#!/usr/bin/env python3
"""
MovieHub Microservices Test & Seed Script
-----------------------------------------
This script verifies that the MovieHub microservices are operational on your EC2 instance
(or locally) and seeds initial dummy data (admin/users, movies, playlists, and reviews).

Features:
- Zero external dependencies (uses standard Python 3 urllib/json libraries).
- Health check validation for API Gateway, User Service, and Movie Service.
- Dummy user registration & authentication test.
- Seeding high-quality dummy movies with genres and posters.
- Seeding playlists and reviews.
- Formatted status reports.

Usage:
  python3 test_and_seed.py
  python3 test_and_seed.py --base-url http://localhost:8080
  python3 test_and_seed.py --base-url http://<EC2_PUBLIC_IP>:8080
  python3 test_and_seed.py --base-url https://moviehub.nostackdev.online
"""

import argparse
import json
import sys
import time
import urllib.error
import urllib.parse
import urllib.request
from typing import Any, Dict, List, Optional, Tuple

# Terminal colors for clean CLI feedback
GREEN = "\033[92m"
RED = "\033[91m"
YELLOW = "\033[93m"
CYAN = "\033[96m"
BOLD = "\033[1m"
RESET = "\033[0m"


def log_step(title: str):
    print(f"\n{BOLD}{CYAN}=== {title} ==={RESET}")


def log_success(msg: str):
    print(f"  {GREEN}✓{RESET} {msg}")


def log_warn(msg: str):
    print(f"  {YELLOW}⚠{RESET} {msg}")


def log_fail(msg: str):
    print(f"  {RED}✗{RESET} {msg}")


def make_request(
    url: str,
    method: str = "GET",
    data: Optional[Dict[str, Any]] = None,
    token: Optional[str] = None,
    timeout: int = 10,
) -> Tuple[int, Any]:
    """Execute an HTTP request and parse JSON response using urllib."""
    headers = {"Content-Type": "application/json", "Accept": "application/json"}
    if token:
        headers["Authorization"] = f"Bearer {token}"

    req_data = None
    if data is not None:
        req_data = json.dumps(data).encode("utf-8")

    req = urllib.request.Request(url=url, data=req_data, headers=headers, method=method)

    try:
        with urllib.request.urlopen(req, timeout=timeout) as response:
            status = response.getcode()
            body_bytes = response.read()
            if not body_bytes:
                return status, {}
            try:
                return status, json.loads(body_bytes.decode("utf-8"))
            except json.JSONDecodeError:
                return status, body_bytes.decode("utf-8")
    except urllib.error.HTTPError as e:
        body_bytes = e.read()
        try:
            parsed = json.loads(body_bytes.decode("utf-8"))
        except Exception:
            parsed = body_bytes.decode("utf-8", errors="ignore")
        return e.code, parsed
    except urllib.error.URLError as e:
        return 0, str(e.reason)
    except Exception as e:
        return 0, str(e)


# =============================================================================
# DUMMY SEED DATA
# =============================================================================

DUMMY_USERS = [
    {
        "username": "admin_user",
        "email": "admin@cinedeck.dev",  # Matches frontend admin email in config.js
        "password": "AdminPassword123!",
        "role": "Admin",
    },
    {
        "username": "johndoe",
        "email": "john.doe@example.com",
        "password": "UserPassword123!",
        "role": "User",
    },
    {
        "username": "sarah_connor",
        "email": "sarah.connor@example.com",
        "password": "UserPassword123!",
        "role": "User",
    },
]

