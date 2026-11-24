# KepKep (ເກັບເກັບ) — Backend Plan

> Send the slip, KepKep writes it down. Open the app, see what you can safely spend today.

This plan covers the **backend only**, from an empty folder to production. It is written for project-based learning: each phase builds one part of the product and teaches the backend concepts that part needs.

---

## 1. Product summary

| Item | Decision |
|---|---|
| Name | KepKep (ເກັບເກັບ) |
| Mascot | Kep the duck, says "kep kep!" |
| First users | Young salaried workers and students in Vientiane who pay by bank QR (BCEL One first, then LDB, JDB) |
| Core promise | No typing for bank payments (slip scan) + "Safe to spend today" on the home screen |
| Signature feature | Savings goal shown as an egg that hatches when the goal is reached |
| Language | Lao first, English second |
| Money | Lao kip, stored as `int64` |

### What the MVP must prove

1. AI reads Lao bank slips correctly (amount right in ≥ 95% of scans). **Tested first, in Phase 1, before building the backend.**
2. Users keep logging after 30 days.
3. "Safe to spend today" helps users stay within budget.
4. Enough users pay for Premium (unlimited scans).

### Two releases

- **Beta** — the smallest product that tests the core promise: login, accounts, manual entry, transfers, slip scan, categories, bills, income, savings goal, safe to spend, history, daily summary, 3 alerts, Free/Premium.
- **v1.1** — added after real users are using the beta: budgets, debt book, bulk import, monthly report, category breakdown, the other 4 alerts.

Getting real users early matters more than having every feature on day one.

---

## 2. Actors

| Actor | Type | Inputs | Outputs |
|---|---|---|---|
| User (Free) | Person | Phone number; income amount, frequency, payday; fixed bills; slip photos (30/month); manual cash entries; categories; budgets; debts; one savings goal | Safe to spend today; history; daily summary; monthly report; alerts |
| User (Premium) | Person | Same + unlimited scans + bulk import | Same + full history + detailed monthly report |
| Admin (you) | Person | Default categories; slip rules per bank; plan prices; premium activation | Admin metrics: active users, scans/day, AI accuracy, AI cost, conversion, errors |
| AI slip reader | System | Slip image | Amount, date, time, bank, receiver, reference number, suggested category, confidence |
| Scheduler | System | Clock + user settings | Daily summaries, alerts, reminders |

---

## 3. MVP features (27)

**Release** column: **Beta** = first release, **1.1** = after beta.

### A. Account and setup
| # | Feature | Input | Output | Release |
|---|---|---|---|---|
| A1 | Sign up / login | Phone, one-time code | Account, tokens | Beta |
| A2 | Profile | Name, language | App in chosen language | Beta |
| A3 | Money accounts | Name (BCEL, LDB, cash), opening balance | Balance per account + total | Beta |
| A4 | Income setup | Amount, frequency (monthly, weekly, daily, irregular), next payday | Days to next income | Beta |
| A5 | Fixed bills | Name, amount, due day | Reserved before safe to spend; reminder | Beta |

### B. Recording money
| # | Feature | Input | Output | Release |
|---|---|---|---|---|
| B1 | Slip scan | One slip image | AI draft transaction → confirm | Beta |
| B2 | Bulk import (Premium) | Up to 50 slips | Review list of drafts | 1.1 |
| B3 | Manual entry | Amount, category, note, optional bill | Saved transaction | Beta |
| B4 | Income entry | Amount, source, account | Updated balance | Beta |
| B5 | Transfer | From, to, amount | Balances move; not spending | Beta |
| B6 | Edit / delete | Changes | Corrected history and balances | Beta |
| B7 | Duplicate check | Automatic (image hash + slip reference number) | "Already saved" warning | Beta |
| B8 | Auto-category learning | Automatic (from user corrections) | Right category next time | Beta |

### C. Categories
| # | Feature | Input | Output | Release |
|---|---|---|---|---|
| C1 | Lao default categories | — | Food, coffee & drinks, transport & fuel, phone top-up, rent, electricity & water, family support, ceremonies & merit-making (ບຸນ, weddings, envelopes), shopping, health, education, debt payment, other | Beta |
| C2 | Edit categories | Add, rename, icon, colour, hide | Categories that fit the user | Beta |

