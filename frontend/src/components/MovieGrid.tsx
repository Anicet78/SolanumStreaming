import type { ParentComponent } from "solid-js";

type MovieGridProps = {
  loading?: boolean;
};

const MovieGrid: ParentComponent<MovieGridProps> = (props) => (
  <div
    class="grid grid-cols-[repeat(auto-fill,minmax(9.375rem,1fr))] gap-x-4 gap-y-5.5 mobile:grid-cols-2 mobile:gap-x-3 mobile:gap-y-4 [&>a]:block [&>a]:h-fit [&>a]:min-w-0 [&>a]:leading-[0]"
    classList={{ "opacity-65": props.loading }}
  >
    {props.children}
  </div>
);

export default MovieGrid;
