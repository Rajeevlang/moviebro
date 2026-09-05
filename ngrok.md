# 🚇 ngrok Field Guide: Architecture, Usage & Common DevOps Patterns

A comprehensive reference for using **ngrok** to expose local services, capture webhooks, debug distributed systems, and bypass NAT/firewalls.

---

## 🏛️ 1. What is ngrok and How Does It Work?

**ngrok** is a reverse proxy tool that establishes a secure, outbound tunnel from your local machine to ngrok's public cloud servers.

```
┌─────────────────────────────────────────────────────────────┐
│ 1. Local Machine runs: `ngrok http 8080`                    │
│                                                             │
│  [ Local Service ] ◄───► [ ngrok Client ]                   │
│  (e.g., localhost:8080)        │                            │
└────────────────────────────────┼────────────────────────────┘
                                 │ Outbound TLS Tunnel
                                 ▼
┌─────────────────────────────────────────────────────────────┐
│ 2. Public Internet                                          │
│                                                             │
│  [ GitHub / Stripe / Client ] ──► [ https://xyz.ngrok.app ] │
│  (Sends webhook or request)       (Public ngrok Edge)       │
└─────────────────────────────────────────────────────────────┘
```

### Why is ngrok so popular?
- **Zero Firewall / Router Configuration:** You don't need port forwarding, dynamic DNS, or a public static IP. The tunnel initiates **outbound** from your machine (which firewalls almost always allow).
- **Free Built-in HTTPS:** Automatic SSL/TLS termination on ngrok's edge.
- **Traffic Inspection UI:** Includes a local web console at `http://127.0.0.1:4040` where you can inspect every incoming HTTP header, payload, and replay webhooks with one click.

---

## 📦 2. Installation & Quick Setup

### On Linux (Ubuntu / Debian)
```bash
# Add ngrok repository key & list
curl -sSL https://ngrok-agent.s3.amazonaws.com/ngrok.asc \
  | sudo tee /etc/apt/trusted.gpg.d/ngrok.asc >/dev/null

echo "deb https://ngrok-agent.s3.amazonaws.com buster main" \
  | sudo tee /etc/apt/sources.list.d/ngrok.list

sudo apt update && sudo apt install ngrok
```

### On macOS
```bash
brew install ngrok/ngrok/ngrok
```

### Connect Your Account (Free Auth Token)
1. Sign up for a free account at [dashboard.ngrok.com](https://dashboard.ngrok.com).
2. Copy your Authtoken and configure ngrok:
```bash
ngrok config add-authtoken YOUR_AUTHTOKEN_HERE
```

---

## 🛠️ 3. Most Common Real-World Use Cases

### 🔹 Use Case 1: Webhook Development & Debugging
The #1 use case for developers. External services (like GitHub, Stripe, PayPal, Twilio, Razorpay) need a public HTTPS URL to notify your app of events:
```bash
# Expose your local backend running on port 5000:
ngrok http 5000
```
- ngrok gives you: `https://abcd-123.ngrok-free.app`
- Put `https://abcd-123.ngrok-free.app/webhook` in your Stripe or GitHub repository settings.
- **Inspect Traffic:** Open `http://127.0.0.1:4040` in your browser. You can inspect the raw JSON payload and click **"Replay"** to test your code without re-triggering the event on GitHub!

---

### 🔹 Use Case 2: Jenkins Local Controller Webhook Trigger (CI/CD)
When running a local Jenkins Master behind home Wi-Fi:
```bash
ngrok http 8080
```
- Webhook URL for GitHub:
  `https://your-ngrok-subdomain.ngrok-free.app/github-webhook/`
- Every `git push` reaches your local Jenkins instantly.

---

### 🔹 Use Case 3: Testing on Real Mobile Devices
When building mobile apps (React Native, Flutter, iOS/Android) that need to query your laptop's backend:
```bash
# Expose your local API server:
ngrok http 8081
```
Use the `https://xyz.ngrok-free.app` URL inside your mobile app's API base URL. Works on cellular data and any Wi-Fi network without requiring your phone and laptop to share a local network.

---

### 🔹 Use Case 4: Client & Teammate Demos
Need to show a client or teammate a feature running on your local machine before pushing to production?
```bash
# Add basic password protection so strangers can't access it:
ngrok http 3000 --basic-auth="client:secretpassword"
```

---

### 🔹 Use Case 5: SSH Tunneling Through Firewalls
Need to SSH into a computer or Raspberry Pi behind a university or corporate firewall?
```bash
# Expose the local SSH port (requires a TCP tunnel):
ngrok tcp 22
```
Output gives you: `tcp://0.tcp.ngrok.io:19425`. Connect from anywhere in the world:
```bash
ssh -p 19425 user@0.tcp.ngrok.io
```

---

## 🧰 4. Essential ngrok Command Cheat Sheet

| Command | Purpose |
| :--- | :--- |
| `ngrok http 8080` | Tunnels HTTP & HTTPS traffic to local port `8080` |
| `ngrok http 8080 --host-header="localhost:8080"` | Rewrites the `Host` header (useful for virtual hosts/Nginx) |
| `ngrok http 8080 --basic-auth="user:pass"` | Adds HTTP Basic Authentication |
| `ngrok http 8080 --domain=your-name.ngrok-free.app` | Uses your reserved free static subdomain |
| `ngrok tcp 22` | Tunnels raw TCP (for SSH, database connections) |
| `http://127.0.0.1:4040` | Local web dashboard to inspect and replay requests |

---

## ⚖️ 5. Tool Comparison: ngrok vs Cloudflare Tunnel

| Feature | ngrok | Cloudflare Tunnel (`cloudflared`) |
| :--- | :--- | :--- |
| **Best For** | Quick development, temporary webhook debugging | Permanent production services, custom domains |
| **Custom Domain** | 1 free static domain per account | Unlimited (uses your own Cloudflare domain) |
| **Session Lifetime** | Ephemeral URL unless configured with static domain | Permanent DNS record on your domain |
| **Traffic Inspection** | Built-in UI (`localhost:4040`) with one-click replay | Handled in Cloudflare Analytics & Logs |
| **Zero Trust Access** | Add-on | Included free (Google / GitHub login barrier) |
| **Cost** | Free tier (with rate limits) | 100% Free |
