# 🎬 Open-Source Movie Database Bulk Seeder Guide
### Populate GinMovieAPI with Rich Movie Data from Free Open-Source APIs & Curated Blockbuster Catalogs

This guide explains how to use the automated movie database seeder to populate your **GinMovieAPI** instance with dozens of real movies, high-resolution posters, release years, and genre tags.

---

## 📑 Table of Contents

1. [🌟 Supported Open Data Sources](#1-🌟-supported-open-data-sources)
2. [🚀 Usage Instructions](#2-🚀-usage-instructions)
   - [Local Docker Compose](#local-docker-compose)
   - [Minikube Cluster](#minikube-cluster)
   - [AWS EKS / EC2 Cloud Deployments](#aws-eks--ec2-cloud-deployments)
3. [⚙️ Available CLI Options & Flags](#3-⚙️-available-cli-options--flags)
4. [🔍 Verifying Seeded Movie Data](#4-🔍-verifying-seeded-movie-data)

---

## 1. 🌟 Supported Open Data Sources

The bulk seeder supports multiple data ingestion modes:

1. **TVMaze Open API** (100% free, no API key required):
   - Fetches live movie and show metadata directly from `https://api.tvmaze.com/shows`.
   - Cleans HTML formatting and maps genres and posters automatically.
2. **Curated Blockbuster Catalog**:
   - Built-in dataset of top-rated movies (*Oppenheimer, Interstellar, The Dark Knight, Dune: Part Two, Blade Runner 2049, Spirited Away, Gladiator, Parasite, Avengers: Endgame, The Matrix*).
3. **Keyword-Based Search**:
   - Allows dynamically searching open movie databases (e.g. `--query "Marvel"`, `--query "Cyberpunk"`).

---

## 2. 🚀 Usage Instructions

### Local Docker Compose
When running the app on `http://localhost:8080`:

```bash
# Seed 30 movies using Python seeder
python3 scripts/seed_movies_opensource.py --target http://localhost:8080 --source all --count 30

# Or using the pure Bash script
bash scripts/seed_movies_opensource.sh http://localhost:8080
```

---

### Minikube Cluster
When running on Minikube:

```bash
# 1. Forward API Gateway port in a separate terminal
kubectl port-forward svc/api-gateway 8080:8080 -n moviehub

# 2. Run seeder
python3 scripts/seed_movies_opensource.py --target http://localhost:8080 --count 40
```

---

### AWS EKS / EC2 Cloud Deployments
Point the seeder to your public AWS Application Load Balancer or EC2 endpoint:

```bash
python3 scripts/seed_movies_opensource.py \
  --target "http://<YOUR_AWS_ALB_DNS_NAME>" \
  --source all \
  --count 50
```

---

## 3. ⚙️ Available CLI Options & Flags

| Flag | Default | Description |
| :--- | :--- | :--- |
| `--target` | `http://localhost:8080` | URL of API Gateway or Movie Service. |
| `--source` | `all` | Data source: `all`, `curated`, `tvmaze`, `search`. |
| `--query` | `""` | Search query for open database (e.g. `--query "Action"`). |
| `--count` | `25` | Maximum number of movies to ingest. |

### Examples:
```bash
# Search and seed 20 Sci-Fi / Space movies
python3 scripts/seed_movies_opensource.py --query "Space" --count 20

# Seed only the curated Oscar-winning catalog
python3 scripts/seed_movies_opensource.py --source curated
```

---

## 4. 🔍 Verifying Seeded Movie Data

### Query via REST API:
```bash
curl -s http://localhost:8080/api/v1/movies | jq .
```

### View in Web UI:
Open **`http://localhost:8080/`** in your browser to browse the visual movie grid with posters and genre badges!
