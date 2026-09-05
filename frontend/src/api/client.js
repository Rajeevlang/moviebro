import axios from "axios";
import { API_BASE, TOKEN_KEY } from "../config";

const client = axios.create({
  baseURL: `${API_BASE}/api/v1`,
  timeout: 15000,
});

client.interceptors.request.use((config) => {
  const token = localStorage.getItem(TOKEN_KEY);
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

let onUnauthorized = null;
export function setUnauthorizedHandler(fn) {
  onUnauthorized = fn;
}

client.interceptors.response.use(
  (res) => res,
  (error) => {
    if (error.response?.status === 401 && onUnauthorized) {
      onUnauthorized();
    }
    return Promise.reject(error);
  }
);

/** Extract a readable message from any backend error shape. */
export function apiError(error) {
  if (error.code === "ECONNABORTED") return "Request timed out. Is the gateway running?";
  if (!error.response) return "Cannot reach the server. Is the gateway running?";
  const data = error.response.data;
  if (!data) return `Request failed (${error.response.status})`;
  if (typeof data.error === "string") return data.error;
  if (data.errors && typeof data.errors === "object") {
    return Object.entries(data.errors)
      .map(([field, msg]) => `${field}: ${msg}`)
      .join(" · ");
  }
  if (Array.isArray(data.errors)) return data.errors.join(" · ");
  if (typeof data.message === "string") return data.message;
  return `Request failed (${error.response.status})`;
}

export default client;
