# 🛠️ The DevOps Troubleshooting Playbook & Field Guide

A collection of real-world diagnostic frameworks, triage heuristics, and case studies developed during the deployment of the MovieHub distributed microservices stack.

---

## 🏛️ 1. The 5-Layer DevOps Diagnostic Stack

Whenever an incident or pipeline failure occurs, walk up this 5-layer hierarchy from bottom to top:

```
┌────────────────────────────────────────────────────────┐
│ Layer 5: Pipeline & Application Logic (Groovy / Code)  │ ◄── Syntax, Exit codes, App bugs
├────────────────────────────────────────────────────────┤
│ Layer 4: Runtimes & Toolchains (JVM / Go / Node / npm) │ ◄── Java 21 vs 17, Go 1.25 vs 1.23
├────────────────────────────────────────────────────────┤
│ Layer 3: Host Permissions & Sockets (/var/run, users)  │ ◄── usermod -aG docker, chmod, chown
├────────────────────────────────────────────────────────┤
│ Layer 2: DNS & Name Resolution (/etc/resolv.conf)      │ ◄── Could not resolve host, NXDOMAIN
├────────────────────────────────────────────────────────┤
│ Layer 1: Network & Firewalls (AWS SG, Routing, Ports)  │ ◄── Connection timed out vs refused
└────────────────────────────────────────────────────────┘
```

---

## 🔬 2. The Golden Rules of Error Interpretation

| Error Message | Meaning | Layer | Diagnostic Action |
| :--- | :--- | :--- | :--- |
| **`Connection timed out`** | Packets sent, but dropped silently with no response. | Layer 1 (Firewall) | Check AWS Security Group inbound rules, Subnet Route Tables, and OS firewalls (`ufw`). |
| **`Connection refused`** | Server received the packet, but **no service is listening** on that port. | Layer 1 / 3 | Check if Docker container crashed or service failed to bind. |
| **`Could not resolve host`** | Server asked its DNS resolver for an IP, and got silence or NXDOMAIN. | Layer 2 (DNS) | Check `/etc/resolv.conf`, test with `dig` or `nslookup`. |
| **`Permission denied` / `EACCES`**| Linux kernel blocked the syscall due to file permissions or group membership. | Layer 3 (Permissions)| Missing group (`usermod -aG docker`), wrong socket permissions (`docker.sock`). |
| **`UnsupportedClassVersionError`**| Bytecode was compiled for a newer JVM than the runtime. | Layer 4 (Runtime) | Bytecode 65 = Java 21; Bytecode 61 = Java 17; Bytecode 55 = Java 11. |
| **`requires go >= X.Y`** | Go module enforces minimum language version. | Layer 4 (Runtime) | Upgrade compiler container image tag (e.g., `golang:alpine`). |

---

## 📚 3. Case Studies: Today's Real-World Incidents & Solutions

###  Incident 1: Jenkins Master Port Not Accessible on `localhost:8080`
- **Symptom:** Container was running, but browser could not connect to `http://localhost:8080`.
- **Root Cause:** Container lost attachment to Docker's default `bridge` network, clearing port mappings in `.NetworkSettings.Ports`.
- **Diagnostic Command:**
  ```bash
  docker inspect jenkins-master --format '{{json .NetworkSettings.Ports}}'
  ```
- **Fix:**
  ```bash
  docker network connect bridge jenkins-master
  ```

---

### Incident 2: SSH Connection Timed Out Between Master and EC2 Agent
- **Symptom:** Jenkins reported `IOException: Connect timed out` attempting to connect to `15.207.161.137:22`.
- **Root Cause:** AWS Security Group blocked port 22 or restricted it to a stale dynamic IP.
- **Diagnostic Command:**
  ```bash
  nc -zv -w 5 15.207.161.137 22
  ```
- **Fix:** Added Inbound Rule in AWS Security Group for SSH (Port 22) from `0.0.0.0/0` (or current public IP).

---

