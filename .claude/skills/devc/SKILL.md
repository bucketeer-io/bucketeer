---
name: devc
description: >-
  Detect the running Bucketeer dev container (local VS Code devcontainer or
  GitHub Codespace) and run commands inside it. Use this whenever a task should
  run in the dev container environment — make targets, builds, tests, kubectl /
  helm / minikube commands, checking whether the container is up — or when the
  user says "devc", "dev container", "devcontainer", "codespace", or "run this
  inside the container". Also use it when a task needs tools the container
  guarantees but the host may lack (protoc 23.4, mockgen, protolock, helm,
  kubectl, minikube). devc-generate and devc-deploy build on this skill.
---

# devc — run commands inside the Bucketeer dev container

The dev container is the canonical Bucketeer development environment: Ubuntu with
docker-in-docker, minikube + helm + kubectl, protoc v23.4, and Go tooling in
`/home/codespace/go-tools/bin` (a persistent volume, NOT on PATH in plain
non-login shells). The workspace is `/workspaces/bucketeer`, the user is
`codespace` (passwordless sudo).

## How to run anything inside it

Always go through the wrapper script — it finds the container and sets up PATH:

```bash
# Where is the container, and is the environment healthy?
bash .claude/skills/devc/scripts/devc-exec.sh status

# Run any command in /workspaces/bucketeer inside the container
bash .claude/skills/devc/scripts/devc-exec.sh 'make build-api'
bash .claude/skills/devc/scripts/devc-exec.sh 'kubectl get pods'
```

Detection order (the script handles all of this):
1. Already inside the container (`/workspaces/bucketeer` exists, user `codespace`) → run directly.
2. Local devcontainer → `docker ps` filtered by label `devcontainer.local_folder=<repo root>`, exec via `docker exec`.
3. GitHub Codespace → `gh codespace list` (needs the `codespace` auth scope), exec via `gh codespace ssh`.

Exit code 2 means no container was found; the script prints how to start one.
Don't fall back to running the command on the host in that case — tell the user
and let them choose, because host tool versions (especially protoc) may differ.

## Local devcontainer vs Codespace — the one difference that matters

- **Local devcontainer**: `/workspaces/bucketeer` is a bind mount of the host
  repo. Files generated inside appear in the host working tree immediately.
- **Codespace**: a separate clone. Generated or edited files stay in the
  codespace. To get them back: commit and push from inside, or
  `gh codespace cp 'remote:/workspaces/bucketeer/<path>' <local-path>`.
  Always tell the user which mode you're in when file changes are involved
  (`status` prints it).

## Environment facts and gotchas

- Long commands (image builds, deploys) can take many minutes — use a generous
  Bash timeout (600000) or `run_in_background`.
- `dockerd` inside the container is started by the post-attach hook, but that
  only fires when an editor attaches. If `status` says it's not running:
  `bash .claude/skills/devc/scripts/devc-exec.sh 'nohup sudo dockerd > /tmp/dockerd.log 2>&1 & sleep 5 && docker info > /dev/null && echo ok'`
- minikube must be started with `make start-minikube`, never `minikube start`
  directly (the make target restores the cluster config and localenv services).
  Note: `make start-minikube` intentionally **exits 1 if minikube is already
  running** — check `minikube status` first instead of treating that as failure.
- `web-gateway.bucketeer.io` / `api-gateway.bucketeer.io` resolve via the
  container's own `/etc/hosts` (pointed at `minikube ip`). Health checks with
  curl against those hosts must run *inside* the container, not on the host.
- If go-tools are missing or permissions look broken, the fix is the setup
  script: `bash .devcontainer/setup.sh` (idempotent, cache-aware).
- **Never run kubectl/helm bare on the host for dev work.** The host's kubectl
  context may point at a real GKE cluster, not minikube — always go through the
  wrapper so commands hit the cluster inside the container.
- The host may also run a docker-compose Bucketeer stack in parallel
  (`docker-compose/compose.yml`). That is a different environment — this skill
  is only about the dev container / minikube world.
