import { For } from "solid-js"
import type { SearchResponse } from "./api/types"
import MovieCard from "./components/MovieCard"

const Search = () => {
	const resp: SearchResponse = {
		page: 1,
		results: [{
			id: 1,
			title: "Backrooms",
			original_title: "The Backrooms",
			popularity: 10,
			adult: false,
			video: false,
			poster_path: "/google/backrooms",
			release_date: "2026",
			imdb_id: "1",
			runtime: 120
		}],
		total_pages: 1,
		total_results: 1
	}

	return (
		<div class="flex align-items">
			<ul class="list bg-base-100 rounded-box shadow-md">

				<For each={resp.results}>
					{(movie) => <li><MovieCard/></li>}
				</For>

			</ul>
		</div>
	)
}

export default Search