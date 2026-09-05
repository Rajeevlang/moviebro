#!/usr/bin/env python3
"""
=============================================================================
GinMovieAPI — Massive Movie Dataset Seeder from GitHub Open Dataset
Downloads public JSON dataset and ingests into GinMovieAPI catalog.
=============================================================================
"""

import argparse
import json
import sys
import urllib.error
import urllib.request

JSON_DATA_URL = "https://raw.githubusercontent.com/erik-sytnyk/movies-list/master/db.json"


def fetch_movies():
    print("📥 Downloading massive movie dataset from GitHub...")
    req = urllib.request.Request(JSON_DATA_URL, headers={"User-Agent": "Mozilla/5.0 (GinMovieAPI Seeder)"})
    try:
        with urllib.request.urlopen(req, timeout=15) as response:
            data = json.loads(response.read().decode("utf-8"))
            return data.get("movies", [])
    except Exception as e:
        print(f"❌ Failed to download dataset: {e}")
        return []


def authenticate(base_url: str) -> str:
    """Registers / logs in admin user to obtain a JWT token."""
    email = "admin@cinedeck.dev"
    password = "AdminPassword2026!"
    username = "admin_massive_seeder"

    # Register (ignore if already exists)
    reg_url = f"{base_url}/api/v1/users/register"
    reg_payload = json.dumps({"username": username, "email": email, "password": password}).encode("utf-8")
    reg_req = urllib.request.Request(reg_url, data=reg_payload, headers={"Content-Type": "application/json"})
    try:
        urllib.request.urlopen(reg_req, timeout=5)
    except Exception:
        pass

    # Login
    login_url = f"{base_url}/api/v1/users/login"
    login_payload = json.dumps({"email": email, "password": password}).encode("utf-8")
    login_req = urllib.request.Request(login_url, data=login_payload, headers={"Content-Type": "application/json"})
    try:
        with urllib.request.urlopen(login_req, timeout=5) as resp:
            data = json.loads(resp.read().decode("utf-8"))
            return data.get("token", "")
    except Exception as e:
        print(f"ℹ️ Auth warning: {e}")
        return ""


def seed_database(base_url: str, count: int):
    api_url = f"{base_url}/api/v1/movies/"
    token = authenticate(base_url)

    movies = fetch_movies()
    if not movies:
        print("❌ No movies found to seed.")
        return

    print(f"🎬 Found {len(movies)} movies. Seeding up to {count} movies into {api_url}...")

    success_count = 0
    fail_count = 0

    headers = {"Content-Type": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"

    for movie in movies[:count]:
        try:
            poster = movie.get("posterUrl") or "https://images.unsplash.com/photo-1485846234645-a62644f84728?w=600"
            if not poster.startswith("http://") and not poster.startswith("https://"):
                poster = "https://images.unsplash.com/photo-1485846234645-a62644f84728?w=600"

            genres = movie.get("genres") or ["Drama"]
            if not isinstance(genres, list) or len(genres) == 0:
                genres = ["Drama"]

            payload = {
                "title": movie.get("title", "Unknown"),
                "description": (movie.get("plot") or "No description available.")[:490],
                "release_year": int(movie.get("year", "2000")),
                "genres": genres,
                "poster": poster,
            }

            req = urllib.request.Request(
                api_url,
                data=json.dumps(payload).encode("utf-8"),
                headers=headers,
                method="POST",
            )

            with urllib.request.urlopen(req, timeout=5) as response:
                if response.status in (200, 201):
                    success_count += 1
                    print(f"   ✅ Added: {payload['title']} ({payload['release_year']})")
                else:
                    fail_count += 1

        except urllib.error.HTTPError as e:
            if e.code == 409:
                success_count += 1
                print(f"   ℹ️ Already exists: {movie.get('title')}")
            else:
                fail_count += 1
        except Exception:
            fail_count += 1

    print(f"\n🎉 Seeding Complete! Successfully added: {success_count} | Skipped/Failed: {fail_count}")


def main():
    parser = argparse.ArgumentParser(description="Seed massive movie catalog from GitHub open dataset.")
    parser.add_argument("--target", default="http://localhost:8080", help="Base URL of API Gateway (default: http://localhost:8080)")
    parser.add_argument("--count", type=int, default=50, help="Number of movies to seed (default: 50)")
    args = parser.parse_args()

    seed_database(args.target.rstrip("/"), args.count)


if __name__ == "__main__":
    main()
