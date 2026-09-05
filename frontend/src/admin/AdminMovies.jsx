import { useCallback, useEffect, useState } from "react";
import { getMovies, createMovie, updateMovie, deleteMovie } from "../api/movies";
import { apiError } from "../api/client";
import { GENRES } from "../config";
import { ErrorBox, Modal, Spinner } from "../components/ui";
import { ErrorInline } from "../components/SaveToPlaylistModal";

const EMPTY_FORM = {
  title: "",
  description: "",
  release_year: new Date().getFullYear(),
  poster: "",
  genres: [],
  directors: "",
  actors: "",
};

export default function AdminMovies() {
  const [movies, setMovies] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");
  const [search, setSearch] = useState("");

  // create / edit
  const [editing, setEditing] = useState(null); // null | "new" | movie
  const [form, setForm] = useState(EMPTY_FORM);
  const [busy, setBusy] = useState(false);
  const [formError, setFormError] = useState("");

  const [confirmDelete, setConfirmDelete] = useState(null);
  const [deleteBusy, setDeleteBusy] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      setMovies(await getMovies({ limit: 0 }));
    } catch (e) {
      setError(apiError(e));
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    load();
  }, [load]);

  function openCreate() {
    setForm(EMPTY_FORM);
    setEditing("new");
    setFormError("");
  }

  function openEdit(movie) {
    setForm({
      title: movie.title,
      description: movie.description,
      release_year: movie.release_year,
      poster: movie.poster,
      genres: [...(movie.genres || [])],
      directors: (movie.directors || []).map((d) => d.name).join(", "),
      actors: (movie.actors || []).map((a) => a.name).join(", "),
    });
    setEditing(movie);
    setFormError("");
  }

  async function handleSubmit(e) {
    e.preventDefault();
    setBusy(true);
    setFormError("");
    try {
      const payload = {
        title: form.title.trim(),
        description: form.description.trim(),
        release_year: Number(form.release_year),
        poster: form.poster.trim(),
        genres: form.genres,
      };
      if (editing === "new") {
        await createMovie(payload);
      } else {
        // PATCH supports partial updates; directors/actors aren't in the PATCH input,
        // so only send the supported fields.
        await updateMovie(editing.id, payload);
      }
      setEditing(null);
      await load();
    } catch (err) {
      setFormError(apiError(err));
    } finally {
      setBusy(false);
    }
  }

  async function handleDelete() {
    setDeleteBusy(true);
    setError("");
    try {
      await deleteMovie(confirmDelete.id);
      setMovies((prev) => prev.filter((m) => String(m.id) !== String(confirmDelete.id)));
      setConfirmDelete(null);
    } catch (e) {
      setError(apiError(e));
      setConfirmDelete(null);
    } finally {
      setDeleteBusy(false);
    }
  }

  function toggleGenre(g) {
    setForm((f) => ({
      ...f,
      genres: f.genres.includes(g) ? f.genres.filter((x) => x !== g) : [...f.genres, g],
    }));
  }

  const filtered = movies.filter(
    (m) =>
      !search.trim() ||
      m.title.toLowerCase().includes(search.toLowerCase()) ||
      m.release_year.toString().includes(search)
  );

  return (
    <div>
      <div className="page-head">
        <h1>Movies</h1>
        <button type="button" className="btn btn-primary" onClick={openCreate}>
          + Add movie
        </button>
      </div>

      <input
        className="table-search"
        type="search"
        placeholder="Filter by title or year…"
        value={search}
        onChange={(e) => setSearch(e.target.value)}
      />

      <ErrorBox message={error} />

      {loading ? (
        <Spinner label="Loading movies…" />
      ) : (
        <div className="card table-card">
          <table className="admin-table">
            <thead>
              <tr>
                <th>Poster</th>
                <th>Title</th>
                <th>Year</th>
                <th>Genres</th>
                <th>Rating</th>
                <th>Reviews</th>
                <th className="right">Actions</th>
              </tr>
            </thead>
            <tbody>
              {filtered.map((m) => (
                <tr key={String(m.id)}>
                  <td>
                    <div className="table-poster">
                      {m.poster ? <img src={m.poster} alt="" loading="lazy" /> : "🎬"}
                    </div>
                  </td>
                  <td>{m.title}</td>
                  <td>{m.release_year}</td>
                  <td>
                    <span className="tag">{(m.genres || []).slice(0, 2).join(", ")}</span>
                    {(m.genres || []).length > 2 && (
                      <span className="muted"> +{(m.genres || []).length - 2}</span>
                    )}
                  </td>
                  <td>{m.rating?.average ? m.rating.average.toFixed(1) : "—"}</td>
                  <td>{m.reviews?.length || 0}</td>
                  <td className="right actions-cell">
                    <button type="button" className="btn btn-outline btn-sm" onClick={() => openEdit(m)}>
                      Edit
                    </button>
                    <button
                      type="button"
                      className="btn btn-danger-outline btn-sm"
                      onClick={() => setConfirmDelete(m)}
                    >
                      Delete
                    </button>
                  </td>
                </tr>
              ))}
              {filtered.length === 0 && (
                <tr>
                  <td colSpan={7} className="center muted">
                    No movies found.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      )}

      <Modal
        open={Boolean(editing)}
        title={editing === "new" ? "Add movie" : `Edit “${editing?.title}”`}
        onClose={() => setEditing(null)}
      >
        <form onSubmit={handleSubmit} className="stack">
          <label className="field">
            <span>Title *</span>
            <input
              required
              minLength={2}
              value={form.title}
              onChange={(e) => setForm((f) => ({ ...f, title: e.target.value }))}
              placeholder="Movie title"
            />
          </label>

          <label className="field">
            <span>Description *</span>
            <textarea
              rows={3}
              required
              value={form.description}
              onChange={(e) => setForm((f) => ({ ...f, description: e.target.value }))}
              placeholder="Short synopsis"
            />
          </label>

          <div className="form-row">
            <label className="field">
              <span>Release year *</span>
              <input
                type="number"
                min={1888}
                max={2100}
                required
                value={form.release_year}
                onChange={(e) => setForm((f) => ({ ...f, release_year: e.target.value }))}
              />
            </label>
            <label className="field grow">
              <span>Poster URL *</span>
              <input
                type="url"
                required
                value={form.poster}
                onChange={(e) => setForm((f) => ({ ...f, poster: e.target.value }))}
                placeholder="https://…"
              />
            </label>
          </div>

          <div className="field">
            <span>Genres * ({form.genres.length} selected)</span>
            <div className="genre-picker">
              {GENRES.map((g) => (
                <button
                  key={g}
                  type="button"
                  className={`tag tag-btn ${form.genres.includes(g) ? "active" : ""}`}
                  onClick={() => toggleGenre(g)}
                >
                  {g}
                </button>
              ))}
            </div>
          </div>

          <p className="muted small">
            Note: cast & crew editing isn't supported by the backend PATCH endpoint yet — only
            title, description, year and poster can be updated.
          </p>

          <ErrorInline message={formError} />
          <div className="modal-actions">
            <button type="button" className="btn btn-ghost" onClick={() => setEditing(null)}>
              Cancel
            </button>
            <button type="submit" className="btn btn-primary" disabled={busy}>
              {busy ? "Saving…" : editing === "new" ? "Create movie" : "Save changes"}
            </button>
          </div>
        </form>
      </Modal>

      <Modal open={Boolean(confirmDelete)} title="Delete movie?" onClose={() => setConfirmDelete(null)}>
        <p>
          Permanently delete <strong>“{confirmDelete?.title}”</strong> ({confirmDelete?.release_year})
          from the catalog? This cannot be undone.
        </p>
        <div className="modal-actions">
          <button type="button" className="btn btn-ghost" onClick={() => setConfirmDelete(null)}>
            Cancel
          </button>
          <button type="button" className="btn btn-danger" disabled={deleteBusy} onClick={handleDelete}>
            {deleteBusy ? "Deleting…" : "Delete"}
          </button>
        </div>
      </Modal>
    </div>
  );
}
