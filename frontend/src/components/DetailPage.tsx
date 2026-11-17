import type { ParentComponent } from "solid-js";

const DetailPage: ParentComponent = (props) => (
  <main class="min-h-screen bg-ink px-[max(24px,5vw)] pt-7 pb-16 max-lg:min-h-[auto] max-lg:pb-8 narrow:px-3.5 narrow:pt-4.5 narrow:pb-8 3xl:px-16 4xl:mx-auto 4xl:max-w-[3600px] 4xl:px-25 4xl:pt-16 4xl:pb-25">
    {props.children}
  </main>
);

export default DetailPage;
