import { createEffect, createSignal, createUniqueId, For, onCleanup, Show, type Component } from "solid-js";
import { clearOverride, setOverride, TIERS, type TierName } from "../api/capabilities";
import { useDisplay } from "../Capabilities";
import { pill, pillIdle } from "./pill";

const TIER_NAMES = Object.keys(TIERS) as TierName[];

const tierLabel = (tier: TierName) => TIERS[tier].label.replace(/\s*\(.*\)$/, "");
const summary = (tier: TierName, hdr: boolean) => `${tierLabel(tier)} · ${hdr ? "HDR" : "SDR"}`;

const DisplayMenu: Component = () => {
  const uid = createUniqueId();
  const { display, override, detected } = useDisplay();

  const [open, setOpen] = createSignal(false);
  const [shift, setShift] = createSignal(0);
  const [tier, setTier] = createSignal<TierName>("fhd");
  const [hdr, setHdr] = createSignal(false);

  let root!: HTMLDivElement;
  let trigger!: HTMLButtonElement;
  let panel: HTMLDivElement | undefined;

  const customised = () => override().tier !== undefined || override().hdr !== undefined;
  const dirty = () => {
    const d = display();
    return !!d && (d.tier !== tier() || d.hdr !== hdr());
  };

  const toggle = () => {
    if (open()) {
      setOpen(false);
      return;
    }
    const d = display();
    if (d) {
      setTier(d.tier);
      setHdr(d.hdr);
    }
    setOpen(true);
  };

  const close = () => {
    setOpen(false);
    trigger.focus();
  };

  const apply = () => {
    const det = detected();
    if (!det) return;
    setOverride({
      tier: tier() === det.tier ? undefined : tier(),
      hdr: hdr() === det.hdr ? undefined : hdr(),
    });
    close();
  };

  const reset = () => {
    clearOverride();
    close();
  };

  createEffect(() => {
    if (!open()) return;
    const onPointerDown = (e: PointerEvent) => {
      if (!root.contains(e.target as Node)) setOpen(false);
    };
    const onKeyDown = (e: KeyboardEvent) => {
      if (e.key === "Escape") close();
    };
    document.addEventListener("pointerdown", onPointerDown);
    document.addEventListener("keydown", onKeyDown);
    onCleanup(() => {
      document.removeEventListener("pointerdown", onPointerDown);
      document.removeEventListener("keydown", onKeyDown);
    });
  });

  createEffect(() => {
    if (!open() || !panel) return;
    const margin = parseFloat(getComputedStyle(document.documentElement).fontSize) * 0.75;
    const left = trigger.getBoundingClientRect().left;
    const max = window.innerWidth - panel.offsetWidth - margin;
    setShift(Math.max(margin, Math.min(left, max)) - left);
    panel.querySelector<HTMLInputElement>('input[type="radio"]:checked')?.focus();
  });

  return (
    <Show when={display()}>
      {(d) => (
        <div class="relative" ref={root}>
          <button
            ref={trigger}
            type="button"
            class={`group flex cursor-pointer items-center gap-[0.5em] transition-[border-color,color] duration-200 hover:border-brand hover:text-brand pointer-coarse:py-3 ${pill}`}
            classList={{ [pillIdle]: !customised(), "border-brand/60 text-brand": customised() }}
            aria-haspopup="dialog"
            aria-expanded={open()}
            aria-controls={`${uid}-panel`}
            title="Display used to pick the stream quality"
            onClick={toggle}
          >
            <svg class="size-[1em] shrink-0" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <rect x="2" y="3" width="20" height="14" rx="2" />
              <path d="M8 21h8M12 17v4" />
            </svg>
            {summary(d().tier, d().hdr)}
            <svg class="size-[0.8em] shrink-0 opacity-70 transition-transform group-aria-expanded:rotate-180" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round" aria-hidden="true">
              <path d="m6 9 6 6 6-6" />
            </svg>
          </button>

          <Show when={open()}>
            <div
              ref={panel}
              id={`${uid}-panel`}
              role="dialog"
              aria-label="Display settings"
              class="absolute top-full z-10 mt-2 flex w-[22em] max-w-[calc(100vw-1.5rem)] flex-col gap-[1em] rounded-[0.9em] border border-line-strong bg-surface p-[1em] text-[0.8rem] font-normal text-fog shadow-[0_0.75rem_1.5rem_rgb(0_0_0/28%)] transition-[opacity,translate] duration-150 starting:-translate-y-1 starting:opacity-0"
              style={{ left: `${shift()}px` }}
            >
              <div>
                <div class="font-semibold text-snow">Your display</div>
                <p class="mt-[0.25em] text-[0.9em] leading-snug text-dim">
                  We pick the stream that fits it. Saved on this device.
                </p>
              </div>

              <div>
                <div class="mb-[0.5em] flex items-baseline justify-between gap-[1em]">
                  <span id={`${uid}-res`} class="text-snow">Resolution</span>
                  <span class="text-[0.85em] text-dim">
                    Detected: {detected() ? tierLabel(detected()!.tier) : "–"}
                  </span>
                </div>
                <div role="radiogroup" aria-labelledby={`${uid}-res`} class="grid grid-cols-2 gap-[0.5em]">
                  <For each={TIER_NAMES}>
                    {(t) => (
                      <label
                        title={`${TIERS[t].width} × ${TIERS[t].height}`}
                        class="flex cursor-pointer items-center justify-center rounded-[0.6em] border border-line-strong px-[0.75em] py-[0.5em] text-center transition-colors pointer-coarse:min-h-11 hover:border-brand/60 has-checked:border-brand has-checked:bg-brand/10 has-checked:text-brand has-focus-visible:outline-2 has-focus-visible:outline-offset-2 has-focus-visible:outline-brand"
                      >
                        <input
                          type="radio"
                          name={`${uid}-tier`}
                          class="sr-only"
                          checked={tier() === t}
                          onChange={() => setTier(t)}
                        />
                        {tierLabel(t)}
                      </label>
                    )}
                  </For>
                </div>
              </div>

              <label class="flex cursor-pointer items-center justify-between gap-[1em]">
                <span>
                  <span class="block text-snow">HDR</span>
                  <span class="block text-[0.85em] text-dim">
                    Detected: {detected()?.hdr ? "on" : "off"}
                  </span>
                </span>
                <input
                  type="checkbox"
                  role="switch"
                  class="peer sr-only"
                  checked={hdr()}
                  onChange={(e) => setHdr(e.currentTarget.checked)}
                />
                <span
                  aria-hidden="true"
                  class="relative h-[1.5em] w-[2.6em] shrink-0 rounded-full border border-line-strong bg-surface-deep transition-colors after:absolute after:top-1/2 after:left-[0.2em] after:size-[1em] after:-translate-y-1/2 after:rounded-full after:bg-fog after:transition-[translate,background-color] after:content-[''] peer-checked:border-brand peer-checked:bg-brand peer-checked:after:translate-x-[1.1em] peer-checked:after:bg-ink peer-focus-visible:outline-2 peer-focus-visible:outline-offset-2 peer-focus-visible:outline-brand"
                />
              </label>

              <div class="flex items-center justify-between gap-[0.75em] border-t border-line-soft pt-[0.9em]">
                <Show when={customised()}>
                  <button type="button" class="cursor-pointer py-[0.4em] text-dim transition-colors hover:text-brand pointer-coarse:min-h-11" onClick={reset}>
                    Use detected
                  </button>
                </Show>
                <button
                  type="button"
                  class="ml-auto cursor-pointer rounded-[0.6em] bg-brand px-[1em] py-[0.5em] font-bold text-ink pointer-coarse:min-h-11 hover:brightness-[1.08] disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:brightness-100"
                  disabled={!dirty()}
                  onClick={apply}
                >
                  Apply
                </button>
              </div>
            </div>
          </Show>
        </div>
      )}
    </Show>
  );
};

export default DisplayMenu;