DUMMY_MOVIES = [
    {
        "title": "Inception",
        "description": "A thief who steals corporate secrets through dream-sharing technology is given the inverse task of planting an idea into the mind of a CEO.",
        "release_year": 2010,
        "genres": ["Action", "Sci-Fi", "Thriller"],
        "poster": "https://image.tmdb.org/t/p/w500/oYuLEt3zVCKq57qu2Dc8dVQI7Ni.jpg",
    },
    {
        "title": "The Dark Knight",
        "description": "When the menace known as the Joker wreaks havoc and chaos on the people of Gotham, Batman must accept one of the greatest psychological and physical tests of his ability to fight injustice.",
        "release_year": 2008,
        "genres": ["Action", "Crime", "Drama"],
        "poster": "https://image.tmdb.org/t/p/w500/qJ2tW6WMUDux911r6m7haRef0WH.jpg",
    },
    {
        "title": "Interstellar",
        "description": "A team of explorers travel through a wormhole in space in an attempt to ensure humanity's survival as Earth faces catastrophic blight.",
        "release_year": 2014,
        "genres": ["Sci-Fi", "Drama", "Adventure"],
        "poster": "https://image.tmdb.org/t/p/w500/gEU2QniE6E77NI6lCU6MxlNBvIx.jpg",
    },
    {
        "title": "Parasite",
        "description": "Greed and class discrimination threaten the newly formed symbiotic relationship between the wealthy Park family and the destitute Kim clan.",
        "release_year": 2019,
        "genres": ["Drama", "Thriller", "Comedy"],
        "poster": "https://image.tmdb.org/t/p/w500/7IiTTgloJzvGI1TAYymCfbfl3vT.jpg",
    },
    {
        "title": "Pulp Fiction",
        "description": "The lives of two mob hitmen, a boxer, a gangster and his wife, and a pair of diner bandits intertwine in four tales of violence and redemption.",
        "release_year": 1994,
        "genres": ["Crime", "Drama"],
        "poster": "https://image.tmdb.org/t/p/w500/d5iIlFn5s0ImszYzBPb8JPIfbXD.jpg",
    },
    {
        "title": "Spirited Away",
        "description": "During her family's move to the suburbs, a sullen 10-year-old girl wanders into a world ruled by gods, witches, and spirits, where humans are changed into beasts.",
        "release_year": 2001,
        "genres": ["Animation", "Family", "Fantasy"],
        "poster": "https://image.tmdb.org/t/p/w500/39wmItIWsg5sZMyRUHLkWBcuVCM.jpg",
    },
    {
        "title": "The Matrix",
        "description": "When a beautiful stranger leads computer hacker Neo to a forbidding underworld, he discovers the shocking truth--the life he knows is the elaborate deception of an evil cyber-intelligence.",
        "release_year": 1999,
        "genres": ["Action", "Sci-Fi"],
        "poster": "https://image.tmdb.org/t/p/w500/f89U3ADr1oiB1s9GkdPOEpXUk5H.jpg",
    },
    {
        "title": "Blade Runner 2049",
        "description": "Young Blade Runner K's discovery of a long-buried secret leads him to track down former Blade Runner Rick Deckard, who's been missing for thirty years.",
        "release_year": 2017,
        "genres": ["Sci-Fi", "Mystery", "Drama"],
        "poster": "https://image.tmdb.org/t/p/w500/gajva2L0rPYkEWjzgFlBXCAVBE5.jpg",
    },
    {
        "title": "Oppenheimer",
        "description": "The story of American scientist J. Robert Oppenheimer and his role in the development of the atomic bomb during World War II.",
        "release_year": 2023,
        "genres": ["Drama", "History"],
        "poster": "https://image.tmdb.org/t/p/w500/8Gxv8gSFCU0XGDykEGv7zR1n2ua.jpg",
    },
    {
        "title": "Whiplash",
        "description": "A promising young drummer enrolls at a cut-throat music conservatory where his dreams of greatness are mentored by an instructor who will stop at nothing to realize a student's potential.",
        "release_year": 2014,
        "genres": ["Drama", "Music"],
        "poster": "https://image.tmdb.org/t/p/w500/7fn624j5lj3xTme2SgiLCeuedmO.jpg",
    },
]


def test_healthchecks(base_url: str, check_direct_ports: bool) -> bool:
    log_step("1. Checking Service Health Endpoints")
    all_healthy = True

    # 1. API Gateway Health
    gw_url = f"{base_url}/healthcheck"
    status, body = make_request(gw_url)
    if status == 200:
        log_success(f"API Gateway is healthy ({gw_url}) -> {body}")
    else:
        log_fail(f"API Gateway health check failed ({gw_url}) [HTTP {status}]: {body}")
        all_healthy = False

    # 2. Check direct internal ports if requested or running directly on EC2
    if check_direct_ports:
        services = [
            ("User Service", "http://localhost:8083/healthcheck"),
            ("Movie Service", "http://localhost:8082/healthcheck"),
        ]
        for name, s_url in services:
            s_status, s_body = make_request(s_url)
            if s_status == 200:
                log_success(f"{name} (Direct) is healthy ({s_url}) -> {s_body}")
            else:
                log_warn(f"{name} (Direct port) not directly reachable: {s_body} (This is normal if port is inside Docker network)")

    return all_healthy


