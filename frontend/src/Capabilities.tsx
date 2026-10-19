import { createSignal, onCleanup, onMount } from "solid-js";
import { ensureCapabilities, getDisplay, getOverride, subscribeCapabilities } from "./api/capabilities";

function useDisplay() {
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

export default useDisplay;