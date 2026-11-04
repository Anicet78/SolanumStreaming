import { createSignal, onCleanup, onMount, Show, type ParentComponent } from "solid-js";
import Sidebar from "./Sidebar";
import { useAuthContext } from "../contexts/AuthContext";

const Layout: ParentComponent = (props) => {
	const auth = useAuthContext();
  const [sidebarOpen, setSidebarOpen] = createSignal(
    typeof window !== "undefined" && window.matchMedia("(min-width: 1024px)").matches
  );

  onMount(() => {
    const mediaQuery = window.matchMedia("(min-width: 1024px)");
    const syncSidebar = (event: MediaQueryListEvent) => setSidebarOpen(event.matches);
    mediaQuery.addEventListener("change", syncSidebar);
    onCleanup(() => mediaQuery.removeEventListener("change", syncSidebar));
  });
  return (
    <div class="min-h-screen bg-ink bg-[radial-gradient(circle_at_80%_-10%,var(--color-glow)_0,transparent_38%)] 3xl:text-[18px] 4xl:text-[25px]">
      <div
        class="ml-0 min-h-screen min-w-0 transition-[margin-left] duration-220 ease-[ease]"
        classList={{
          "lg:ml-62 3xl:ml-75": sidebarOpen(),
          "lg:pt-14.5 3xl:pt-19": !sidebarOpen(),
        }}
      >
        <div class="hidden h-16 items-center gap-3.5 border-b border-line-soft bg-[rgb(10_17_24/92%)] px-4.5 max-lg:flex">
          <button class="grid size-8.5 cursor-pointer place-items-center rounded-[9px] border border-line-strong text-brand" type="button" aria-label="Open menu" onClick={() => setSidebarOpen(true)}>☰</button>
          <span class="font-heading text-[0.95rem] tracking-[0.14em] text-snow">SOLANUM</span>
        </div>

        <main class="min-w-0">
          <Show when={auth.user()?.lang} keyed>
            {props.children}
          </Show>
        </main>
      </div>

      <div
        class="fixed inset-y-0 left-0 z-30 w-62 transition-transform duration-220 ease-[ease] max-lg:w-[min(84vw,300px)] 3xl:w-75"
        classList={{
          "translate-x-0 pointer-events-auto": sidebarOpen(),
          "-translate-x-full pointer-events-none": !sidebarOpen(),
        }}
      >
        <button
          class="z-0 hidden border-0 max-lg:fixed max-lg:cursor-pointer max-lg:inset-0 max-lg:h-screen max-lg:w-screen max-lg:bg-[rgb(3_8_12/65%)]"
          classList={{ "max-lg:block": sidebarOpen() }}
          type="button"
          aria-label="Close menu"
          onClick={() => setSidebarOpen(false)}
        />
        <button class="absolute top-4.5 right-3.5 z-2 hidden size-8.5 cursor-pointer place-items-center rounded-[9px] border border-line-strong bg-surface-deep text-[1rem] text-brand hover:border-brand hover:bg-surface-hover lg:grid" type="button" aria-label="Hide sidebar" onClick={() => setSidebarOpen(false)}>←</button>
        <Sidebar/>
      </div>
      <Show when={!sidebarOpen()}>
        <button class="fixed top-4.5 left-4.5 z-20 hidden size-8.5 cursor-pointer place-items-center rounded-[9px] border border-line-strong bg-surface-deep text-[1rem] text-brand hover:border-brand hover:bg-surface-hover lg:grid" type="button" aria-label="Show sidebar" onClick={() => setSidebarOpen(true)}>→</button>
      </Show>
    </div>
  );
};

export default Layout;
