import { Link } from "react-router-dom";
import { Stars } from "./ui";

export default function MovieCard({ movie, onSave, savedIn }) {
  return (
    <article className="movie-card">
      <Link to={`/movies/${movie.id}`} className="movie-poster-link">
        <div className="movie-poster">
          {movie.poster ? (
            <img src={movie.poster} alt={`${movie.title} poster`} loading="lazy" />
          ) : (
            <div className="poster-fallback">🎬</div>
          )}
          {movie.rating?.count > 0 && (
            <span className="rating-badge">{movie.rating.average.toFixed(1)}</span>
          )}
        </div>
      </Link>
      <div className="movie-card-body">
        <h3 className="movie-title">
          <Link to={`/movies/${movie.id}`}>{movie.title}</Link>
        </h3>
        <p className="movie-meta">
          <span>{movie.release_year}</span>
          <Stars average={movie.rating?.average} count={movie.rating?.count} />
        </p>
        <div className="genre-tags">
          {(movie.genres || []).slice(0, 3).map((g) => (
            <span key={g} className="tag">
              {g}
            </span>
          ))}
        </div>
        {onSave && (
          <button
            className={`btn btn-save ${savedIn ? "btn-saved" : ""}`}
            onClick={() => onSave(movie)}
            title={savedIn ? `In ${savedIn.length} playlist(s)` : "Save to playlist"}
          >
            {savedIn ? `✓ Saved (${savedIn.length})` : "+ Save"}
          </button>
        )}
      </div>
    </article>
  );
}
