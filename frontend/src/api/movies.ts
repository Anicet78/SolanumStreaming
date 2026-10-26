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
}

export const moviesApi = {
  search: (filters: { title?: string; page?: number }) => {
    const params = new URLSearchParams();
    if (filters.title) params.set("title", filters.title);
    if (filters.page) params.set("page", String(filters.page));

    return api.get<SearchResponse>(`:8082/search?${params.toString()}`);
  },
  getMovieInfos: (movieId: number) => api.get<MovieT>(`:8082/${movieId}`),
  searchTorrent: (movieId: number, caps: string) => api.get<CollectionMovieT>(`:8082/${movieId}/torrent?caps=${caps.toString()}`),
  getCollection: () => api.get<CollectionMovieT[]>(":8082/collection"),
  getInCollection: (movieId: number) => api.get<CollectionMovieT>(`:8082/collection/${movieId}`),
  addToCollection: (data: { movie_id: number }) => api.post<CollectionMovieT>(`:8082/collection`, data),
  removeFromCollection: (movieId: number) => api.delete<void>(`:8082/collection/${movieId}`),
};
