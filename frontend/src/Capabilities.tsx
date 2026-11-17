import { createSignal, onCleanup, onMount } from "solid-js";
import { ensureCapabilities, getCapabilities, getDisplay, getOverride, subscribeCapabilities } from "./api/capabilities";

export function useDisplay() {
  const [display, setDisplay] = createSignal(getDisplay());
  const [override, setOverride] = createSignal(getOverride());
  const [detected, setDetected] = createSignal(getCapabilities()?.detected ?? null);

  const sync = () => {
    setDisplay(getDisplay());
    setOverride({ ...getOverride() });
    setDetected(getCapabilities()?.detected ?? null);
  };
  onCleanup(subscribeCapabilities(sync));
  onMount(() => void ensureCapabilities());

  return { display, override, detected };
}

const Capabilities = () => {
  const display = useDisplay().display

  return (
    <div>
      <p>{display()?.tier}</p>
      <p>{display()?.height}</p>
      <p>{display()?.width}</p>
    </div>
  )
}

export default Capabilities
