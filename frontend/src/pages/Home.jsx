import { useCallback, useEffect, useMemo, useState } from "react";
import { getMovies } from "../api/movies";
import { getUserPlaylists } from "../api/playlists";
import { useAuth } from "../auth/AuthContext";
import { apiError } from "../api/client";
import { GENRES } from "../config";
import MovieCard from "../components/MovieCard";
import SaveToPlaylistModal from "../components/SaveToPlaylistModal";
import { ErrorBox, EmptyState, Spinner } from "../components/ui";

const PAGE_SIZE = 12;

export default function Home() {
  const { isAuthenticated, user } = useAuth();
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  // server-side filters
  const [filters, setFilters] = useState({ year: "", actor: "", director: "", min_rating: "" });
  const [genre, setGenre] = useState("");
  const [search, setSearch] = useState(""); // client-side title search
  const [page, setPage] = useState(1);
  const [showFilters, setShowFilters] = useState(false);

  const [saveTarget, setSaveTarget] = useState(null); // movie being saved
  const [myPlaylists, setMyPlaylists] = useState([]);

  const loadMovies = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const list = await getMovies({ ...filters, limit: 0 }); // limit=0 → all
      setMovies(list);
    } catch (e) {
      setError(apiError(e));
    } finally {
      setLoading(false);
    }
  }, [filters]);

  useEffect(() => {
    loadMovies();
  }, [loadMovies]);

  useEffect(() => {
    setPage(1);
  }, [filters, genre, search]);

  useEffect(() => {
    async function loadPlaylists() {
      if (!isAuthenticated || !user?.id) return;
      try {
        setMyPlaylists(await getUserPlaylists(user.id));
      } catch {
        /* non-critical */
      }
    }
    loadPlaylists();
  }, [isAuthenticated, user?.id]);

  function playlistsContaining(movieId) {
    return myPlaylists.filter((pl) =>
      (pl.movie_ids || []).some((id) => String(id) === String(movieId))
    );
  }

  const visible = useMemo(() => {
    let list = movies;
    if (genre) list = list.filter((m) => (m.genres || []).includes(genre));
    if (search.trim()) {
      const q = search.trim().toLowerCase();
      list = list.filter(
        (m) =>
          m.title.toLowerCase().includes(q) ||
          m.description.toLowerCase().includes(q)
      );
    }
    return list;
  }, [movies, genre, search]);

  const totalPages = Math.max(1, Math.ceil(visible.length / PAGE_SIZE));
  const pageItems = visible.slice((page - 1) * PAGE_SIZE, page * PAGE_SIZE);

  function onFilter(e) {
    e.preventDefault();
    const data = new FormData(e.target);
    setFilters({
      year: data.get("year")?.trim() || "",
      actor: data.get("actor")?.trim() || "",
      director: data.get("director")?.trim() || "",
      min_rating: data.get("min_rating")?.trim() || "",
    });
    setShowFilters(false);
  }

  return (
    <div className="page">
      <section className="hero">
        <h1>
          Discover your next <span className="accent">favorite film</span>
        </h1>
        <p className="muted">
          Browse the catalog, save movies to playlists, and share your ratings.
        </p>

        <form className="search-bar" onSubmit={(e) => e.preventDefault()}>
          <input
            type="search"
            placeholder="Search by title or description…"
            value={search}
            onChange={(e) => setSearch(e.target.value)}
          />
          <button type="button" className="btn btn-outline" onClick={() => setShowFilters((v) => !v)}>
            Filters {showFilters ? "▲" : "▼"}
          </button>
        </form>

        <div className="genre-strip">
          <button
            type="button"
            className={`tag tag-btn ${genre === "" ? "active" : ""}`}
            onClick={() => setGenre("")}
          >
            All
          </button>
          {GENRES.map((g) => (
            <button
              key={g}
              type="button"
              className={`tag tag-btn ${genre === g ? "active" : ""}`}
              onClick={() => setGenre(g)}
            >
              {g}
            </button>
          ))}
        </div>

        {showFilters && (
          <form className="card filter-card" onSubmit={onFilter}>
            <label>
              Year
              <input name="year" type="number" defaultValue={filters.year} placeholder="e.g. 2019" />
            </label>
            <label>
              Min rating
              <input
                name="min_rating"
                type="number"
                step="0.1"
                min="0"
                max="10"
                defaultValue={filters.min_rating}
                placeholder="e.g. 8"
              />
            </label>
            <label>
              Actor
              <input name="actor" defaultValue={filters.actor} placeholder="e.g. DiCaprio" />
            </label>
            <label>
              Director
              <input name="director" defaultValue={filters.director} placeholder="e.g. Nolan" />
            </label>
            <div className="filter-actions">
              <button className="btn btn-primary btn-sm" type="submit">
                Apply
              </button>
              <button
                type="button"
                className="btn btn-ghost btn-sm"
                onClick={() =>
                  setFilters({ year: "", actor: "", director: "", min_rating: "" })
                }
              >
                Reset
              </button>
            </div>
          </form>
        )}
      </section>

      <section>
        <p className="results-count muted">
          {loading ? "Loading…" : `${visible.length} movie${visible.length === 1 ? "" : "s"} found`}
        </p>

        <ErrorBox message={error} onRetry={loadMovies} />

        {loading && movies.length === 0 ? (
          <Spinner label="Fetching the catalog…" />
        ) : !loading && visible.length === 0 && !error ? (
          <EmptyState icon="🍿" title="No movies match" hint="Try clearing filters or searching for something else." />
        ) : (
          <div className="movie-grid">
            {pageItems.map((m) => (
              <MovieCard
                key={String(m.id)}
                movie={m}
                onSave={isAuthenticated ? (mv) => setSaveTarget(mv) : null}
                savedIn={playlistsContaining(String(m.id))}
              />
            ))}
          </div>
        )}

        {!loading && totalPages > 1 && (
          <div className="pagination">
            <button type="button" className="btn btn-ghost" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>
              ← Prev
            </button>
            <span>
              Page {page} of {totalPages}
            </span>
            <button
              type="button"
              className="btn btn-ghost"
              disabled={page >= totalPages}
              onClick={() => setPage((p) => p + 1)}
            >
              Next →
            </button>
          </div>
        )}
      </section>

      <SaveToPlaylistModal
        movie={saveTarget}
        open={Boolean(saveTarget)}
        onClose={() => setSaveTarget(null)}
      />
    </div>
  );
}
