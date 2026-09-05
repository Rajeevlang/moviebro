# API Endpoints Reference — MovieHub

All public traffic goes through the **api-gateway**. It reverse-proxies to the two
backend services and injects `X-User-ID` (from the JWT) into proxied requests.

```
Gateway :8080 ── /api/v1/movies, /api/v1/list  → movie-service :8081 (MongoDB)
             ├─ /api/v1/users                  → user-service  :8082 (PostgreSQL)
             ├─ /healthcheck                   → gateway itself
             └─ /metrics                       → Prometheus
```

## Authentication model

- **JWT Bearer tokens** issued by user-service on login/register.
- Gateway middleware (`internal/gateway/middleware/auth.go`):
  - **Public without token:** any `GET` under `/api/v1/movies*`
  - **Everything else** (`POST/PATCH/DELETE`, all playlists, users) requires
    `Authorization: Bearer <token>`
- Rate limit: **100 req/s, burst 200 per IP**, applied to both proxies.
- CORS: only for origins listed in `CORS_ALLOWED_ORIGINS` (Cloudflare Pages domain in prod).

## Base URLs

| Environment | URL |
|---|---|
| Local dev | `http://localhost:8080/api/v1` |
| Production | `https://api.<your-domain>/api/v1` |
| Direct service access (dev only) | movie `http://localhost:8082`, user `http://localhost:8083` |

---

## 1. User Service — `/api/v1/users`

| Method | Path | Auth | Description |
|---|---|---|---|
| POST | `/api/v1/users/register` | ❌ | Create account, returns JWT |
| POST | `/api/v1/users/login` | ❌ | Email+password → JWT |
| POST | `/api/v1/users/google` | ❌ | Google Sign-In via OIDC `id_token` |

### Register
```bash
curl -s -X POST http://localhost:8080/api/v1/users/register \
  -H "Content-Type: application/json" \
  -d '{"username":"raj","email":"raj@example.com","password":"SuperSecret123"}'
# → { "token": "<jwt>", ... }
```

### Login
```bash
curl -s -X POST http://localhost:8080/api/v1/users/login \
  -H "Content-Type: application/json" \
  -d '{"email":"raj@example.com","password":"SuperSecret123"}'
# → { "token": "<jwt>" }
```

Validation: username 3–50 chars, valid email, password ≥ 8 chars.
Google flow: `POST /api/v1/users/google` with `{"id_token": "<google-id-token>"}`.

> ⚠️ Note: the OAuth redirect endpoints (`/googlelogin`, `/auth/callback`) are registered at
> the **root of user-service**, not under `/api/v1/users`. The gateway only proxies
> `/api/v1/users/*`, so those two routes are reachable only by hitting the service directly
> (dev: port 8083). The gateway OAuth callback path in `.env.example`
> (`/api/v1/users/auth/google/callback`) will 404 until these routes are moved under the
> proxy prefix or added to the gateway mux.

## 2. Movie Catalog — `/api/v1/movies`

| Method | Path | Auth | Description |
|---|---|---|---|
| GET | `/api/v1/movies` | ❌ public | List movies |
| GET | `/api/v1/movies/:id` | ❌ public | Single movie incl. reviews & rating |
| POST | `/api/v1/movies` | ✅ | Add movie |
| PATCH | `/api/v1/movies/:id` | ✅ | Update movie |
| DELETE | `/api/v1/movies/:id` | ✅ | Delete movie |
| POST | `/api/v1/movies/:id/reviews` | ✅ | Add a review (updates aggregate rating) |

```bash
TOKEN="<paste-jwt>"
AUTH="Authorization: Bearer $TOKEN"

# Browse (no token needed)
curl -s http://localhost:8080/api/v1/movies

# Create (auth required)
curl -s -X POST http://localhost:8080/api/v1/movies \
  -H "$AUTH" -H "Content-Type: application/json" \
  -d '{
    "title": "Interstellar",
    "description": "Space farmers.",
    "release_year": 2014,
    "genres": ["Sci-Fi", "Drama"],
    "directors": [{"name":"Christopher Nolan"}],
    "actors":    [{"name":"Matthew McConaughey"}]
  }'

# Review a movie
curl -s -X POST http://localhost:8080/api/v1/movies/<movieId>/reviews \
  -H "$AUTH" -H "Content-Type: application/json" \
  -d '{"comment":"Masterpiece.","score":9.5}'
```

## 3. Playlists — `/api/v1/list`

All playlist routes require auth; ownership is tied to the JWT subject.

| Method | Path | Description |
|---|---|---|
| POST | `/api/v1/list` | Create playlist |
| GET | `/api/v1/list/:id` | Get one playlist |
| GET | `/api/v1/list/user/:user_id` | All playlists of a user |
| DELETE | `/api/v1/list/:id` | Delete playlist |
| POST | `/api/v1/list/:id/movies/:movie_id` | Add movie to playlist |
| DELETE | `/api/v1/list/:id/movies/:movie_id` | Remove movie from playlist |

```bash
curl -s -X POST http://localhost:8080/api/v1/list \
  -H "$AUTH" -H "Content-Type: application/json" \
  -d '{"name":"Watch Later"}'

curl -s -X POST http://localhost:8080/api/v1/list/<playlistId>/movies/<movieId> -H "$AUTH"
curl -s http://localhost:8080/api/v1/list/user/<userId> -H "$AUTH"
```

## 4. Ops endpoints (gateway, user & movie services each expose their own)

| Method | Path | Purpose |
|---|---|---|
| GET | `/healthcheck` | LB/Docker/CD smoke test → `{"status":"ok",...}` |
| GET | `/metrics` | Prometheus metrics (RED metrics per route) |

```bash
curl -sf http://localhost:8080/healthcheck   # gateway
curl -sf http://localhost:8082/healthcheck   # movie-service direct (dev)
curl -sf http://localhost:8083/healthcheck   # user-service direct (dev)
curl -s  http://localhost:8080/metrics | head
```

> 🔒 In production `/metrics` is publicly reachable through the proxy — consider blocking it
> at Cloudflare/proxy level or adding auth.

## Error shapes

The frontend's `apiError()` helper understands these backend error formats:
- `{"error": "string"}` — simple message
- `{"errors": {"field": "msg"}}` — validation map
- `{"errors": ["msg", ...]}` — validation list
- `{"message": "string"}` — generic

Auth failures return plain-text `401 Unauthorized: ...` from the gateway middleware.

## Full endpoint dump (quick test script)

```bash
GW=http://localhost:8080
for p in /healthcheck "/api/v1/movies" ; do echo "== $p"; curl -si "$GW$p" | head -1; done
echo "== unauthorized write blocked"; curl -si -X POST $GW/api/v1/list -H 'Content-Type: application/json' -d '{}' | head -1   # expect 401
```

A ready-made Postman collection exists: `Movie_API.postman_collection.json`.
