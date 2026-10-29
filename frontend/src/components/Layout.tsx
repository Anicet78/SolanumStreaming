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
    <div class="app-shell" classList={{ "sidebar-is-open": sidebarOpen() }}>
      <div class="app-main">
        <div class="mobile-header">
          <button class="menu-button" type="button" aria-label="Open menu" onClick={() => setSidebarOpen(true)}>☰</button>
          <span class="brand-mark">SOLANUM</span>
        </div>

        <main class="app-content">
          <Show when={auth.user()?.lang} keyed>
            {props.children}
          </Show>
        </main>
      </div>

      <div class="sidebar-panel" classList={{ "sidebar-visible": sidebarOpen() }}>
        <button class="sidebar-overlay" type="button" aria-label="Close menu" onClick={() => setSidebarOpen(false)} />
        <button class="desktop-sidebar-close" type="button" aria-label="Hide sidebar" onClick={() => setSidebarOpen(false)}>←</button>
        <Sidebar/>
      </div>
      <Show when={!sidebarOpen()}>
        <button class="desktop-sidebar-open" type="button" aria-label="Show sidebar" onClick={() => setSidebarOpen(true)}>→</button>
      </Show>
    </div>
  );
};

export default Layout;
