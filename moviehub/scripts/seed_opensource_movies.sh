#!/usr/bin/env bash
# =============================================================================
# GinMovieAPI — Open-Source Movie Database Seeder (Shell Script)
# Fetches movie data from free Open-Source APIs (TVMaze API + Blockbuster Catalog)
# and seeds them into the GinMovieAPI microservices database.
#
# Usage:
#   bash seed_opensource_movies.sh [TARGET_URL]
#   e.g. bash seed_opensource_movies.sh http://127.0.0.1
# =============================================================================

set -euo pipefail

TARGET_URL="${1:-http://127.0.0.1}"
ADMIN_EMAIL="admin@cinedeck.dev"
ADMIN_PASS="AdminPassword2026!"
ADMIN_USER="admin_cinedeck"

echo "================================================================="
echo "  🎬 GinMovieAPI Open-Source Movie Seeder"
echo "  Target: $TARGET_URL"
echo "================================================================="

# ── 1. Install jq if missing ──────────────────────────────────────────────────
if ! command -v jq &>/dev/null; then
    echo "📦 Installing jq for JSON processing..."
    sudo apt-get update -qq && sudo apt-get install -y -qq jq || true
fi

# ── 2. Authenticate & Obtain Admin JWT Token ──────────────────────────────────
echo "👤 1. Authenticating Admin User ($ADMIN_EMAIL)..."

# Ensure user exists (ignore error if already registered)
curl -s -X POST "$TARGET_URL/api/v1/users/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\": \"$ADMIN_USER\", \"email\": \"$ADMIN_EMAIL\", \"password\": \"$ADMIN_PASS\"}" > /dev/null 2>&1 || true

# Login to get JWT
LOGIN_RESP=$(curl -s -X POST "$TARGET_URL/api/v1/users/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\": \"$ADMIN_EMAIL\", \"password\": \"$ADMIN_PASS\"}")

if command -v jq &>/dev/null; then
    JWT_TOKEN=$(echo "$LOGIN_RESP" | jq -r '.token // empty')
