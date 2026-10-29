import { createResource, For, Show } from 'solid-js';
import { moviesApi } from './api/movies.ts';
import SearchBar from './components/SearchBar.tsx';
import { A } from '@solidjs/router';
import MovieCard from './components/MovieCard.tsx';

const Home = () => {
  const [popular] = createResource(moviesApi.getPopular)

  return (
    <main class="page">
      <div class="page-heading">
        <div><div class="eyebrow">Discover</div><h1 class="page-title">Find your next story.</h1><p class="page-subtitle">A hand-picked stream of popular movies, ready for your watchlist.</p></div>
        <div class="search-wrap"><SearchBar/></div>
      </div>
      <Show when={popular()}>
          <div class="movie-grid">
            <For each={popular()?.results}>
              {(movie, idx) =>
                <Show when={idx() < 18}>
                  <A href={`/${movie.id}`}>
                    <MovieCard {...movie} />
                  </A>
                </Show>
              }
            </For>
          </div>
      </Show>
    </main>
  );
}

export default Home
