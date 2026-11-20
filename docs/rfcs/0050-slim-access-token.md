# Slim Access Token

This RFC reduces the console access token to the caller's identity. Roles,
system admin status, the display name and the organization are resolved on
the server for each request. Services and public API requests, which share
one token with users today, get their own credentials.

- Status: Draft
- Authors: Alessandro Yuichi Okimoto
- Related: RFC 0049 (centralized authorization, #3006), PR #3015 (authorization engine)

## 1. Background

The console signs in with Google and receives two JWTs signed with the same
key: an access token (10 minutes by default) and a refresh token (7 days by
default). When this RFC was written the access token carried more than the
caller's identity. PR #3021 has since removed `organization_id` and `name`
from the access token and `organization_id` from the refresh token (§6, step
3); the other claims are unchanged:

| Claim | Read by | Problem |
| --- | --- | --- |
| `email` | Every authorization check | The only identity. There is no `sub`. |
| `organization_id` | No backend code; only the dashboard's organization picker | `RefreshToken` writes the first organization with an enabled account, not the one the user selected, so the claim can name a different organization from the one in use. |
| `name` | `pkg/role`, for the editor in audit logs | Names are stored per account. The token carries the name from the organization chosen at sign-in, which may not be the one being edited. |
| `is_system_admin` | `pkg/role` and a few handlers | An authorization decision fixed for the lifetime of the token. Its derivation also looks wrong (see below). |
| `is_service_token` | `pkg/role` and a few handlers | A privilege flag inside the same token type users get. |

An account is one row of `account_v2`, keyed by `(email, organization_id)`,
so one email may hold accounts in several organizations. Anything in the
token that depends on the organization is therefore either a guess or has to
be reissued on every switch. The backend never read the claim: every request
names its organization or environment, and the account is looked up for that
request.

Three further problems were found while reviewing the token:

- **System admin derivation.** `getMyOrganizations` checks for a system
  admin organization before it leaves out disabled accounts, and sign-in then
  accepts an enabled account in any organization. A user whose system admin
  account is disabled but who has an enabled account elsewhere appears to
  receive `is_system_admin: true`. This was found by reading the code and has
  not been reproduced.
- **Service token.** `hack/generate-service-token` issues an access token
  valid for 100 years with both flags set. It has no identifier, so
  cancelling one means rotating the signing key, which signs out every user.
  One token is shared by the web, API, batch and subscriber services, the
  batch CronJobs and the operator tools, so a leak cannot be traced to a
  holder and replacing it touches all of them. Every internal call uses it,
  including the calls a service makes while handling a user's request.
- **Public API relay.** The API service checks the API key, then calls the
  web services with the shared service token. On every write and a few reads
  it adds three metadata fields: the obfuscated key, the maintainer's email
  and the key name. The web services accept those fields on any service
  token and record the maintainer as the editor, so every holder of the
  service token can attribute a change to any user who has the role for it.
  The other reads carry the service token alone, so the web services give
  them full access whoever the maintainer is. The key's own role and
  environment are checked by the API service only.
- **Session lifetime.** `SwitchOrganization` accepts an access token and
  returns a new refresh token, so an access token that leaks on its own can
  be turned into lasting access. Each refresh returns a refresh token valid
  for another 7 days and the previous one stays valid, so a stolen refresh
  token can be renewed for as long as the account is enabled.

RFC 0049 replaces `pkg/role` with one authorization engine. In PR #3015 its
caller is the email, a service flag and the relayed API key, and roles, the
name and system admin status are read from the database. This RFC is the
token side of that change.

## 2. Goals / Non-Goals

Goals

- The access token says who the caller is and when they signed in, and
  nothing else.
- Service identities get their own token type. User tokens carry no
  privilege flags.
- Each service has its own token, which names the service and can be
  cancelled without rotating the signing key.
- A service token works only inside the cluster. It is refused when it
  arrives through the public address.
- Of the deployments, only web holds the key that signs user tokens.
- A public API request reaches the web services with one signed credential
  that names the API key and its maintainer. No request metadata is trusted
  for identity.
- A change of role, system admin status or name applies on the next request,
  not when the token expires.
- An access token that leaks on its own stops working when it expires.
- A stolen refresh token stops working after a maximum sign-in age.
- No new storage. Users sign in and switch organization as today, and sign
  in again once the maximum sign-in age has passed.
- A rolling upgrade and a rollback keep working at every step.

Non-Goals

- Server-side sessions: sign-out that revokes, immediate revocation and
  refresh token rotation (§7).
- Moving the refresh token out of local storage into an HttpOnly cookie (§7).
- A user table with a stable ID. `sub` is the email until one exists.
- Changing `exp` and `iat` to numeric dates (they are written as RFC 3339
  strings today).
- The authorization rules themselves, which are RFC 0049. That includes
  limiting which RPCs each service may call: this RFC gives every service a
  name, and a service token keeps full access until RFC 0049 has rules for
  it.
- Network controls other than refusing service tokens on the public
  listener: removing the public service routes, restricting the service
  ports, and mutual TLS between services (§7).

## 3. Design Overview

```mermaid
flowchart LR
    Token[Token] -->|"who is calling, and its kind"| Caller[Caller]
    Request[Request] -->|"environment_id or organization_id"| Scope[ScopeResolver]
    Caller --> Authorizer
    Scope --> Authorizer
    Accounts[("account_v2")] -->|"roles, name, system admin"| Authorizer
    Authorizer -->|allow| Handler
    Authorizer -->|deny| Denied[PermissionDenied]
```

1. **The token identifies the caller.** It supplies the email of a user, the
   name of a service, or the API key and its maintainer.
2. **The request names the organization.** RFC 0049's Table 2 says which
   field of each RPC holds the environment or organization.
3. **The database supplies the rest.** Roles, the name and system admin
   status come from `account_v2` for that request.

Every request that reaches a web service carries exactly one of three
tokens. All are verified by the same interceptor on the gRPC server, so the
REST gateway, direct gRPC and calls from other services are checked in one
place.

```mermaid
flowchart LR
    Console[Console] -->|"access token"| Web[Web services]
    Client[Public API client] -->|"API key"| API[API service]
    API -->|"delegated token"| Web
    API -->|"service token of api"| Web
    Jobs[Batch and subscriber] -->|"service token of each"| Web
    Tools[Operator tools] -->|"service token of each"| Web
```

| Caller | Before this RFC | Proposed |
| --- | --- | --- |
| A user in the console | Access token with the organization, name and two flags | Access token with the identity only |
| A service acting on its own, including web services calling each other | One shared access token with both flags, valid for 100 years | A service token per service |
| The API service acting for an API key | The shared service token plus metadata naming the key and its maintainer | A delegated token issued for that key |

The access token changes as follows; the other two are described in §4.

| Claim | Before this RFC | Proposed |
| --- | --- | --- |
| `iss`, `aud`, `exp`, `iat` | yes | yes |
| `sub` | no | the email |
| `email` | yes | yes |
| `token_use` | `access` | `access` |
| `auth_time` | no | the time of the original sign-in |
| `organization_id` | yes | removed (done in PR #3021) |
| `name` | yes | removed (done in PR #3021) |
| `is_system_admin` | yes | removed |
| `is_service_token` | yes | removed; service tokens are their own type |

## 4. Design Details

### Access token

```go
type AccessToken struct {
	Issuer   string    `json:"iss"`
	Subject  string    `json:"sub"`
	Audience string    `json:"aud"`
	Expiry   time.Time `json:"exp"`
	IssuedAt time.Time `json:"iat"`
	Email    string    `json:"email"`
	AuthTime time.Time `json:"auth_time"`
	TokenUse string    `json:"token_use"`
}
```

`sub` identifies the person. It is the email today and becomes a stable user
ID if a user table is added. `email` stays because verifiers of the previous
version require it and handlers read it. There is no `jti`: nothing stores or
looks up user tokens.

`auth_time` is the time of the original sign-in, copied unchanged on every
refresh. It is not used to authorize. It is there because
`SwitchOrganization` issues a refresh token from an access token, and
without it that RPC would start the sign-in age again (see Refresh token).

### Service token

```go
type ServiceToken struct {
	Issuer   string    `json:"iss"`
	Subject  string    `json:"sub"`
	Audience string    `json:"aud"`
	Expiry   time.Time `json:"exp"`
	IssuedAt time.Time `json:"iat"`
	Email    string    `json:"email,omitempty"`
	ID       string    `json:"jti"`
	TokenUse string    `json:"token_use"` // "service"
}
```

- `sub` names the service the token was issued to: `web`, `api`, `batch`,
  `subscriber`, or a name chosen for an operator tool.
  `hack/generate-service-token` gains a `--service` flag for it.
- `email` is what `--email` sets today. It stays the editor in audit logs,
  so existing entries and the way they are shown do not change.
- **One token per service.** Each Helm chart already reads its own
  `serviceToken.secret`, so a deployment gives each chart a different Secret.
  The Docker Compose and development setups generate one file per service. A
  CronJob uses the token of its chart.
- **Cancelling a token.** A server option gives, for each service name, a
  time before which its tokens are not accepted, and the verifier rejects a
  service token whose `iat` is earlier. A leaked token is cancelled by
  issuing a new token to that one service and then setting the time to the
  new token's `iat`. Users stay signed in and the other services are not
  touched. One value per service covers every copy of the old token, needs
  no token ID and never grows. Rotation is the same two actions, so it also
  ends the old token. The order matters: the time is set only after every
  pod of that service runs with the new token. Both tokens are valid until
  then, so a rolling restart interrupts no call; setting the time earlier
  rejects the calls of the pods that still hold the old token, which is the
  price of cutting off a leaked token at once. A separate value applies to
  legacy tokens, which carry no name. `jti` stays as the identifier of one
  token in logs and is not used to cancel.
- **Lifetime.** The generator gains a `--ttl` flag. A token issued with
  `--service` lasts one year by default, so the tokens issued in step 6 of
  §6 do not inherit today's 100 years. A service reads its token at startup,
  so rotating it means replacing the Secret and restarting. From the release
  in which the per-service tokens are issued, a service logs a warning and
  reports a metric when its own token is close to expiry, so none expires
  unnoticed. Legacy tokens already issued keep their expiry until the legacy
  form is dropped.
- **Signing key.** A deployment may give service tokens their own key pair.
  The generator then signs with that private key, and every verifier checks
  service tokens against its public key: a token verified with the service
  token key is accepted only as a service token, and a `service` token
  signed with the shared key is rejected. Legacy service tokens are still
  verified with the shared key. When no service token key is configured,
  service tokens use the shared key, so a deployment that does not need the
  separation adds no Secret.
  - The private key is needed only where tokens are generated. It does not
    have to be stored in the cluster, and the key that signs user tokens no
    longer has to leave the web deployment to issue a service token.
  - Replacing it means issuing the service tokens again. Users stay signed
    in. A verifier accepts more than one service token public key, so the
    key is replaced without interrupting calls between services: add the new
    public key to every verifier, issue and deploy the new tokens, then
    remove the old public key.
  - It does not limit what a leaked key allows: whoever holds it can sign a
    service token with full access, so it is protected like the shared key.
  - A deployment configures it before issuing the per-service tokens.
    Adding it later means issuing them again.
- **Access.** A service token keeps the full access it has today. Naming the
  service is what allows RFC 0049 to later give each one only the RPCs it
  calls.
- The verifier dispatches on `token_use`. An access token that has both
  `is_system_admin` and `is_service_token` is a legacy service token and is
  treated as a service identity, whether its `token_use` is `access` (the
  current generator stamps it) or absent (older tokens). Otherwise the access
  token verifier rejects a service token, and the reverse, as the existing
  token kinds already do.
- While `pkg/role` exists, the auth interceptor presents a service token to
  it as an access token with both legacy flags set, so services that are not
  yet migrated to RFC 0049 behave as today.

### Delegated token

A public API request acts for the maintainer of its API key. Today the API
service says so in request metadata next to the shared service token. The
metadata is replaced by a token that the auth service signs for that key:

```go
type DelegatedToken struct {
	Issuer         string    `json:"iss"`
	Subject        string    `json:"sub"`   // the maintainer
	Audience       string    `json:"aud"`
	Expiry         time.Time `json:"exp"`
	IssuedAt       time.Time `json:"iat"`
	Email          string    `json:"email"` // the maintainer
	TokenUse       string    `json:"token_use"` // "delegated"
	APIKeyID       string    `json:"api_key_id"`
	APIKeyName     string    `json:"api_key_name"`
	APIKeyToken    string    `json:"api_key_token"` // obfuscated, as in audit logs
	APIKeyRole     string    `json:"api_key_role"`
	EnvironmentID  string    `json:"environment_id"`
	OrganizationID string    `json:"organization_id"`
}
```

```mermaid
sequenceDiagram
    participant P as PublicAPIClient
    participant G as APIService
    participant A as AuthService
    participant W as WebServices

    P->>G: Request with the API key
    G->>G: Check the key, its role and its environment
    G->>A: Exchange the API key, as the api service
    A-->>G: Delegated token for the maintainer, valid 5 minutes
    Note over G: Kept per key until it expires
    G->>W: Request with the delegated token
    W-->>G: Response, authorized and logged as the maintainer
    G-->>P: Response
```

- **Issuing.** A new auth service RPC takes the API key and requires the
  `api` service token. It looks the key up and refuses a disabled key, a
  disabled environment and a key whose role is not a public API role. The
  signing key stays in the web deployment.
- **Caching.** The API service keeps the token per key until shortly before
  it expires, next to the API keys it already caches. That is about one
  extra call per key every 5 minutes.
- **Verifying.** The web interceptor builds the same caller that
  `Caller.Relayed` builds from the metadata in PR #3015: the maintainer,
  with the API key attached. The maintainer's account and role are still
  read from the database on every request, and the editor in audit logs is
  the same as today.
- **Key role and scope.** A delegated token is a bearer credential, so the
  web services do not rely on the API service having checked the key. Each
  RPC that the public API uses is listed with the lowest key role that may
  call it and with its scope. An RPC that is not listed is denied to a
  delegated token, so the token reaches nothing the public API does not
  offer, whatever the maintainer's own role is.
  - *Environment scope:* the request's environment must be the token's
    `environment_id`.
  - *Organization scope:* the request must resolve to the token's
    `organization_id`. This covers RPCs that take an organization, and RPCs
    that take an environment only to find the organization, such as
    `GetAccountV2ByEnvironmentID`, where the public API accepts any
    environment of the key's organization today.

  The list is filled from the checks the API service makes today. Those
  checks stay, and a test fails when the two differ. When several public API
  methods use the same RPC, the lowest of their roles applies. The list
  becomes columns of RFC 0049's Table 2 when the API key roles move there.
- **Which calls use it.** Every call that carries out a public API request,
  reads included. The API service's own token is left for the calls it makes
  for itself: loading flags and segments for the SDK endpoints, the
  exchange, and the lookups behind its own checks, such as finding the
  organization of an environment named in a request.
- **Reads change behaviour.** The reads that carry no metadata today pass
  whoever the maintainer is. With a delegated token they are judged as the
  maintainer, like writes: a key whose maintainer is no longer an enabled
  member, or has no access to what the request names, is denied. This follows
  RFC 0049's rule that a key never grants more than its maintainer has.
  Nothing updates a key when its maintainer's account is disabled or
  removed, so some keys in use may be affected. Those reads are reported
  before they are denied (§6, steps 11 and 12).

The access token carries no organization because a user can switch, and the
claim would be a guess. A delegated token differs: an API key belongs to one
environment of one organization, the token cannot be refreshed, and it lives
5 minutes, so both claims are exact for as long as the token is valid.
Disabling a key or changing its role applies when the cached token expires.

Until the exchange exists, the metadata is accepted only from the `api`
service token (and from the legacy token while it is still accepted), not
from every service token.

### Refresh token

The refresh token gains `auth_time`, the time of the original sign-in. Its
`organization_id` is already removed (PR #3021). The request does not
change: it carries only the refresh token.

`RefreshToken` does not act on any organization, so it needs none. It
answers one question, whether this person may stay signed in:

1. Verify the refresh token and take the email.
2. Require at least one enabled account in any organization, as today.
3. Reject the token when `auth_time` is older than the maximum sign-in age.
   The console then sends the user to sign in again, as it does for any
   rejected refresh.
4. Return a new access token and a new refresh token with the same
   `auth_time`.

A refresh token without `auth_time`, issued before this change, is treated
as signed in at its own `iat`, but only in the release that introduces
`auth_time`. A pod of an older release still issues refresh tokens without
it, so during a rolling update or after a rollback a refresh through such a
pod would start the sign-in age again. From the next release, a refresh
token without `auth_time` is rejected and the user signs in again. That
release is deployed no sooner than the refresh token lifetime after the
first one, so the tokens issued before the change have expired by then.

`SwitchOrganization` also returns a refresh token, and it is given only an
access token. While the RPC exists (§6, step 4) it follows the same rules:
it copies `auth_time` from the access token into the tokens it returns, and
rejects the request when `auth_time` is older than the maximum sign-in age.
An access token without `auth_time` is handled like a refresh token without
it: treated as signed in at its own `iat` in the first release, and rejected
from the next. Without this, a stolen refresh token could be refreshed and
the new access token traded for a refresh token with a new sign-in time.

The maximum sign-in age is a new option next to `accessTokenTTL` and
`refreshTokenTTL`.

### One claims type for every kind

The demo-creation token is already marked `token_use: "demo"`. What it still
has of its own is a Go type and a sign and a verify function, and the service
and delegated tokens would add two more sets. Once the access token is
slimmed, the kinds carry almost the same claims, so they share one type:

```go
type Claims struct {
	Issuer   string    `json:"iss"`
	Subject  string    `json:"sub"`
	Audience string    `json:"aud"`
	Expiry   time.Time `json:"exp"`
	IssuedAt time.Time `json:"iat"`
	Email    string    `json:"email,omitempty"`
	ID       string    `json:"jti,omitempty"`
	AuthTime time.Time `json:"auth_time,omitempty"` // access tokens only
	TokenUse string    `json:"token_use"`
}
```

One function signs and one verifies, each taking the kind. The verifier
checks the signature, issuer, audience and expiry for every kind, and then
the rule of that kind, so "a token of one kind does not verify as another"
is decided in one place:

| Kind | `token_use` | Required claims | Accepted by |
| --- | --- | --- | --- |
| access | `access`, or absent on a token issued before the claim existed | `email` | Every RPC except `CreateDemoOrganization` |
| service | `service`, or the legacy form described under Service token | `sub`, `jti` (not required in the legacy form) | The same RPCs, as a service identity, but not through the public listener |
| delegated | `delegated` | `email`, the API key claims, `environment_id`, `organization_id` | Only the RPCs listed for the public API, within the key's role and scope, as the maintainer acting through the API key, but not through the public listener |
| demo | `demo` | `email` | `CreateDemoOrganization` only |

The demo token on the wire does not change, except that it also carries
`sub`. Demo tokens live 5 minutes, so no older form has to be accepted.

`demo` stays a kind of its own and is not replaced by an access token. It
proves something an access token does not: `ExchangeDemoToken` issues it
only after checking that the email has no organization to sign in to, that
is, no enabled account in an organization that is neither disabled nor
archived. `CreateDemoOrganization` does not repeat that check. If an access
token were accepted there, any existing user could create demo
organizations.

Three things keep their own shape:

- **Go types for handlers.** `AccessToken`, `ServiceToken`, `DelegatedToken`
  and `DemoCreationToken` remain distinct named types over `Claims`, so the
  compiler still rejects a demo token where an access token is expected. The
  structs in the sections above list the claims each kind carries.
- **The delegated token's extra claims.** Its type embeds `Claims` and adds
  the API key and environment claims.
- **The refresh token.** It omits `iss` and `aud` on purpose, for verifiers
  of earlier versions, and would become an opaque value with server-side
  sessions (§7).

The access token keeps its own struct until the last old claims are removed
(§6, step 13), and moves to `Claims` then.

### Which organization a request is for

| Request | Organization comes from | Check |
| --- | --- | --- |
| Carries `environment_id` | The organization that owns the environment | The caller's account in that organization, and its role in the environment |
| Carries `organization_id` | The request | The caller's account in that organization |
| The caller's own data (`GetMyOrganizations`) | None | A verified caller |
| System admin only | None | An enabled account in a system admin organization |

The table describes a user. Three callers are treated differently, as today
and as RFC 0049 keeps them:

| Caller | Check |
| --- | --- |
| System admin | Reads need no account in the organization. Writes need one, like any user. |
| Service token | No account is looked up. The service has full access. |
| Delegated token | The maintainer's account in the organization, as for a user. The RPC must be listed for the key's role, and the request must be within the token's environment or organization, as that RPC's scope says. |

A user may put any organization in a request, and it is accepted only where
the user has an enabled account, apart from system admin reads. A user whose
account is disabled in one organization and enabled in another keeps a valid
token and is rejected by the first organization's check, as today.

```mermaid
sequenceDiagram
    participant C as Console
    participant A as AuthService
    participant W as WebAPI
    participant D as Database

    C->>A: ExchangeToken with the Google code
    A->>D: Accounts of the email
    A-->>C: Access token and refresh token (email, auth_time)
    Note over C: The selected organization is kept in local storage
    C->>W: Request with organization_id or environment_id
    W->>D: Account of the email in that organization
    W-->>C: Response, or denied when not an active member
    Note over C: Switching organization changes local storage only
    C->>A: RefreshToken with the refresh token
    A->>D: Any enabled account of the email
    A-->>C: New tokens, or Unauthenticated past the maximum sign-in age
```

### Switching organization

The console already keeps the selected organization in local storage and
sends it with each request. Switching becomes: store the new organization,
then call `GetMe` for it. The console stops decoding the access token, and
the `SwitchOrganization` RPC is removed.

### System admin

System admin status is read by email from the caller's account in a system
admin organization, and that account must be enabled. PR #3015 already does
this in the engine. Until every service uses the engine, `pkg/role` still
trusts the claim, so its derivation is fixed first: disabled accounts are
left out before the system admin check.

### Service and delegated tokens on the public listener

The Ingress sends every request from outside to the public listener of the
web deployment's Envoy (port 9003). That listener serves the REST gateway,
and it also routes the `/bucketeer.<package>.<Service>` paths, `BatchService`
included, straight to the services. A service token that leaks therefore
works from the internet today.

Service and delegated tokens are for calls inside the cluster, so they are
refused when they arrive through the public listener:

- **Marking.** The public listener sets a request header on everything it
  forwards and removes any value a client sent for it. The interceptor on
  the gRPC server refuses a service token, in either form, and a delegated
  token on a request that carries the header. Access and demo tokens are not
  affected.
- **Every way in.** A request from the public listener reaches a token check
  by four paths, and the rule has to hold on each:
  - *Direct gRPC:* the listener routes to the service, and the interceptor
    sees the header.
  - *The REST gateway:* it forwards only the `authorization` header to the
    services today, so it would drop the mark. It forwards the mark as well.
  - *HTTP handlers that verify the token themselves:* the AI chat endpoint
    does not pass the gRPC interceptor. It makes the same check, and so does
    any handler added later.
  - *The auth service:* its server has no token interceptor, because signing
    in needs no token. Two of its RPCs verify a token themselves.
    `SwitchOrganization` takes an access token in the request body. It
    accepts a user's access token only and refuses a service token in either
    form, wherever the request comes from, so a leaked service token cannot
    be traded for user tokens that the public listener accepts. The
    delegated token exchange requires the `api` service token and checks the
    mark like the interceptor.
- **Callers inside the cluster.** Calls between services use the gRPC
  listener (port 9000), which the Ingress does not target, and are not
  affected. Two callers send a service token to the public listener today:
  the batch CronJobs call port 9003 inside the cluster, and the operator
  tools (the `hack/` commands such as `create-api-key`, which the test setup
  also runs) go through the public address. They move first, to a listener
  that the Ingress does not target and that does not set the header: the
  gRPC listener, or a new one for callers that need the REST form. The tools
  reach it from inside the cluster or through a port-forward. The tests
  themselves use access tokens and are not affected.
- **Report before refusing.** One release logs and counts the service and
  delegated tokens that arrive through the public listener, on every path.
  A later release refuses them, once the report is empty. A test sends a
  service token through each path of the public listener and expects it to
  be reported, then refused, so a path that loses the mark is found before
  the report is trusted. Another test gives a service token to
  `SwitchOrganization` and expects it to be refused.

This does not stop a pod inside the cluster from calling a service directly
with a leaked token. That needs the network controls in §7.

### What a stolen token gives

| | Today | Proposed |
| --- | --- | --- |
| Access token leaks on its own | Every organization of the user; can be traded for a refresh token through `SwitchOrganization` | Every organization of the user, until it expires. While `SwitchOrganization` exists it can still be traded for a refresh token, which now ends at the maximum sign-in age. |
| Both tokens are stolen | Every organization of the user, renewable without limit | Every organization of the user, until the maximum sign-in age |
| Role or system admin status is removed | Applies when the token expires | Applies on the next request |
| A service token leaks | Full access for up to 100 years. The holder cannot be identified, and cancelling the token rotates the signing key, which signs out every user and replaces the token everywhere. | Full access from inside the cluster only, until the tokens of its service are cancelled or it expires, one year at most by default. The token names its service, and only that service gets a new one. |
| The key that signs service tokens leaks | It is the key that signs user tokens. Replacing it signs out every user and replaces every service token. | With a service token key: full access until it is replaced, which replaces the service tokens only. Without one: as before. |
| A service token is used to name a user as the editor | Any holder can, on any request, through the relay metadata | Not on a request of its own. A delegated token is issued only for a valid API key and names that key's maintainer. An indirect way remains, described below. |
| A delegated token leaks | No such token exists | What the API key itself allows: the RPCs open to its role, in its environment or organization, as its maintainer, for at most 5 minutes |

A leaked service token still gives full access to a holder inside the
cluster until it is cancelled. That includes creating an API key with
another user as its maintainer and then using the key through the public
API, which acts as that user. Unlike the relay metadata, this leaves a
record: creating the key is written to the audit log with the service token
as the editor. What narrows a service token further is outside this RFC:
rules per service in RFC 0049, so that only the services that need to can
manage API keys, and the network controls in §7.

## 5. Alternatives Considered

- **`sub` as `organization_id:email`.** It matches the `account_v2` key, but
  a system admin viewing an organization has no account there, the subject
  would change on every switch, and every reader would have to parse it.
- **An enforced `organization_id` claim.** The server would reject a request
  for any other organization. It helps only when the access token leaks on
  its own, and then narrows a few minutes of access from every organization
  to one. When the browser's storage is stolen the refresh token is stolen
  with it and can switch organization, unless switching requires signing in
  again, which costs every user who works in several organizations.
- **Server-side sessions now.** They add revocation, but also a table,
  rotation logic and a two-release rollout. The maximum sign-in age bounds
  the same theft without storage, and sessions can follow without changing
  the claims (§7).
- **A `typ` header instead of `token_use`.** `typ: at+jwt` is the standard
  marker for access tokens. `token_use` is as safe, since both are signed,
  and was introduced recently with its own compatibility handling. Nothing
  outside our verifier reads these tokens, so replacing it has no benefit
  today.
- **A separate audience for service tokens.** PR #3015 considered it. The
  verifier already dispatches on `token_use`, so an audience would be a
  second marker for the same purpose, and the two would have to agree.
- **Checking tokens only at the gateway.** The REST gateway would
  authenticate, and the gRPC services behind it would trust every call. The
  gateway is not the only way in: the public listener also routes the
  service paths straight to the services, which the tests, the operator
  tools and the batch CronJobs use, and the charts do not limit the service
  ports to Envoy. One routing mistake would expose an administrative API
  with no check. The check stays on the gRPC server, where every route
  passes through it.
- **Signing the delegated token in the API service.** It saves the exchange
  call, but the signing key would have to be mounted in the API service,
  which faces the internet. Only the web deployment needs it. The batch
  deployment also mounts it today, without using it, and stops doing so
  (§6, step 15).
- **Cancelling service tokens by ID.** A list of cancelled `jti` values
  cancels one token at a time. It needs the ID of the leaked token, grows
  with every entry, and leaves a rotated token valid until it expires. A
  time per service has none of these limits.
- **A service token key required in every deployment.** Every installation
  would need a second key pair, including those that run with the defaults.
  Making it optional gives the separation to the deployments that want it.
- **Issuing service tokens inside the cluster.** No operator would handle a
  signing key, but a component in the cluster would have to hold the key and
  be allowed to write Secrets. This RFC keeps tokens issued by the
  generator; Kubernetes service account tokens are the longer-term direction
  (§7).
- **Forwarding the API key to the web services.** They could look the key up
  themselves, with no new token. The long-lived key would then travel to
  every service and its logs, where it is deliberately obfuscated today.
- **Keeping the relay metadata.** Limiting it to the `api` service token is
  the interim step of this RFC. On its own it leaves the identity of the
  user in unsigned request metadata, and the web services unable to check
  the key's environment.

## 6. Rollout Plan

Old claims are still written until no running version reads them, and a new
form is accepted one release before anything sends it. `organization_id` and
`name` are the exception: they go in step 3, for the reasons given there.

Split the user and service tokens:

1. **Fix the system admin derivation.** Independent of everything else.
2. **Read both forms.** The signer adds `sub` to the access token and
   `auth_time` to the access and refresh tokens. The verifier accepts tokens
   with or without them, and service tokens in the new and the legacy form.
   It rejects a service token issued before the time set for its service,
   and checks service tokens against the service token key when one is
   configured. The generator gains `--service` and `--ttl`, with one year as
   the default for a token that names a service. `RefreshToken` and
   `SwitchOrganization` enforce the maximum sign-in age, and
   `SwitchOrganization` refuses a service token in either form. The service
   and demo tokens use the shared claims type, and every kind is signed and
   verified through one path. The next release rejects a refresh token, and
   an access token given to `SwitchOrganization`, without `auth_time`. It is
   deployed no sooner than the refresh token lifetime after this one.
3. **Console, `organization_id` and `name`.** The console stops decoding the
   token and stops calling `SwitchOrganization`. Mark the RPC deprecated.
   The signer stops writing `organization_id` and `name` in the same
   release, and `pkg/role` takes the editor name from the account it loads.
   - `organization_id`: no backend version reads it. A console of the
     previous release calls `SwitchOrganization` when the claim is missing,
     and the RPC is still served.
   - `name`: the previous release copies it into the editor of audit logs.
     During a rolling update, or after a rollback, a pod of that release
     records an empty editor name until the token expires. The console shows
     the email when the name is empty, so this is accepted.
   - Where `pkg/role` loads no account (a service token, or a system admin
     acting without membership), the editor name is empty until RFC 0049's
     authorizer replaces `pkg/role`.
4. **Remove `SwitchOrganization`** no earlier than one release after step
   3, and only once its request metric shows no calls. The console does not
   reload itself when a new version is deployed, so a tab left open keeps
   calling the RPC until it is reloaded.

Authorize in one place:

5. **Build the caller from the token** in RFC 0049's interceptor on the gRPC
   server: an access token gives the email, a service token gives the
   service name.

Give each service its own token:

6. **Issue one token per service** once the release after step 2 is
   deployed. The release before step 2 rejects the new type; waiting one
   release means a rollback lands on a version that accepts it. The new
   tokens last one year. From this release a service warns when its own
   token is close to expiry, and the rotation steps are documented. A
   deployment that wants a service token key configures it before issuing
   the tokens. The legacy token stays valid, and putting it back is the way
   to roll back further.
7. **Accept the relay metadata from `api` only.** The legacy token may still
   relay until it is dropped.

Replace the relay metadata:

8. **Accept delegated tokens.** The auth service gains the exchange RPC and
   the web services verify delegated tokens, including the key role and
   scope of each RPC. The API service's Envoy gets a route to the auth
   service. Nothing sends the token yet.
9. **Send delegated tokens for the calls that relay today**, one release
   later. A rollback of the API service returns to the metadata, which is
   still accepted.
10. **Stop reading the relay metadata** one release after step 9.

Judge public API reads as the maintainer:

11. **Report the reads that would be denied.** For a read from a key whose
    maintainer would not pass the check, the API service writes a log entry
    and counts it in a metric, and still sends its own service token. The
    web services cannot report this, because they do not see the maintainer
    on those reads. The release notes ask operators to change the maintainer
    of the reported keys. This step depends on nothing else and can ship in
    any release.
12. **Send delegated tokens for reads** no earlier than one release after
    steps 9 and 11, once the report shows no key in use that would be
    denied, or its owners have been told.

Clean up:

13. **Stop writing the remaining old claims** (`is_system_admin`,
    `is_service_token`) in the release after the one that deletes
    `pkg/role`, so a rollback target does not need them. The access token
    then moves to the shared claims type.
14. **Drop the legacy service token** in a later release, announced in the
    release notes. A deployment can stop accepting it earlier by setting the
    time for legacy tokens.

Keep the signing key and the service tokens inside. These steps do not wait
for the steps above:

15. **Mount only the public key outside the web deployment.** The batch
    chart mounts the whole key Secret and reads only `public.pem`. This step
    depends on nothing else and can ship in any release.
16. **Report service and delegated tokens on the public listener.** The
    public listener marks what it forwards, the REST gateway forwards the
    mark, and the AI chat endpoint and the delegated token exchange check it
    like the interceptor. The CronJobs and the operator tools move to a
    listener that the Ingress does not target. The service and delegated
    tokens that still arrive through the public listener are logged and
    counted on every path.
17. **Refuse them** no earlier than one release after step 16, once the
    report is empty.

## 7. Future Work

- **Server-side sessions.** A session table in the existing database gives
  sign-out and immediate revocation, and refresh token rotation with reuse
  detection.
- **HttpOnly cookie for the refresh token**, so a script on the page cannot
  read it. Nothing in the backend sets cookies today.
- **Rules per service.** With a name in every service token, RFC 0049 can
  list the RPCs each service may call instead of giving all of them full
  access.
- **Remove the public service routes.** The public listener routes
  `/bucketeer.<package>.<Service>` paths straight to the services. After
  step 17 of §6 they no longer accept a service token, but they still offer
  every RPC to a user token next to the REST gateway. They can be removed
  once nothing outside the cluster uses them.
- **Restrict the service ports.** The Go servers listen on every interface
  and the charts define no NetworkPolicy, so nothing in the charts stops
  another pod from reaching them without passing Envoy.
- **Mutual TLS between services.** TLS authenticates only the server today:
  every pod mounts the same certificate, and the Go clients present none. A
  private CA with one certificate per workload would identify the calling
  service at the connection, and could then replace the service token for
  calls inside the cluster.
- **Kubernetes service account tokens** as the service credential on
  Kubernetes. The kubelet issues one per workload and rotates it about every
  hour, so no long-lived token is stored. Docker Compose would keep the
  signed service token.

## 8. Open Questions

1. The default maximum sign-in age. The proposal is 30 days.
2. Whether the service names are a fixed list in the code or free text. The
   relay rule and later the rules per service need at least `api` to be a
   known name.
3. How many API keys in use have a maintainer who would be denied on reads.
   It is not known until step 11 reports them, and it decides how long to
   wait before step 12.
4. Where the key role and scope of each RPC live until RFC 0049 moves the
   API key roles into its Table 2: in a list next to the delegated token
   verifier, or in Table 2 from the start.
5. PR #2061 (password login) changes the same auth service and has not been
   checked against this design.
6. How the operator tools reach a listener inside the cluster (§6, step
   16): through a port-forward, or by running inside the cluster.
