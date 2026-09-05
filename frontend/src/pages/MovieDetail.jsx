import { useCallback, useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { getMovie, addReview } from "../api/movies";
import { getUserPlaylists, addMovieToPlaylist } from "../api/playlists";
import { useAuth } from "../auth/AuthContext";
import { apiError } from "../api/client";
import { ErrorBox, Spinner, Stars } from "../components/ui";
import SaveToPlaylistModal from "../components/SaveToPlaylistModal";
import { ErrorInline } from "../components/SaveToPlaylistModal";

export default function MovieDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { isAuthenticated, user } = useAuth();

  const [movie, setMovie] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState("");

  const [saveOpen, setSaveOpen] = useState(false);
  const [playlists, setPlaylists] = useState([]);
  const [savedIn, setSavedIn] = useState([]);

  const [score, setScore] = useState("8");
  const [comment, setComment] = useState("");
  const [reviewBusy, setReviewBusy] = useState(false);
  const [reviewError, setReviewError] = useState("");
  const [reviewDone, setReviewDone] = useState(false);

  const load = useCallback(async () => {
    setLoading(true);
    setError("");
    try {
      setMovie(await getMovie(id));
    } catch (e) {
      setError(apiError(e));
    } finally {
      setLoading(false);
    }
  }, [id]);

  useEffect(() => {
    load();
  }, [load]);

  useEffect(() => {
    async function checkSaved() {
      if (!isAuthenticated || !user?.id || !movie?.id) return;
      try {
        const pls = await getUserPlaylists(user.id);
        setPlaylists(pls);
        setSavedIn(pls.filter((pl) => (pl.movie_ids || []).some((mid) => String(mid) === String(movie.id))));
      } catch {
        /* non-critical */
      }
    }
    checkSaved();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [isAuthenticated, user?.id, movie?.id]);

  async function handleQuickAdd() {
    if (!isAuthenticated) return navigate("/login");
    if (playlists.length === 0) {
      setSaveOpen(true); // modal handles "create first playlist" flow
      return;
    }
    try {
      await addMovieToPlaylist(playlists[0].id, movie.id);
      setSavedIn([{ id: playlists[0].id, name: playlists[0].name }]);
    } catch (e) {
      setSaveOpen(true);
    }
  }

  async function handleReview(e) {
    e.preventDefault();
    setReviewError("");
    setReviewBusy(true);
    try {
      await addReview(id, { userId: user.id, comment: comment.trim(), score: Number(score) });
      setComment("");
      setReviewDone(true);
      await load(); // refresh rating + reviews
    } catch (err) {
      setReviewError(apiError(err));
    } finally {
      setReviewBusy(false);
    }
  }

  if (loading && !movie) return <Spinner label="Loading movie…" />;
  if (error)
    return (
      <div className="page">
        <ErrorBox message={error} onRetry={load} />
        <p className="center">
          <Link to="/" className="btn btn-ghost">
            ← Back to movies
          </Link>
        </p>
      </div>
    );
  if (!movie) return null;

  const alreadySaved = savedIn.length > 0;

  return (
    <div className="page detail-page">
      <button type="button" className="btn btn-ghost btn-sm back-btn" onClick={() => navigate(-1)}>
        ← Back
      </button>

      <div className="detail-hero">
        <div className="detail-poster">
          {movie.poster ? (
            <img src={movie.poster} alt={`${movie.title} poster`} />
          ) : (
            <div className="poster-fallback big">🎬</div>
          )}
        </div>

        <div className="detail-info">
          <h1>{movie.title}</h1>
          <div className="detail-meta">
            <span className="tag">{movie.release_year}</span>
            {(movie.genres || []).map((g) => (
              <span key={g} className="tag">
                {g}
              </span>
            ))}
          </div>

          <Stars average={movie.rating?.average} count={movie.rating?.count} />

          <p className="detail-desc">{movie.description}</p>

          {isAuthenticated ? (
            <div className="detail-actions">
              <button
                type="button"
                className={`btn ${alreadySaved ? "btn-saved" : "btn-primary"}`}
                onClick={alreadySaved ? () => setSaveOpen(true) : handleQuickAdd}
                disabled={alreadySaved}
                title={alreadySaved ? `In playlist: ${savedIn.map((p) => p.name).join(", ")}` : "Save to your playlist"}
              >
                {alreadySaved ? `✓ Saved to “${savedIn[0]?.name || "playlist"}”` : "＋ Save to playlist"}
              </button>
              {!alreadySaved && (
                <button type="button" className="btn btn-outline" onClick={() => setSaveOpen(true)}>
                  Choose playlist…
                </button>
              )}
            </div>
          ) : (
            <p className="muted">
              <Link to="/login">Log in</Link> to save this movie and leave a review.
            </p>
          )}

          {(movie.directors?.length > 0 || movie.actors?.length > 0) && (
            <div className="cast-grid">
              {movie.directors?.length > 0 && (
                <div>
                  <h4>Director{movie.directors.length > 1 ? "s" : ""}</h4>
                  <ul>
                    {movie.directors.map((d, i) => (
                      <li key={i}>{d.name}</li>
                    ))}
                  </ul>
                </div>
              )}
              {movie.actors?.length > 0 && (
                <div>
                  <h4>Cast</h4>
                  <ul>
                    {movie.actors.slice(0, 8).map((a, i) => (
                      <li key={i}>
                        {a.name}
                        {a.role ? <em className="muted"> — {a.role}</em> : null}
                      </li>
                    ))}
                  </ul>
                </div>
              )}
            </div>
          )}
        </div>
      </div>

      <section className="reviews-section">
        <h2>Reviews {movie.reviews?.length ? `(${movie.reviews.length})` : ""}</h2>

        {isAuthenticated && (
          <form className="card review-form" onSubmit={handleReview}>
            <div className="review-form-row">
              <label className="field score-field">
                <span>Your score</span>
                <select value={score} onChange={(e) => setScore(e.target.value)}>
                  {Array.from({ length: 10 }, (_, i) => 10 - i).map((n) => (
                    <option key={n} value={n}>
                      {n}/10
                    </option>
                  ))}
                </select>
              </label>
              <label className="field grow">
                <span>Comment (optional)</span>
                <input
                  value={comment}
                  onChange={(e) => setComment(e.target.value)}
                  placeholder="What did you think?"
                  maxLength={500}
                />
              </label>
            </div>
            <ErrorInline message={reviewError} />
            {reviewDone && !reviewError && <div className="toast-success">Review submitted!</div>}
            <button type="submit" className="btn btn-primary" disabled={reviewBusy}>
              {reviewBusy ? "Submitting…" : "Submit review"}
            </button>
          </form>
        )}

        {movie.reviews?.length ? (
          <ul className="review-list">
            {[...movie.reviews].reverse().map((r) => (
              <li key={String(r.id)} className="card review-item">
                <div className="review-head">
                  <strong>{r.user_id === user?.id ? "You" : r.user_id}</strong>
                  <span className="review-score">{r.score}/10</span>
                  <time className="muted">{new Date(r.created_at).toLocaleDateString()}</time>
                </div>
                {r.comment && <p>{r.comment}</p>}
              </li>
            ))}
          </ul>
        ) : (
          <p className="muted">No reviews yet. Be the first!</p>
        )}
      </section>

      <SaveToPlaylistModal movie={movie} open={saveOpen} onClose={() => setSaveOpen(false)} />
    </div>
  );
}
