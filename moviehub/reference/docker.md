# Docker Reference — MovieHub

Two compose files drive everything:

| File | Purpose | Services |
|---|---|---|
| `docker-compose.yml` | Local dev — builds images, exposes DB ports on loopback | redis, postgres, mongo, movie, user, gateway |
| `docker-compose.prod.yml` | EC2 prod — hardened (limits, log rotation, healthcheck-gated startup), DBs **not** exposed, gateway on `127.0.0.1:8080` only | same |

## Day-to-day commands

The Makefile wraps almost everything:

```bash
make up          # start everything (detached)
make down        # stop, keep data volumes
make restart     # down + up
make ps          # what's running
make logs        # tail all logs
make build       # rebuild images
make clean       # ⚠ down -v — DESTROYS all database data
```

Infrastructure only (native `go run` dev workflow):

```bash
make db-up       # just mongo + postgres + redis
make run-movie   # go run ./cmd/movie  (port 8081)
make run-user    # go run ./cmd/user   (port 8082)
make run-gateway # go run ./cmd/gateway(port 8080)
```

Health & shells:

```bash
make health              # curl every service /healthcheck + pg/redis probes
make db-shell-mongo
make db-shell-postgres
make db-shell-redis
```

Raw compose equivalents:

```bash
docker compose up -d --build
docker compose ps
docker compose down            # keep volumes
docker compose down -v         # ⚠ wipe volumes too
docker compose restart api-gateway
```

## Checking logs (the important part)

```bash
# All services, follow
docker compose logs -f

# One service, follow
docker compose logs -f api-gateway
docker compose logs -f movie-service
docker compose logs -f user-service

# Last 100 lines
docker compose logs --tail 100 api-gateway

# Since a point in time
docker compose logs --since 10m api-gateway      # last 10 minutes
docker compose logs --since 2026-08-23T10:00:00 api-gateway

# Prod file (different project name / env file!)
docker compose -f docker-compose.prod.yml --env-file .env logs -f api-gateway

# Bypass compose — plain docker logs works too
docker logs -f moviehub-api-gateway          # prod container names:
docker logs -f moviehub-movie-service        #   moviehub-* prefix
docker logs -f moviehub-user-service
docker logs -f moviehub-postgres
docker logs -f moviehub-mongodb
docker logs -f moviehub-redis
```

Prod containers use the json-file driver with rotation (`max-size: 10m`, `max-file: 3`) — logs never grow unbounded.

What to grep for:

```bash
docker compose logs api-gateway | grep -i "unauthorized\|rate limit"  # auth/ratelimit hits
docker compose logs movie-service | grep -i "error\|mongo"
docker compose logs user-service | grep -i "error\|postgres\|jwt"
```

## Production deploy flow (EC2)

Normally done by GitHub Actions (`.github/workflows/cd-ec2.yml`), but manually:

```bash
cd /opt/moviehub                       # repo must be present, .env must exist
export DOCKER_REGISTRY=<dockerhub-user>
export APP_VERSION=<tag>               # e.g. 42-abc1234; or 'latest'

docker compose -f docker-compose.prod.yml --env-file .env pull \
    movie-service user-service api-gateway

docker compose -f docker-compose.prod.yml --env-file .env up -d --no-deps \
    movie-service user-service api-gateway    # --no-deps: don't touch databases

curl -sf http://127.0.0.1:8080/healthcheck
```

Rollback = re-run with previous `APP_VERSION`.

## Inspection & debugging

```bash
docker ps                                    # running containers
docker ps -a                                 # include crashed/exited
docker inspect moviehub-api-gateway | jq '.[0].State.Health'   # Dockerfile HEALTHCHECK status
docker stats                                 # live CPU/mem (prod limits: 128M per app svc)
docker exec -it moviehub-api-gateway sh      # shell into a container
docker exec -it moviehub-postgres psql -U app_user -d app_db
docker network ls && docker network inspect <project>_moviehub-net
docker volume ls                             # postgres_data, mongodb_data, redis_data
```

Image housekeeping:

```bash
docker image ls --filter dangling=true
docker image prune -f                        # safe cleanup after deploys
docker system df                             # disk usage report
```

## Port map cheat sheet

| Service | Container port | Dev host port | Prod host port |
|---|---|---|---|
| api-gateway | 8080 | `127.0.0.1:8080` | `127.0.0.1:8080` (behind tunnel/proxy) |
| movie-service | 8081 | `127.0.0.1:8082` | not exposed |
| user-service | 8082 | `127.0.0.1:8083` | not exposed |
| postgresdb | 5432 | `127.0.0.1:5434` | not exposed |
| mongodb | 27017 | `127.0.0.1:27018` | not exposed |
| redisdb | 6379 | `127.0.0.1:6379` | not exposed |

> Everything binds to `127.0.0.1` — nothing is reachable from outside the machine unless you publish it explicitly.

## Common problems

| Symptom | Fix |
|---|---|
| `up` fails: port already allocated | Something else owns the port → `lsof -i :8080`, stop it or change `HOST_*_PORT` in `.env` |
| App services keep restarting in prod | Missing/unhealthy DBs — they wait on `condition: service_healthy`; check `docker compose ps` and DB container logs |
| Gateway 502 on some routes | Backend service down — `docker compose logs movie-service` |
| Images won't pull on EC2 | `DOCKER_REGISTRY` wrong / `docker login` missing on host |
| Data gone after redeploy | Someone ran `down -v` — volumes were deleted |
