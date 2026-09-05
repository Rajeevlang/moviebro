# 👥 User Database Seeder Guide (Admin & Dummy Users)
### Automatically Provision Administrator and Dummy User Accounts in GinMovieAPI PostgreSQL

This guide explains how to use the automated user database seeder script to populate your **GinMovieAPI** instance with a primary administrator account and 10 realistic test users.

---

## 📑 Table of Contents

1. [🔑 Generated User Accounts & Credentials](#1-🔑-generated-user-accounts--credentials)
2. [🚀 Usage Instructions](#2-🚀-usage-instructions)
   - [Local Docker Compose](#local-docker-compose)
   - [Minikube Cluster](#minikube-cluster)
   - [AWS EKS / EC2 Deployments](#aws-eks--ec2-deployments)
3. [🧪 Verifying User Authentication (JWT)](#3-🧪-verifying-user-authentication-jwt)

---

## 1. 🔑 Generated User Accounts & Credentials

### Primary Administrator:
| Field | Value |
| :--- | :--- |
| **Role** | 👑 Primary System Administrator |
| **Username** | `admin` |
| **Email** | `admin@cinedeck.dev` |
| **Password** | `AdminPassword2026!` |

### Dummy Users Catalog:
| Role / Persona | Username | Email | Password |
| :--- | :--- | :--- | :--- |
| **Film Critic** | `alex_morgan` | `alex.morgan@cinedeck.dev` | `UserPass2026!` |
| **VIP Member** | `sarah_connor` | `sarah.connor@cyberdyne.io` | `UserPass2026!` |
| **Premium Member** | `john_wick` | `john.wick@continental.org` | `UserPass2026!` |
| **Developer** | `neo_anderson` | `neo.anderson@matrix.dev` | `UserPass2026!` |
| **Executive Producer** | `bruce_wayne` | `bruce.wayne@waynecorp.com` | `UserPass2026!` |
| **Tech Lead** | `tony_stark` | `tony.stark@starkenterprises.com` | `UserPass2026!` |
| **Movie Enthusiast** | `elena_rostova` | `elena.rostova@cinema.net` | `UserPass2026!` |
| **Director** | `david_fincher` | `david.fincher@directors.org` | `UserPass2026!` |
| **Member** | `clara_oswald` | `clara.oswald@tardis.space` | `UserPass2026!` |
| **Historian** | `marcus_aurelius` | `marcus.aurelius@rome.ancient` | `UserPass2026!` |

---

## 2. 🚀 Usage Instructions

### Local Docker Compose
When running the app on `http://localhost:8080`:

```bash
# Seed admin and all 10 dummy users
python3 scripts/seed_users.py --target http://localhost:8080

# Or using the pure Bash script
bash scripts/seed_users.sh http://localhost:8080
```

---

### Minikube Cluster
When running on Minikube:

```bash
# 1. Forward API Gateway port in a separate terminal
kubectl port-forward svc/api-gateway 8080:8080 -n moviehub

# 2. Run seeder
python3 scripts/seed_users.py --target http://localhost:8080
```

---

### AWS EKS / EC2 Deployments
Point the seeder to your public AWS Application Load Balancer or EC2 instance:

```bash
python3 scripts/seed_users.py --target "http://<YOUR_AWS_ALB_DNS_NAME_OR_EC2_IP>"
```

---

## 3. 🧪 Verifying User Authentication (JWT)

You can log in with any seeded user using `curl`:

### Login as Admin:
```bash
curl -s -X POST "http://localhost:8080/api/v1/users/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@cinedeck.dev", "password": "AdminPassword2026!"}' | jq .
```

**Expected Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### Login as Dummy User (`alex_morgan`):
```bash
curl -s -X POST "http://localhost:8080/api/v1/users/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "alex.morgan@cinedeck.dev", "password": "UserPass2026!"}' | jq .
```
