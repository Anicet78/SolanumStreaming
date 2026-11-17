import type { ParentComponent } from "solid-js";
import Eyebrow from "./Eyebrow";

type AuthPanelProps = {
  eyebrow: string;
  tagline: string;
  switchText: string;
  switchLabel: string;
  onSwitch: () => void;
};

const AuthPanel: ParentComponent<AuthPanelProps> = (props) => (
  <div class="grid min-h-screen place-items-center bg-ink bg-[radial-gradient(circle_at_50%_0,var(--color-glow)_0,transparent_48%)] p-6">
    <div class="w-[min(100%,430px)] rounded-[18px] border border-line bg-[rgb(17_29_39/92%)] p-8.5 shadow-[0_24px_70px_rgb(0_0_0/28%)]">
      <div class="mb-6.5 text-center">
        <Eyebrow>{props.eyebrow}</Eyebrow>
        <strong class="mt-[7px] block font-heading text-[1.65rem] text-snow">Solanum Streaming</strong>
        <p class="mt-2 text-[0.88rem] text-dim">{props.tagline}</p>
      </div>

      {props.children}

      <div class="mt-5.5 text-center text-[0.82rem] text-dim">
        <span>{props.switchText}</span>
        <button class="btn btn-xs btn-link text-brand" onClick={props.onSwitch}>
          {props.switchLabel}
        </button>
      </div>
    </div>
  </div>
);

export default AuthPanel;
