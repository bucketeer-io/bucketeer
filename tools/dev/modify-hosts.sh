#!/bin/sh
# Point the local gateway names at the minikube node in /etc/hosts.
#
# Usage: modify-hosts.sh <minikube-ip>
#
# The local gateways use the reserved .test TLD (RFC 2606) so local development
# never names the production domain. Idempotent, so deploy-bucketeer runs it on
# every deploy and an existing dev container picks up renamed hosts or a new
# minikube IP without a restart:
#
# - When web-gateway.bucketeer.test and api-gateway.bucketeer.test already map
#   to the given IP and nothing maps the former *.bucketeer.io names, the file is
#   left untouched and sudo is not invoked.
# - Otherwise only the gateway hostname tokens (current and former names) are
#   removed from the rows they appear on; the address, other aliases and the
#   comment of a shared row are kept, and a row that listed nothing but gateway
#   names is dropped. The two current entries are then appended.
#
# The new content is computed first and written back in place (tee under sudo
# when the file is not writable), because /etc/hosts is a bind mount inside the
# dev container that sed -i cannot replace. HOSTS_FILE overrides the path for
# tests.
set -eu

ip="${1:?usage: $0 <minikube-ip>}"
hosts="${HOSTS_FILE:-/etc/hosts}"

current='web-gateway.bucketeer.test api-gateway.bucketeer.test'
former='web-gateway.bucketeer.io api-gateway.bucketeer.io'
managed="$current $former"

# Address/name pairs the file currently holds for managed names (comments ignored).
have="$(awk -v managed="$managed" '
  BEGIN { n = split(managed, m, " "); for (i = 1; i <= n; i++) is[m[i]] = 1 }
  {
    code = $0; sub(/#.*/, "", code)
    nt = split(code, t, /[ \t]+/); c = 0
    for (i = 1; i <= nt; i++) if (t[i] != "") tk[++c] = t[i]
    for (i = 2; i <= c; i++) if (tk[i] in is) print tk[1], tk[i]
  }
' "$hosts" | sort)"
want="$(for name in $current; do echo "$ip $name"; done | sort)"

if [ "$have" = "$want" ]; then
  echo "$hosts already maps the local gateways to $ip"
  exit 0
fi

if [ -n "$have" ]; then
  echo "$hosts: replacing local gateway entries:"
  echo "$have" | sed 's/^/  /'
fi
echo "$hosts: mapping web-gateway.bucketeer.test and api-gateway.bucketeer.test to $ip"

new="$(awk -v managed="$managed" -v ip="$ip" -v current="$current" '
  BEGIN { n = split(managed, m, " "); for (i = 1; i <= n; i++) is[m[i]] = 1 }
  {
    idx = index($0, "#")
    if (idx > 0) { code = substr($0, 1, idx - 1); comment = substr($0, idx) } else { code = $0; comment = "" }
    nt = split(code, t, /[ \t]+/); c = 0
    for (i = 1; i <= nt; i++) if (t[i] != "") tk[++c] = t[i]
    hit = 0; kept = ""
    for (i = 2; i <= c; i++) {
      if (tk[i] in is) hit = 1
      else kept = kept (kept == "" ? "" : " ") tk[i]
    }
    if (!hit) { print; next }      # no gateway name on this row: keep it byte for byte
    if (kept == "") next           # the row listed nothing but gateway names: drop it
    print tk[1] "\t" kept (comment == "" ? "" : "\t" comment)
  }
  END { n = split(current, names, " "); for (i = 1; i <= n; i++) print ip "\t" names[i] }
' "$hosts")"

if [ -w "$hosts" ]; then
  printf '%s\n' "$new" > "$hosts"
else
  printf '%s\n' "$new" | sudo tee "$hosts" >/dev/null
fi
