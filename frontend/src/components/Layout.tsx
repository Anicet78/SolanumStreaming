import { Show, type ParentComponent } from "solid-js";
import Sidebar from "./Sidebar";
import { useAuthContext } from "../contexts/AuthContext";

const Layout: ParentComponent = (props) => {
	const auth = useAuthContext();
  return (
    <div class="drawer lg:drawer-open">
      <input id="app-drawer" type="checkbox" class="drawer-toggle" />

      <div class="drawer-content flex flex-col">
        <div class="navbar bg-base-100 lg:hidden">
          <label for="app-drawer" class="btn btn-square btn-ghost" aria-label="Open menu">
            ☰
          </label>
        </div>

        <main>
          <Show when={auth.user()?.lang} keyed>
            {props.children}
          </Show>
        </main>
      </div>

      {/* Sidebar */}
      <div class="drawer-side z-10">
        <label for="app-drawer" aria-label="Close menu" class="drawer-overlay" />
        <Sidebar/>
      </div>
    </div>
  );
};

export default Layout;
