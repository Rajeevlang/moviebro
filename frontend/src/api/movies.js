import client from "./client";

export async function getMovies(params = {}) {
  const clean = Object.fromEntries(
    Object.entries(params).filter(([, v]) => v !== "" && v !== null && v !== undefined)
  );
  const { data } = await client.get("/movies/", { params: clean });
  return Array.isArray(data) ? data : [];
}

export async function getMovie(id) {
  const { data } = await client.get(`/movies/${id}`);
  return data;
}

export async function createMovie(payload) {
  const { data } = await client.post("/movies/", payload);
  return data;
}

export async function updateMovie(id, payload) {
  const { data } = await client.patch(`/movies/${id}`, payload);
  return data;
}

export async function deleteMovie(id) {
  await client.delete(`/movies/${id}`);
}

export async function addReview(movieId, { userId, comment, score }) {
  const { data } = await client.post(`/movies/${movieId}/reviews`, {
    user_id: userId,
    comment,
    score,
  });
  return data;
}
