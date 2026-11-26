import type { ParentComponent } from "solid-js";

const Eyebrow: ParentComponent = (props) => (
  <div class="text-[0.72rem] font-bold tracking-[0.16em] text-brand uppercase">
    {props.children}
  </div>
);

export default Eyebrow;
