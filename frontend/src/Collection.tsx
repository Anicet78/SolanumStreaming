import { moviesApi } from "./api/movies";
import { createMemo, createResource, For } from "solid-js";
import { A } from "@solidjs/router";
import CollectionMovieCard from "./components/CollectionMovieCard";
import SearchBar from "./components/SearchBar";

const Collection = () => {
  const [results] = createResource(moviesApi.getCollection);

  const sortedMovies = createMemo(() =>
    [...(results() ?? [])].sort((a, b) => new Date(b.last_watched).getTime() - new Date(a.last_watched).getTime())
  );

  return (
    <>
      <div class="flex items-center justify-center mt-4">
        <SearchBar />
      </div>
      <div class="min-h-screen flex flex-col items-center p-4">
        <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4 bg-base-100 rounded-box shadow-md p-4">
          <For each={sortedMovies()}>
            {(movie) => (
              <A href={`/collection/${movie.movie_id}`}>
                <CollectionMovieCard {...movie} />
              </A>
            )}
          </For>
        </div>
      </div>
    </>
  );
};

export default Collection
