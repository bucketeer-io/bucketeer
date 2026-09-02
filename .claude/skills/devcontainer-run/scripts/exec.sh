#!/usr/bin/env bash
# Locate the Bucketeer dev container and run a command inside it.
#
# Usage:
#   exec.sh status              Report where the dev container is and its health
#   exec.sh <command ...>       Run a command in /workspaces/bucketeer inside it
#
# Exit codes: 0 success, 2 no dev container found, otherwise the command's exit code.
set -euo pipefail

WORKDIR=/workspaces/bucketeer
# go-tools live in a persistent volume; not on PATH in non-login shells
SETUP_PATH='export PATH=/home/codespace/go-tools/bin:$PATH'

MODE=""
CID=""
CODESPACE=""

detect() {
  # Case 1: this shell is already inside the dev container
  if [ -d "$WORKDIR" ] && [ "$(id -un)" = "codespace" ]; then
    MODE=inside
    return
  fi

  # Case 2: local devcontainer (VS Code "Reopen in Container" / devcontainer CLI).
  # Derive the repo root from this script's location (<repo>/.claude/skills/devcontainer-run/scripts/)
  # so the exact label match works from any cwd — no fuzzy fallback that could pick
  # the wrong container when multiple checkouts are running.
  local repo_root
  repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../../../.." && pwd)"
  CID="$(docker ps -q --filter "label=devcontainer.local_folder=$repo_root" 2>/dev/null | head -1 || true)"
  if [ -n "$CID" ]; then
    MODE=local
    return
  fi

  # Case 3: GitHub Codespace (requires `gh` with the codespace scope).
  # Never silently pick one: a fork named `bucketeer`, or several codespaces on
  # this repo, would otherwise run commands in the wrong clone/branch/cluster.
  local candidates
  candidates="$(gh codespace list --json name,repository,state \
    -q '.[] | select((.repository | endswith("/bucketeer")) and .state == "Available") | .name' 2>/dev/null \
    || true)"
  if [ -n "${BUCKETEER_CODESPACE:-}" ]; then
    if ! echo "$candidates" | grep -qx -- "$BUCKETEER_CODESPACE"; then
      echo "BUCKETEER_CODESPACE='$BUCKETEER_CODESPACE' is not an available Bucketeer codespace." >&2
      echo "Available:" >&2
      echo "$candidates" | sed 's/^/  - /' >&2
      exit 2
    fi
    CODESPACE="$BUCKETEER_CODESPACE"
    MODE=codespace
    return
  fi
  local count
  count="$(echo "$candidates" | grep -c . || true)"
  if [ "$count" -gt 1 ]; then
    echo "Multiple available Bucketeer codespaces found — refusing to guess:" >&2
    echo "$candidates" | sed 's/^/  - /' >&2
    echo "" >&2
    echo "Pick one explicitly:" >&2
    echo "  export BUCKETEER_CODESPACE=<name>" >&2
    exit 2
  fi
  if [ "$count" -eq 1 ]; then
    CODESPACE="$candidates"
    MODE=codespace
    return
  fi
}

run_inside() {
  local cmd="$1"
  case "$MODE" in
    inside)
      bash -c "cd $WORKDIR && $SETUP_PATH && $cmd"
      ;;
    local)
      docker exec -u codespace -w "$WORKDIR" "$CID" bash -c "$SETUP_PATH && $cmd"
      ;;
    codespace)
      gh codespace ssh -c "$CODESPACE" -- "cd $WORKDIR && $SETUP_PATH && $cmd"
      ;;
  esac
}

not_found() {
  echo "No running Bucketeer dev container found." >&2
  echo "" >&2
  echo "Checked: this shell, local devcontainers (docker label devcontainer.local_folder)," >&2
  echo "and GitHub Codespaces (gh codespace list)." >&2
  echo "" >&2
  echo "To start one:" >&2
  echo "  - VS Code: 'Dev Containers: Reopen in Container' on this repo" >&2
  echo "  - CLI:     devcontainer up --workspace-folder ." >&2
  echo "  - Codespace: gh codespace create -R bucketeer-io/bucketeer" >&2
  if ! gh codespace list >/dev/null 2>&1; then
    echo "" >&2
    echo "Note: 'gh codespace list' failed — if you use Codespaces, grant the scope with:" >&2
    echo "  gh auth refresh -h github.com -s codespace" >&2
  fi
  exit 2
}

status_report() {
  case "$MODE" in
    inside)    echo "mode: inside (this shell is already in the dev container)" ;;
    local)     echo "mode: local devcontainer (docker exec, container $CID)"
               echo "workspace: bind-mounted from the host — file changes appear in the host repo directly" ;;
    codespace) echo "mode: GitHub Codespace '$CODESPACE' (gh codespace ssh)"
               echo "workspace: SEPARATE clone — changes made inside do NOT appear in the host repo" ;;
  esac
  echo "---"
  run_inside '
    echo "user: $(id -un)  workdir: $(pwd)"
    echo "protoc: $(protoc --version 2>/dev/null || echo MISSING)"
    command -v mockgen >/dev/null && echo "go-tools: OK" || echo "go-tools: MISSING (run bash .devcontainer/setup.sh)"
    docker info >/dev/null 2>&1 && echo "dockerd: running" || echo "dockerd: NOT running (start with: nohup sudo dockerd > /tmp/dockerd.log 2>&1 &)"
    if minikube status >/dev/null 2>&1; then
      # Always query the minikube context explicitly: the active context may point
      # somewhere else entirely (a real GKE cluster), and reporting its pods here
      # would be exactly the confusion this wrapper exists to prevent.
      ctx=$(kubectl config current-context 2>/dev/null || echo unknown)
      [ "$ctx" != "minikube" ] && echo "kube-context: WARNING active context is \"$ctx\", not minikube — bare kubectl/helm commands would hit that cluster; pass --context minikube / --kube-context minikube"
      # Ignore transient states (Pending/ContainerCreating/Init) — batch CronJobs
      # constantly spawn short-lived pods and would make the count flap.
      if pods=$(kubectl --context minikube get pods --no-headers 2>/dev/null); then
        total=$(echo "$pods" | grep -c . || true)
        failing=$(echo "$pods" | grep -cE "CrashLoopBackOff|ImagePull|ErrImage|Error|OOMKilled|Evicted" || true)
        echo "minikube: running ($total pods, $failing failing)"
        [ "$failing" -gt 0 ] && echo "$pods" | grep -E "CrashLoopBackOff|ImagePull|ErrImage|Error|OOMKilled|Evicted"
      else
        echo "minikube: running, but kubectl failed to list pods in the minikube context — check kubectl config (kubectl config get-contexts)"
      fi
    else
      echo "minikube: NOT running (start with: make start-minikube — never minikube start)"
    fi
    echo "git: $(git status --porcelain | wc -l | tr -d " ") modified files on branch $(git branch --show-current)"
  '
}

detect
[ -z "$MODE" ] && not_found

if [ "$#" -eq 0 ] || [ "$1" = "status" ]; then
  status_report
else
  run_inside "$*"
fi