### Incident 3: Java Bytecode Mismatch (`UnsupportedClassVersionError: 65.0 vs 61.0`)
- **Symptom:** SSH connection succeeded, but agent launch failed immediately.
- **Root Cause:** Jenkins Master was compiled for Java 21 (bytecode 65.0), while the EC2 host had Java 17 (bytecode 61.0).
- **Diagnostic Command:** Looked up Java class file versions:
  - 65.0 = Java 21
  - 61.0 = Java 17
- **Fix:**
  ```bash
  sudo apt install -y openjdk-21-jre-headless
  ```

---

### Incident 4: Jenkins Container Unable to Resolve `github.com`
- **Symptom:** `git ls-remote` returned `fatal: unable to access ... Could not resolve host: github.com`.
- **Root Cause:** Docker inherited stale/unreachable nameservers from host's `/run/systemd/resolve/resolv.conf`.
- **Diagnostic Command:**
  ```bash
  docker exec jenkins-master cat /etc/resolv.conf
  ```
- **Fix:** Injected public DNS resolvers into container:
  ```bash
  docker exec -u 0 jenkins-master sh -c 'printf "nameserver 8.8.8.8\nnameserver 1.1.1.1\n" > /etc/resolv.conf'
  ```

---

### Incident 5: Declarative Pipeline Syntax Error (`Invalid option type "ansiColor"`)
- **Symptom:** Build #2 failed at startup before executing any stages.
- **Root Cause:** `ansiColor('xterm')` was declared in `options { ... }`, but the optional `ansicolor` plugin was not installed.
- **Fix:** Removed third-party option from `Jenkinsfile` and adhered to Jenkins core built-in options (`timestamps`, `timeout`, `disableConcurrentBuilds`).

---

### Incident 6: Go Backend Test Compilation Failure
- **Symptom:** Stage 2 failed with: `go: go.mod requires go >= 1.25.8 (running go 1.23.12; GOTOOLCHAIN=local)`.
- **Root Cause:** Pipeline container was pinned to `golang:1.23-alpine`, while `go.mod` required `>= 1.25.8`.
- **Fix:** Updated test container image to `golang:alpine` (Go 1.27) and mounted a persistent volume for dependency caching (`-v go-mod-cache:/go/pkg/mod`).

---

### Incident 7: Cloudflare Error 522 (Host Down) & The 1-Digit DNS Typo
- **Symptom:** Browser showed Cloudflare `Error 522: Connection timed out`.
- **Diagnostic Process (The Scientific Elimination):**
  1. Tested inside EC2: `curl -k -I https://localhost` ➔ Returned **`HTTP 200 OK`** (Proved app was healthy).
  2. Tested directly from outside internet: `curl -k -I https://15.207.161.137/` ➔ Returned **`HTTP 200 OK`** (Proved AWS firewall and ports were open).
  3. Switched Cloudflare to **DNS-only (Gray Cloud)** to unmask the raw IP.
  4. Ran verbose curl:
     ```bash
     curl -Iv https://moviehub.nostackdev.online
     # Output: Trying 15.207.61.137:443... (Timed out)
     ```
- **The Smoking Gun:**
  - Real EC2 Elastic IP: `15.207.161.137`
  - Cloudflare DNS Record: `15.207.61.137` *(Missing the digit `1`!)*
- **Fix:** Corrected the typo in Cloudflare DNS.

---

## 🧰 4. Essential DevOps CLI Cheat Sheet

```bash
# --- Network & Connectivity ---
nc -zv -w 5 <IP> <PORT>                  # Test TCP port reachability
curl -Iv https://yourdomain.com          # Inspect full TLS handshake & response
curl -4 https://ifconfig.me              # Get true external public IPv4
dig @1.1.1.1 yourdomain.com +short       # Check DNS records against Cloudflare DNS

# --- Docker Debugging ---
docker ps --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
docker inspect <container> --format '{{json .NetworkSettings.Ports}}'
docker logs --tail 50 -f <container>
docker run --rm -it -v $(pwd):/app -w /app <image> sh  # Interactive test shell

# --- System & Resource Monitoring ---
free -m                                  # Check RAM & Swap usage
df -h                                    # Check disk capacity
dmesg -T | grep -i oom                   # Verify if Linux Out-Of-Memory killer fired
```
