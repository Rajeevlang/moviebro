# PostgreSQL Reference — User Service

PostgreSQL stores everything for the **user-service**: accounts, credentials, Google IDs.

## Where it lives in this project

| What | Value |
|---|---|
| Image | `postgres:15-alpine` |
| Dev container name | `postgres` (`docker-compose.yml`) |
| Prod container name | `moviehub-postgres` (`docker-compose.prod.yml`) |
| Host port (dev only) | `127.0.0.1:5434` → container `5432` |
| Host port (prod) | **not exposed** — internal network only |
| Database / User / Password | `app_db` / `app_user` / `$POSTGRES_PASSWORD` (default `secure_pass_123`, set real value in `.env`) |
| Schema init | `internal/user/schema.sql` auto-runs on **first** startup via `/docker-entrypoint-initdb.d/01-schema.sql` |
| Data volume | `postgres_data` |

## Quick health checks

```bash
make health                      # runs pg_isready among other checks
docker compose exec postgresdb pg_isready -U app_user -d app_db
curl -sf http://localhost:8083/healthcheck   # user-service (proves it can reach PG)
```

## Connecting

```bash
# Easiest — Makefile shell alias
make db-shell-postgres

# Manual (dev compose)
docker compose exec postgresdb psql -U app_user -d app_db

# From host machine (dev: port 5434)
psql "postgresql://app_user:secure_pass_123@127.0.0.1:5434/app_db?sslmode=disable"

# Prod EC2 (no exposed port — must exec)
docker compose -f docker-compose.prod.yml --env-file .env exec postgresdb \
    psql -U "${POSTGRES_USER:-app_user}" -d "${POSTGRES_DB:-app_db}"
```

## Schema (what actually exists)

```
users
├── id            UUID PK, default gen_random_uuid()
├── username      VARCHAR(50) UNIQUE NOT NULL
├── email         VARCHAR(255) UNIQUE NOT NULL     ← login identifier
├── password_hash VARCHAR(255) NULL                ← NULL for Google-only accounts
├── google_id     VARCHAR(255) UNIQUE NULL
├── created_at    TIMESTAMPTZ DEFAULT now()
└── updated_at    TIMESTAMPTZ DEFAULT now()

Indexes: idx_users_email, idx_users_username, idx_users_google_id
```

## Useful SQL for this project

```sql
-- Count users
SELECT count(*) FROM users;

-- Recent signups
SELECT username, email, created_at,
       (google_id IS NOT NULL) AS is_google_account
FROM users
ORDER BY created_at DESC
LIMIT 20;

-- Find one user by email (login problems)
SELECT * FROM users WHERE email = 'someone@example.com';

-- Check who has a password vs Google-only
SELECT count(*) FILTER (WHERE password_hash IS NOT NULL) AS password_users,
       count(*) FILTER (WHERE google_id IS NOT NULL)     AS google_users
FROM users;

-- Delete a test user (cascades nothing — playlists live in MongoDB)
DELETE FROM users WHERE email = 'test@example.com';

-- Table sizes
\dt+            -- psql meta command
```

Handy `psql` meta commands: `\dt` list tables · `\d users` describe · `\x` expanded output · `\q` quit.

## Backups & restore

```bash
# Backup (prod EC2)
docker compose -f docker-compose.prod.yml --env-file .env exec -T postgresdb \
  pg_dump -U app_user app_db > backup_$(date +%F).sql

# Restore
cat backup_2026-08-23.sql | docker compose exec -T postgresdb \
  psql -U app_user -d app_db
```

## Troubleshooting

| Symptom | Check |
|---|---|
| user-service won't start | `docker compose logs user-service` — usually bad `POSTGRES_URI` |
| Schema missing | Only loads on **first** boot of an empty volume. Fix: `make clean && make up` (destroys data!) or run schema.sql manually |
| Password auth fails | `.env` changed after first start? Postgres only reads `POSTGRES_PASSWORD` when initializing a fresh volume |
| Connection refused | Service names are docker-network DNS: use `postgresdb:5432` inside containers, `127.0.0.1:5434` from dev host |

> ⚠️ Changing DB credentials requires recreating the volume: `docker compose down -v && up -d`.
