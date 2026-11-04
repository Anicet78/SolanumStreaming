import { Show, type Component } from "solid-js";
import type { moviesApi } from "../api/movies";
import Eyebrow from "./Eyebrow";

type MovieInfos = Awaited<ReturnType<typeof moviesApi.getMovieInfos>>;

type MovieDetailsProps = {
  movie: MovieInfos;
  eyebrow: string;
  compact?: boolean;
};

const stat =
  "rounded-lg border border-line-strong bg-[rgb(11_17_24/60%)] px-[11px] py-2 text-[0.8rem] text-haze 3xl:text-[0.9rem] 4xl:px-4.5 4xl:py-[13px] 4xl:text-[1.5rem]";

const MovieDetails: Component<MovieDetailsProps> = (props) => (
  <>
    <Eyebrow>{props.eyebrow}</Eyebrow>
    <h1
      class="mt-2 font-heading leading-[0.96] tracking-[-0.04em] text-snow"
      classList={{
        "text-[clamp(2.2rem,6vw,5rem)] 4xl:text-[clamp(5rem,5vw,8rem)]": !props.compact,
        "text-[clamp(1.8rem,4vw,3rem)]": props.compact,
      }}
    >
      {props.movie.title}
    </h1>
    <Show when={props.movie.title !== props.movie.original_title}>
      <div class="mt-[9px] text-[0.9rem] text-dim 3xl:text-[1rem] 4xl:text-[1.75rem]">
        Original title: {props.movie.original_title}
      </div>
    </Show>
    <div class="my-6 flex flex-wrap gap-2.5">
      <span class={stat}>{props.movie.release_date.slice(0, 4)}</span>
      <span class={stat}>{Math.floor(props.movie.runtime / 60)}h {props.movie.runtime % 60}m</span>
      <span class={stat}>★ {props.movie.vote_average.toFixed(1)} / 10</span>
    </div>
    <Show when={props.movie.tagline || props.movie.overview}>
      <Show when={props.movie.tagline}>
        <div class="mt-6 text-[0.82rem] font-bold tracking-[0.12em] text-brand uppercase 3xl:text-[0.9rem] 4xl:text-[1.8rem]">
          {props.movie.tagline}
        </div>
      </Show>
      <Show when={props.movie.overview}>
        <p class="mt-3 max-w-[650px] leading-[1.65] text-haze 3xl:text-[1.05rem] 4xl:max-w-[1050px] 4xl:text-[2rem] 4xl:leading-[1.5]">
          {props.movie.overview}
        </p>
      </Show>
    </Show>
  </>
);

export default MovieDetails;
