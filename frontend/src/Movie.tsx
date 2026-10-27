import { action, useAction, useParams, useSubmission } from "@solidjs/router";
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
  const [movieInfos] = createResource(id, moviesApi.getMovieInfos)

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
    <Show when={movieInfos()}>
      <main
        class="relative min-h-screen bg-cover bg-center"
        style={{
          "background-image": `linear-gradient(to bottom, rgb(0 0 0 / 10%), rgb(0 0 0 / 90%)), url(https://image.tmdb.org/t/p/original${movieInfos()?.backdrop_path})`
        }}
      >
        <div class="flex items-center justify-center pt-4">
          <SearchBar/>
        </div>
        <div class="flex flex-col justify-center items-center mt-50">
          <div class="flex flex-col self-end p-5 rounded-md mr-20 backdrop-blur-xs bg-gray-500/15">
            <div class="flex flex-col justify-center w-md">
              <div class="flex flex-col items-center pt-2">
                <span class="font-black text-2xl">{movieInfos()!.title}</span>
                <Show when={movieInfos()!.title !== movieInfos()!.original_title}>
                  <span class="text-sm text-gray-400">{movieInfos()!.original_title}</span>
                </Show>
              </div>
              <div class="flex flex-col items-center pt-2">
                <span class="text-2xl">🗓️ </span>
                <span class="font-black text-2xl">{movieInfos()!.release_date.slice(0, 4)}</span>
              </div>
              <div class="flex flex-col items-center pt-2">
                <span class="text-2xl">⏳ </span>
                <span class="font-black text-2xl">{Math.floor(movieInfos()!.runtime / 60)}h {movieInfos()!.runtime % 60}m</span>
              </div>
              <div class="flex flex-col items-center pt-2 pb-2">
                <span class="text-2xl">⭐ </span>
                <span class="font-black text-2xl">{movieInfos()!.vote_average}</span>
              </div>
            </div>
            <Switch fallback={<p>An error has occured</p>}>
              <Match when={movieInCollection.loading || addSubmission.pending || removeSubmission.pending}>
                <button class="btn btn-soft btn-primary w-full p-0">
                  <span class="loading loading-spinner loading-sm"></span>
                </button>
              </Match>
              <Match when={inCollection()}>
                <button class="btn btn-soft btn-primary w-full p-0" onClick={handleRemoveFromCollection}>Remove from collection</button>
              </Match>
              <Match when={!inCollection()}>
                <button class="btn btn-soft btn-primary w-full p-0" onClick={handleAddToCollection}>Add to collection</button>
              </Match>
            </Switch>
          </div>
        </div>
      </main>
    </Show>
  );
}

export default Movie
