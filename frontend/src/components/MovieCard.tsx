import type { MovieT } from "../api/types"

const MovieCard = (movie: MovieT) => {
	return (
		<div class="hover-3d">
		<figure class="max-w-100 rounded-2xl">
		<img src={movie.poster_path} alt={movie.title} />
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