#!/usr/bin/env python3
"""
=============================================================================
GinMovieAPI — Open-Source Movie Database Bulk Seeder
Fetches movie data from free open-source APIs (TVMaze API / Curated Catalog)
and ingests them into the GinMovieAPI microservices stack.
=============================================================================
"""

import argparse
import json
import re
import sys
import urllib.request
import urllib.error

# ── Curated Catalog of Blockbuster & Award-Winning Movies ─────────────────────
CURATED_MOVIES = [
    {
        "title": "Oppenheimer",
        "description": "The story of American scientist J. Robert Oppenheimer and his role in the development of the atomic bomb during World War II.",
        "release_year": 2023,
        "genres": ["Biography", "Drama", "History"],
        "poster": "https://images.unsplash.com/photo-1485846234645-a62644f84728?auto=format&fit=crop&q=80&w=800"
    },
    {
        "title": "Interstellar",
        "description": "When Earth becomes uninhabitable, a farmer and ex-NASA pilot is tasked to pilot a spacecraft along with a team of researchers to find a new planet for humans.",
        "release_year": 2014,
        "genres": ["Adventure", "Drama", "Sci-Fi"],
        "poster": "https://images.unsplash.com/photo-1446776811953-b23d57bd21aa?auto=format&fit=crop&q=80&w=800"
    },
    {
        "title": "Inception",
        "description": "A thief who steals corporate secrets through the use of dream-sharing technology is given the inverse task of planting an idea into the mind of a CEO.",
        "release_year": 2010,
        "genres": ["Action", "Sci-Fi", "Thriller"],
        "poster": "https://images.unsplash.com/photo-1536440136628-849c177e76a1?auto=format&fit=crop&q=80&w=800"
    },
    {
        "title": "The Dark Knight",
        "description": "When the menace known as the Joker wreaks havoc and chaos on Gotham City, Batman must accept one of the greatest psychological tests to fight injustice.",
        "release_year": 2008,
        "genres": ["Action", "Crime", "Drama"],
        "poster": "https://images.unsplash.com/photo-1509347528160-9a9e33742cdb?auto=format&fit=crop&q=80&w=800"
    },
    {
        "title": "Dune: Part Two",
        "description": "Paul Atreides unites with Chani and the Fremen while seeking revenge against the conspirators who destroyed his family.",
        "release_year": 2024,
        "genres": ["Action", "Adventure", "Sci-Fi"],
        "poster": "https://images.unsplash.com/photo-1506744038136-46273834b3fb?auto=format&fit=crop&q=80&w=800"
    },
    {
        "title": "Blade Runner 2049",
        "description": "Young Blade Runner K's discovery of a long-buried secret leads him to track down former Blade Runner Rick Deckard, who's been missing for thirty years.",
        "release_year": 2017,
        "genres": ["Action", "Drama", "Sci-Fi"],
        "poster": "https://images.unsplash.com/photo-1518709268805-4e9042af9f23?auto=format&fit=crop&q=80&w=800"
    },
    {
        "title": "Spider-Man: Across the Spider-Verse",
        "description": "Miles Morales catapults across the Multiverse, where he encounters a team of Spider-People charged with protecting its very existence.",
        "release_year": 2023,
        "genres": ["Animation", "Action", "Adventure"],
        "poster": "https://images.unsplash.com/photo-1607604276583-eef5d076aa5f?auto=format&fit=crop&q=80&w=800"
    },
    {
        "title": "Parasite",
        "description": "Greed and class discrimination threaten the newly formed symbiotic relationship between the wealthy Park family and the destitute Kim clan.",
        "release_year": 2019,
        "genres": ["Drama", "Thriller"],
        "poster": "https://images.unsplash.com/photo-1517604931442-7e0c8ed2963c?auto=format&fit=crop&q=80&w=800"
    },
    {
        "title": "The Matrix",
        "description": "When a beautiful stranger leads computer hacker Neo to a forbidding underworld, he discovers the shocking truth--the life he knows is the elaborate deception of an evil cyber-intelligence.",
        "release_year": 1999,
        "genres": ["Action", "Sci-Fi"],
        "poster": "https://images.unsplash.com/photo-1526374965328-7f61d4dc18c5?auto=format&fit=crop&q=80&w=800"
    },
    {
        "title": "Spirited Away",
        "description": "During her family's move to the suburbs, a sullen 10-year-old girl wanders into a world ruled by gods, witches and spirits, where humans are changed into beasts.",
        "release_year": 2001,
        "genres": ["Animation", "Adventure", "Family"],
        "poster": "https://images.unsplash.com/photo-1534447677768-be436bb09401?auto=format&fit=crop&q=80&w=800"
    },
    {
        "title": "Gladiator",
        "description": "A former Roman General sets out to exact vengeance against the corrupt emperor who murdered his family and sent him into slavery.",
        "release_year": 2000,
        "genres": ["Action", "Adventure", "Drama"],
        "poster": "https://images.unsplash.com/photo-1568605117036-5fe5e7bab0b7?auto=format&fit=crop&q=80&w=800"
    },
    {
        "title": "Whiplash",
        "description": "A promising young drummer enrolls at a cut-throat music conservatory where his dreams of greatness are mentored by an instructor who will stop at nothing to realize a student's potential.",
        "release_year": 2014,
        "genres": ["Drama", "Music"],
        "poster": "https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?auto=format&fit=crop&q=80&w=800"
    },
    {
        "title": "Avengers: Endgame",
        "description": "After the devastating events of Infinity War, the universe is in ruins. With the help of remaining allies, the Avengers assemble once more to reverse Thanos' actions.",
        "release_year": 2019,
        "genres": ["Action", "Adventure", "Drama"],
        "poster": "https://images.unsplash.com/photo-1534809027769-b00d750a6bac?auto=format&fit=crop&q=80&w=800"
    },
    {
        "title": "Pulp Fiction",
        "description": "The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.",
        "release_year": 1994,
        "genres": ["Crime", "Drama"],
        "poster": "https://images.unsplash.com/photo-1594909122845-11baa439b7bf?auto=format&fit=crop&q=80&w=800"
    },
    {
        "title": "Interstellar 2: Beyond Horizons",
        "description": "Humanity embarks on a journey through deep spacetime to establish new civilization beyond known gravitational anomalies.",
        "release_year": 2026,
        "genres": ["Sci-Fi", "Adventure", "Drama"],
        "poster": "https://images.unsplash.com/photo-1451187580459-43490279c0fa?auto=format&fit=crop&q=80&w=800"
    }
]

