import client from "./client";

export async function getUserPlaylists(userId) {
  const { data } = await client.get(`/list/user/${userId}`);
  return Array.isArray(data) ? data : [];
}

export async function getPlaylist(id) {
  const { data } = await client.get(`/list/${id}`);
  return data;
}

export async function createPlaylist({ userId, name, description, isPublic }) {
  const { data } = await client.post("/list/", {
    user_id: userId,
    playlist_name: name,
    description: description || undefined,
    visiblity: Boolean(isPublic), // matches backend field spelling
  });
  return data;
}

export async function deletePlaylist(id) {
  await client.delete(`/list/${id}`);
}

export async function addMovieToPlaylist(playlistId, movieId) {
  const { data } = await client.post(`/list/${playlistId}/movies/${movieId}`);
  return data;
}

export async function removeMovieFromPlaylist(playlistId, movieId) {
  await client.delete(`/list/${playlistId}/movies/${movieId}`);
}