### D. Planning
| # | Feature | Input | Output | Release |
|---|---|---|---|---|
| D1 | Safe to spend today | Uses balance, bills, savings, payday | Home-screen number | Beta |
| D2 | Monthly budget per category | Category, limit | Progress; alerts at 80% and 100% | 1.1 |
| D3 | Savings goal | Name, target, date | Amount per day/month; egg progress | Beta |

#### Safe to spend — definition

```
available   = total balance − unpaid bills due before next payday − savings set aside
daily_limit = available ÷ days until next payday
remaining   = daily_limit − spent today
```

- **Fixed at the start of the day.** `daily_limit` is calculated from the state at 00:00 in the user's timezone (today's transactions are ignored). Spending during the day lowers `remaining`, not `daily_limit`. Otherwise "over today's limit" could never be detected.
- **Exception:** income received today recalculates `daily_limit` straight away.
- **Total balance** = sum of all account balances.
- **Unpaid bill** = an active bill whose due day falls before the next payday and that has no transaction with its `bill_id` in the current bill period.
- **Savings set aside** = `savings_goals.saved_amount`. This money is still inside the user's accounts, so it is subtracted once. *(Product decision — can be changed.)*
- **Days until next payday** = whole days from today (user's timezone), counting today. Payday today counts as 1 day, never 0.
- **Several income sources:** use the earliest next payday.
- **Payday in the past** (user did not update it): move it forward by its frequency until it is in the future.
- **Irregular income or no income:** divide by 7.
- **Negative `available`:** `daily_limit` is 0, and the API also returns the shortfall so the app can show "you are short X kip before payday".

### E. Debt book
| # | Feature | Input | Output | Release |
|---|---|---|---|---|
| E1 | Add debt | Direction (I owe / owed to me), person, amount, due date, optional account | Totals both ways | 1.1 |
| E2 | Repayment | Amount, date, optional account | Remaining balance; closes at zero | 1.1 |

When an account is chosen, the debt or repayment also creates a transaction and moves the account balance. Without an account, it is only a record. *(Product decision — can be changed.)*

### F. History
| # | Feature | Input | Output | Release |
|---|---|---|---|---|
| F1 | History list | Scroll, search | Grouped by day with totals | Beta |
| F2 | Filters | Date range, category, account, type | Filtered list + total | Beta |

### G. Reports and alerts
| # | Feature | Input | Output | Release |
|---|---|---|---|---|
| G1 | Daily summary | Notification time | Yesterday's spending + today's safe to spend | Beta |
| G2 | Monthly report | — | Income, spending, savings, top categories, vs last month | 1.1 |
| G3 | Alerts | On/off per type | See section 4 | Beta (3) / 1.1 (4) |

### H. Plans
| # | Feature | Input | Output | Release |
|---|---|---|---|---|
| H1 | Free | — | 30 scans/month, 3 months history | Beta |
| H2 | Premium | Bank QR payment (admin activates in MVP) | Unlimited scans, bulk import, full history | Beta |

Free history limit **hides** older data. It never deletes it. When a user upgrades, everything is visible again.

---

## 4. Reports and alerts

### Reports (5)
| Report | When | Content | Release |
|---|---|---|---|
| Home summary | Every app open | Safe to spend, spent today, next bill, savings egg (+ budget bars in 1.1) | Beta |
| Daily summary | Each morning (user's time) | Yesterday by category + today's safe to spend | Beta |
| Category breakdown | Any date range | Totals and share per category | 1.1 |
| Monthly report | 1st of month | Income, spending, savings, top 5 categories, 5 biggest payments, budgets kept/missed, debt status, one small win | 1.1 |
| Debt summary | Any time | I owe / owed to me / next due dates | 1.1 |

### Alerts (7)
| Alert | Trigger | Dedupe key (fires once per…) | Release |
|---|---|---|---|
| Bill reminder | 3 days before due | `bill:<bill_id>:<period>` | Beta |
| Over today's limit | Today's spending > daily limit | `overlimit:<date>` | Beta |
| Low balance | Balance won't cover bills before payday | `lowbal:<next_payday>` | Beta |
| Budget warning | Category reaches 80% | `budget80:<month>:<category_id>` | 1.1 |
| Budget reached | Category reaches 100% | `budget100:<month>:<category_id>` | 1.1 |
| Debt reminder | 1 day before due | `debt:<debt_id>:<due_date>` | 1.1 |
| Small win | 3 days in a row under safe to spend | `win:<date>` | 1.1 |

Free users also get a notice at 25/30 scans (`scans25:<month>`).

---

## 5. Data entities (19 tables)

Tables are **not** all created at once. Each phase adds the migrations it needs (column "Phase").

| Table | Key fields | Phase |
|---|---|---|
| `users` | id, phone (E.164, e.g. `+8562055512345`), name, language, role (user, admin), plan (free, premium — cache, see below), plan_expires_at, timezone, daily_summary_time, created_at | 2 |
| `otp_codes` | id, phone, code_hash, expires_at, attempts, used_at | 2 |
| `refresh_tokens` | id, user_id, family_id, token_hash, expires_at, revoked_at, replaced_by | 2 |
| `categories` | id, user_id (null = default), name_lo, name_en, icon, color, kind (expense, income), hidden | 2 |
| `accounts` | id, user_id, name, type (bank, cash), bank_code, opening_balance, balance | 4 |
| `transactions` | id, user_id, account_id, to_account_id, category_id, bill_id, kind (expense, income, transfer, adjustment), amount, occurred_at, receiver, note, source (slip, manual, debt), slip_id | 4 |
| `idempotency_keys` | user_id, key, request_hash, response_status, response_body, created_at | 4 |
| `income_sources` | id, user_id, amount, frequency, next_payday | 5 |
| `bills` | id, user_id, name, amount, due_day, category_id, active | 5 |
| `savings_goals` | id, user_id, name, target_amount, target_date, saved_amount | 5 |
| `slips` | id, user_id, status (pending, extracted, confirmed, failed, discarded), image_sha256, bank, reference_no, extracted_json, confidence, object_key, image_deleted_at | 6 |
| `receiver_rules` | id, user_id, receiver, category_id, times_used | 6 |
| `devices` | id, user_id, push_token, platform, last_seen_at | 7 |
| `notification_settings` | user_id, type, enabled | 7 |
| `notifications` | id, user_id, type, dedupe_key, payload, sent_at, read_at | 7 |
| `subscriptions` | id, user_id, amount, paid_at, verified_by, expires_at | 8 |
| `budgets` | id, user_id, category_id, month (first day of month, user's timezone), limit_amount | 11 |
| `debts` | id, user_id, direction, person_name, amount, due_date, status, account_id, transaction_id | 11 |
| `debt_payments` | id, debt_id, amount, paid_at, account_id, transaction_id | 11 |

River creates its own job tables with its own migration (Phase 6).

### Data rules

- Money is `int64` kip. Never `float64`.
- Times are stored in UTC (`timestamptz`); daily jobs use each user's timezone (default `Asia/Vientiane`).
- Every query on user data filters by `user_id`. This is your main security rule.
- All user tables use `ON DELETE CASCADE` from `users`, so `DELETE /me` removes everything.

### Balances

- `accounts.balance` is a cached value. The truth is: `balance = opening_balance + sum of the account's transactions`. A reconciliation test checks this.
- Editing or deleting a transaction first reverses its old effect, then applies the new one, inside one DB transaction.
- When the user changes an account balance by hand, the backend creates an `adjustment` transaction for the difference. History always adds up.

### Constraints on `transactions`

- `amount > 0`, except `kind = 'adjustment'` where `amount <> 0` (can be negative).
- `kind = 'transfer'` ⇔ `to_account_id IS NOT NULL`.
- `to_account_id <> account_id`.
- Transfers and adjustments have no `category_id`.

### Duplicate slips

- **Before calling the AI:** check `image_sha256`. The same photo uploaded twice costs nothing.
- **After extraction:** partial unique index on `slips (user_id, bank, reference_no) WHERE status IN ('extracted', 'confirmed')`. Failed or discarded slips do not block a re-upload.
- Deleting a transaction that came from a slip sets the slip to `discarded`.

### Other rules

- **Plans:** `subscriptions` is the source of truth. `users.plan` and `plan_expires_at` are a cache, updated in the same DB transaction as the subscription.
- **Refresh tokens:** every login starts a new `family_id`. Reusing an old token revokes the whole family.
- **Notifications:** unique index on `notifications (user_id, type, dedupe_key)`. Inserting first, then sending, is what makes each alert fire once.
- **Slip images:** deleted after the slip is confirmed or after 24 hours, whichever comes first (River job). Account numbers are masked in `extracted_json`.

---

## 6. Tech stack

| Concern | Choice | Why |
|---|---|---|
| Language | Go (latest stable) | Your main language |
| HTTP framework | **Echo v4** | You already know it from Evently |
| Database | **PostgreSQL** | Transactions, constraints, strong for money data |
| DB driver | **pgx v5** | Fastest, most complete Postgres driver for Go |
| Queries | **sqlc** (no ORM) | You write SQL, sqlc generates type-safe Go. See decision below |
| Migrations | **goose** | Simple SQL migrations, can embed in the binary |
| Validation | **go-playground/validator v10** | Standard struct-tag validation, works with Echo |
| Auth tokens | **golang-jwt/jwt v5** | Short access token (JWT) + opaque refresh token in DB |
| Authorization | **Own middleware** (no Casbin) | See decision below |
| Config | Environment variables + **caarlos0/env** | Twelve-factor style, no config files with secrets |
| Logging | **log/slog** (standard library) | Structured JSON logs, no extra dependency |
| Background jobs | **River** | Job queue stored in Postgres: no Redis needed |
| File storage | **MinIO** (local) → **S3-compatible storage** in production (e.g. Cloudflare R2) via **minio-go** | Temporary slip images |
| AI slip reading | A vision-capable LLM API behind a `SlipReader` interface | Swap providers after your accuracy test |
| Push notifications | **Firebase Cloud Messaging** (firebase-admin Go SDK) | iOS + Android pushes |
| Rate limiting | **golang.org/x/time/rate** | Protect OTP and scan endpoints. In-memory: fine for one server |
| Testing | `testing` + table-driven tests + **testcontainers-go** | Real Postgres in integration tests |
| Linting | **golangci-lint** | Catches bugs and style issues |
| API docs | **OpenAPI** via swaggo/swag (+ echo-swagger) | Your mobile developer needs a contract |
| Load testing | **vegeta** (written in Go) | Check performance before launch |
| Containers | **Docker + Docker Compose** | Same setup locally and on the server |
| CI | **GitHub Actions** | Test + lint on every push |
| Hosting | One VPS in Singapore (close to Laos) + **Caddy** for HTTPS | Cheap, simple, enough for thousands of users |
| Errors (production) | **Sentry** (sentry-go) | See crashes with stack traces |

### Decision: ORM or not? → **No ORM. Use sqlc.**

- An ORM (like GORM) hides the SQL. You are learning backend, and SQL is the skill employers test.
- Money needs exact control: transactions, `SELECT ... FOR UPDATE`, unique constraints. With sqlc you write that SQL yourself and still get typed Go functions.
- Reports (`GROUP BY`, date ranges, sums) are easy in SQL and painful in an ORM.
- sqlc catches SQL mistakes at generate time, not at runtime.

### Decision: Casbin or not? → **Not for the MVP.**

- KepKep has only two roles (`user`, `admin`), and almost every rule is "users can only touch their own rows."
- That rule belongs in your SQL (`WHERE user_id = $1`), not in a policy engine.
- Admin routes use one small middleware: `RequireRole("admin")`.
- Revisit Casbin only if you add shared family accounts with many roles.

### Decision: architecture → **Modular monolith, hexagonal only at the edges. Not full clean architecture.**

- "Monolith vs microservices" is about deployment. "Layered / clean / hexagonal" is about code inside. They are separate choices.
- **One binary:** one developer, one database. A slip confirm changes 4 tables in one DB transaction — easy in one process, hard across services.
- **By feature:** `internal/<feature>` with handler → service → repository. All slip code lives in one folder.
- **Interfaces only at the edges** (things you swap or fake in tests): `SlipReader`, `SMSSender`, `ObjectStore`, `Pusher`. The interface is defined in the package that **uses** it.
- **No interface per repository.** You never swap Postgres. Repositories are tested against real Postgres (testcontainers), which catches real SQL bugs.
- **Kept from clean architecture:** dependencies point inward (services never import Echo), business rules are pure functions that take `now time.Time` as a parameter.
- **Revisit** only if one part must scale on its own. Slip processing is the only candidate, and River workers already cover most of that.

### Decision: JWT design

- **Access token:** JWT, 15 minutes, contains `user_id` and `role`.
- **Refresh token:** random string, stored **hashed** in `refresh_tokens`, 30 days, **rotated** on every refresh. Reuse of an old token = revoke the whole family.
- **Login:** phone + one-time code. Phone numbers are normalised to E.164 before hashing or lookup. During development, a fake SMS sender logs that a code was sent (never the code itself in production). Choose a real provider (a Lao SMS gateway or Firebase Phone Auth) before beta.

---

## 7. Architecture

A **modular monolith**: one Go binary, organised **by feature**, with handler → service → repository inside each feature (the same layering you used in Evently).

```
kepkep-backend/          # module github.com/Phaphonn12okjo/kepkep-backend
├── cmd/
│   ├── api/            # HTTP server + River workers
│   └── slipeval/       # CLI: test AI accuracy on a folder of slips (Phase 1)
├── internal/
│   ├── config/
│   ├── platform/       # db, logger, storage, ai, push, jobs, middleware
│   ├── auth/           # otp, jwt, refresh tokens
│   ├── user/
│   ├── account/
│   ├── category/
│   ├── transaction/
│   ├── slip/
│   ├── planning/       # income, bills, savings, safe-to-spend (+ budgets in 1.1)
│   ├── debt/           # 1.1
│   ├── report/
│   ├── notify/
│   └── admin/
├── db/
│   ├── migrations/     # goose SQL files, added phase by phase
│   └── queries/        # sqlc .sql files
├── sqlc.yaml
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── CLAUDE.md
└── .github/workflows/ci.yml
```

Each feature folder has `handler.go`, `service.go`, `repository.go` (sqlc wrapper), `model.go`, and `*_test.go`.

### One DB transaction across features

Some actions touch several features at once. Confirming a slip changes `slips`, `transactions`, `accounts` and `receiver_rules`. These must succeed or fail together.

- `internal/platform/db` has one helper, for example `InTx(ctx, pool, func(q *sqlc.Queries) error) error`. It begins the transaction, commits on success and rolls back on error.
- Inside it, sqlc's `queries.WithTx(tx)` gives every repository the same transaction.
- Services accept the queries they should use, so the caller decides whether work runs inside a transaction.

Decide and build this in Phase 4, the first time balances change.

### Idempotency

Phones lose network and retry. `POST /transactions`, `POST /transfers` and `POST /slips/:id/confirm` accept an `Idempotency-Key` header. The first response is saved in `idempotency_keys`. A retry with the same key returns the saved response instead of creating a second row.

---

## 8. API endpoints

All routes under `/v1`. All except auth require `Authorization: Bearer <access token>`. Routes marked *(1.1)* come after beta.

| Area | Endpoints |
|---|---|
| Health | `GET /healthz` · `GET /readyz` |
| Auth | `POST /auth/otp/request` · `POST /auth/otp/verify` · `POST /auth/refresh` · `POST /auth/logout` |
| Me | `GET /me` · `PATCH /me` · `DELETE /me` (delete account + data) |
| Accounts | `GET/POST /accounts` · `PATCH/DELETE /accounts/:id` (changing the balance creates an adjustment) |
| Income & bills | `GET/POST /income-sources` · `PATCH/DELETE /income-sources/:id` · `GET/POST /bills` · `PATCH/DELETE /bills/:id` |
| Categories | `GET/POST /categories` · `PATCH/DELETE /categories/:id` |
| Transactions | `GET/POST /transactions` · `GET/PATCH/DELETE /transactions/:id` · `POST /transfers` |
| Slips | `POST /slips` (upload) · `GET /slips/:id` (status + draft) · `POST /slips/:id/confirm` · `DELETE /slips/:id` · `POST /slips/bulk` *(1.1)* |
| Planning | `GET /safe-to-spend` · `GET/POST /savings-goal` · `PATCH /savings-goal` · `GET/POST /budgets` *(1.1)* · `PATCH/DELETE /budgets/:id` *(1.1)* |
| Debts *(1.1)* | `GET/POST /debts` · `PATCH/DELETE /debts/:id` · `POST /debts/:id/payments` |
| Reports | `GET /reports/home` · `GET /reports/categories?from&to` *(1.1)* · `GET /reports/monthly?month` *(1.1)* · `GET /reports/debts` *(1.1)* |
| Notifications | `POST /devices` (push token) · `GET /notifications` · `GET/PATCH /notification-settings` |
| Admin | `GET /admin/metrics` · `GET /admin/users` · `POST /admin/subscriptions` · `GET/POST /admin/categories` |

---

## 9. Phases: from zero to production

Rough time estimates assume about 10–15 hours a week. Beginners often need about 1.5× this; that is normal. Every phase ends with tests passing and a short note in `LEARNING_LOG.md` about what you learned.

Each phase that needs new tables writes its own migrations (see section 5).

### Phase 0 — Foundations (1 week)
- **Build:** Go module, Echo server, `/healthz`, config from env, slog JSON logger, request-ID middleware, graceful shutdown, Docker Compose (**Postgres only**), Makefile (`run`, `test`, `lint`, `migrate-up`, `sqlc`), GitHub Actions CI.
- **Learn:** project layout, environment config, Docker, graceful shutdown, CI.
- **Done when:** `docker compose up` + `make run` gives a 200 on `/healthz`, and CI is green.

### Phase 1 — Slip accuracy test (1 week)
This tests the riskiest idea in the product before building everything else.
- **Build:** `cmd/slipeval`. Put 30+ real slips (BCEL One, LDB, JDB) in a folder with a file of expected values. The CLI sends each slip to the AI, parses the JSON answer, compares it with the expected values and prints accuracy per field and per bank. Define the `SlipReader` interface here. Try 2–3 providers or prompts.
- **Learn:** interfaces, calling an external API with `net/http` (timeouts, errors), JSON parsing, `context`, measuring cost per call.
- **Done when:** accuracy and cost per slip are recorded for each provider in `LEARNING_LOG.md`, and one provider + prompt is chosen. **If no option reaches 95% on amount, stop and rethink the slip feature before Phase 2.**

### Phase 2 — Database and auth tables (1 week)
- **Build:** pgx pool, goose setup, sqlc setup, migrations for `users`, `otp_codes`, `refresh_tokens`, `categories`, seed default Lao categories, first sqlc queries.
- **Learn:** SQL types, primary/foreign keys, unique and check constraints, indexes, `ON DELETE CASCADE`, why `int64` for money, UTC timestamps.
- **Done when:** `make migrate-up` and `make migrate-down` both work from zero, and you can explain every index.

### Phase 3 — Auth (1–2 weeks)
- **Build:** phone normalisation (E.164), OTP request/verify (fake SMS sender), hashed OTP with expiry and attempt limit, JWT access tokens, rotating refresh tokens with `family_id`, auth middleware, `GET /me`, `PATCH /me`, rate limit on OTP routes (per phone and per IP), `RequireRole` middleware.
- **Learn:** authentication vs authorization, hashing, token rotation, rate limiting, 401 vs 403.
- **Done when:** a full login → refresh → logout flow works, reusing an old refresh token revokes the whole family, and tests cover each failure case.

### Phase 4 — Core money (2 weeks)
- **Build:** migrations for `accounts`, `transactions` (with check constraints), `idempotency_keys`. The `InTx` helper. Accounts with opening balance, categories (create, rename, hide), transactions CRUD (edit and delete reverse the old effect), transfers, adjustments, idempotency keys, history with filters and cursor pagination, consistent error responses.
- **Learn:** ACID transactions, row locks (`SELECT ... FOR UPDATE`), check constraints, idempotency, validation, error handling patterns, pagination, integration tests with testcontainers.
- **Done when:**
  - 50 concurrent requests on one account never produce a wrong balance (test it).
  - A reconciliation test shows `balance = opening_balance + sum(transactions)` after create, edit and delete.
  - The same `Idempotency-Key` sent twice creates one transaction (test it).
  - A user can never read another user's data (test it).

### Phase 5 — Planning (1–2 weeks)
- **Build:** migrations for `income_sources`, `bills`, `savings_goals`. Income sources, bills (paid via `transactions.bill_id`), savings goal, and the safe-to-spend calculation as a **pure function** following the definition in section 3.
- **Learn:** keeping business logic in pure functions, table-driven tests, date and timezone handling (start of day, month boundaries, payday in the past).
- **Done when:** safe to spend has at least 15 table-driven test cases and all pass. Include at least: no income, irregular income, payday today, payday in the past, two income sources, overspent (negative available), bill paid today, income received today, spending during the day does not change `daily_limit`.

### Phase 6 — Slip scan with AI (2 weeks)
- **Build:** add MinIO to Docker Compose. Migrations for `slips`, `receiver_rules`, River tables. Upload endpoint (size and type limits) → image hash duplicate check → store in MinIO → River job → `SlipReader` (from Phase 1) with timeout and retry → validate JSON → reference number duplicate check → draft transaction → confirm endpoint (inside `InTx`, with idempotency key) → `receiver_rules` learning → delete image → monthly scan quota for Free users.
- **Learn:** async jobs and queues, calling external APIs (timeouts, retries, cost), idempotency, interfaces for testing (fake `SlipReader`), handling private data, race conditions on quotas.
- **Done when:**
  - The same photo uploaded twice never calls the AI twice.
  - Two different photos of the same slip are rejected by reference number.
  - A failed slip can be uploaded again.
  - Images are gone after confirm or 24 hours.
  - Two uploads at the same moment cannot go over the Free quota (test it).
  - Tests run without calling the real AI.

### Phase 7 — Home summary and notifications (1–2 weeks)
- **Build:** migrations for `devices`, `notification_settings`, `notifications`. Home summary, River periodic jobs for the daily summary and the 3 beta alerts (bill reminder, over today's limit, low balance) + the 25/30 scan notice, FCM push, notification settings and log.
- **Learn:** `GROUP BY` and date ranges, indexes for reports, scheduled jobs, per-user timezones, avoiding duplicate notifications with a unique dedupe key.
- **Done when:** the home summary matches hand-calculated numbers in tests, and each alert fires once (not twice) when the scheduler runs twice.

### Phase 8 — Plans and admin (1 week)
- **Build:** migration for `subscriptions`. Free/Premium limits (scan quota, history hidden after 3 months), manual premium activation (updates `subscriptions` and the `users.plan` cache in one DB transaction), admin metrics (active users, scans/day, correction rate, AI cost, conversion), admin category management.
- **Learn:** feature gating, role-based access, metrics queries, source of truth vs cache.
- **Done when:** a Free user hits the scan limit with a clear error, old history is hidden (not deleted) for Free users, and non-admins get 403 on every admin route.

### Phase 9 — Production readiness (2 weeks)
- **Build:**
  - Security: rate limits on all sensitive routes, request body limits, CORS, secure headers, secrets only in env, `DELETE /me` that removes all user data, `govulncheck` in CI.
  - Reliability: Sentry, Prometheus metrics endpoint, structured access logs, `/readyz` checks DB.
  - Data: daily `pg_dump` backups to object storage **and one restore test**.
  - Deploy: production Dockerfile (multi-stage), VPS with Docker Compose + Caddy (HTTPS), staging environment, CI deploy on tag.
  - Docs: OpenAPI published for the mobile developer; privacy policy covering slip images, the AI provider and data deletion.
  - Load test with vegeta: home summary and slip upload.
- **Learn:** deployment, HTTPS, backups, monitoring, security basics, performance testing.
- **Done when:** production runs on HTTPS, a backup has been restored successfully, `DELETE /me` leaves no rows for that user (test it), and the load test results are recorded.

### Phase 10 — Beta (4 weeks, ongoing)
- **Do:** 20–50 real users, watch Sentry and metrics daily, measure the success metrics below, and fix what hurts most.

### Phase 11 — v1.1 (3–4 weeks, after beta feedback)
- **Build:** migrations for `budgets`, `debts`, `debt_payments`. Budgets with 80%/100% alerts, debt book with repayments (optional account creates a transaction), bulk slip import, category breakdown, monthly report, debt summary, debt reminder and small-win alerts.
- **Learn:** report SQL (`GROUP BY`, window functions), batch jobs, changing a live schema safely.
- **Done when:** each report matches hand-calculated numbers in tests, and every new alert fires once.
- Let beta feedback reorder this list. If users ask for something else first, build that first.

### Success metrics (first 3 months)
| Metric | Target |
|---|---|
| Slip amount read correctly | ≥ 95% |
| Transactions added by slip | ≥ 50% |
| Users still logging after 30 days | ≥ 40% |
| Days per week a user opens the app | ≥ 4 |
| Free → Premium | ≥ 5% |
| AI cost per active user per month | Lower than Premium price |

---

## 10. Out of scope for the MVP

Voice entry · shared family accounts · multiple savings goals · THB/USD currencies · export to Excel/PDF · direct bank connections · automatic premium payment · Casbin · microservices.

---

## 11. How to run each phase with Claude Code

1. Start a session with: "We are starting Phase N. Read KEPKEP_BACKEND_PLAN.md and CLAUDE.md, explain the concepts for this phase briefly, then give me the first small task."
2. Write the code yourself first. Ask Claude Code to review it, not to write it for you.
3. Ask "why does this work?" after each task.
4. Before closing a phase: run tests and lint, check the "Done when" line, and update `LEARNING_LOG.md`.
5. Commit at the end of each task with a clear message.
