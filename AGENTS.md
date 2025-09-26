## AGENTS Guidelines (Operational Playbook)

This repository hosts a layered Go application built around: clean(ish) layering (handler → service → repository → DB), dependency injection via Google Wire, modular servers (`cmd/server`, `cmd/migration`, `cmd/task`), and supporting packages under `pkg/` (logging, config, JWT, SID generation, HTTP/GRPC bootstrap). This document tells an automated agent (and humans) exactly how to work safely and productively in this codebase.

---
### 1. Golden Rule: Use the Dev Workflow, NOT the Prod Build
* Always start / restart with: `make bootstrap`
  - Spins up dependent services (via Docker Compose if defined)
  - Loads config (default `config/local.yml`)
  - Runs migrations (through migrate server logic)
  - Starts HTTP server with hot reload (`nunu`)
* Do NOT invoke `make build` during interactive development inside the agent session. That produces a production artifact and slows iteration.
* If a production image or binary is needed, perform it manually outside of the agent automation scope.

---
### 2. Architecture Map (Read Before Changing Code)
Layer / Area | Path(s) | Purpose / Notes
-------------|---------|----------------
API Schema & Errors | `api/v1` | Request/response DTOs + central error registry (`errors.go`). Extend by adding new typed request structs & using `newError` pattern.
HTTP Routing | `internal/router` | Composition of Gin routes. New domain routes go in a dedicated file (mirroring `user.go`).
Handlers (Transport) | `internal/handler` | Translate HTTP <-> service calls. Keep them thin (validation, binding, response formatting).
Services (Business Logic) | `internal/service` | Core application orchestration. No transport code here. Return domain models or simple structs.
Repositories (Persistence) | `internal/repository` | DB abstraction (GORM). Add interfaces + concrete implementations. Inject via Wire.
Models | `internal/model` | GORM model definitions (e.g. `user.go`). Keep persistence-focused fields here.
Background Jobs | `internal/job` | Long-running / scheduled job definitions.
Tasks / CLI | `cmd/task` + `internal/task` | One-off or recurring maintenance / batch tasks.
Migrations | `cmd/migration` + `internal/server/migration.go` | Run schema migrations at bootstrap.
Server Bootstrapping | `pkg/app`, `internal/server`, `pkg/server/http` | App lifecycle & HTTP server wiring.
Logging | `pkg/log`, `pkg/zapgorm2` | Unified logging. Use injected `*log.Logger`—never create ad-hoc loggers.
Security / Auth | `pkg/jwt`, middleware in `internal/middleware/jwt.go` | JWT creation & validation; keep auth changes centralized.
ID Generation | `pkg/sid` | Structured unique ID helpers.
Configuration | `pkg/config`, `config/*.yml` | Viper-based config. Accept external path with `-conf` flag.
Tests & Mocks | `test/` + `test/mocks` | Organized by layer; mocks generated via `make mock`.

---
### 3. Adding a New Domain (Example: Article)
1. Model: Create `internal/model/article.go` (define GORM struct + tags).
2. Repository: Add interface + implementation in `internal/repository/article.go`.
3. Service: Add `internal/service/article.go` with business logic.
4. Handler: Add `internal/handler/article.go` (Gin handlers) + request/response DTOs in `api/v1` if externally exposed.
5. Router: Register routes in `internal/router/article.go` and include from `router.go` init logic.
6. Wire: Add constructor(s) to `internal/repository` / `internal/service` and include them in the appropriate `wire.go` set (e.g. `cmd/server/wire/wire.go`). Regenerate with `wire` (usually handled externally—avoid running full build; if needed, run `go generate ./...`).
7. Tests: Add service tests under `test/server/service`, repository tests under `test/server/repository`, and handler tests under `test/server/handler`.
8. Swagger: Update DTOs in `api/v1` and run `make swag` to regenerate docs.

Keep handlers transport-only; do not embed SQL queries or business logic there. All data access must go through repositories.

---
### 4. Error Handling & Response Pattern
* Define new semantic errors using `newError(code, message)` in `api/v1/errors.go`.
* Always return responses through `v1.HandleSuccess` / `v1.HandleError`.
* Reserve HTTP status codes for protocol-level meaning (e.g., 400/401/404/500) and use `.Code` in response JSON for business codes.
* Never leak raw DB error strings to clients—wrap or map them.

---
### 5. Dependency Injection (Wire)
* Wire sets exist per command in `cmd/*/wire`.
* To introduce a new injectable component:
  - Provide a constructor `NewXxx(...) *Xxx` (add a short comment summarizing purpose).
  - Add it to the relevant `wire.NewSet()`.
  - Regenerate wire output (outside agent if possible). The generated `wire_gen.go` must NOT be manually edited.
* Do not introduce global singletons; prefer DI wiring.

---
### 6. Configuration & Environments
* Default local config: `config/local.yml`
* Prod example: `config/prod.yml`
* Pass custom config path with: `-conf path/to/file.yml`
* Never hardcode secrets—use environment variables that Viper can bind or external secret management.

