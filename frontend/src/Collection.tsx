import { moviesApi } from "./api/movies";
import { createMemo, createResource, For } from "solid-js";
import { A } from "@solidjs/router";
import CollectionMovieCard from "./components/CollectionMovieCard";
import MovieGrid from "./components/MovieGrid";
import Page from "./components/Page";
import PageHeading from "./components/PageHeading";

const Collection = () => {
  const [results] = createResource(moviesApi.getCollection);

  const sortedMovies = createMemo(() =>
    [...(results() ?? [])].sort((a, b) => new Date(b.last_watched).getTime() - new Date(a.last_watched).getTime())
  );

  return (
    <Page>
      <PageHeading eyebrow="Your library" title="Collection" subtitle="The movies you saved, sorted by your latest activity." />
      <MovieGrid>
        <For each={sortedMovies()}>
          {(movie) => (
            <A href={`/collection/${movie.movie_id}`}>
              <CollectionMovieCard {...movie} />
            </A>
          )}
        </For>
      </MovieGrid>
    </Page>
  );
};

export default Collection
