import { createResource, Show, type Component } from "solid-js"
import { moviesApi, type CollectionMovieT } from "../api/movies"
import MovieCard from "./MovieCard"

const CollectionMovieCard: Component<CollectionMovieT> = (props) => {
  const [movieInfos] = createResource(() => props.movie_id, moviesApi.getMovieInfos);

  return (
    <Show when={!movieInfos.loading && !movieInfos.error}>
      <MovieCard {...movieInfos()!}/>
    </Show>
  )
}

export default CollectionMovieCard
