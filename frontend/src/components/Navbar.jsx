import { useState } from "react";
import { Link, NavLink, useNavigate } from "react-router-dom";
import { useAuth } from "../auth/AuthContext";

export default function Navbar() {
  const { user, isAuthenticated, isAdmin, logout } = useAuth();
  const navigate = useNavigate();
  const [open, setOpen] = useState(false);

  function handleLogout() {
    logout();
    setOpen(false);
    navigate("/");
  }

  return (
    <header className="navbar">
      <div className="container navbar-inner">
        <Link to="/" className="brand" onClick={() => setOpen(false)}>
          <span className="brand-mark">▶</span> CineDeck
        </Link>

        <button
          type="button"
          className="nav-toggle"
          aria-label="Toggle menu"
          onClick={() => setOpen((v) => !v)}
        >
          ☰
        </button>

        <nav className={`nav-links ${open ? "open" : ""}`}>
          <NavLink to="/" end onClick={() => setOpen(false)}>
            Movies
          </NavLink>
          {isAuthenticated && (
            <NavLink to="/playlists" onClick={() => setOpen(false)}>
              My Playlists
            </NavLink>
          )}
          {isAuthenticated && (
            <NavLink to="/profile" onClick={() => setOpen(false)}>
              Profile
            </NavLink>
          )}
          {isAdmin && (
            <NavLink to="/admin" onClick={() => setOpen(false)}>
              Admin
            </NavLink>
          )}

          <div className="nav-spacer" />

          {isAuthenticated ? (
            <>
              <Link to="/profile" className="user-chip" title="Open profile">
                <span className="avatar">{(user.username || "U")[0].toUpperCase()}</span>
                <span className="username">{user.username}</span>
              </Link>
              <button className="btn btn-ghost btn-sm" onClick={handleLogout}>
                Logout
              </button>
            </>
          ) : (
            <>
              <Link className="btn btn-ghost btn-sm" to="/login" onClick={() => setOpen(false)}>
                Login
              </Link>
              <Link className="btn btn-primary btn-sm" to="/register" onClick={() => setOpen(false)}>
                Sign up
              </Link>
            </>
          )}
        </nav>
      </div>
    </header>
  );
}
