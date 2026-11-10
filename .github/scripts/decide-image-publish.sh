#!/usr/bin/env bash
# Decide whether a hack-tool image should be published.
#
# Usage: decide-image-publish.sh <makefile>
# Reads TAG from the Makefile unless INPUT_TAG is set.
# PACKAGE_NAME is the ghcr.io package name (without the registry or owner).
# Writes tag and push ("true" or "false") to GITHUB_OUTPUT.
#
# Publishes when the tag is a newer vX.Y.Z than the newest published release.
# A matching tag is left unchanged. An older tag fails.
set -euo pipefail

makefile=${1:?makefile path is required}
: "${PACKAGE_NAME:?PACKAGE_NAME is required}"
: "${GITHUB_OUTPUT:?GITHUB_OUTPUT is required}"
: "${GITHUB_REPOSITORY_OWNER:?GITHUB_REPOSITORY_OWNER is required}"

if [ -n "${INPUT_TAG:-}" ]; then
  TAG="$INPUT_TAG"
else
  TAG=$(grep -E '^TAG=' "$makefile" | sed 's/^TAG=//')
fi
if ! [[ "$TAG" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
  echo "Tag must look like v1.2.3, got: ${TAG}" >&2
  exit 1
fi

err=$(mktemp)
set +e
published=$(gh api --paginate \
  "/orgs/${GITHUB_REPOSITORY_OWNER}/packages/container/${PACKAGE_NAME}/versions" \
  --jq '.[] | .metadata.container.tags[]?' 2>"$err")
status=$?
set -e
if [ "$status" -ne 0 ]; then
  if grep -q 'HTTP 404' "$err"; then
    published=""
  else
    cat "$err" >&2
    exit 1
  fi
fi
latest=$(printf '%s\n' "$published" | grep -E '^v[0-9]+\.[0-9]+\.[0-9]+$' | sort -V | tail -n 1 || true)

push=true
if [ -z "$latest" ]; then
  echo "No published release tag. Publishing ${TAG}."
elif [ "$TAG" = "$latest" ]; then
  echo "Published tag ${latest} matches ${TAG}. Leaving the image unchanged."
  push=false
else
  newer=$(printf '%s\n%s\n' "$latest" "$TAG" | sort -V | tail -n 1)
  if [ "$newer" != "$TAG" ]; then
    echo "Tag ${TAG} is older than the published tag ${latest}." >&2
    exit 1
  fi
  echo "Tag ${TAG} is newer than the published tag ${latest}."
fi
echo "tag=${TAG}" >> "$GITHUB_OUTPUT"
echo "push=${push}" >> "$GITHUB_OUTPUT"
