import { A } from "@solidjs/router";
import type { Component } from "solid-js";
import SearchBar from "./SearchBar";

type DetailBackbarProps = {
  href: string;
  label: string;
};

const DetailBackbar: Component<DetailBackbarProps> = (props) => (
  <div class="mb-8 flex items-center justify-between gap-5 narrow:flex-col-reverse narrow:items-stretch">
    <A class="-my-2 py-2 text-[0.85rem] text-dim hover:text-brand pointer-coarse:-my-3 pointer-coarse:py-3" href={props.href}>
      {props.label}
    </A>
    <div class="w-[min(100%,30rem)] narrow:w-full">
      <SearchBar />
    </div>
  </div>
);

export default DetailBackbar;
