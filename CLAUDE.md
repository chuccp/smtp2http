# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Common Commands

### Backend (Go)

```bash
go build -o http2smtp ./        # Build
go build ./...                  # Compile-check everything
go vet ./...
go test ./...                   # Run all tests
go test ./util -run TestParserCron   # Run a single test
```

The go.mod carries `replace github.com/chuccp/go-web-frame => ../go-web-frame`, so the
sibling checkout is what actually gets compiled — editing it takes effect immediately,
and a version bump in go.mod alone changes nothing.

Run the server:

```bash
./http2smtp                                    # Defaults: 12566 (mgmt) + 12567 (API)
./http2smtp -web_port 12566 -api_port 12567    # Custom ports
./http2smtp -storage_root /path/to/storage     # Custom storage
```

**A port given by flag or environment means Docker** (`core.isdocker`): the container's run
command fixed it and its port mapping depends on it, so the UI may not move it. Both UIs
therefore only *display* the ports there — the wizard always, the settings page when
`isDocker` — and `RestartService` refuses to change them. The flag is written on every
start, true or false: it lives in the config file, and a stale `true` from an earlier
container run would otherwise lock the ports for good.

Port resolution is flag/env → `config.ini` → the 12566/12567 constants, read through
`model.GetConfig` so the INI's strings survive. `config.ini` matters here: the wizard and
the settings page both write it, and ignoring it used to mean a port saved in the UI was
silently dropped on the next start. The wizard echoes the ports back in `/dbInit` even
though it does not edit them, because that handler binds onto `model.DefaultConfig()` and an
omitted port would be persisted as 12566/12567.

Stop a running instance:

```bash
taskkill /F /IM http2smtp.exe                   # Windows
lsof -ti :12566 | xargs -I {} kill -9 {}        # macOS / Linux
```

Cross-compile and Docker:

```bash
make release-linux      # → build_assets/http2smtp
make release-osx        # → build_assets/http2smtp
cd docker && docker compose up -d
docker build -t http2smtp .
```

### Frontend (`view/`, Vue 3 + Vite)

```bash
cd view && pnpm install
cd view && pnpm run dev          # Dev server on :3000, proxies /api → localhost:12566
cd view && pnpm run type-check   # vue-tsc --noEmit
cd view && pnpm run build        # vue-tsc && vite build — must stay green
cd view && pnpm run build:force  # vite build only
```

pnpm (11.25.0, pinned in `packageManager`) is the package manager; `pnpm-lock.yaml` is the
only lockfile. Two things about this setup that fail loudly if mishandled:

- pnpm 11 keeps project config in `pnpm-workspace.yaml` (which otherwise has no `packages:`
  field here — it is config only). It requires a yes/no on every dependency install script;
  an undecided entry is written as a placeholder string, and that placeholder makes *every*
  `pnpm run` fail on the pre-run deps check. The one entry here denies `@parcel/watcher`'s
  build, which `scripts/build-from-source.js` runs only when no prebuilt binary shipped.
- **element-plus stays at 2.14.0 — do not take 2.14.5.** Its el-table slot props got typed
  as `DefaultRow`, which `vue-tsc` rejects at every `#default="{ row }"` destructure (18
  errors across `views/{smtp,token,user}/index.vue`), and `pnpm run build` runs `vue-tsc`
  first, so the build fails rather than the release.

## Architecture

HTTP2SMTP is an HTTP-to-SMTP gateway: configure SMTP servers through the web UI, then send
mail over HTTP from your own code.

### Dual server

- **Management server** (`manage.port`, default 12566, context path `/api`) — web UI plus
  admin APIs. Authenticated by a session cookie.
- **Public API server** (`api.port`, default 12567, no context path) — only `/sendMail`.
  Authenticated by a token from the `t_token` table, passed as a query/body param. No
  session cookie is involved, so the two auth paths never mix.

`rest/API` is the only REST type on the public server; everything else lives on the
management server. Keeping them separate is what lets you expose mail-sending publicly
while the admin surface stays private.

### Web UI hosting

The management server serves the built frontend itself, so no nginx or reverse proxy is
needed. `manage.webPath` (default `web`) points at the built bundle and `main.go` hands it
to the server group as `web.ServerConfig{Locations, Page404: "index.html"}`. API routes
stay under `/api`; everything else falls through to the static handler, which serves the
bundle from the site root and returns `index.html` for unmatched HTML requests so
client-side routes survive a refresh or deep link.

The static handler runs as a no-route fallback, so it sits outside the auth filter — the
bundle (login page included) must load unauthenticated. Every `/api/*` route is still
authorised normally.

### Backend layout

| Path | Role |
|---|---|
| `main.go` | Wires both rest groups, models, services and the cron runner |
| `rest/` | HTTP handlers. `Init(context)` registers routes; protect with `.WithMeta(auth2.WithLogin())` |
| `service/` | Business logic: `TokenService`, `ScheduleService`, `LogService`, `SmtpService`, `UserService`, `RestartService` |
| `model/` | GORM models plus the `Config` struct mapped to `config.ini` |
| `entity/` | Request/response DTOs and the log status constants |
| `auth/` | Session cookie (`authentication.go`) and the AES-CBC helpers (`token.go`) |
| `smtp/` | Delivery via `wneessen/go-mail`; `*2` suffixed functions are the live ones |
| `runner/` | `ScheduleRunner` — reloads scheduled tasks into cron every minute |
| `db/` | Builds the GORM connection from config (sqlite or mysql) |
| `util/` | String, time, template and HTTP helpers |

Tables: `t_user`, `t_SMTP`, `t_mail`, `t_token`, `t_schedule`, `t_log`.

