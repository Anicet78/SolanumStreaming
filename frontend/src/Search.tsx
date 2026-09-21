import { createResource, For, Show } from "solid-js"
import MovieCard from "./components/MovieCard"
import { moviesApi } from "./api/movies"
import { useSearchParams } from "@solidjs/router";

function toSingleString(v: string | string[] | undefined): string | undefined {
  return Array.isArray(v) ? v[0] : v;
}

const Search = () => {
  const [searchParams] = useSearchParams();

  const [results] = createResource(
    () => ({
      title: toSingleString(searchParams.title),
      page: toSingleString(searchParams.page) ? Number(toSingleString(searchParams.page)) : undefined,
    }),
    moviesApi.search
  );

  return (
    <div class="flex align-items justify-center h-screen flex-col">
      <Show when={results.loading}>
        <For each={Array.from({ length: 30 })}>
          {(_, i) => <div class="skeleton w-50 h-70 rounded-2xl"></div>}
        </For>
      </Show>

      <Show when={results.error}>
        <span class="text-2xl mb-0 m-auto">
          Cannot load your search
        </span>
        <span class="text-xl mt-0 m-auto text-gray-300">
          {results.error?.message}
        </span>
      </Show>

      <Show when={!results.loading && !results.error}>
        <ul class="list bg-base-100 rounded-box shadow-md">
          <For each={results()?.results}>
            {(movie) => <li><MovieCard {...movie} /></li>}
          </For>
        </ul>
      </Show>
    </div>
  );
};

export default Search