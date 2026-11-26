import { api } from "./client";

export interface SearchResponse {
  page: number;
  results: MovieT[];
  total_pages: number;
  total_results: number;
}

export interface MovieT {
  id: number;
  title: string;
  original_title: string;
  popularity: number;
  adult: boolean;
  video: boolean;
  poster_path: string;
  backdrop_path: string;
  release_date: string;
  imdb_id: string;
  runtime: number;
  tagline: string;
  overview: string;
  vote_average: number;
}

export interface CollectionMovieT {
  movie_id: number,
  torrent_link: string
  length: number
  progression: number
  added_on: string
  last_watched: string
}

const BASE = import.meta.env.VITE_MOVIES_PATH;

export const moviesApi = {
  search: (filters: { title?: string; page?: number }) => {
    const params = new URLSearchParams();
    if (filters.title) params.set("title", filters.title);
    if (filters.page) params.set("page", String(filters.page));

    return api.get<SearchResponse>(`${BASE}/search?${params.toString()}`);
  },
  getNowPlaying: () => api.get<SearchResponse>(`${BASE}/list/now_playing`),
  getPopular: () => api.get<SearchResponse>(`${BASE}/list/popular`),
  getMovieInfos: (movieId: number) => api.get<MovieT>(`${BASE}/${movieId}`),
  searchTorrent: (movieId: number, caps: string) => api.get<CollectionMovieT>(`${BASE}/${movieId}/torrent?caps=${caps.toString()}`),
  getCollection: () => api.get<CollectionMovieT[]>(`${BASE}/collection`),
  getInCollection: (movieId: number) => api.get<CollectionMovieT>(`${BASE}/collection/${movieId}`),
  addToCollection: (data: { movie_id: number }) => api.post<CollectionMovieT>(`${BASE}/collection`, data),
  removeFromCollection: (movieId: number) => api.delete<void>(`${BASE}/collection/${movieId}`),
};