---
### 7. Logging
* Use injected logger: `logger.Infof(...)`, etc.
* For DB logging, `pkg/zapgorm2` integrates with GORM—avoid custom verbose SQL prints.
* Ensure sensitive values (passwords, tokens) are redacted.

---
### 8. Auth & Security
* JWT creation lives in `pkg/jwt`; validation is middleware in `internal/middleware/jwt.go`.
* Protect new routes by adding the middleware group in the router.
* For password handling (see user flow), always hash before persistence—never log plaintext.

---
### 9. ID & Data Consistency
* Use `pkg/sid` for generating sortable, prefixed IDs if required.
* Enforce uniqueness constraints (e.g., email) at both application (pre-check) and DB levels when feasible.

---
### 10. Testing Guidance
Layer | What to Test | Hints
------|--------------|------
Repository | DB CRUD paths | Use a test DB / transaction rollback.
Service | Business rules, branching | Mock repositories from `test/mocks`.
Handler | HTTP contract | Spin up Gin with only required routes; assert JSON.
Auth | Token issuance & middleware | Mock time where helpful.

Run tests with: `make test`
Regenerate mocks (when interfaces change): `make mock`

Aim for fast, deterministic tests—avoid sleeping; prefer context deadlines.

---
### 11. Swagger / API Docs
* Update request/response structs in `api/v1` with `json` + `example` tags.
* Run `make swag` to regenerate `docs/swagger.*`.
* Serve docs if integrated (check router for `/swagger` mounting logic if added later).

---
### 12. Performance & Resource Hygiene
* Avoid N+1 queries—prefer `Preload` or explicit joins.
* Keep allocations low in hot paths (services). Reuse buffers only when safe.
* Wrap long DB operations with context timeouts if added in future.

---
### 13. Concurrency Notes
* Do not share mutable state across goroutines without synchronization.
* Background jobs (`internal/job`) should be idempotent and stoppable via context.

---
### 14. Making Safe Changes (Agent Checklist)
1. Read relevant files (model/service/repository/handler) before editing.
2. Create minimal, targeted diff—avoid formatting unrelated regions.
3. Add a short comment to any newly created function (per repository rule).
4. Run `go vet` / `go mod tidy` (implicitly via tests or manual) if dependencies change.
5. Run `make test`.
6. If API shape changes, regenerate swagger.
7. Update this document if process-impacting conventions change.

---
### 15. Common Pitfalls
Pitfall | Avoid By
--------|---------
Calling `make build` mid-iteration | Use `make bootstrap`.
Embedding SQL in handlers | Put logic in repository/service.
Returning raw errors to clients | Map via `HandleError`.
Duplicating user/password handling | Centralize in service layer.
Adding global state | Inject with Wire.

---
### 16. Useful Commands (Quick Reference)
| Command          | Description |
|------------------|-------------|
| `make bootstrap` | Start dev stack (services + migrations + hot reload). |
| `make test`      | Run tests with coverage. |
| `make mock`      | Regenerate mocks for changed interfaces. |
| `make swag`      | Regenerate Swagger documentation. |
| `go mod tidy`    | Reconcile module dependencies. |
| `wire` / `go generate ./...` | Regenerate DI code (avoid unless sets changed). |

---
### 17. When to Ask for Human Review
* Schema changes (new tables / destructive migrations)
* Authentication logic changes
* Introduction of external services / dependencies
* Changes to error code semantics

---
### 18. Migration Flow
* Migration server: `cmd/migration` uses Wire to assemble a migrate-only app.
* Ensure migrations are idempotent and additive.
* Keep destructive operations (drops) behind explicit human review.

---
### 19. Background Jobs / Tasks
* Jobs: Place logic in `internal/job` and register in server startup.
* Tasks / batch runs: Use `cmd/task` with its own Wire graph.

---
### 20. Logging Patterns
Good:
  logger.Infof("creating user", log.String("email", safeEmail))
Bad:
  logger.Infof("creating user %s with password %s", email, password)

Mask secrets / PII where appropriate.

---
### 21. Code Style Quick Notes
* Follow Go naming: exported types/functions start uppercase when part of public surface.
* Keep functions short; extract helpers for clarity.
* Use context-aware function signatures where cancellation matters.
* Add comments to all new public functions AND any new private helpers (repository policy: "Always add a comment to created functions").

---
### 22. Security Hygiene
* Never commit real secrets to `config/*.yml`.
* Validate all user input with Gin `binding` tags.
* Use prepared statements implicitly (GORM handles this) – avoid manual string concatenation for queries.

---
### 23. If Something Breaks
1. Re-run `make bootstrap` (clears some transient states).
2. Check logs in `storage/logs/server.log`.
3. Run focused tests for failing layer.
4. Inspect DI graph changes (Wire sets) if nil pointer panics at startup.

---
### 24. Minimal Pull / Change Description Template
Subject: <scope>: short imperative summary

Body:
* Motivation:
* Changes:
* Risk & Mitigation:
* Follow-up:

---
### 25. Final Reminder
Prefer incremental, well-tested changes over large refactors. If in doubt: read, trace, test, THEN modify.

