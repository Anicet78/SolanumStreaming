import { createResource, For, Show } from "solid-js"
import MovieCard from "./components/MovieCard"
import { moviesApi } from "./api/movies"
import { A, useSearchParams } from "@solidjs/router";
import SearchBar from "./components/SearchBar";

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
    <>
      <div class="flex flex-col m-4 items-center justify-center">
        <SearchBar/>
      </div>
      <div class="min-h-screen flex flex-col items-center p-4">
        <Show when={results.loading}>
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4 bg-base-100 rounded-box shadow-md p-4">
            <For each={Array.from({ length: 18 })}>
              {() => <div class="skeleton w-60 aspect-2/3 rounded-2xl"></div>}
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
              {(movie) =>
                <A href={`/${movie.id}`}>
                  <MovieCard {...movie} />
                </A>
              }
            </For>
          </div>
          <Show when={results()?.total_pages! > 1}>
            <div class="join flex justify-end w-full mr-7">
              <For each={Array.from({ length: results()?.total_pages!})}>
                {(_, i) =>
                  <Show when={i() == 0 || i() + 1 == results()?.total_pages || (results()?.page! - i() < 4 && results()?.page! - i() >= 0) || (i() - results()?.page! < 4 && i() - results()?.page! >= 0)}>
                    <A href={`/search?title=${searchParams.title}&page=${i() + 1}`}>
                      <input
                        class="join-item btn btn-square"
                        type="radio"
                        name="options"
                        aria-label={(i() + 1).toString()}
                        checked={i() + 1 == results()?.page ? true : false}
                      />
                    </A>
                  </Show>
                }
              </For>
            </div>
          </Show>
        </Show>
      </div>
    </>
  );
};

export default Search
