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
    <div class="min-h-screen flex flex-col items-center p-4">
      <Show when={results.loading}>
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4 bg-base-100 rounded-box shadow-md p-4">
          <For each={Array.from({ length: 30 })}>
            {() => <div class="skeleton max-w-100 w-full aspect-5/7 rounded-2xl"></div>}
          </For>
        </div>
      </Show>

      <Show when={results.error}>
        <div class="flex flex-1 flex-col items-center justify-center">
          <span class="text-2xl">Cannot load your search</span>
          <span class="text-xl text-gray-300">{results.error?.message}</span>
        </div>
      </Show>

      <Show when={!results.loading && !results.error}>
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4 bg-base-100 rounded-box shadow-md p-4">
          <For each={results()?.results}>
            {(movie) => <MovieCard {...movie} />}
          </For>
        </div>
      </Show>
    </div>
  );
};

export default Search