def clean_html(raw_html: str) -> str:
    """Removes HTML tags from API description strings."""
    if not raw_html:
        return "No description available."
    clean_r = re.compile('<.*?>')
    cleaned = re.sub(clean_r, '', raw_html).strip()
    return (cleaned[:490] + '...') if len(cleaned) > 490 else cleaned

def fetch_from_tvmaze(query: str = None, count: int = 20) -> list:
    """Fetches movies/shows from the free open TVMaze API."""
    if query:
        url = f"https://api.tvmaze.com/search/shows?q={urllib.parse.quote(query)}"
    else:
        url = "https://api.tvmaze.com/shows?page=0"

    print(f"🌐 Fetching movie catalog from TVMaze Open API ({url})...")
    
    req = urllib.request.Request(url, headers={'User-Agent': 'Mozilla/5.0 (GinMovieAPI Seeder)'})
    try:
        with urllib.request.urlopen(req, timeout=10) as response:
            data = json.loads(response.read().decode('utf-8'))
    except Exception as e:
        print(f"⚠️ TVMaze API fetch failed: {e}. Falling back to curated catalog.")
        return []

    movies = []
    items = data if not query else [item['show'] for item in data if 'show' in item]

    for item in items[:count]:
        title = item.get('name')
        summary = clean_html(item.get('summary', ''))
        premiered = item.get('premiered', '')
        
        # Extract year
        year = 2020
        if premiered and len(premiered) >= 4 and premiered[:4].isdigit():
            year = int(premiered[:4])

        genres = item.get('genres') or ["Drama"]
        
        # Poster image
        image_obj = item.get('image') or {}
        poster = image_obj.get('original') or image_obj.get('medium') or "https://images.unsplash.com/photo-1485846234645-a62644f84728?w=800"

        if title and summary:
            movies.append({
                "title": title,
                "description": summary,
                "release_year": year,
                "genres": genres,
                "poster": poster
            })

    return movies

