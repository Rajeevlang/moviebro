# MongoDB Reference — Movie Service

MongoDB stores everything for the **movie-service**: the catalog, reviews, and playlists.

## Where it lives in this project

| What | Value |
|---|---|
| Image | `mongo:7.0` |
| Dev container name | `mongodb` (`docker-compose.yml`) |
| Prod container name | `moviehub-mongodb` (`docker-compose.prod.yml`) |
| Host port (dev only) | `127.0.0.1:27018` → container `27017` |
| Host port (prod) | **not exposed** — internal network only |
| Root credentials | `$MONGO_ROOT_USER` / `$MONGO_ROOT_PASSWORD` (defaults `admin` / `password123`, set real values in `.env`) |
| Database name | `movie_catalog` (hardcoded in `cmd/movie/main.go:42`) |
| Collections | `movies`, `playlists` (`cmd/movie/main.go:43-44`) |
| Data volume | `mongodb_data` |

Connection string used by movie-service:
```
mongodb://${MONGO_ROOT_USER}:${MONGO_ROOT_PASSWORD}@mongodb:27017/?authSource=admin
```

## Quick health checks

```bash
make health                                   # includes redis/pg; mongo via compose below
docker compose exec mongodb mongosh --eval "db.adminCommand('ping')"
curl -sf http://localhost:8082/healthcheck    # movie-service (proves it can reach Mongo)
```

## Connecting

```bash
# Easiest — Makefile shell alias
make db-shell-mongo

# Manual (dev compose) — authenticates against admin DB
docker compose exec mongodb mongosh -u admin -p password123 --authenticationDatabase admin

# From host machine (dev: port 27018)
mongosh "mongodb://admin:password123@127.0.0.1:27018/?authSource=admin"

# Prod EC2
docker compose -f docker-compose.prod.yml --env-file .env exec mongodb \
  mongosh -u "$MONGO_ROOT_USER" -p "$MONGO_ROOT_PASSWORD" --authenticationDatabase admin
```

## Data model

**`movie_catalog.movies`**
```js
{
  _id: ObjectId,
  title: String,
  description: String,
  release_year: Number,
  genres: [String],                    // e.g. ["Action", "Sci-Fi"]
  directors: [ { name: String, role: String } ],
  actors:    [ { name: String, role: String } ],
  rating: { average: Number, count: Number },
  reviews: [
    { _id: ObjectId, user_id: String, comment: String, score: Number, created_at: Date }
  ]
}
```

**`movie_catalog.playlists`**
```js
{
  _id: ObjectId,
  user_id: String,        // JWT subject from gateway (X-User-ID)
  name: String,
  movie_ids: [ObjectId],
  created_at: Date,
  updated_at: Date
}
```

## Useful mongosh commands for this project

```js
show dbs
use movie_catalog
show collections

// Counts
db.movies.countDocuments()
db.playlists.countDocuments()

// Browse
db.movies.find().limit(5).pretty()
db.movies.find({}, { title: 1, release_year: 1, "rating.average": 1 })

// Filter examples
db.movies.find({ genres: "Sci-Fi" })
db.movies.find({ release_year: { $gte: 2020 } })
db.movies.find({ title: { $regex: /matrix/i } })

// Top rated
db.movies.find().sort({ "rating.average": -1 }).limit(10)

// A user's playlists
db.playlists.find({ user_id: "<uuid-from-users-table>" })

// Reviews of a movie with at least one review
db.movies.find({ "reviews.0": { $exists: true } }, { title: 1 })

// Delete test data
db.playlists.deleteMany({ user_id: "test-user-id" })
db.movies.deleteOne({ _id: ObjectId("<id>") })

// Drop everything and reseed (careful!)
db.dropDatabase()
```

Seed data afterwards: `make seed` or `./scripts/seed_opensource_movies.sh`.

## Backups & restore

```bash
# Dump (prod EC2)
docker compose -f docker-compose.prod.yml --env-file .env exec -T mongodb \
  mongodump --authenticationDatabase admin \
  -u "$MONGO_ROOT_USER" -p "$MONGO_ROOT_PASSWORD" \
  --archive --gzip > mongo_backup_$(date +%F).archive.gz

# Restore
cat mongo_backup_2026-08-23.archive.gz | docker compose exec -T mongodb \
  mongorestore --authenticationDatabase admin \
  -u "$MONGO_ROOT_USER" -p "$MONGO_ROOT_PASSWORD" --archive --gzip
```

## Troubleshooting

| Symptom | Check |
|---|---|
| movie-service won't start | `docker compose logs movie-service` — usually bad `MONGO_URI` / auth source missing |
| `auth failed` in logs | Credentials changed after volume was initialized — recreate volume or create user manually |
| Empty API responses | No seed data — run `make seed` |
| Connection refused | Inside containers use host `mongodb:27017`; from dev host use `127.0.0.1:27018` |

> ⚠️ Changing `MONGO_ROOT_PASSWORD` requires recreating the volume: `docker compose down -v && up -d`.
