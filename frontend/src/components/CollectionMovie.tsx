import { createResource, Show } from "solid-js"
import { moviesApi } from "../api/movies"
import MovieCard from "./MovieCard"
import { useParams } from "@solidjs/router"

const CollectionMovie = () => {
  const params = useParams()

  if (!params.movie_id)
    return <span>Not Found</span>
  const id: number = Number(params.movie_id)

  const [movieInCollection] = createResource(() => id, moviesApi.getInCollection);
  const [movieInfos] = createResource(() => id, moviesApi.getMovieInfos);


  return (
    <Show when={!movieInfos.loading && !movieInfos.error}>
      <MovieCard {...movieInfos()!}/>
    </Show>
  )
}

export default CollectionMovie