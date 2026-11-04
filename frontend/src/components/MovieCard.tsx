import type { Component } from "solid-js"
import type { MovieT } from "../api/movies"

const MovieCard: Component<MovieT> = (props) => {
  return (
    <div class="relative aspect-2/3 overflow-hidden rounded-[14px] border border-line bg-[#15212c] shadow-[0_12px_26px_rgb(0_0_0/18%)] transition-[translate,box-shadow] duration-200 ease-[ease] after:pointer-events-none after:absolute after:inset-x-0 after:bottom-0 after:h-[35%] after:bg-[linear-gradient(transparent,rgb(5_10_15/70%))] after:content-[''] hover:-translate-y-[5px] hover:shadow-[0_0_0_1px_var(--color-brand),0_18px_34px_rgb(0_0_0/30%)]">
      <img class="block size-full object-cover" src={`https://image.tmdb.org/t/p/w500${props.poster_path}`} alt={props.title} loading="lazy" />
    </div>
  )
}

export default MovieCard
