import { Show, type ParentComponent } from "solid-js";
import Eyebrow from "../Eyebrow";

type AuthPanelProps = {
  eyebrow: string;
  tagline: string;
  submitLabel: string;
  pendingLabel: string;
  pending?: boolean;
  error?: string;
  onSubmit: (e: SubmitEvent) => void;
  switchText: string;
  switchLabel: string;
  onSwitch: () => void;
};

const AuthPanel: ParentComponent<AuthPanelProps> = (props) => (
  <div class="grid min-h-dvh place-items-center bg-ink bg-[radial-gradient(circle_at_50%_0,var(--color-glow)_0,transparent_48%)] p-6 mobile:p-3">
    <div class="w-full max-w-[26.875rem] rounded-[1.125rem] border border-line bg-[rgb(17_29_39/92%)] p-8.5 shadow-[0_1.5rem_4.375rem_rgb(0_0_0/28%)] mobile:p-5">
      <div class="mb-6.5 text-center">
        <Eyebrow>{props.eyebrow}</Eyebrow>
        <h1 class="mt-[0.4375rem] font-heading text-[1.65rem] font-bold break-words text-snow">Solanum Streaming</h1>
        <p class="mt-2 text-[0.88rem] text-dim">{props.tagline}</p>
      </div>

      <form class="grid gap-3" onSubmit={props.onSubmit}>
        {props.children}
        <button
          class="btn btn-primary mt-1 w-full text-[0.9rem] font-bold pointer-coarse:[--size:3rem]"
          type="submit"
          disabled={props.pending}
        >
          {props.pending ? props.pendingLabel : props.submitLabel}
        </button>
        <Show when={props.error}>
          <p role="alert" class="rounded-lg border border-danger-line px-3 py-2 text-[0.85rem] break-words text-danger">
            {props.error}
          </p>
        </Show>
      </form>

      <div class="mt-5.5 text-center text-[0.82rem] text-dim">
        <span>{props.switchText}</span>
        <button
          class="ml-1 cursor-pointer rounded-sm px-1.5 py-2 font-semibold text-brand underline underline-offset-4 pointer-coarse:py-3"
          type="button"
          onClick={props.onSwitch}
        >
          {props.switchLabel}
        </button>
      </div>
    </div>
  </div>
);

export default AuthPanel;
