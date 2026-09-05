# 🛡️ Cloudflare WARP Guide: Installation, Usage & Bypass

A quick reference for using **Cloudflare WARP (1.1.1.1 with WARP)** on Linux to bypass network filters, fix DNS issues, and secure outbound traffic.

---

## 🏛️ 1. What is Cloudflare WARP?

**Cloudflare WARP** uses the modern **WireGuard** protocol to route and encrypt all your computer's internet traffic through Cloudflare’s global network (over 300+ cities).

```
[ Your Computer ]
       │
       ▼ (Encrypted WireGuard Tunnel over UDP / Anycast)
[ Campus / School Wi-Fi Firewall ] (Cannot inspect or block destinations)
       │
       ▼
[ Cloudflare Global Edge ] 
       │
       ▼
[ GitHub / Docker Hub / Internet Services ]
```

### Why it solves campus & ISP blocks:
- The local network only sees an encrypted WireGuard connection to Cloudflare.
- Domain inspections (SNI sniffing) and DNS poisoning used by schools/campuses fail.
- All blocked developer sites (`github.com`, `docker.io`, package registries) work immediately.

---

## 📦 2. Installation on Linux (Ubuntu / Debian)

```bash
# 1. Add Cloudflare WARP GPG Key
curl -fsSL https://pkg.cloudflareclient.com/pubkey.gpg | sudo gpg --yes --dearmor --output /usr/share/keyrings/cloudflare-warp-archive-keyring.gpg

# 2. Add the Cloudflare repository to APT sources
echo "deb [arch=amd64 signed-by=/usr/share/keyrings/cloudflare-warp-archive-keyring.gpg] https://pkg.cloudflareclient.com/ $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/cloudflare-client.list

# 3. Update package list & install
sudo apt update && sudo apt install -y cloudflare-warp
```

---

## ⚡ 3. Initial Setup & Registration (One-Time)

Before connecting for the first time, register your device client:

```bash
warp-cli registration new
```

To view your registration details and license:
```bash
warp-cli registration show
```

---

## 🚀 4. Everyday Usage & Commands

### Connect to WARP
```bash
warp-cli connect
```

### Check Connection Status
```bash
warp-cli status
```
*Expected output when working:*
```text
Status update: Connected
Network: healthy
```

### Disconnect from WARP
```bash
warp-cli disconnect
```

### Toggle Modes
- **WARP mode (Default - Full tunnel for DNS & Traffic):**
  ```bash
  warp-cli mode warp
  ```
- **DoH mode (DNS-Only over HTTPS, keeps normal IP):**
  ```bash
  warp-cli mode doh
  ```

---

## 🧪 5. Verifying That WARP is Working

Run these quick checks:

```bash
# 1. Verify Cloudflare sees you on the WARP network
curl -s https://www.cloudflare.com/cdn-cgi/trace | grep warp
# Output should be: warp=on

# 2. Test GitHub connectivity
curl -I https://github.com | head -n 1
# Output should be: HTTP/2 200
```

---

## 💡 Troubleshooting WARP

| Issue | Cause | Solution |
| :--- | :--- | :--- |
| `Status: Connecting (stuck)` | Campus blocks UDP/WireGuard on port 2408 | Switch to fallback DNS mode: `warp-cli mode doh` |
| `warp-cli: command not found` | Package not installed | Run the installation script above |
| Local services (`localhost:8080`) slow down | WARP routing loop | Run: `warp-cli tunnel host add localhost` |
