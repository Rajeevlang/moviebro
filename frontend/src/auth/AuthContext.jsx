import { createContext, useCallback, useContext, useEffect, useMemo, useState } from "react";
import { TOKEN_KEY, USER_KEY, ADMIN_EMAILS } from "../config";
import { setUnauthorizedHandler } from "../api/client";

const AuthContext = createContext(null);

function decodeJwt(token) {
  try {
    const payload = token.split(".")[1];
    return JSON.parse(atob(payload.replace(/-/g, "+").replace(/_/g, "/")));
  } catch {
    return {};
  }
}

function isExpired(claims) {
  if (!claims?.exp) return false;
  return Date.now() >= claims.exp * 1000;
}

function loadUser(token) {
  const claims = decodeJwt(token);
  if (isExpired(claims)) return null;
  let email = claims.email || "";
  if (!email) {
    try {
      email = JSON.parse(localStorage.getItem(USER_KEY) || "{}").email || "";
    } catch {
      /* ignore */
    }
  }
  return {
    id: claims.id || claims.sub || "",
    username: claims.username || claims.name || "User",
    email,
  };
}

export function AuthProvider({ children }) {
  const [token, setToken] = useState(() => localStorage.getItem(TOKEN_KEY));
  const [user, setUser] = useState(() => {
    const stored = localStorage.getItem(TOKEN_KEY);
    if (!stored) return null;
    const loaded = loadUser(stored);
    if (!loaded) localStorage.removeItem(TOKEN_KEY);
    return loaded;
  });

  const logout = useCallback(() => {
    localStorage.removeItem(TOKEN_KEY);
    setToken(null);
    setUser(null);
  }, []);

  useEffect(() => {
    setUnauthorizedHandler(logout);
    return () => setUnauthorizedHandler(null);
  }, [logout]);

  /**
   * Store session from a JWT returned by /login.
   * `profile` may carry { email } captured from the login form
   * (the backend token does not embed the email).
   */
  const loginWithToken = useCallback((jwtToken, profile = {}) => {
    localStorage.setItem(TOKEN_KEY, jwtToken);
    if (profile.email) {
      localStorage.setItem(USER_KEY, JSON.stringify({ email: profile.email }));
    }
    const claims = decodeJwt(jwtToken);
    setToken(jwtToken);
    setUser({
      id: claims.id || claims.sub || "",
      username: profile.username || claims.username || claims.name || "User",
      email: profile.email || claims.email || "",
    });
  }, []);

  const value = useMemo(() => {
    const isAdmin = Boolean(user) && ADMIN_EMAILS.includes((user?.email || "").toLowerCase());
    const authed = Boolean(token) && Boolean(user);
    return { user, token, isAdmin, isAuthenticated: authed, loginWithToken, logout };
  }, [token, user, loginWithToken, logout]);

  return <AuthContext.Provider value={value}>{children}</AuthContext.Provider>;
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error("useAuth must be used inside <AuthProvider>");
  return ctx;
}
