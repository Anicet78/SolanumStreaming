import type { ParentComponent } from "solid-js";

const Page: ParentComponent = (props) => (
  <main class="mx-auto w-[min(100%_-_3rem,90rem)] pt-8.5 pb-16 mobile:w-[min(100%_-_1.75rem,90rem)] mobile:pt-6">
    {props.children}
  </main>
);

export default Page;
