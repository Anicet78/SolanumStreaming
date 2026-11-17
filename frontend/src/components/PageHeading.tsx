import type { Component } from "solid-js";
import Eyebrow from "./Eyebrow";
import SearchBar from "./SearchBar";

type PageHeadingProps = {
  eyebrow: string;
  title: string;
  subtitle: string;
};

const PageHeading: Component<PageHeadingProps> = (props) => (
  <div class="mb-7.5 flex items-end justify-between gap-6 mobile:mb-5.5 mobile:block 4xl:mb-14">
    <div>
      <Eyebrow>{props.eyebrow}</Eyebrow>
      <h1 class="mt-1.5 font-heading text-[clamp(2rem,4vw,3.5rem)] leading-none tracking-[-0.03em] text-snow 3xl:text-[clamp(3rem,3.5vw,5rem)] 4xl:text-[clamp(4rem,4vw,7rem)]">
        {props.title}
      </h1>
      <p class="mt-3 max-w-[520px] text-[0.98rem] text-dim 3xl:text-[1.1rem] 4xl:max-w-[900px] 4xl:text-[1.4rem] 4xl:leading-[1.55]">
        {props.subtitle}
      </p>
    </div>
    <div class="w-[min(100%,560px)] 4xl:w-[min(100%,820px)]">
      <SearchBar />
    </div>
  </div>
);

export default PageHeading;
