import { useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { getUserPlaylists, createPlaylist, addMovieToPlaylist } from "../api/playlists";
import { useAuth } from "../auth/AuthContext";
import { apiError } from "../api/client";
import { Modal, Spinner } from "./ui";

/**
 * Modal that saves a movie into one of the user's playlists.
 * Also supports creating a brand-new playlist inline.
 */
export default function SaveToPlaylistModal({ movie, open, onClose }) {
  const { user } = useAuth();
  const [playlists, setPlaylists] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [status, setStatus] = useState("");
  const [creating, setCreating] = useState(false);
  const [newName, setNewName] = useState("");
  const [busyId, setBusyId] = useState(null);

  async function refresh() {
    setLoading(true);
    setError("");
    try {
      setPlaylists(await getUserPlaylists(user.id));
    } catch (e) {
      setError(apiError(e));
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    if (open && user) refresh();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [open, user?.id]);

  useEffect(() => {
    if (!open) {
      setStatus("");
      setCreating(false);
      setNewName("");
      setError("");
    }
  }, [open]);

  if (!movie) return null;

  const alreadyIn = (pl) => (pl.movie_ids || []).some((id) => String(id) === String(movie.id));

  async function addTo(playlist) {
    setBusyId(playlist.id);
    setError("");
    try {
      await addMovieToPlaylist(playlist.id, movie.id);
      setStatus(`Saved to “${playlist.name}”`);
      setPlaylists((prev) =>
        prev.map((p) =>
          p.id === playlist.id
            ? { ...p, movie_ids: [...(p.movie_ids || []), { toString: () => movie.id }] }
            : p
        )
      );
    } catch (e) {
      setError(apiError(e));
    } finally {
      setBusyId(null);
    }
  }

  async function handleCreate(e) {
    e.preventDefault();
    if (!newName.trim()) return;
    setError("");
    try {
      const pl = await createPlaylist({ userId: user.id, name: newName.trim() });
      setNewName("");
      setCreating(false);
      await addMovieToPlaylist(pl.id, movie.id);
      setStatus(`Saved to new playlist “${pl.name}”`);
      await refresh();
    } catch (err) {
      setError(apiError(err));
    }
  }

  return (
    <Modal open={open} title={`Save “${movie.title}” to…`} onClose={onClose}>
      <ErrorInline message={error} />
      {status && <div className="toast-success">{status}</div>}

      {loading ? (
        <Spinner label="Loading playlists…" />
      ) : playlists.length === 0 && !creating ? (
        <p className="muted">You have no playlists yet. Create your first one below.</p>
      ) : (
        <ul className="playlist-pick-list">
          {playlists.map((pl) => {
            const saved = alreadyIn(pl);
            return (
              <li key={pl.id} className="playlist-pick-item">
                <span>
                  🎞️ {pl.name}
                  {saved && <em className="already-in"> · already saved</em>}
                </span>
                <button
                  type="button"
                  className="btn btn-sm"
                  disabled={saved || busyId === pl.id}
                  onClick={() => addTo(pl)}
                >
                  {saved ? "✓" : busyId === pl.id ? "…" : "Add"}
                </button>
              </li>
            );
          })}
        </ul>
      )}

      {creating ? (
        <form className="inline-form" onSubmit={handleCreate}>
          <input
            autoFocus
            value={newName}
            onChange={(e) => setNewName(e.target.value)}
            placeholder="New playlist name"
          />
          <button className="btn btn-primary btn-sm" type="submit" disabled={!newName.trim()}>
            Create & save
          </button>
          <button type="button" className="btn btn-ghost btn-sm" onClick={() => setCreating(false)}>
            Cancel
          </button>
        </form>
      ) : (
        <button type="button" className="btn btn-outline btn-block" onClick={() => setCreating(true)}>
          + New playlist
        </button>
      )}

      <p className="muted center">
        Manage everything in <Link to="/playlists">My Playlists</Link>
      </p>
    </Modal>
  );
}

export function ErrorInline({ message }) {
  if (!message) return null;
  return <div className="error-inline">{message}</div>;
}
