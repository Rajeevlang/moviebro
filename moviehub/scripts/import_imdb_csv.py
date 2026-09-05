#!/usr/bin/env python3
"""
=============================================================================
IMDb Top 1000 CSV to MongoDB Direct Seeder for GinMovieAPI
Reads 'IMDB top 1000.csv' and imports directly into MongoDB (movie_catalog.movies).
=============================================================================
"""

import csv
import json
import os
import re
import subprocess
import sys
from datetime import datetime, timezone

GENRE_POSTERS = {
    "action": "https://images.unsplash.com/photo-1509347528160-9a9e33742cdb?w=800",
    "adventure": "https://images.unsplash.com/photo-1506744038136-46273834b3fb?w=800",
    "animation": "https://images.unsplash.com/photo-1534447677768-be436bb09401?w=800",
    "biography": "https://images.unsplash.com/photo-1485846234645-a62644f84728?w=800",
    "comedy": "https://images.unsplash.com/photo-1514306191717-452ec28c7814?w=800",
    "crime": "https://images.unsplash.com/photo-1594909122845-11baa439b7bf?w=800",
    "drama": "https://images.unsplash.com/photo-1517604931442-7e0c8ed2963c?w=800",
    "fantasy": "https://images.unsplash.com/photo-1518709268805-4e9042af9f23?w=800",
    "history": "https://images.unsplash.com/photo-1461360370896-922624d12aa1?w=800",
    "horror": "https://images.unsplash.com/photo-1509248961158-e54f6934749c?w=800",
    "music": "https://images.unsplash.com/photo-1511671782779-c97d3d27a1d4?w=800",
    "mystery": "https://images.unsplash.com/photo-1519074069444-1ba4eae16e61?w=800",
    "romance": "https://images.unsplash.com/photo-1518199266791-5375a83190b7?w=800",
    "sci-fi": "https://images.unsplash.com/photo-1446776811953-b23d57bd21aa?w=800",
    "thriller": "https://images.unsplash.com/photo-1536440136628-849c177e76a1?w=800",
    "western": "https://images.unsplash.com/photo-1568605117036-5fe5e7bab0b7?w=800",
    "default": "https://images.unsplash.com/photo-1489599849927-2ee91cede3ba?w=800"
}

def parse_title_and_year(raw_title):
    """
    Parses '1. The Shawshank Redemption (1994)' into:
    title = 'The Shawshank Redemption'
    year = 1994
    """
    cleaned = raw_title.strip()
    # Strip leading ranking digits like '1. ', '100. '
    cleaned = re.sub(r'^\d+\.\s*', '', cleaned)
    
    # Extract year from end like '(1994)'
    year_match = re.search(r'\((\d{4})\)$', cleaned)
    if year_match:
        year = int(year_match.group(1))
        title = cleaned[:year_match.start()].strip()
    else:
        year = 2000
        title = cleaned.strip()
        
    return title, year

def parse_cast(cast_str):
    """
    Parses 'Director: Christopher Nolan | Stars: Christian Bale, Heath Ledger'
    """
    directors = []
    actors = []
    if not cast_str:
        return directors, actors
        
    parts = cast_str.split('|')
    for part in parts:
        part = part.strip()
        if part.startswith("Director:") or part.startswith("Directors:"):
            names = re.sub(r'^Directors?:\s*', '', part).split(',')
            for n in names:
                if n.strip():
                    directors.append({"name": n.strip(), "role": "Director"})
        elif part.startswith("Star:") or part.startswith("Stars:"):
            names = re.sub(r'^Stars?:\s*', '', part).split(',')
            for n in names:
                if n.strip():
                    actors.append({"name": n.strip(), "role": "Lead Actor"})
                    
    return directors, actors

def get_poster_for_genres(genres):
    for g in genres:
        key = g.lower().strip()
        if key in GENRE_POSTERS:
            return GENRE_POSTERS[key]
    return GENRE_POSTERS["default"]

