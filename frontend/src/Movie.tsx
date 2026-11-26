import { A, action, useAction, useParams, useSubmission } from "@solidjs/router";
import { moviesApi } from "./api/movies";
import { createEffect, createResource, createSignal, Match, Show, Switch } from "solid-js";
import DetailPage from "./components/DetailPage";
import DetailBackbar from "./components/DetailBackbar";
import MovieDetails from "./components/MovieDetails";

export const addToCollectionAction = action(async (movie_id: number) => {
  return moviesApi.addToCollection({ movie_id });
});

export const removeFromCollectionAction = action(async (movie_id: number) => {
  return moviesApi.removeFromCollection(movie_id);
});

const actionButton =
  "inline-flex min-h-10.5 cursor-pointer items-center justify-center rounded-[0.5625rem] border px-4 text-[0.85rem] font-bold hover:brightness-[1.08] narrow:flex-[1_1_100%] pointer-coarse:min-h-12";
const actionPrimary = "border-brand bg-brand text-ink";
const actionSecondary = "border-[#496273] bg-[rgb(11_17_24/65%)] text-[#d7e2eb]";

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
        <DetailPage>
          <DetailBackbar href="/" label="← Back to discover" />
          <section
            class="relative flex min-h-[min(40.625rem,calc(100vh_-_8.125rem))] items-end overflow-hidden rounded-[1.25rem] border border-line bg-cover bg-center shadow-[0_1.5rem_4.375rem_rgb(0_0_0/28%)] before:absolute before:inset-0 before:bg-[linear-gradient(90deg,rgb(6_12_18/96%)_0%,rgb(6_12_18/78%)_42%,rgb(6_12_18/18%)_100%),linear-gradient(0deg,rgb(6_12_18/96%)_0%,transparent_58%)] before:content-[''] max-lg:min-h-0 narrow:bg-[position:64%_center] narrow:before:bg-[linear-gradient(0deg,rgb(6_12_18/98%)_8%,rgb(6_12_18/50%)_64%,rgb(6_12_18/20%)_100%)]"
            style={{ "background-image": `url(https://image.tmdb.org/t/p/original${movie().backdrop_path})` }}
          >
            <div class="relative w-[min(47.5rem,100%)] p-[clamp(1.5rem,5vw,4rem)] narrow:px-5 narrow:py-6">
              <MovieDetails movie={movie()} eyebrow="Movie details" />
              <div class="mt-7 flex flex-wrap gap-2.5">
              <Switch fallback={<p>Something went wrong.</p>}>
                <Match when={movieInCollection.loading || addSubmission.pending || removeSubmission.pending}>
                  <button class={`${actionButton} ${actionPrimary}`} disabled>
                    <span class="loading loading-spinner loading-sm"></span>
                  </button>
                </Match>
                <Match when={inCollection()}>
                  <button class={`${actionButton} ${actionSecondary}`} onClick={handleRemoveFromCollection}>Remove from collection</button>
                  <A class={`${actionButton} ${actionPrimary}`} href={`/collection/${movie().id}`}>Watch from collection</A>
                </Match>
                <Match when={!inCollection()}>
                  <button class={`${actionButton} ${actionPrimary}`} onClick={handleAddToCollection}>+ Add to collection</button>
                </Match>
              </Switch>
              </div>
            </div>
          </section>
        </DetailPage>
      )}
    </Show>
  );
}

export default Movie
