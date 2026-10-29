import type { Component } from "solid-js"
import type { MovieT } from "../api/movies"

const MovieCard: Component<MovieT> = (props) => {
  return (
    <div class="movie-card">
      <img src={`https://image.tmdb.org/t/p/w500${props.poster_path}`} alt={props.title} loading="lazy" />
    </div>
  )
}

export default MovieCard
