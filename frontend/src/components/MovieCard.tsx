import type { Component } from "solid-js"
import type { MovieT } from "../api/movies"

const MovieCard: Component<MovieT> = (props) => {
  return (
    <div class="hover-3d w-full aspect-5/7 rounded-2xl">
    <figure class="max-w-100 rounded-2xl">
    <img src={`https://image.tmdb.org/t/p/original${props.poster_path}`} alt={props.title} />
    </figure>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    <div></div>
    </div>
  )
}

export default MovieCard