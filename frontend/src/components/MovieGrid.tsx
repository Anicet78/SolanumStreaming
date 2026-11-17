import type { ParentComponent } from "solid-js";

type MovieGridProps = {
  loading?: boolean;
};

const MovieGrid: ParentComponent<MovieGridProps> = (props) => (
  <div
    class="grid grid-cols-[repeat(auto-fill,minmax(150px,1fr))] gap-x-4 gap-y-5.5 mobile:grid-cols-2 mobile:gap-x-3 mobile:gap-y-4 3xl:grid-cols-[repeat(auto-fill,minmax(230px,1fr))] 3xl:gap-x-6.5 3xl:gap-y-8 4xl:grid-cols-[repeat(auto-fill,minmax(300px,1fr))] 4xl:gap-x-9 4xl:gap-y-11.5 [&>a]:block [&>a]:h-fit [&>a]:min-w-0 [&>a]:leading-[0]"
    classList={{ "opacity-65": props.loading }}
  >
    {props.children}
  </div>
);

export default MovieGrid;
