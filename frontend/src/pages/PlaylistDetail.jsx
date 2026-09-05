import { useCallback, useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { getPlaylist, removeMovieFromPlaylist, deletePlaylist } from "../api/playlists";
import { getMovie } from "../api/movies";
import { apiError } from "../api/client";
import { useAuth } from "../auth/AuthContext";
import { ErrorBox, EmptyState, Modal, Spinner } from "../components/ui";

export default function PlaylistDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { user } = useAuth();

  const [playlist, setPlaylist] = useState(null);
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [removingId, setRemovingId] = useState(null);
  const [confirmDelete, setConfirmDelete] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      const pl = await getPlaylist(id);
      setPlaylist(pl);
      if (pl.movie_ids?.length) {
        // Playlist stores only ObjectIDs — fetch each movie doc for display
        const all = await Promise.all(
          pl.movie_ids.map((mid) => getMovie(String(mid)).catch(() => null))
        );
        setMovies(all.filter(Boolean));
      } else {
        setMovies([]);
      }
    } catch (e) {
      setError(apiError(e));
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    load();
  }, [load]);

  async function handleRemove(movie) {
    setRemovingId(String(movie.id));
    setError("");
    try {
      await removeMovieFromPlaylist(id, movie.id);
      setMovies((prev) => prev.filter((m) => String(m.id) !== String(movie.id)));
      setPlaylist((pl) => ({
        ...pl,
        movie_ids: (pl.movie_ids || []).filter((mid) => String(mid) !== String(movie.id)),
      }));
    } catch (e) {
      setError(apiError(e));
    } finally {
      setRemovingId(null);
    }
  }

  async function handleDeletePlaylist() {
    try {
      await deletePlaylist(id);
      navigate("/playlists");
    } catch (e) {
      setError(apiError(e));
      setConfirmDelete(false);
    }
  }

  if (loading && !playlist) return <Spinner label="Loading playlist…" />;

  return (
    <div className="page">
      <button type="button" className="btn btn-ghost btn-sm back-btn" onClick={() => navigate("/playlists")}>
        ← My Playlists
      </button>

      <ErrorBox message={error} />

      {playlist && (
        <>
          <div className="page-head">
            <div>
              <h1>{playlist.name}</h1>
              <p className="muted">
                {playlist.description || "No description"} ·{" "}
                <span className={playlist.is_public ? "public" : "private"}>
                  {playlist.is_public ? "🌐 Public" : "🔒 Private"}
                </span>{" "}
                · {(playlist.movie_ids || []).length} movies
              </p>
            </div>
            {user?.id === playlist.user_id && (
              <button type="button" className="btn btn-danger-outline" onClick={() => setConfirmDelete(true)}>
                Delete playlist
              </button>
            )}
          </div>

          {movies.length === 0 ? (
            <EmptyState
              icon="🍿"
              title="This playlist is empty"
              hint={
                <>
                  Browse the catalog and hit “Save” to add movies here.{" "}
                  <Link to="/">Go to movies →</Link>
                </>
              }
            />
          ) : (
            <ul className="playlist-movie-list">
              {movies.map((m) => (
                <li key={String(m.id)} className="card playlist-movie-item">
                  <Link to={`/movies/${m.id}`} className="mini-poster">
                    {m.poster ? <img src={m.poster} alt="" loading="lazy" /> : "🎬"}
                  </Link>
                  <div className="pmi-info">
                    <Link to={`/movies/${m.id}`} className="pmi-title">
                      {m.title}
                    </Link>
                    <p className="muted">
                      {m.release_year} · {(m.genres || []).slice(0, 3).join(", ")}
                    </p>
                  </div>
                  {user?.id === playlist.user_id && (
                    <button
                      type="button"
                      className="btn btn-ghost btn-sm danger"
                      disabled={removingId === String(m.id)}
                      onClick={() => handleRemove(m)}
                    >
                      {removingId === String(m.id) ? "…" : "Remove"}
                    </button>
                  )}
                </li>
              ))}
            </ul>
          )}
        </>
      )}

      <Modal open={confirmDelete} title="Delete playlist?" onClose={() => setConfirmDelete(false)}>
        <p>
          Delete <strong>“{playlist?.name}”</strong>? The saved movies stay in the catalog but the
          playlist itself will be gone. This cannot be undone.
        </p>
        <div className="modal-actions">
          <button type="button" className="btn btn-ghost" onClick={() => setConfirmDelete(false)}>
            Cancel
          </button>
          <button type="button" className="btn btn-danger" onClick={handleDeletePlaylist}>
            Delete
          </button>
        </div>
      </Modal>
    </div>
  );
}
