export function Spinner({ label = "Loading…" }) {
  return (
    <div className="spinner-wrap">
      <div className="spinner" />
      <p>{label}</p>
    </div>
  );
}

export function ErrorBox({ message, onRetry }) {
  if (!message) return null;
  return (
    <div className="error-box">
      <span>⚠️ {message}</span>
      {onRetry && (
        <button className="btn btn-ghost btn-sm" onClick={onRetry}>
          Retry
        </button>
      )}
    </div>
  );
}

export function EmptyState({ icon = "🎬", title, hint }) {
  return (
    <div className="empty-state">
      <div className="empty-icon">{icon}</div>
      <h2>{title}</h2>
      {hint && <p>{hint}</p>}
    </div>
  );
}

export function Stars({ average, count }) {
  const value = Number(average) || 0;
  const full = Math.round(value / 2); // score is out of 10
  return (
    <span className="stars" title={`${value.toFixed(1)} / 10 (${count} votes)`}>
      <span className="stars-icons" aria-hidden="true">
        {"★".repeat(full)}
        {"☆".repeat(Math.max(0, 5 - full))}
      </span>
      <span className="stars-num">{value ? value.toFixed(1) : "—"}</span>
      {count ? <span className="stars-count">({count})</span> : null}
    </span>
  );
}

export function Modal({ open, title, onClose, children }) {
  if (!open) return null;
  return (
    <div className="modal-backdrop" onMouseDown={(e) => e.target === e.currentTarget && onClose()}>
      <div className="modal" role="dialog" aria-modal="true" aria-label={title}>
        <div className="modal-head">
          <h3>{title}</h3>
          <button type="button" className="modal-close" aria-label="Close" onClick={onClose}>
            ✕
          </button>
        </div>
        <div className="modal-body">{children}</div>
      </div>
    </div>
  );
}