def register_and_login_users(base_url: str) -> Dict[str, Dict[str, Any]]:
    log_step("2. Registering & Authenticating Dummy Users")
    user_sessions = {}

    for u in DUMMY_USERS:
        email = u["email"]
        username = u["username"]
        password = u["password"]

        # 1. Register User
        reg_url = f"{base_url}/api/v1/users/register"
        reg_payload = {"username": username, "email": email, "password": password}
        status, reg_resp = make_request(reg_url, method="POST", data=reg_payload)

        user_id = None
        if status in [200, 201]:
            user_id = reg_resp.get("id")
            log_success(f"Registered '{username}' ({email}) -> ID: {user_id}")
        elif status == 400 or status == 500:
            log_warn(f"User '{username}' already exists or registration note: {reg_resp}")
        else:
            log_warn(f"Registration response for '{username}' [HTTP {status}]: {reg_resp}")

        # 2. Login User to get JWT
        login_url = f"{base_url}/api/v1/users/login"
        login_payload = {"email": email, "password": password}
        log_status, log_resp = make_request(login_url, method="POST", data=login_payload)

        if log_status == 200 and isinstance(log_resp, dict) and "token" in log_resp:
            token = log_resp["token"]
            log_success(f"Logged in '{username}' successfully (JWT token obtained)")
            user_sessions[email] = {
                "username": username,
                "email": email,
                "token": token,
                "role": u["role"],
                "user_id": user_id,
            }
        else:
            log_fail(f"Failed to log in '{username}' [HTTP {log_status}]: {log_resp}")

    return user_sessions


def seed_movies(base_url: str, admin_token: str) -> List[Dict[str, Any]]:
    log_step("3. Seeding Dummy Movies Catalog")
    created_movies = []

    # First, let's check existing movies
    list_url = f"{base_url}/api/v1/movies/"
    status, existing = make_request(list_url, method="GET")
    existing_titles = set()
    if status == 200 and isinstance(existing, list):
        existing_titles = {m.get("title", "").strip().lower() for m in existing if isinstance(m, dict)}
        log_success(f"Found {len(existing_titles)} existing movie(s) in catalog")

    for movie in DUMMY_MOVIES:
        if movie["title"].strip().lower() in existing_titles:
            log_warn(f"Movie '{movie['title']}' already exists, skipping duplicate")
            continue

        create_url = f"{base_url}/api/v1/movies/"
        c_status, c_resp = make_request(create_url, method="POST", data=movie, token=admin_token)

        if c_status in [200, 201] and isinstance(c_resp, dict):
            m_id = c_resp.get("id") or c_resp.get("_id")
            log_success(f"Added movie: '{movie['title']}' ({movie['release_year']}) [ID: {m_id}]")
            created_movies.append(c_resp)
        else:
            log_fail(f"Failed to add '{movie['title']}' [HTTP {c_status}]: {c_resp}")

    # Fetch updated list
    status, updated_list = make_request(list_url, method="GET")
    if status == 200 and isinstance(updated_list, list):
        log_success(f"Total movies now in database: {len(updated_list)}")
        return updated_list
    return created_movies


