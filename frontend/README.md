# CineDeck — React Frontend

A movie catalog UI for the Gin Movie API gateway (port 8080).

**Features**
- 🔐 Register / Login / Logout (JWT stored in `localStorage`, auto-logout on expiry)
- 🎬 Browse movies with search, genre chips, year/rating/actor/director filters + pagination
- ⭐ Movie details with cast, reviews, and score submission
- 🎞️ Playlists: create, delete, public/private, save/remove movies
- 🛠️ Admin panel (`/admin`): dashboard stats + full movie CRUD

## Run (dev)

```bash
cd frontend
npm install
npm run dev        # http://localhost:5173
```

`/api` is proxied to `http://localhost:8080` (your gateway) — start the Go gateway first.
No CORS setup needed.

## Admin access

There are no roles in the backend yet, so admin gating is UI-only via an email allowlist:

```bash
# frontend/.env
VITE_ADMIN_EMAILS=you@example.com,boss@example.com
```

> ⚠️ Anyone can still call POST/PATCH/DELETE `/api/v1/movies/...` directly with a valid JWT.
> Add role checks to the gateway/user-service before exposing this publicly.

## Production build

```bash
npm run build      # outputs frontend/dist/
```

To serve it from the Go gateway, copy `dist/*` into the gateway's `./public` folder:

```bash
rm -rf public/* && cp -r frontend/dist/* public/
```

Routing uses a `HashRouter`, so deep links work behind the plain `http.FileServer`
without adding SPA-fallback logic to the gateway.

### Optional: switch to clean URLs (BrowserRouter)

1. In `src/App.jsx`, replace `HashRouter` with `BrowserRouter`.
2. Add an SPA fallback in `cmd/gateway/main.go`:
   ```go
   mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
       if strings.HasPrefix(r.URL.Path, "/api") {
           http.NotFound(w, r)
           return
       }
       http.ServeFile(w, r, "./public/index.html")
   })
   ```

## API mapping

| Frontend | Backend |
|---|---|
| Register / Login | `POST /api/v1/users/register` · `POST /api/v1/users/login` |
| List movies (+filters) | `GET /api/v1/movies/?year=&min_rating=&actor=&director=` |
| Movie detail / reviews | `GET /api/v1/movies/:id` · `POST /api/v1/movies/:id/reviews` |
| Create / edit / delete movie (admin) | `POST / PATCH / DELETE /api/v1/movies/:id` |
| Playlists CRUD + membership | `POST /api/v1/list/`, `GET /api/v1/list/user/:uid`, `DELETE /api/v1/list/:id`, `POST/DELETE /api/v1/list/:id/movies/:movie_id` |
# moviedeskf
# moviedesk-frontend