def convert_csv_to_mongo_docs(csv_filepath, limit=None):
    if not os.path.exists(csv_filepath):
        print(f"❌ Error: CSV file not found at '{csv_filepath}'")
        sys.exit(1)

    movies = []
    seen_titles = set()

    with open(csv_filepath, mode='r', encoding='utf-8-sig') as f:
        reader = csv.DictReader(f)
        for row in reader:
            raw_title = row.get("Title", "")
            if not raw_title:
                continue

            title, year = parse_title_and_year(raw_title)
            if not title or title.lower() in seen_titles:
                continue
            seen_titles.add(title.lower())

            # Genres
            raw_genre = row.get("Genre", "Drama")
            genres = [g.strip() for g in raw_genre.split(",") if g.strip()]
            if not genres:
                genres = ["Drama"]

            # Description (capped at 490 chars)
            desc = row.get("Description", "No description available.").strip()
            if len(desc) > 490:
                desc = desc[:487] + "..."

            # Rating
            try:
                rate = float(row.get("Rate", 7.5))
            except ValueError:
                rate = 7.5

            # Cast
            directors, actors = parse_cast(row.get("Cast", ""))

            # Poster
            poster = get_poster_for_genres(genres)

            movie_doc = {
                "title": title,
                "description": desc,
                "release_year": year,
                "genres": genres,
                "directors": directors,
                "actors": actors,
                "poster": poster,
                "rating": {
                    "average": rate,
                    "count": 100
                },
                "reviews": [],
                "created_at": "__NEW_DATE__",
                "updated_at": "__NEW_DATE__"
            }
            movies.append(movie_doc)

            if limit and len(movies) >= limit:
                break

    return movies

def main():
    script_dir = os.path.dirname(os.path.abspath(__file__))
    default_csv = os.path.join(script_dir, "imbd.csv")
    csv_file = sys.argv[1] if len(sys.argv) > 1 else default_csv

    print(f"📖 Reading IMDb dataset from '{csv_file}'...")
    movies = convert_csv_to_mongo_docs(csv_file)
    print(f"🎬 Parsed {len(movies)} unique movies from CSV.")

    # Write temporary JSON file for mongoimport
    temp_json = os.path.join(script_dir, "movies_import.json")
    with open(temp_json, "w", encoding="utf-8") as f:
        json.dump(movies, f, indent=2)

    print(f"💾 Generated '{temp_json}'. Importing directly into MongoDB container...")

    # Find the running MongoDB container ID/name
    try:
        # Try to find minikube container first, fallback to standard docker-compose name
        container_cmd = subprocess.run(
            ['docker', 'ps', '-qf', 'name=k8s_mongodb_mongodb'],
            capture_output=True, text=True, check=True
        )
        container_id = container_cmd.stdout.strip().split('\n')[0]
        
        if not container_id:
            container_cmd = subprocess.run(
                ['docker', 'ps', '-qf', 'name=mongodb'],
                capture_output=True, text=True, check=True
            )
            # Find the first one that is NOT a k8s pause container
            ids = [i for i in container_cmd.stdout.strip().split('\n') if i]
            for cid in ids:
                # check image name
                img_cmd = subprocess.run(['docker', 'inspect', '-f', '{{.Config.Image}}', cid], capture_output=True, text=True)
                if 'pause' not in img_cmd.stdout:
                    container_id = cid
                    break
            if not container_id:
                print("❌ No running MongoDB container found! Run 'docker compose up -d' or minikube first.")
                sys.exit(1)

        print(f"🐳 Found MongoDB container: {container_id}")

        mongo_user = os.getenv("MONGO_ROOT_USER", "admin")
        mongo_pass = os.getenv("MONGO_ROOT_PASSWORD", "password123")

        # Execute direct import by piping script via stdin
        raw_json = json.dumps(movies)
        raw_json = raw_json.replace('"__NEW_DATE__"', 'new Date()')
        mongosh_script = f"db.movies.insertMany({raw_json});"
        proc = subprocess.Popen(
            [
                'docker', 'exec', '-i', container_id,
                'mongosh', '-u', mongo_user, '-p', mongo_pass,
                '--authenticationDatabase', 'admin', 'movie_catalog', '--quiet'
            ],
            stdin=subprocess.PIPE,
            stdout=subprocess.PIPE,
            stderr=subprocess.PIPE,
            text=True
        )
        stdout, stderr = proc.communicate(input=mongosh_script)

        if proc.returncode == 0:
            print(f"🎉 Successfully imported {len(movies)} movies into MongoDB database 'movie_catalog.movies'!")
        else:
            print(f"⚠️ Insertion output: {stdout}\n{stderr}")

    except Exception as e:
        print(f"❌ Error during import: {e}")
    finally:
        if os.path.exists(temp_json):
            os.remove(temp_json)

if __name__ == "__main__":
    main()
