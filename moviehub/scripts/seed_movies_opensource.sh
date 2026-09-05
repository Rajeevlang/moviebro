#!/usr/bin/env bash
# =============================================================================
# GinMovieAPI — Open-Source Movie Database Bulk Seeder (Bash / Curl)
# Usage: ./scripts/seed_movies_opensource.sh [TARGET_URL]
# =============================================================================

set -euo pipefail

TARGET_URL="${1:-http://127.0.0.1:8080}"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "================================================================="
echo "  🎬 Seeding GinMovieAPI via Open-Source APIs -> $TARGET_URL"
echo "================================================================="

# 1. Run Python Seeder if python3 is available
if command -v python3 &>/dev/null; then
    python3 "$SCRIPT_DIR/seed_movies_opensource.py" --target "$TARGET_URL" --source all --count 30
    exit 0
fi

# 2. Fallback pure-Bash Seeder
ADMIN_EMAIL="admin@cinedeck.dev"
ADMIN_PASS="AdminSecurePass2026!"
ADMIN_USER="admin_bash_seeder"

echo "👤 1. Authenticating with user-service..."
curl -s -X POST "$TARGET_URL/api/v1/users/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\": \"$ADMIN_USER\", \"email\": \"$ADMIN_EMAIL\", \"password\": \"$ADMIN_PASS\"}" > /dev/null 2>&1 || true

LOGIN_RESP=$(curl -s -X POST "$TARGET_URL/api/v1/users/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\": \"$ADMIN_EMAIL\", \"password\": \"$ADMIN_PASS\"}")

JWT_TOKEN=$(echo "$LOGIN_RESP" | grep -o '"token":"[^"]*' | cut -d'"' -f4 || true)

AUTH_HEADER=""
if [ -n "$JWT_TOKEN" ]; then
    echo "✅ JWT Token obtained!"
    AUTH_HEADER="Authorization: Bearer $JWT_TOKEN"
fi

echo "🌐 2. Fetching open-source show catalog from TVMaze API..."
SHOWS_JSON=$(curl -s "https://api.tvmaze.com/shows?page=0" | head -c 50000)

echo "📽️ 3. Adding Movies..."
# Adding curated blockbusters
MOVIES=(
  '{"title":"Oppenheimer","description":"The story of American scientist J. Robert Oppenheimer and his role in the Manhattan Project.","release_year":2023,"genres":["Biography","Drama"],"poster":"https://images.unsplash.com/photo-1485846234645-a62644f84728?w=800"}'
  '{"title":"Interstellar","description":"A team of explorers travel through a wormhole in space to ensure humanity survival.","release_year":2014,"genres":["Adventure","Sci-Fi"],"poster":"https://images.unsplash.com/photo-1446776811953-b23d57bd21aa?w=800"}'
  '{"title":"Inception","description":"A thief steals corporate secrets through dream-sharing technology.","release_year":2010,"genres":["Action","Sci-Fi"],"poster":"https://images.unsplash.com/photo-1536440136628-849c177e76a1?w=800"}'
  '{"title":"The Dark Knight","description":"Batman battles the Joker to save Gotham City.","release_year":2008,"genres":["Action","Crime"],"poster":"https://images.unsplash.com/photo-1509347528160-9a9e33742cdb?w=800"}'
  '{"title":"Dune: Part Two","description":"Paul Atreides unites with Chani and the Fremen to avenge his family.","release_year":2024,"genres":["Sci-Fi","Action"],"poster":"https://images.unsplash.com/photo-1506744038136-46273834b3fb?w=800"}'
)

COUNT=0
for m in "${MOVIES[@]}"; do
    TITLE=$(echo "$m" | grep -o '"title":"[^"]*' | cut -d'"' -f4)
    if [ -n "$AUTH_HEADER" ]; then
        HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$TARGET_URL/api/v1/movies/" \
            -H "Content-Type: application/json" \
            -H "$AUTH_HEADER" \
            -d "$m")
    else
        HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$TARGET_URL/api/v1/movies/" \
            -H "Content-Type: application/json" \
            -d "$m")
    fi
    if [ "$HTTP_CODE" -eq 200 ] || [ "$HTTP_CODE" -eq 201 ]; then
        echo "   ✅ Added: $TITLE"
        COUNT=$((COUNT + 1))
    else
        echo "   ℹ️ ($HTTP_CODE): $TITLE"
    fi
done

echo ""
echo "================================================================="
echo "  🎉 Seeding Finished! $COUNT movies processed."
echo "  Explore catalog at: $TARGET_URL/api/v1/movies"
echo "================================================================="
