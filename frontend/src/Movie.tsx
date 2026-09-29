import { action, useAction, useParams, useSubmission } from "@solidjs/router";
import { moviesApi } from "./api/movies";
import { createResource, Show } from "solid-js";

export const addToCollectionAction = action(async (movie_id: number) => {
  return moviesApi.addToCollection({ movie_id });
});

export const removeFromCollectionAction = action(async (movie_id: number) => {
  return moviesApi.removeFromCollection(movie_id);
});

const Movie = () => {
  const params = useParams();

  if (!params.title)
    return <span>Not Found</span>

  const sep: number = params.title.lastIndexOf("-")
  const title: string = params.title.substring(0, sep)
  const id: number = Number(params.title.substring(sep + 1))

  const [movieInCollection] = createResource(() => id, moviesApi.getInCollection);

  const addToCollection = useAction(addToCollectionAction);
  const addSubmission = useSubmission(addToCollectionAction);

  const removeFromCollection = useAction(removeFromCollectionAction);
  const removeSubmission = useSubmission(removeFromCollectionAction);

  const handleAddToCollection = async () => {
    await addToCollection(id);
  };

  const handleRemoveFromCollection = async () => {
    await removeFromCollection(id);
  };

  return (
    <main class="flex flex-col justify-center items-center">
      <h1 class="mt-10">{decodeURIComponent(title)}</h1>
      <h1 class="m-2 mb-10">{id}</h1>
      <Show when={movieInCollection.loading || addSubmission.pending || removeSubmission.pending}>
        <button class="btn btn-soft btn-primary btn-wide"><span class="loading loading-spinner loading-sm"></span></button>
      </Show>
      <Show when={movieInCollection.error }>
        <button class="btn btn-soft btn-primary btn-wide" onClick={handleAddToCollection}>Add to collection</button>
      </Show>
      <Show when={!movieInCollection.loading && !movieInCollection.error}>
        <button class="btn btn-soft btn-primary btn-wide" onClick={handleRemoveFromCollection}>Remove from collection</button>
      </Show>
    </main>
  );
}

export default Movie