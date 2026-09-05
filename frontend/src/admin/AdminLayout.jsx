import { NavLink, Outlet } from "react-router-dom";

export default function AdminLayout() {
  return (
    <div className="page admin-page">
      <div className="admin-shell">
        <aside className="admin-sidebar">
          <h2>Admin</h2>
          <nav>
            <NavLink to="/admin" end>
              📊 Dashboard
            </NavLink>
            <NavLink to="/admin/movies">🎬 Movies</NavLink>
          </nav>
          <p className="muted small">
            UI-level gating only. Add server-side role checks before exposing this publicly.
          </p>
        </aside>
        <main className="admin-content">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
