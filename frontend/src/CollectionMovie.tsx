import { useParams } from "@solidjs/router";
import { createResource, Show } from "solid-js";
import { moviesApi } from "./api/movies";
import MovieCard from "./components/MovieCard";

const CollectionMovie = () => {
  const params = useParams()

  const id = Number(params.movie_id)
  if (!id) return <span>Not Found</span>

  const [movieInCollection] = createResource(() => id, moviesApi.getInCollection);
  const [movieInfos] = createResource(() => id, moviesApi.getMovieInfos);

  return (
    <Show when={!movieInfos.loading && !movieInfos.error}>
      <div class="min-h-screen flex flex-col items-center p-4">
        <span class="text-2xl">{movieInfos()?.title}</span>
        <Show when={movieInfos()?.title !== movieInfos()?.original_title}><span class="text-lg">{movieInfos()?.original_title}</span></Show>
      </div>
    </Show>
  )
}

export default CollectionMovie