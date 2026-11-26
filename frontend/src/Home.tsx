import { createResource, For, Show } from 'solid-js';
import { moviesApi } from './api/movies.ts';
import { A } from '@solidjs/router';
import MovieCard from './components/MovieCard';
import MovieGrid from './components/MovieGrid';
import Page from './components/Page';
import PageHeading from './components/PageHeading';

const Home = () => {
  const [popular] = createResource(moviesApi.getPopular)

  return (
    <Page>
      <PageHeading eyebrow="Discover" title="Explore the movie catalog" subtitle="Trending and popular movies" />
      <Show when={popular()}>
        <MovieGrid>
          <For each={popular()?.results}>
            {(movie) =>
              <A href={`/${movie.id}`}>
                <MovieCard {...movie} />
              </A>
            }
          </For>
        </MovieGrid>
      </Show>
    </Page>
  );
}

export default Home
