import { createSignal, onCleanup, onMount, Show, type ParentComponent } from "solid-js";
import Sidebar from "./Sidebar";
import { useAuthContext } from "../contexts/AuthContext";

const desktopQuery = "(min-width: 64rem)";

const Layout: ParentComponent = (props) => {
	const auth = useAuthContext();
  const [sidebarOpen, setSidebarOpen] = createSignal(
    typeof window !== "undefined" && window.matchMedia(desktopQuery).matches
  );

  onMount(() => {
    const mediaQuery = window.matchMedia(desktopQuery);
    const syncSidebar = (event: MediaQueryListEvent) => setSidebarOpen(event.matches);
    mediaQuery.addEventListener("change", syncSidebar);
    onCleanup(() => mediaQuery.removeEventListener("change", syncSidebar));
  });
  return (
    <div class="min-h-dvh bg-ink bg-[radial-gradient(circle_at_80%_-10%,var(--color-glow)_0,transparent_38%)]">
      <div
        class="ml-0 min-h-dvh min-w-0 transition-[margin-left] duration-220 ease-[ease]"
        classList={{
          "lg:ml-62": sidebarOpen(),
          "lg:pt-14.5": !sidebarOpen(),
        }}
      >
        <div class="z-20 hidden h-16 items-center gap-3.5 border-b border-line-soft bg-[rgb(10_17_24/92%)] px-4.5 max-lg:flex tall:sticky tall:top-0">
          <button class="grid size-8.5 cursor-pointer place-items-center rounded-[0.5625rem] border border-line-strong text-brand pointer-coarse:size-11" type="button" aria-label="Open menu" onClick={() => setSidebarOpen(true)}>☰</button>
          <span class="font-heading text-[0.95rem] tracking-[0.14em] text-snow">SOLANUM</span>
        </div>

        <main class="min-w-0">
          <Show when={auth.user()?.lang} keyed>
            {props.children}
          </Show>
        </main>
      </div>

      <div
        class="fixed inset-y-0 left-0 z-30 w-62 transition-transform duration-220 ease-[ease] max-lg:w-[min(84vw,18.75rem)]"
        classList={{
          "translate-x-0 pointer-events-auto": sidebarOpen(),
          "-translate-x-full pointer-events-none": !sidebarOpen(),
        }}
      >
        <button
          class="z-0 hidden border-0 max-lg:fixed max-lg:cursor-pointer max-lg:inset-0 max-lg:h-dvh max-lg:w-screen max-lg:bg-[rgb(3_8_12/65%)]"
          classList={{ "max-lg:block": sidebarOpen() }}
          type="button"
          aria-label="Close menu"
          onClick={() => setSidebarOpen(false)}
        />
        <button class="absolute top-4.5 right-3.5 z-2 hidden size-8.5 cursor-pointer place-items-center rounded-[0.5625rem] border border-line-strong bg-surface-deep text-[1rem] text-brand hover:border-brand hover:bg-surface-hover pointer-coarse:size-11 lg:grid" type="button" aria-label="Hide sidebar" onClick={() => setSidebarOpen(false)}>←</button>
        <Sidebar/>
      </div>
      <Show when={!sidebarOpen()}>
        <button class="fixed top-4.5 left-4.5 z-20 hidden size-8.5 cursor-pointer place-items-center rounded-[0.5625rem] border border-line-strong bg-surface-deep text-[1rem] text-brand hover:border-brand hover:bg-surface-hover pointer-coarse:size-11 lg:grid" type="button" aria-label="Show sidebar" onClick={() => setSidebarOpen(true)}>→</button>
      </Show>
    </div>
  );
};

export default Layout;
