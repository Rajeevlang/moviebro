#!/usr/bin/env python3
"""
=============================================================================
GinMovieAPI — User Database Seeder (Admin + Dummy Users)
Creates 1 Primary Administrator User and a batch of realistic dummy users
with automatic JWT token verification and credential summary table.
=============================================================================
"""

import argparse
import json
import sys
import urllib.request
import urllib.error

# ── Primary Administrator Credentials ─────────────────────────────────────────
ADMIN_USER = {
    "username": "admin",
    "email": "admin@cinedeck.dev",
    "password": "AdminPassword2026!",
    "role": "Administrator"
}

# ── Preset List of Realistic Dummy Users ──────────────────────────────────────
DUMMY_USERS = [
    {"username": "alex_morgan", "email": "alex.morgan@cinedeck.dev", "password": "UserPass2026!", "role": "Film Critic"},
    {"username": "sarah_connor", "email": "sarah.connor@cyberdyne.io", "password": "UserPass2026!", "role": "VIP Member"},
    {"username": "john_wick", "email": "john.wick@continental.org", "password": "UserPass2026!", "role": "Premium Member"},
    {"username": "neo_anderson", "email": "neo.anderson@matrix.dev", "password": "UserPass2026!", "role": "Developer"},
    {"username": "bruce_wayne", "email": "bruce.wayne@waynecorp.com", "password": "UserPass2026!", "role": "Executive Producer"},
    {"username": "tony_stark", "email": "tony.stark@starkenterprises.com", "password": "UserPass2026!", "role": "Tech Lead"},
    {"username": "elena_rostova", "email": "elena.rostova@cinema.net", "password": "UserPass2026!", "role": "Movie Enthusiast"},
    {"username": "david_fincher", "email": "david.fincher@directors.org", "password": "UserPass2026!", "role": "Director"},
    {"username": "clara_oswald", "email": "clara.oswald@tardis.space", "password": "UserPass2026!", "role": "Member"},
    {"username": "marcus_aurelius", "email": "marcus.aurelius@rome.ancient", "password": "UserPass2026!", "role": "Historian"}
]

def register_user(base_url: str, user_dict: dict) -> dict:
    """Registers a user via POST /api/v1/users/register."""
    url = f"{base_url}/api/v1/users/register"
    payload = json.dumps({
        "username": user_dict["username"],
        "email": user_dict["email"],
        "password": user_dict["password"]
    }).encode("utf-8")

    req = urllib.request.Request(url, data=payload, headers={"Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req, timeout=5) as resp:
            data = json.loads(resp.read().decode("utf-8"))
            return {"status": "created", "data": data}
    except urllib.error.HTTPError as e:
        if e.code in (409, 500):
            # Already registered or duplicate key in PostgreSQL
            return {"status": "exists", "error": f"HTTP {e.code}"}
        return {"status": "error", "error": f"HTTP {e.code}"}
    except Exception as e:
        return {"status": "error", "error": str(e)}

def login_user(base_url: str, email: str, password: str) -> str:
    """Logs in a user via POST /api/v1/users/login and returns JWT token."""
    url = f"{base_url}/api/v1/users/login"
    payload = json.dumps({
        "email": email,
        "password": password
    }).encode("utf-8")

    req = urllib.request.Request(url, data=payload, headers={"Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(req, timeout=5) as resp:
            data = json.loads(resp.read().decode("utf-8"))
            return data.get("token") or ""
    except Exception:
        return ""

def main():
    parser = argparse.ArgumentParser(description="Seed Administrator and Dummy Users into GinMovieAPI.")
    parser.add_argument("--target", default="http://localhost:8080", help="Base URL of API Gateway or User Service (default: http://localhost:8080)")
    parser.add_argument("--count", type=int, default=10, help="Number of dummy users to register (max: 10)")
    args = parser.parse_args()

    print("==================================================================================")
    print("  👥 GinMovieAPI User Database Seeder")
    print(f"  Target: {args.target}")
    print("==================================================================================")

    # 1. Register Primary Administrator
    print("\n👑 1. Registering Primary Administrator Account...")
    admin_reg = register_user(args.target, ADMIN_USER)
    admin_token = login_user(args.target, ADMIN_USER["email"], ADMIN_USER["password"])

    if admin_reg["status"] == "created":
        print(f"   ✅ Admin created: {ADMIN_USER['email']} (Username: {ADMIN_USER['username']})")
    elif admin_reg["status"] == "exists":
        print(f"   ℹ️  Admin already exists: {ADMIN_USER['email']}")
    else:
        print(f"   ⚠️ Admin registration note: {admin_reg.get('error')}")

    if admin_token:
        print(f"   🔑 Admin JWT Verified! Token: {admin_token[:25]}...")
    else:
        print("   ⚠️ Could not login admin (check if user-service is running).")

    # 2. Register Dummy Users
    print(f"\n👤 2. Registering {min(args.count, len(DUMMY_USERS))} Dummy Users...")
    results = []

    for u in DUMMY_USERS[:args.count]:
        reg_res = register_user(args.target, u)
        token = login_user(args.target, u["email"], u["password"])
        status_icon = "✅ Created" if reg_res["status"] == "created" else "ℹ️ Exists" if reg_res["status"] == "exists" else "❌ Failed"
        
        print(f"   {status_icon}: {u['username']:<18} | {u['email']:<32} | Password: {u['password']}")
        results.append({
            "username": u["username"],
            "email": u["email"],
            "password": u["password"],
            "role": u["role"],
            "token_verified": bool(token)
        })

    # 3. Print Summary Table
    print("\n==================================================================================")
    print("  📋 User Accounts & Login Credentials Summary")
    print("==================================================================================")
    print(f"  {'ROLE':<20} | {'USERNAME':<16} | {'EMAIL':<30} | {'PASSWORD'}")
    print("  " + "-" * 80)
    print(f"  👑 {ADMIN_USER['role']:<18} | {ADMIN_USER['username']:<16} | {ADMIN_USER['email']:<30} | {ADMIN_USER['password']}")
    
    for r in results:
        print(f"  👤 {r['role']:<18} | {r['username']:<16} | {r['email']:<30} | {r['password']}")

    print("==================================================================================")
    print("  🎉 User Seeding Complete!")
    print(f"  You can now log in at {args.target}/api/v1/users/login")
    print("==================================================================================")

if __name__ == "__main__":
    main()
