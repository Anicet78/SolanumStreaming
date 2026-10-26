import { createSignal, onCleanup, onMount } from "solid-js";
import { ensureCapabilities, getDisplay, getOverride, subscribeCapabilities } from "./api/capabilities";

export function useDisplay() {
  const [display, setDisplay] = createSignal(getDisplay());
  const [override, setOverride] = createSignal(getOverride());

  const sync = () => {
    setDisplay(getDisplay());
    setOverride({ ...getOverride() });
  };
  onCleanup(subscribeCapabilities(sync));
  onMount(() => void ensureCapabilities());

  return { display, override };
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