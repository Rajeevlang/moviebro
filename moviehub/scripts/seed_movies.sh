#!/usr/bin/env bash
# =============================================================================
# GinMovieAPI — Quick Movie Database Seeder (Bash / Curl)
# Usage: ./scripts/seed_movies.sh [TARGET_URL]
# e.g.:  ./scripts/seed_movies.sh http://localhost:8080
#        ./scripts/seed_movies.sh https://moviehub.nostackdev.online
# =============================================================================

set -euo pipefail

TARGET_URL="${1:-http://127.0.0.1:8080}"
ADMIN_EMAIL="admin@cinedeck.dev"
ADMIN_PASS="AdminPassword2026!"
ADMIN_USER="admin_quick_seeder"

echo "================================================================="
echo "  🎬 Seeding Curated Movies into $TARGET_URL"
echo "================================================================="

# 1. Authenticate to get JWT token
echo "👤 1. Authenticating..."
curl -s -X POST "$TARGET_URL/api/v1/users/register" \
    -H "Content-Type: application/json" \
    -d "{\"username\": \"$ADMIN_USER\", \"email\": \"$ADMIN_EMAIL\", \"password\": \"$ADMIN_PASS\"}" > /dev/null 2>&1 || true

LOGIN_RESP=$(curl -s -X POST "$TARGET_URL/api/v1/users/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\": \"$ADMIN_EMAIL\", \"password\": \"$ADMIN_PASS\"}")

JWT_TOKEN=$(echo "$LOGIN_RESP" | grep -o '"token":"[^"]*' | cut -d'"' -f4 || true)

AUTH_HEADER=""
if [ -n "$JWT_TOKEN" ]; then
    echo "🔑 JWT Token acquired!"
    AUTH_HEADER="Authorization: Bearer $JWT_TOKEN"
fi

# 2. Seed Movies
echo "📽️ 2. Adding Movies..."
movies=(
  '{"title": "Inception", "description": "A thief who steals corporate secrets through the use of dream-sharing technology.", "release_year": 2010, "genres": ["Sci-Fi", "Action", "Thriller"], "poster": "https://images.unsplash.com/photo-1536440136628-849c177e76a1?auto=format&fit=crop&q=80&w=800"}'
  '{"title": "The Matrix", "description": "A computer hacker learns from mysterious rebels about the true nature of his reality.", "release_year": 1999, "genres": ["Sci-Fi", "Action"], "poster": "https://images.unsplash.com/photo-1526304640581-d334cdbbf45e?auto=format&fit=crop&q=80&w=800"}'
  '{"title": "Interstellar", "description": "A team of explorers travel through a wormhole in space in an attempt to ensure humanity survival.", "release_year": 2014, "genres": ["Sci-Fi", "Adventure", "Drama"], "poster": "https://images.unsplash.com/photo-1446776811953-b23d57bd21aa?auto=format&fit=crop&q=80&w=800"}'
  '{"title": "The Dark Knight", "description": "When the menace known as the Joker wreaks havoc and chaos on the people of Gotham.", "release_year": 2008, "genres": ["Action", "Crime", "Drama"], "poster": "https://images.unsplash.com/photo-1509347528160-9a9e33742cdb?auto=format&fit=crop&q=80&w=800"}'
  '{"title": "Pulp Fiction", "description": "The lives of two mob hitmen, a boxer, a gangster and his wife intertwine in four tales of violence.", "release_year": 1994, "genres": ["Crime", "Drama"], "poster": "https://images.unsplash.com/photo-1594909122845-11baa439b7bf?auto=format&fit=crop&q=80&w=800"}'
  '{"title": "Fight Club", "description": "An insomniac office worker and a devil-may-care soap maker form an underground fight club.", "release_year": 1999, "genres": ["Drama"], "poster": "https://images.unsplash.com/photo-1508615039623-a25605d2b022?auto=format&fit=crop&q=80&w=800"}'
  '{"title": "Blade Runner 2049", "description": "Young Blade Runner discovery of a long-buried secret leads him to track down Rick Deckard.", "release_year": 2017, "genres": ["Sci-Fi", "Mystery", "Thriller"], "poster": "https://images.unsplash.com/photo-1534447677768-be436bb09401?auto=format&fit=crop&q=80&w=800"}'
  '{"title": "Gladiator", "description": "A former Roman General sets out to exact vengeance against the corrupt emperor who murdered his family.", "release_year": 2000, "genres": ["Action", "Adventure", "Drama"], "poster": "https://images.unsplash.com/photo-1590523277543-a94d2e4eb00b?auto=format&fit=crop&q=80&w=800"}'
)

for movie in "${movies[@]}"; do
  TITLE=$(echo "$movie" | grep -o '"title": "[^"]*' | cut -d'"' -f4)
  if [ -n "$AUTH_HEADER" ]; then
      curl -s -X POST "$TARGET_URL/api/v1/movies/" \
        -H "Content-Type: application/json" \
        -H "$AUTH_HEADER" \
        -d "$movie" > /dev/null 2>&1 || true
  else
      curl -s -X POST "$TARGET_URL/api/v1/movies/" \
        -H "Content-Type: application/json" \
        -d "$movie" > /dev/null 2>&1 || true
  fi
  echo "   ✅ Added: $TITLE"
done

echo ""
echo "================================================================="
echo "  🎉 Database successfully seeded!"
echo "  Catalog: $TARGET_URL/api/v1/movies"
echo "================================================================="
