import { A, action, useAction, useParams, useSubmission } from "@solidjs/router";
import { moviesApi } from "./api/movies";
import { createEffect, createResource, createSignal, Match, Show, Switch } from "solid-js";
import SearchBar from "./components/SearchBar";

export const addToCollectionAction = action(async (movie_id: number) => {
  return moviesApi.addToCollection({ movie_id });
});

export const removeFromCollectionAction = action(async (movie_id: number) => {
  return moviesApi.removeFromCollection(movie_id);
});

const Movie = () => {
  const params = useParams()

  if (!params.movie_id)
    return <span>Not Found</span>
  const id: number = Number(params.movie_id)

  const [movieInCollection] = createResource(() => id, moviesApi.getInCollection);
  const [movie] = createResource(id, moviesApi.getMovieInfos)

  const [inCollection, setInCollection] = createSignal(false);

  createEffect(() => {
    if (movieInCollection.state === "ready") setInCollection(true);
    else if (movieInCollection.state === "errored") setInCollection(false);
  });

  const addToCollection = useAction(addToCollectionAction);
  const addSubmission = useSubmission(addToCollectionAction);

  const removeFromCollection = useAction(removeFromCollectionAction);
  const removeSubmission = useSubmission(removeFromCollectionAction);

  const handleAddToCollection = async () => {
    await addToCollection(id);
    if (!addSubmission.error) setInCollection(true)
  };

  const handleRemoveFromCollection = async () => {
    await removeFromCollection(id);
    if (!removeSubmission.error) setInCollection(false)
  };

  return (
    <Show when={movie()}>
      {(movie) => (
        <main class="detail-page">
          <div class="detail-backbar">
            <A class="detail-back" href="/">← Back to discover</A>
            <div class="detail-search"><SearchBar /></div>
          </div>
          <section class="detail-hero" style={{ "background-image": `url(https://image.tmdb.org/t/p/original${movie().backdrop_path})` }}>
            <div class="detail-content">
              <div class="eyebrow">Movie details</div>
              <h1 class="detail-title">{movie().title}</h1>
              <Show when={movie().title !== movie().original_title}>
                <div class="detail-original">Original title: {movie().original_title}</div>
              </Show>
              <div class="detail-stats">
                <span class="detail-stat">{movie().release_date.slice(0, 4)}</span>
                <span class="detail-stat">{Math.floor(movie().runtime / 60)}h {movie().runtime % 60}m</span>
                <span class="detail-stat">★ {movie().vote_average.toFixed(1)} / 10</span>
              </div>
              <Show when={movie().tagline || movie().overview}>
                <Show when={movie().tagline}><div class="detail-tagline">{movie().tagline}</div></Show>
                <Show when={movie().overview}><p class="detail-overview">{movie().overview}</p></Show>
              </Show>
              <div class="detail-actions">
              <Switch fallback={<p>Something went wrong.</p>}>
                <Match when={movieInCollection.loading || addSubmission.pending || removeSubmission.pending}>
                  <button class="action-button" disabled>
                    <span class="loading loading-spinner loading-sm"></span>
                  </button>
                </Match>
                <Match when={inCollection()}>
                  <button class="action-button secondary" onClick={handleRemoveFromCollection}>Remove from collection</button>
                  <A class="action-button" href={`/collection/${movie().id}`}>Watch from collection</A>
                </Match>
                <Match when={!inCollection()}>
                  <button class="action-button" onClick={handleAddToCollection}>+ Add to collection</button>
                </Match>
              </Switch>
              </div>
            </div>
          </section>
        </main>
      )}
    </Show>
  );
}

export default Movie
