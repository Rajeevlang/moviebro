# 🚀 MovieHub DevOps Journey — Milestone Progress (Phase 1 & 2)

## 📌 Project Overview
**Application:** MovieHub Microservices (Go API Gateway, User Service, Movie Service, PostgreSQL, MongoDB, Redis, and React Vite Frontend with Nginx).  
**Deployment Target:** AWS EC2 Instance (`t3.medium`, Ubuntu Linux) with Elastic IP.  
**CI/CD Engine:** Distributed Jenkins Architecture (Local Controller / Master + Cloud EC2 Worker Node).  
**Live Production URL:** [https://moviehub.nostackdev.online](https://moviehub.nostackdev.online)

---

## 🏛️ System Architecture

```
[ Developer Local Machine ]
    ├── Jenkins Controller (Docker container: `jenkins-master`)
    ├── Git Source Code & Declarative Pipeline (`Jenkinsfile`)
    └── SSH Tunnel (Port 22 with private key credentials)
             │
             ▼
[ AWS EC2 Instance (15.207.161.137) ]
    ├── Jenkins Worker Agent (`remoting.jar` on Java 21)
    ├── Docker Engine & Docker Compose v2
    └── Microservices Stack:
         ├── moviehub-frontend-1    (Port 80 & 443 — React + Nginx SSL)
         ├── moviehub-api-gateway-1 (Port 8080 — Reverse Proxy)
         ├── moviehub-user-service-1(Port 8082 — Go + PostgreSQL)
         ├── moviehub-movie-service-1 (Port 8081 — Go + MongoDB + Redis)
         ├── postgresdb-1           (PostgreSQL 15)
         ├── mongodb-1              (MongoDB 7.0)
         └── redisdb-1              (Redis 7)
```

---

## ✅ Milestones Accomplished

### 1. Repository Cleanliness & Public Security Hardening
- Implemented a comprehensive [`.gitignore`](.gitignore) preventing leaks to public GitHub:
  - Blocked all `.env`, `.env.local`, and backup files.
  - Excluded TLS/SSL private keys (`certs/*.pem`, `certs/*.key`).
  - Filtered out ~88MB of pre-compiled Go binaries (`moviehub/bin/*`).
  - Preserved configuration templates ([`.env.example`](moviehub/app.env.example)).

### 2. Distributed Jenkins Master-Agent Cluster
- Launched local Jenkins Controller via Docker with persistent volume storage.
- Provisioned and hardened AWS EC2 instance with an Elastic IP (`15.207.161.137`).
- Connected local master to the EC2 worker node using the `ssh-slaves` remoting protocol.
- Aligned runtime bytecode to **OpenJDK 21** across both nodes.

### 3. Production-Ready Declarative Pipeline (`Jenkinsfile`)
- **Node Targeting:** Set to execute on `agent { label 'ec2-agent' }`.
- **Zero-Dependency Host Testing:**
  - Backend unit tests and static analysis run inside throwaway `golang:alpine` containers with persistent module caching (`-v go-mod-cache:/go/pkg/mod`).
  - Frontend production builds run inside throwaway `node:20-alpine` containers.
  - **Parallelism:** Tests run simultaneously to cut CI time in half.
- **Continuous Deployment:** Zero-downtime rolling updates via `docker compose -f docker-compose.prod.yml up -d --remove-orphans`.
- **Automated SSL Injection:** Dynamically mounts SSL certificates from host storage (`/home/ubuntu/certs`).
- **Disk Hygiene:** Automatically prunes dangling Docker build layers on every run (`docker image prune -f`).

### 4. Live DNS & SSL Routing
- Configured Cloudflare DNS with Full SSL encryption.
- Routed traffic to the EC2 Elastic IP serving valid HTTPS responses on port 443.

---

## 📊 Live Verification Checklist

| Service | Port | Internal / External | Status |
| :--- | :--- | :--- | :--- |
| **Frontend UI** | `80`, `443` | External (Cloudflare) | ✅ Active (HTTP 200 OK) |
| **API Gateway** | `8080` | External / Internal | ✅ Active |
| **Movie Service** | `8081` | Internal (Bridge Network) | ✅ Healthy |
| **User Service** | `8082` | Internal (Bridge Network) | ✅ Healthy |
| **PostgreSQL** | `5432` | Internal (Bridge Network) | ✅ Healthy |
| **MongoDB** | `27017`| Internal (Bridge Network) | ✅ Healthy |
| **Redis** | `6379` | Internal (Bridge Network) | ✅ Healthy |
