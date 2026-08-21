---
name: devcontainer-generate
description: >-
  Regenerate Bucketeer protobuf Go bindings, OpenAPI/Swagger specs, and gomock
  files inside the dev container, where protoc is guaranteed to be exactly
  v23.4. Use this whenever a .proto file changed, generated *.pb.go /
  *.pb.gw.go / swagger files need regenerating, a mocked Go interface changed
  (mockgen), or the user says "generate proto", "regen protos", "make
  proto-all", "make mockgen", or "devcontainer-generate". Prefer this over running
  protoc or make proto-all on the host — a host protoc version mismatch
  silently rewrites every generated file's header.
---

# devcontainer-generate — code generation inside the dev container

Generated files are committed to the repo, and their headers record the protoc
version (`protoc v4.23.4`). The dev container ships exactly protoc 23.4, so
generation must happen there; a different host protoc churns every `.pb.go`
file and the PR becomes unreviewable. (Human-facing background:
`DEVELOPMENT.md` § "Working with the dev container from the host".)

All commands below go through the devcontainer-run wrapper (see the `devcontainer-run` skill for how
detection works):

```bash
DEVC="bash .claude/skills/devcontainer-run/scripts/exec.sh"
```

## 1. Pick the right target

| What changed | Target |
|---|---|
| `.proto` files | `make proto-all` |
| A Go interface that has generated mocks in a `mock/` dir | `make mockgen` |
| Both, or unsure | `make generate-all` |

## 2. Run it

```bash
$DEVC 'make proto-all'      # or make mockgen / make generate-all
```

This is minutes-long; use a generous Bash timeout (600000). If the run fails
on protolock, the `.proto` change broke backward compatibility — read the
error; don't force it without flagging the compatibility break to the user.

## 3. Verify before declaring success

- `git status --porcelain` (on the host for a local devcontainer; inside via
  `$DEVC 'git status --porcelain'` for a codespace) — the changed files should
  be only the ones related to your proto/interface change plus their generated
  outputs. **A diff touching every `.pb.go` in the repo means a wrong protoc
  version — abort and check `$DEVC 'protoc --version'`.**
- Spot-check one regenerated file's header still says `protoc v4.23.4`:
  `grep -m1 "protoc " proto/<domain>/<file>.pb.go`
- Build still compiles: `$DEVC 'make build-go'` (or the affected
  `make build-<service>`), and `$DEVC 'make gofmt'` after any Go changes.

## Codespace caveat

In a codespace the regenerated files land in the codespace's clone, not the
host repo. Commit/push from inside, or copy back with `gh codespace cp`. The
`status` subcommand of the devcontainer-run script tells you which mode you're in.