def authenticate(base_url: str) -> str:
    """Registers and logs into user-service to acquire a JWT token."""
    email = "admin@cinedeck.dev"
    password = "AdminPassword2026!"
    username = "admin_seeder"

    # 1. Register user (ignore error if user already exists)
    reg_url = f"{base_url}/api/v1/users/register"
    reg_payload = json.dumps({"username": username, "email": email, "password": password}).encode('utf-8')
    reg_req = urllib.request.Request(reg_url, data=reg_payload, headers={'Content-Type': 'application/json'})
    try:
        urllib.request.urlopen(reg_req, timeout=5)
    except Exception:
        pass # User may already exist

    # 2. Login to get JWT
    login_url = f"{base_url}/api/v1/users/login"
    login_payload = json.dumps({"email": email, "password": password}).encode('utf-8')
    login_req = urllib.request.Request(login_url, data=login_payload, headers={'Content-Type': 'application/json'})
    
    try:
        with urllib.request.urlopen(login_req, timeout=5) as resp:
            resp_data = json.loads(resp.read().decode('utf-8'))
            token = resp_data.get('token') or resp_data.get('data', {}).get('token')
            if token:
                print("🔑 Acquired JWT Authentication Token from user-service!")
                return token
    except Exception as e:
        print(f"ℹ️ Auth login failed: {e}. Posting directly to movie service...")

    return ""

def post_movie(base_url: str, movie: dict, auth_token: str) -> bool:
    """Posts a single movie object to the GinMovieAPI."""
    url = f"{base_url}/api/v1/movies/"
    payload = json.dumps(movie).encode('utf-8')
    
    headers = {'Content-Type': 'application/json'}
    if auth_token:
        headers['Authorization'] = f"Bearer {auth_token}"

    req = urllib.request.Request(url, data=payload, headers=headers)
    try:
        with urllib.request.urlopen(req, timeout=5) as resp:
            if resp.status in (200, 201):
                return True
    except urllib.error.HTTPError as e:
        if e.code == 409: # Already exists
            print(f"   ℹ️ Already exists: {movie['title']}")
            return True
        print(f"   ⚠️ HTTP {e.code} for {movie['title']}: {e.read().decode('utf-8', errors='ignore')}")
    except Exception as e:
        print(f"   ❌ Network error for {movie['title']}: {e}")
    return False

def main():
    parser = argparse.ArgumentParser(description="Seed movies into GinMovieAPI from open-source databases.")
    parser.add_argument("--target", default="http://localhost:8080", help="Base URL of API Gateway or Movie Service")
    parser.add_argument("--source", choices=["all", "curated", "tvmaze", "search"], default="all", help="Data source to use")
    parser.add_argument("--query", default="", help="Search query for TVMaze open API (e.g. 'Star Wars', 'Batman')")
    parser.add_argument("--count", type=int, default=25, help="Number of movies to seed")
    
    args = parser.parse_args()

    print("=================================================================")
    print("  🎬 GinMovieAPI Open-Source Database Seeder")
    print(f"  Target: {args.target} | Source: {args.source}")
    print("=================================================================")

    # 1. Acquire JWT Token
    auth_token = authenticate(args.target)

    # 2. Gather movies
    movies_to_seed = []
    
    if args.source in ("curated", "all"):
        movies_to_seed.extend(CURATED_MOVIES)

    if args.source in ("tvmaze", "all") or args.query:
        tvmaze_movies = fetch_from_tvmaze(query=args.query if args.query else None, count=args.count)
        movies_to_seed.extend(tvmaze_movies)

    # De-duplicate by title
    seen = set()
    unique_movies = []
    for m in movies_to_seed:
        if m['title'].lower() not in seen:
            seen.add(m['title'].lower())
            unique_movies.append(m)

    print(f"📽️ Seeding {len(unique_movies)} movies into database...")

    success_count = 0
    for movie in unique_movies[:args.count]:
        print(f"   ➕ Adding: {movie['title']} ({movie['release_year']}) [{', '.join(movie['genres'])}]")
        if post_movie(args.target, movie, auth_token):
            success_count += 1

    print("")
    print("=================================================================")
    print(f"  🎉 Seeding Finished! {success_count}/{len(unique_movies)} movies processed.")
    print(f"  Check your catalog at: {args.target}/api/v1/movies")
    print("=================================================================")

if __name__ == "__main__":
    main()
