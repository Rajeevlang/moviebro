#!/usr/bin/env bash
# =============================================================================
# GinMovieAPI — User Database Seeder (Admin + Dummy Users) [Bash / Curl]
# Usage: ./scripts/seed_users.sh [TARGET_URL]
# =============================================================================

set -euo pipefail

TARGET_URL="${1:-http://127.0.0.1:8080}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Run Python script if python3 is available
if command -v python3 &>/dev/null; then
    python3 "$SCRIPT_DIR/seed_users.py" --target "$TARGET_URL"
    exit 0
fi

echo "================================================================="
echo "  👥 Seeding Users into GinMovieAPI -> $TARGET_URL"
echo "================================================================="

# 1. Register Admin
ADMIN_USER="admin"
ADMIN_EMAIL="admin@cinedeck.dev"
ADMIN_PASS="AdminPassword2026!"

echo "👑 1. Registering Admin ($ADMIN_EMAIL)..."
curl -s -X POST "$TARGET_URL/api/v1/users/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\": \"$ADMIN_USER\", \"email\": \"$ADMIN_EMAIL\", \"password\": \"$ADMIN_PASS\"}" > /dev/null 2>&1 || true

LOGIN_RESP=$(curl -s -X POST "$TARGET_URL/api/v1/users/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\": \"$ADMIN_EMAIL\", \"password\": \"$ADMIN_PASS\"}")

TOKEN=$(echo "$LOGIN_RESP" | grep -o '"token":"[^"]*' | cut -d'"' -f4 || true)
if [ -n "$TOKEN" ]; then
    echo "   ✅ Admin registered and JWT login verified!"
else
    echo "   ℹ️ Admin ready."
fi

# 2. Register Dummy Users
USERS=(
  "alex_morgan:alex.morgan@cinedeck.dev:UserPass2026!"
  "sarah_connor:sarah.connor@cyberdyne.io:UserPass2026!"
  "john_wick:john.wick@continental.org:UserPass2026!"
  "neo_anderson:neo.anderson@matrix.dev:UserPass2026!"
  "bruce_wayne:bruce.wayne@waynecorp.com:UserPass2026!"
)

echo "👤 2. Registering Dummy Users..."
for entry in "${USERS[@]}"; do
    IFS=':' read -r UNAME EMAIL PASS <<< "$entry"
    curl -s -X POST "$TARGET_URL/api/v1/users/register" \
        -H "Content-Type: application/json" \
        -d "{\"username\": \"$UNAME\", \"email\": \"$EMAIL\", \"password\": \"$PASS\"}" > /dev/null 2>&1 || true
    echo "   ✅ User: $UNAME ($EMAIL)"
done

echo ""
echo "================================================================="
echo "  🎉 User Seeding Complete!"
echo "  Admin Login: $ADMIN_EMAIL / $ADMIN_PASS"
echo "================================================================="
