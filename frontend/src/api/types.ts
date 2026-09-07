export interface SearchResponse {
  page: number;
  results: Movie[];
  total_pages: number;
  total_results: number;
}

export interface Movie {
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