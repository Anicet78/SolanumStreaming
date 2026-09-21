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
  release_date: string;
  imdb_id: string;
  runtime: number;
}

export const moviesApi = {
  search: (filters: { title?: string; page?: number }) => {
    const params = new URLSearchParams();
    if (filters.title) params.set("title", filters.title);
    if (filters.page) params.set("page", String(filters.page));

    return api.get<SearchResponse>(`:8082/search?${params.toString()}`);
  },
};