else
    JWT_TOKEN=$(echo "$LOGIN_RESP" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')
fi

if [ -z "$JWT_TOKEN" ]; then
    echo "⚠️ Warning: Could not obtain JWT token. Checking if services are online..."
    echo "Response was: $LOGIN_RESP"
    AUTH_HEADER=""
else
    echo "🔑 JWT Token acquired successfully!"
    AUTH_HEADER="Authorization: Bearer $JWT_TOKEN"
fi

# ── 3. Seed Curated Top Blockbusters ──────────────────────────────────────────
echo ""
echo "📽️ 2. Ingesting Curated Award-Winning Movies..."

CURATED_MOVIES=(
  '{"title":"Oppenheimer","description":"The story of American scientist J. Robert Oppenheimer and his role in the Manhattan Project.","release_year":2023,"genres":["Biography","Drama","History"],"poster":"https://images.unsplash.com/photo-1485846234645-a62644f84728?w=800"}'
  '{"title":"Interstellar","description":"A team of explorers travel through a wormhole in space in an attempt to ensure humanity survival.","release_year":2014,"genres":["Adventure","Drama","Sci-Fi"],"poster":"https://images.unsplash.com/photo-1446776811953-b23d57bd21aa?w=800"}'
  '{"title":"Inception","description":"A thief who steals corporate secrets through dream-sharing technology is given the inverse task of planting an idea.","release_year":2010,"genres":["Action","Sci-Fi","Thriller"],"poster":"https://images.unsplash.com/photo-1536440136628-849c177e76a1?w=800"}'
  '{"title":"The Dark Knight","description":"When the menace known as the Joker wreaks havoc on Gotham, Batman must fight injustice.","release_year":2008,"genres":["Action","Crime","Drama"],"poster":"https://images.unsplash.com/photo-1509347528160-9a9e33742cdb?w=800"}'
  '{"title":"Dune: Part Two","description":"Paul Atreides unites with Chani and the Fremen while seeking revenge against conspirators.","release_year":2024,"genres":["Action","Adventure","Sci-Fi"],"poster":"https://images.unsplash.com/photo-1506744038136-46273834b3fb?w=800"}'
  '{"title":"Blade Runner 2049","description":"Young Blade Runner K discovery of a long-buried secret leads him to track down Rick Deckard.","release_year":2017,"genres":["Action","Drama","Sci-Fi"],"poster":"https://images.unsplash.com/photo-1518709268805-4e9042af9f23?w=800"}'
  '{"title":"The Matrix","description":"A computer hacker learns from mysterious rebels about the true nature of his reality.","release_year":1999,"genres":["Action","Sci-Fi"],"poster":"https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?w=800"}'
  '{"title":"Pulp Fiction","description":"The lives of two mob hitmen, a boxer, a gangster and his wife intertwine in four tales of violence.","release_year":1994,"genres":["Crime","Drama"],"poster":"https://images.unsplash.com/photo-1594909122845-11baa439b7bf?w=800"}'
  '{"title":"Gladiator","description":"A former Roman General sets out to exact vengeance against the corrupt emperor who murdered his family.","release_year":2000,"genres":["Action","Adventure","Drama"],"poster":"https://images.unsplash.com/photo-1568605117036-5fe5e7bab0b7?w=800"}'
  '{"title":"Spider-Man: Across the Spider-Verse","description":"Miles Morales catapults across the Multiverse to protect its very existence with fellow Spider-Heroes.","release_year":2023,"genres":["Animation","Action","Adventure"],"poster":"https://images.unsplash.com/photo-1607604276583-eef5d076aa5f?w=800"}'
  '{"title":"Parasite","description":"Greed and class discrimination threaten the newly formed symbiotic relationship between two families.","release_year":2019,"genres":["Drama","Thriller"],"poster":"https://images.unsplash.com/photo-1517604931442-7e0c8ed2963c?w=800"}'
  '{"title":"Spirited Away","description":"A 10-year-old girl wanders into a world ruled by gods, witches, and spirits where humans become beasts.","release_year":2001,"genres":["Animation","Adventure","Family"],"poster":"https://images.unsplash.com/photo-1534447677768-be436bb09401?w=800"}'
)

SUCCESS_COUNT=0

for movie in "${CURATED_MOVIES[@]}"; do
    TITLE=$(echo "$movie" | grep -o '"title":"[^"]*' | cut -d'"' -f4)
    
    if [ -n "$AUTH_HEADER" ]; then
        HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$TARGET_URL/api/v1/movies/" \
            -H "Content-Type: application/json" \
            -H "$AUTH_HEADER" \
            -d "$movie")
    else
        HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$TARGET_URL/api/v1/movies/" \
            -H "Content-Type: application/json" \
            -d "$movie")
    fi

    if [ "$HTTP_CODE" -eq 200 ] || [ "$HTTP_CODE" -eq 201 ]; then
        echo "   ✅ Added: $TITLE"
        SUCCESS_COUNT=$((SUCCESS_COUNT + 1))
    elif [ "$HTTP_CODE" -eq 409 ]; then
        echo "   ℹ️ Already exists: $TITLE"
    else
        echo "   ⚠️ HTTP $HTTP_CODE: $TITLE"
    fi
done

# ── 4. Fetch Dynamic Movies from Open TVMaze API ──────────────────────────────
echo ""
echo "🌐 3. Fetching live movie catalog from Open-Source TVMaze API..."

if command -v python3 &>/dev/null; then
    # Use python helper for parsing TVMaze JSON & sanitizing HTML tags
    python3 - << 'PYEOF' "$TARGET_URL" "$JWT_TOKEN"
import sys, json, re, urllib.request, urllib.error

target_url = sys.argv[1]
jwt_token = sys.argv[2]

try:
    req = urllib.request.Request("https://api.tvmaze.com/shows?page=0", headers={"User-Agent": "GinMovieAPI-Seeder/1.0"})
    with urllib.request.urlopen(req, timeout=10) as resp:
        data = json.loads(resp.read().decode('utf-8'))
except Exception as e:
    print(f"   ⚠️ Could not reach TVMaze API: {e}")
    sys.exit(0)

added = 0
for show in data[:20]:
    title = show.get("name", "")
    summary_raw = show.get("summary", "") or "No description provided."
    summary = re.sub('<.*?>', '', summary_raw).strip()
    if len(summary) > 450:
        summary = summary[:450] + "..."
    
    prem = show.get("premiered", "")
    year = int(prem[:4]) if prem and len(prem) >= 4 and prem[:4].isdigit() else 2021
    genres = show.get("genres") or ["Drama"]
    
    img = show.get("image") or {}
    poster = img.get("original") or img.get("medium") or "https://images.unsplash.com/photo-1485846234645-a62644f84728?w=800"

    payload = json.dumps({
        "title": title,
        "description": summary,
        "release_year": year,
        "genres": genres,
        "poster": poster
    }).encode('utf-8')

    headers = {"Content-Type": "application/json"}
    if jwt_token:
        headers["Authorization"] = f"Bearer {jwt_token}"

    req = urllib.request.Request(f"{target_url}/api/v1/movies/", data=payload, headers=headers)
    try:
        with urllib.request.urlopen(req, timeout=5) as r:
            if r.status in (200, 201):
                print(f"   ✅ Added from Open API: {title} ({year})")
                added += 1
    except urllib.error.HTTPError as e:
        if e.code == 409:
            print(f"   ℹ️ Already exists: {title}")
        else:
            print(f"   ⚠️ HTTP {e.code}: {title}")
    except Exception:
        pass

print(f"\n   🎉 Added {added} additional movies from Open-Source API!")
PYEOF
fi

# ── 5. Completion Summary ─────────────────────────────────────────────────────
echo ""
echo "================================================================="
echo "  🎉 Seeding Complete!"
echo "  Check your catalog at: $TARGET_URL/api/v1/movies"
echo "  Browse UI at:          $TARGET_URL/"
echo "================================================================="