def seed_playlists_and_reviews(
    base_url: str, user_sessions: Dict[str, Dict[str, Any]], movies: List[Dict[str, Any]]
):
    log_step("4. Seeding Sample Playlists and Reviews")
    if not user_sessions or not movies:
        log_warn("Skipping playlists and reviews (requires at least 1 active user and 1 movie)")
        return

    # Extract test user
    user_info = next(iter(user_sessions.values()))
    token = user_info["token"]
    user_id = user_info.get("user_id") or "test-user-id"

    # 1. Create a Playlist
    pl_url = f"{base_url}/api/v1/list/"
    pl_payload = {
        "user_id": user_id,
        "playlist_name": "Must-Watch Masterpieces",
        "description": "Essential movies to watch before you die",
        "visiblity": True,
    }
    p_status, p_resp = make_request(pl_url, method="POST", data=pl_payload, token=token)
    if p_status in [200, 201] and isinstance(p_resp, dict):
        pl_id = p_resp.get("id") or p_resp.get("_id")
        log_success(f"Created playlist 'Must-Watch Masterpieces' [ID: {pl_id}]")

        # Add first 2 movies to the playlist
        if pl_id:
            for m in movies[:2]:
                m_id = m.get("id") or m.get("_id")
                if m_id:
                    add_url = f"{base_url}/api/v1/list/{pl_id}/movies/{m_id}"
                    a_status, a_resp = make_request(add_url, method="POST", token=token)
                    if a_status in [200, 201]:
                        log_success(f"  Added '{m.get('title')}' to playlist {pl_id}")
    else:
        log_warn(f"Playlist creation note [HTTP {p_status}]: {p_resp}")

    # 2. Add Sample Reviews to first 2 movies
    sample_reviews = [
        {"comment": "Mind-bending visual masterpiece! Christopher Nolan at his finest.", "score": 9.5},
        {"comment": "Heath Ledger gives an unforgettable performance. 10/10.", "score": 10.0},
    ]

    for idx, review_data in enumerate(sample_reviews):
        if idx < len(movies):
            target_movie = movies[idx]
            m_id = target_movie.get("id") or target_movie.get("_id")
            if m_id:
                rev_url = f"{base_url}/api/v1/movies/{m_id}/reviews"
                payload = {
                    "user_id": user_id,
                    "comment": review_data["comment"],
                    "score": review_data["score"],
                }
                r_status, r_resp = make_request(rev_url, method="POST", data=payload, token=token)
                if r_status in [200, 201]:
                    log_success(f"Added review to '{target_movie.get('title')}': {review_data['score']}/10")
                else:
                    log_warn(f"Review note for '{target_movie.get('title')}' [HTTP {r_status}]: {r_resp}")


def print_summary(base_url: str, user_sessions: Dict[str, Dict[str, Any]], movies_count: int):
    log_step("5. Verification & Seed Summary")
    print(f"{BOLD}Target Base URL:{RESET} {base_url}")
    print(f"{BOLD}Total Movies in Catalog:{RESET} {movies_count}")
    print(f"\n{BOLD}Seeded User Accounts:{RESET}")
    print(f"{'Role':<10} | {'Username':<15} | {'Email':<25} | {'Password':<18}")
    print("-" * 75)
    for u in DUMMY_USERS:
        print(f"{u['role']:<10} | {u['username']:<15} | {u['email']:<25} | {u['password']:<18}")

    print(f"\n{GREEN}{BOLD}✓ Microservices testing and initial dummy seeding completed successfully!{RESET}\n")


def main():
    parser = argparse.ArgumentParser(description="Test and seed MovieHub microservices.")
    parser.add_argument(
        "--base-url",
        default="http://localhost:8080",
        help="Base URL of API Gateway (default: http://localhost:8080)",
    )
    parser.add_argument(
        "--check-direct-ports",
        action="store_true",
        help="Also probe direct microservice ports (8082, 8083)",
    )
    args = parser.parse_args()

    base_url = args.base_url.rstrip("/")
    print(f"{BOLD}{GREEN}MovieHub Microservice Test & Seeder{RESET}")
    print(f"Connecting to Gateway at: {base_url}\n")

    # Step 1: Healthcheck
    is_healthy = test_healthchecks(base_url, args.check_direct_ports)
    if not is_healthy:
        print(f"\n{RED}{BOLD}ERROR: One or more critical services are not responding.{RESET}")
        print(f"Check your docker containers with: {YELLOW}docker compose ps{RESET} or {YELLOW}docker compose logs -f{RESET}")
        sys.exit(1)

    # Step 2: Register & Login Users
    user_sessions = register_and_login_users(base_url)
    if not user_sessions:
        print(f"\n{RED}{BOLD}ERROR: Could not log in any user. Check PostgreSQL database logs.{RESET}")
        sys.exit(1)

    # Choose admin token for movie seeding
    admin_session = user_sessions.get("admin@cinedeck.dev") or next(iter(user_sessions.values()))
    admin_token = admin_session["token"]

    # Step 3: Seed Movies
    movies = seed_movies(base_url, admin_token)

    # Step 4: Seed Playlists & Reviews
    seed_playlists_and_reviews(base_url, user_sessions, movies)

    # Step 5: Summary
    print_summary(base_url, user_sessions, len(movies))


if __name__ == "__main__":
    main()
