import { action, useAction, useParams, useSubmission } from "@solidjs/router";
import { moviesApi } from "./api/movies";
import { createEffect, createResource, createSignal, Match, Show, Switch } from "solid-js";
import MovieCard from "./components/MovieCard";

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
        class="relative flex min-h-screen flex-col justify-center items-center bg-cover bg-center"
        style={{
          "background-image": `linear-gradient(to bottom, rgb(0 0 0 / 10%), rgb(0 0 0 / 90%)), url(https://image.tmdb.org/t/p/original${movieInfos()?.backdrop_path})`
        }}
      >
        <div class="flex items-center">
          <div class="w-2xs rounded-box shadow-md p-4 mr-5">
            <MovieCard {...movieInfos()!}/>
          </div>
          <div class="flex flex-col w-2xs">
            <span class="font-black text-2xl">{movieInfos()!.title}</span>
            <Show when={movieInfos()!.title !== movieInfos()!.original_title}>
              <span class="self-center text-sm">{movieInfos()!.original_title}</span>
            </Show>
          </div>
        </div>
        <Switch fallback={<p>An error has occured</p>}>
          <Match when={movieInCollection.loading || addSubmission.pending || removeSubmission.pending}>
            <button class="btn btn-soft btn-primary btn-wide"><span class="loading loading-spinner loading-sm"></span></button>
          </Match>
          <Match when={inCollection()}>
            <button class="btn btn-soft btn-primary btn-wide" onClick={handleRemoveFromCollection}>Remove from collection</button>
          </Match>
          <Match when={!inCollection()}>
            <button class="btn btn-soft btn-primary btn-wide" onClick={handleAddToCollection}>Add to collection</button>
          </Match>
        </Switch>
      </main>
    </Show>
  );
}

export default Movie