### Sessions and authorisation

The session is an AES-CBC-encrypted `LoginUser` carried in an **HttpOnly cookie named
`user_token`**, validated against `t_user.salt`. `Salt` is a session *epoch*, not part of
password hashing (bcrypt carries its own salt): the token embeds the salt it was minted
with, and `auth.User` rejects it once the stored value no longer matches.

Salt is therefore rotated whenever sessions must end — **sign-out, password change, and
user disable**. That means signing out ends every session for that user, on every device.
`auth.User` also accepts the token via an `Authorization` header, which the browser never
uses but external callers can.

Ownership: non-admin users only see their own rows, filtered by `user_id`. Admin-only
endpoints re-check `IsAdmin` in the handler — `/reSet` included, since it rewrites the
database and port configuration and hot-swaps the live connection under the running
models. `PUT` handlers verify ownership and preserve the original owner, so an admin
editing someone else's row does not take it over. Row-level reads (`GET /smtp/:id`,
`/mail/:id`, …) answer `data: null` for a row the caller does not own rather than
reporting a permission error, so existence does not leak; the list endpoints are what
admins use to reach other users' rows, and those exclude nothing by default
(`adminOnly=true` narrows them to admin-owned rows).

Two references are checked as well, because resolving them by id alone would hand out
another account's credentials: a token may only point at the owner's own SMTP server and
its own contacts. And a token stops sending once its owner is disabled — the public API
has no session, so `TokenService` re-reads the owner (`t_user.is_use`) on every send.

Mail logs are per-owner: `LogService.Log` records the `user_id` it is given (the token
owner for a token send, the schedule owner for a scheduled call), because the log list
and the attachment download both filter on it.

### Restart

`POST /api/restart` (login + `IsAdmin`) restarts in place through `WebFrame.ReStart`, which
cancels the run context and lets the framework re-init and serve again in the same process.
The settings page's save-and-restart button saves through `/reSet` first, then calls it.

The catch: a restart re-runs the server groups `main.go` built, and those carry the listen
ports captured at startup — a static `Port()` would freeze the old value forever. So
`main.go` hands each group a `*web.ServerConfig` it keeps, and `service.RestartService`
rewrites `Port` on those two before triggering the restart — except under `core.isdocker`,
where those ports are the container's and stay put. It also refuses a change the app could
not come back up with (port in use, both ports equal), leaving the running instance alone;
the two must differ because groups sharing a port are merged onto one gin engine and the
public API would inherit the management server's auth filter. The response carries the ports
now in effect — the client's own origin may be gone by the time it arrives, which is what
the warning on the settings page reports.

### Configuration

Precedence is defaults → `config.ini` → CLI flags → environment variables.

> **INI values load as strings, and the decoder only weak-converts when the target is a
> struct.** Two consequences that have bitten this codebase:
>
> - Read config through `model.GetConfig()` / `model.CoreState()`. `IConfig.GetBoolOrDefault("core.init")`
>   only accepts a real bool, so against an INI-sourced `init = true` it silently returns the
>   default — which once left the setup guards dead and stopped the database from attaching
>   at startup. Pass the `*Config` itself, never `&cfg`: a `**Config` makes the decoder fall
>   back to `encoding/json`, which cannot coerce `"true"` → bool or `"12566"` → int.
> - A `[]string` field cannot be set from an INI file at all (`cannot convert string to
>   []string`), and a bad value crashes startup rather than reporting a config error.

## Frontend (`view/`)

Vue 3 + Element Plus + Pinia + vue-i18n, four locales (`zh-cn`, `en`, `zh-tw`, `ja`) that
must stay in key-for-key sync.

- `src/api/request.ts` — axios on `/api`, `withCredentials: true`, no credentials are
  attached by hand: the session cookie rides along automatically. Non-`0`/`200` codes and
  HTTP errors both reject, and 401 clears the local state and redirects to `/login`.
- `src/store/auth.ts` — holds **display state only** (username, `isAdmin`, a `loggedIn`
  flag). The session credential is HttpOnly and deliberately unreadable from script, so
  this flag drives navigation only; every request is authorised server-side.
  `views/Index.vue` hydrates it from `GET /api/set`, which reports `hasLogin`,
  `username` and `isAdmin`.
- `src/router/index.ts` — guard on `meta.requiresAuth` / `meta.requiresAdmin`.
- `src/views/settings/index.vue` — admin-only; ports and database config, plus the
  save-and-restart button that makes a saved port take effect.
- `vite.config.ts` — dev proxy for `/api`. Keep `VITE_API_BASE_URL` relative so requests
  stay same-origin; pointing it at the backend directly is cross-origin and the session
  cookie will be dropped.

### Release packaging

`view/` is the frontend that ships. The release workflow installs pnpm 11.25.0, then runs
`pnpm install --frozen-lockfile` and `pnpm run build` in `view/`, and packages `view/dist`
as `web/` in the release tarball. Because `pnpm run build` runs `vue-tsc` first, a type error
fails the release rather than shipping.
(Older releases pulled a prebuilt frontend from `chuccp/d-mail-view`; that dependency is
gone, so `d-mail-view` is no longer part of the build.)

## Known gaps

- **Signing out ends every session for that user**, on all devices — it rotates the session
  salt (see above). Per-device revocation would need a server-side session store.
- **The public `/sendMail` endpoint has no rate limiting.** The framework ships a
  `ratelimit` component; it is not wired up here.
- **INI cannot express a `[]string` field.** `web.server.locations` therefore has to be set
  in code (see `main.go`); putting it in `config.ini` makes startup fail with
  `cannot convert string to []string`.
