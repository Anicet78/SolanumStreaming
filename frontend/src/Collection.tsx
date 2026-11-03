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
    <main class="page">
      <div class="page-heading"><div><div class="eyebrow">Your library</div><h1 class="page-title">Collection</h1><p class="page-subtitle">The movies you saved, sorted by your latest activity.</p></div><div class="search-wrap"><SearchBar /></div></div>
        <div class="movie-grid">
          <For each={sortedMovies()}>
            {(movie) => (
              <A href={`/collection/${movie.movie_id}`}>
                <CollectionMovieCard {...movie} />
              </A>
            )}
          </For>
        </div>
    </main>
  );
};

export default Collection
