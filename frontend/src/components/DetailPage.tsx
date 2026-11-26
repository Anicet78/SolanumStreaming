import type { ParentComponent } from "solid-js";

const DetailPage: ParentComponent = (props) => (
  <main class="min-h-dvh bg-ink px-[max(1.5rem,5vw)] pt-7 pb-16 max-lg:min-h-[auto] max-lg:pb-8 narrow:px-3.5 narrow:pt-4.5 narrow:pb-8">
    {props.children}
  </main>
);

export default DetailPage;
