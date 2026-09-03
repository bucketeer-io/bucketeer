---
name: devcontainer-deploy
description: >-
  Build and deploy Bucketeer to the minikube cluster inside the dev container,
  or redeploy/restart a single service there. Use this whenever the user wants
  to deploy locally, run "make deploy-bucketeer", start minikube, get their
  code changes running in the dev cluster, restart a crashing pod, or says
  "devcontainer-deploy", "deploy to minikube", "redeploy the backend". Also use it to
  check deployment health (pods not ready, gateway not responding) in the dev
  container environment.
---

# devcontainer-deploy — deploy Bucketeer inside the dev container

Deployment target is the minikube cluster *inside* the dev container (helm
charts in `manifests/`), not the host docker-compose stack. (Human-facing
docs for these flows: `DEVELOPMENT.md` § "Deploy Bucketeer".) All commands go
through the devcontainer-run wrapper (see the `devcontainer-run` skill):

```bash
DEVC="bash .claude/skills/devcontainer-run/scripts/exec.sh"
```

## 1. Preflight

```bash
$DEVC status
```

- `dockerd` not running → start it (command is in the status output) — image
  builds need it.
- `minikube` not running → `$DEVC 'make start-minikube'`. Never `minikube start`
  directly. If minikube IS already running, skip this: the target intentionally
  exits 1 with "minikube is already running" — that is not an error to fix.

## 2. Full deploy

```bash
$DEVC 'make deploy-bucketeer'
```

What it does, so failures are diagnosable: uninstalls the existing `bucketeer`
helm release → regenerates cert/token/oauth secrets → builds all Go binaries →
builds docker images with `TAG=localenv` → loads them into minikube → helm
install/upgrade `localenv` (MySQL, Redis, Pub/Sub emulator, optionally
Postgres/BigQuery emulator) → helm install `bucketeer` with
`manifests/bucketeer/values.dev.yaml`.

This takes many minutes. Run it with a 600000 timeout or `run_in_background`
and monitor. Postgres/BigQuery enablement is auto-detected from
`dataWarehouse` in `values.dev.yaml` — don't set it manually, but remember the
invariant: `web` and `subscriber` must use the same event store, so data
warehouse changes belong in `values.dev.yaml`, not ad-hoc helm flags.

## 3. Single service — faster than a full deploy

For a code change to one service (e.g. backend):

```bash
$DEVC 'make build-go-embed && TAG=localenv make build-docker-images && TAG=localenv make minikube-load-images'
$DEVC 'kubectl --context minikube rollout restart deployment <name> && kubectl --context minikube rollout status deployment <name>'
```

The Bucketeer deployments are `api`, `web`, `batch-server`, and `subscriber`
(confirm with `$DEVC 'kubectl --context minikube get deployments'`). Chart-level
changes instead:

```bash
$DEVC 'helm upgrade bucketeer manifests/bucketeer/ --kube-context minikube --values manifests/bucketeer/values.dev.yaml'
```

Always go through `$DEVC` and always name the context: a bare `helm`/`kubectl`
runs against whatever context is currently active, which on the host is often a
real cluster. `$DEVC status` warns when the active context is not `minikube`.

## 4. Verify

```bash
$DEVC 'kubectl --context minikube get pods'   # everything Running/Completed, restarts not climbing
$DEVC 'curl -sk https://api-gateway.bucketeer.io/health'   # must run INSIDE the container
```

The `*.bucketeer.io` hosts entries live in the container's `/etc/hosts`
(pointed at `minikube ip`) — curl from the host proves nothing. For a failing
pod: `$DEVC 'kubectl --context minikube logs deploy/<name> --tail=100'` and
`$DEVC 'kubectl --context minikube describe pod <pod>'`.

## Related dev-cluster chores

- Bootstrap e2e accounts (after a fresh deploy, before e2e tests):
  `$DEVC 'make create-dev-container-e2e-accounts'`
- Wipe e2e data: `$DEVC 'make delete-dev-container-mysql-data'` (or the
  `-postgres-` variant). These are destructive — confirm with the user first.
- MySQL from inside the container: host `$(minikube ip)`, port 32000,
  user/pass `bucketeer`, db `bucketeer`.
