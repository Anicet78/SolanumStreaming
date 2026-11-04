import type { ParentComponent } from "solid-js";

const Page: ParentComponent = (props) => (
  <main class="mx-auto w-[min(100%_-_48px,1440px)] pt-8.5 pb-16 mobile:w-[min(100%_-_28px,1440px)] mobile:pt-6 3xl:w-[min(calc(100%_-_80px),2100px)] 3xl:pt-13 4xl:w-[min(calc(100%_-_160px),3600px)] 4xl:pt-19">
    {props.children}
  </main>
);

export default Page;
