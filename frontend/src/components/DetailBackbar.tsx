import { A } from "@solidjs/router";
import type { Component } from "solid-js";
import SearchBar from "./SearchBar";

type DetailBackbarProps = {
  href: string;
  label: string;
};

const DetailBackbar: Component<DetailBackbarProps> = (props) => (
  <div class="mb-8 flex items-center justify-between gap-5 narrow:flex-col-reverse narrow:items-stretch">
    <A class="text-[0.85rem] text-dim hover:text-brand 3xl:text-[0.95rem] 4xl:text-[1.75rem]" href={props.href}>
      {props.label}
    </A>
    <div class="w-[min(100%,480px)] narrow:w-full">
      <SearchBar />
    </div>
  </div>
);

export default DetailBackbar;
