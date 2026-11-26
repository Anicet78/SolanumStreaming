import type { Component } from "solid-js";
import Eyebrow from "./Eyebrow";
import SearchBar from "./SearchBar";

type PageHeadingProps = {
  eyebrow: string;
  title: string;
  subtitle: string;
};

const PageHeading: Component<PageHeadingProps> = (props) => (
  <div class="mb-7.5 flex items-end justify-between gap-6 mobile:mb-5.5 mobile:flex-col mobile:items-stretch mobile:gap-4">
    <div>
      <Eyebrow>{props.eyebrow}</Eyebrow>
      <h1 class="mt-1.5 font-heading text-[clamp(2rem,4vw,3.5rem)] leading-none tracking-[-0.03em] break-words text-snow">
        {props.title}
      </h1>
      <p class="mt-3 max-w-[32.5rem] text-[0.98rem] text-dim">
        {props.subtitle}
      </p>
    </div>
    <div class="w-[min(100%,35rem)]">
      <SearchBar />
    </div>
  </div>
);

export default PageHeading;
