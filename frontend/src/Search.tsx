import { createResource, For, Show } from "solid-js"
import MovieCard from "./components/MovieCard"
import { moviesApi } from "./api/movies"
import { A, useNavigate, useSearchParams } from "@solidjs/router";
import SearchBar from "./components/SearchBar";

function toSingleString(v: string | string[] | undefined): string | undefined {
  return Array.isArray(v) ? v[0] : v;
}

const Search = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  if (!searchParams.title || searchParams.title == "")
    navigate("/");

  const [results] = createResource(
    () => ({
      title: toSingleString(searchParams.title),
      page: toSingleString(searchParams.page) ? Number(toSingleString(searchParams.page)) : undefined,
    }),
    moviesApi.search
  );

  return (
    <main class="page">
        <div class="page-heading"><div><div class="eyebrow">Library search</div><h1 class="page-title">Search results</h1><p class="page-subtitle">Explore titles and find something worth pressing play on.</p></div><div class="search-wrap"><SearchBar/></div></div>
        <Show when={results.loading}>
          <div class="movie-grid loading-grid">
            <For each={Array.from({ length: 18 })}>
              {() => <div class="skeleton-card"></div>}
            </For>
          </div>
        </Show>

        <Show when={results.error || (results() && results()?.total_results == 0)}>
          <div class="empty-state">
            <strong>Nothing found</strong><br />
            <span>{results.error?.message ?? "No result has been found"}</span>
          </div>
        </Show>

        <Show when={!results.loading && !results.error}>
          <div class="movie-grid">
            <For each={results()?.results}>
              {(movie) =>
                <A href={`/${movie.id}`}>
                  <MovieCard {...movie} />
                </A>
              }
            </For>
          </div>
          <Show when={results()?.total_pages! > 1}>
            <div class="pagination">
              <For each={Array.from({ length: results()?.total_pages!})}>
                {(_, i) =>
                  <Show when={i() == 0 || i() + 1 == results()?.total_pages || (results()?.page! - i() < 4 && results()?.page! - i() >= 0) || (i() - results()?.page! < 4 && i() - results()?.page! >= 0)}>
                    <A classList={{ "pagination-current": i() + 1 === results()?.page }} href={`/search?title=${searchParams.title}&page=${i() + 1}`}>
                      {i() + 1}
                    </A>
                  </Show>
                }
              </For>
            </div>
          </Show>
        </Show>
    </main>
  );
};

export default Search
