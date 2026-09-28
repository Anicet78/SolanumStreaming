import type { MovieT } from "../api/movies"

const MovieCard = (movie: MovieT) => {
  return (
    <div class="hover-3d w-full aspect-5/7 rounded-2xl">
    <figure class="max-w-100 rounded-2xl">
    <img src={`https://image.tmdb.org/t/p/w342${movie.poster_path}`} alt={movie.title} />
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