import { createResource, For, Show } from 'solid-js';
import { moviesApi } from './api/movies.ts';
import SearchBar from './components/SearchBar.tsx';
import { A } from '@solidjs/router';
import MovieCard from './components/MovieCard.tsx';

const Home = () => {
  const [popular] = createResource(moviesApi.getPopular)

  return (
    <>
      <div class="flex items-center justify-center pt-4">
        <SearchBar/>
      </div>
      <div class="flex flex-col min-h-screen items-center justify-center">
        <Show when={popular()}>
          <div class="grid grid-cols-2 sm:grid-cols-3 md:grid-cols-4 lg:grid-cols-5 xl:grid-cols-6 gap-4 bg-base-100 rounded-box shadow-md p-4">
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
      </div>
    </>
  );
}

export default Home
