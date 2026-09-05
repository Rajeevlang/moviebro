export const API_BASE = import.meta.env.VITE_API_BASE || "";

// Emails allowed into the admin panel.
// NOTE: the backend has no role system yet — this only gates the UI.
// Real enforcement must be added server-side before production use.
export const ADMIN_EMAILS = (import.meta.env.VITE_ADMIN_EMAILS || "admin@cinedeck.dev")
  .split(",")
  .map((e) => e.trim().toLowerCase());

export const TOKEN_KEY = "cinedeck_token";
export const USER_KEY = "cinedeck_user";

export const GENRES = [
  "Action",
  "Adventure",
  "Animation",
  "Comedy",
  "Crime",
  "Documentary",
  "Drama",
  "Family",
  "Fantasy",
  "History",
  "Horror",
  "Music",
  "Mystery",
  "Romance",
  "Sci-Fi",
  "Thriller",
  "War",
  "Western",
];
