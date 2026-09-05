import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { getMovies } from "../api/movies";
import { apiError } from "../api/client";
import { useAuth } from "../auth/AuthContext";
import { ErrorBox, Spinner } from "../components/ui";

export default function Dashboard() {
  const { user } = useAuth();
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    async function load() {
      try {
        setMovies(await getMovies({ limit: 0 }));
      } catch (e) {
        setError(apiError(e));
      } finally {
        setLoading(false);
      }
    }
    load();
  }, []);

  const stats = useMemo(() => {
    const genreCount = {};
    let totalReviews = 0;
    for (const m of movies) {
      for (const g of m.genres || []) genreCount[g] = (genreCount[g] || 0) + 1;
      totalReviews += m.reviews?.length || 0;
    }
    const topGenres = Object.entries(genreCount)
      .sort((a, b) => b[1] - a[1])
      .slice(0, 5);
    const recent = [...movies]
      .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
      .slice(0, 5);
    return {
      movieCount: movies.length,
      reviewCount: totalReviews,
      avgRating:
        movies.length > 0
          ? (
              movies.reduce((s, m) => s + (m.rating?.average || 0), 0) / movies.length
            ).toFixed(1)
          : "—",
      topGenres,
      recent,
    };
  }, [movies]);

  if (loading) return <Spinner label="Crunching numbers…" />;

  return (
    <div>
      <h1>Dashboard</h1>
      <p className="muted">
        Signed in as <strong>{user.username}</strong>
      </p>

      <ErrorBox message={error} />

      <div className="stat-grid">
        <div className="card stat-card">
          <span className="stat-num">{stats.movieCount}</span>
          <span className="stat-label">Movies</span>
        </div>
        <div className="card stat-card">
          <span className="stat-num">{stats.reviewCount}</span>
          <span className="stat-label">Reviews</span>
        </div>
        <div className="card stat-card">
          <span className="stat-num">{stats.avgRating}</span>
          <span className="stat-label">Avg rating</span>
        </div>
      </div>

      <div className="dash-columns">
        <section className="card">
          <h3>Top genres</h3>
          {stats.topGenres.length === 0 ? (
            <p className="muted">No data yet.</p>
          ) : (
            stats.topGenres.map(([g, count]) => (
              <div key={g} className="bar-row">
                <span className="bar-label">{g}</span>
                <div className="bar-track">
                  <div
                    className="bar-fill"
                    style={{ width: `${(count / stats.movieCount) * 100}%` }}
                  />
                </div>
                <span className="bar-num">{count}</span>
              </div>
            ))
          )}
        </section>

        <section className="card">
          <h3>Recently added</h3>
          {stats.recent.length === 0 ? (
            <p className="muted">No movies yet.</p>
          ) : (
            <ul className="recent-list">
              {stats.recent.map((m) => (
                <li key={String(m.id)}>
                  <Link to={`/movies/${m.id}`}>{m.title}</Link>
                  <span className="muted"> {new Date(m.created_at).toLocaleDateString()}</span>
                </li>
              ))}
            </ul>
          )}
          <Link to="/admin/movies" className="btn btn-outline btn-block">
            Manage movies →
          </Link>
        </section>
      </div>
    </div>
  );
}
