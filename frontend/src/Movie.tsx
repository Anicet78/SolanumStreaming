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
        <main
          class="relative min-h-screen bg-cover bg-center flex flex-col justify-between "
          style={{
            "background-image": `linear-gradient(to bottom, rgb(0 0 0 / 10%), rgb(0 0 0 / 90%)), url(https://image.tmdb.org/t/p/original${movie()?.backdrop_path})`
          }}
        >
          <div class="flex items-center justify-center pt-4">
            <SearchBar/>
          </div>
          <div class="w-full p-10">
            <div class="flex items-center justify-evenly">
              <div class="flex flex-col w-fit">
                <span class="font-black text-2xl">{movie().title}</span>
                <Show when={movie().title !== movie().original_title}>
                  <span class="self-center text-sm text-gray-300">{movie().original_title}</span>
                </Show>
              </div>
              <div class="flex flex-col items-center p-2">
                <span class="text-2xl">🗓️ </span>
                <span class="font-black text-2xl">{movie().release_date.slice(0, 4)}</span>
              </div>
              <div class="flex flex-col items-center p-2">
                <span class="text-2xl">⏳ </span>
                <span class="font-black text-2xl">{Math.floor(movie().runtime / 60)}h {movie().runtime % 60}m</span>
              </div>
              <div class="flex flex-col items-center p-2">
                <span class="text-2xl">⭐ </span>
                <span class="font-black text-2xl">{movie().vote_average}</span>
              </div>
            </div>
            <Show when={movie().overview !== "" || movie().tagline !== ""}>
              <div class="mt-5 flex flex-col">
                <span class="text-gray-400">{movie().tagline.toUpperCase()}</span>
                <Show when={movie().overview !== "" && movie().tagline !== ""}><div class="mb-1"/></Show>
                <span>{movie().overview}</span>
              </div>
            </Show>
            <div class="mt-5">
              <Switch fallback={<p>An error has occured</p>}>
                <Match when={movieInCollection.loading || addSubmission.pending || removeSubmission.pending}>
                  <button class="btn btn-outline btn-primary w-full p-0">
                    <span class="loading loading-spinner loading-sm"></span>
                  </button>
                </Match>
                <Match when={inCollection()}>
                  <div class="flex justify-around">
                    <button class="btn btn-outline btn-primary w-1/2" onClick={handleRemoveFromCollection}>Remove from collection</button>
                    <A class="btn btn-outline btn-primary w-1/2" href={`/collection/${movie().id}`}>See in collection</A>
                  </div>
                </Match>
                <Match when={!inCollection()}>
                  <button class="btn btn-outline btn-primary w-full p-0" onClick={handleAddToCollection}>Add to collection</button>
                </Match>
              </Switch>
            </div>
          </div>
        </main>
      )}
    </Show>
  );
}

export default Movie
