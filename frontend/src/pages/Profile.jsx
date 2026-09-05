import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { getUserPlaylists } from "../api/playlists";
import { apiError } from "../api/client";
import { useAuth } from "../auth/AuthContext";
import { TOKEN_KEY } from "../config";

function decodeClaims(token) {
  try {
    return JSON.parse(atob(token.split(".")[1].replace(/-/g, "+").replace(/_/g, "/")));
  } catch {
    return {};
  }
}

export default function Profile() {
  const { user, logout } = useAuth();
  const [playlists, setPlaylists] = useState([]);
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(true);

  const claims = decodeClaims(localStorage.getItem(TOKEN_KEY) || "");
  const loginAt = claims.iat ? new Date(claims.iat * 1000) : null;
  const expiresAt = claims.exp ? new Date(claims.exp * 1000) : null;

  useEffect(() => {
    async function load() {
      try {
        setPlaylists(await getUserPlaylists(user.id));
      } catch (e) {
        setError(apiError(e));
      } finally {
        setLoading(false);
      }
    }
    if (user?.id) load();
    else setLoading(false);
  }, [user?.id]);

  if (!user) return null;

  return (
    <div className="page">
      <div className="profile-card card">
        <span className="avatar avatar-lg">{(user.username || "U")[0].toUpperCase()}</span>
        <h1>{user.username}</h1>
        <p className="muted">{user.email || "Email not stored (Google account?)"}</p>

        <dl className="profile-meta">
          <div>
            <dt>User ID</dt>
            <dd>
              <code>{user.id || "—"}</code>
            </dd>
          </div>
          <div>
            <dt>Session started</dt>
            <dd>{loginAt ? loginAt.toLocaleString() : "—"}</dd>
          </div>
          <div>
            <dt>Token expires</dt>
            <dd>{expiresAt ? expiresAt.toLocaleTimeString() : "—"}</dd>
          </div>
          <div>
            <dt>Playlists</dt>
            <dd>{loading ? "…" : playlists.length}</dd>
          </div>
        </dl>

        {error && <div className="error-inline">{error}</div>}

        <div className="profile-actions">
          <Link to="/playlists" className="btn btn-outline">
            🎞️ My Playlists ({playlists.length})
          </Link>
          <Link to="/" className="btn btn-outline">
            🎬 Browse movies
          </Link>
          <button type="button" className="btn btn-danger" onClick={logout}>
            Logout
          </button>
        </div>
      </div>

      {!loading && playlists.length > 0 && (
        <section style={{ marginTop: 24 }}>
          <h2>Your playlists</h2>
          <ul className="playlist-movie-list">
            {playlists.map((pl) => (
              <li key={String(pl.id)} className="card playlist-movie-item">
                <div className="pmi-info">
                  <Link to={`/playlists/${pl.id}`} className="pmi-title">
                    {pl.name}
                  </Link>
                  <p className="muted">
                    {(pl.movie_ids || []).length} movies ·{" "}
                    {pl.is_public ? "🌐 Public" : "🔒 Private"}
                  </p>
                </div>
              </li>
            ))}
          </ul>
        </section>
      )}
    </div>
  );
}
