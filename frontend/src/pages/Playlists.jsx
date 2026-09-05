import { useCallback, useEffect, useState } from "react";
import { Link } from "react-router-dom";
import { getUserPlaylists, createPlaylist, deletePlaylist } from "../api/playlists";
import { useAuth } from "../auth/AuthContext";
import { apiError } from "../api/client";
import { ErrorBox, EmptyState, Modal, Spinner } from "../components/ui";
import { ErrorInline } from "../components/SaveToPlaylistModal";

export default function Playlists() {
  const { user } = useAuth();
  const [playlists, setPlaylists] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [createOpen, setCreateOpen] = useState(false);
  const [confirmDelete, setConfirmDelete] = useState(null);
  const [form, setForm] = useState({ name: "", description: "", isPublic: false });
  const [busy, setBusy] = useState(false);
  const [formError, setFormError] = useState("");

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      setPlaylists(await getUserPlaylists(user.id));
    } catch (e) {
      setError(apiError(e));
    } finally {
      setLoading(false);
    }
  }, [user?.id]);

  useEffect(() => {
    if (user?.id) load();
  }, [user?.id, load]);

  async function handleCreate(e) {
    e.preventDefault();
    if (!form.name.trim()) return;
    setBusy(true);
    setFormError("");
    try {
      await createPlaylist({
        userId: user.id,
        name: form.name.trim(),
        description: form.description.trim(),
        isPublic: form.isPublic,
      });
      setForm({ name: "", description: "", isPublic: false });
      setCreateOpen(false);
      await load();
    } catch (err) {
      setFormError(apiError(err));
    } finally {
      setBusy(false);
    }
  }

  async function handleDelete() {
    try {
      await deletePlaylist(confirmDelete.id);
      setConfirmDelete(null);
      await load();
    } catch (e) {
      setError(apiError(e));
      setConfirmDelete(null);
    }
  }

  return (
    <div className="page">
      <div className="page-head">
        <h1>My Playlists</h1>
        <button type="button" className="btn btn-primary" onClick={() => setCreateOpen(true)}>
          + New playlist
        </button>
      </div>

      <ErrorBox message={error} onRetry={load} />

      {loading ? (
        <Spinner label="Loading playlists…" />
      ) : !error && playlists.length === 0 ? (
        <EmptyState
          icon="🎞️"
          title="No playlists yet"
          hint="Create one and start saving movies to it."
        />
      ) : (
        <div className="playlist-grid">
          {playlists.map((pl) => (
            <div key={String(pl.id)} className="card playlist-card">
              <Link to={`/playlists/${pl.id}`} className="playlist-card-link">
                <h3>{pl.name}</h3>
                <p className="muted">{pl.description || "No description"}</p>
              </Link>
              <div className="playlist-card-foot">
                <span className="tag">{(pl.movie_ids || []).length} movies</span>
                <span className={`visibility ${pl.is_public ? "public" : "private"}`}>
                  {pl.is_public ? "🌐 Public" : "🔒 Private"}
                </span>
                <span className="foot-spacer" />
                <button
                  type="button"
                  className="btn btn-ghost btn-sm danger"
                  title="Delete playlist"
                  onClick={() => setConfirmDelete(pl)}
                >
                  🗑
                </button>
              </div>
            </div>
          ))}
        </div>
      )}

      <Modal open={createOpen} title="New playlist" onClose={() => setCreateOpen(false)}>
        <form onSubmit={handleCreate} className="stack">
          <label className="field">
            <span>Name *</span>
            <input
              autoFocus
              value={form.name}
              onChange={(e) => setForm((f) => ({ ...f, name: e.target.value }))}
              placeholder="e.g. Weekend watchlist"
              required
            />
          </label>
          <label className="field">
            <span>Description</span>
            <textarea
              rows={3}
              value={form.description}
              onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
              placeholder="What's this playlist about?"
            />
          </label>
          <label className="check-field">
            <input
              type="checkbox"
              checked={form.isPublic}
              onChange={(e) => setForm((f) => ({ ...f, isPublic: e.target.checked }))}
            />
            <span>Make it public (visible to others)</span>
          </label>
          <ErrorInline message={formError} />
          <button type="submit" className="btn btn-primary btn-block" disabled={busy || !form.name.trim()}>
            {busy ? "Creating…" : "Create playlist"}
          </button>
        </form>
      </Modal>

      <Modal open={Boolean(confirmDelete)} title="Delete playlist?" onClose={() => setConfirmDelete(null)}>
        <p>
          Delete <strong>“{confirmDelete?.name}”</strong>? This cannot be undone.
        </p>
        <div className="modal-actions">
          <button type="button" className="btn btn-ghost" onClick={() => setConfirmDelete(null)}>
            Cancel
          </button>
          <button type="button" className="btn btn-danger" onClick={handleDelete}>
            Delete
          </button>
        </div>
      </Modal>
    </div>
  );
}
