import { A } from "@solidjs/router";
import { Match, Switch, type Component } from "solid-js";

interface SidebarProps {
  username: string;
  locale: string;
}

const Sidebar: Component<SidebarProps> = (props) => {
  return (
    <aside class="flex min-h-full w-64 flex-col gap-4 bg-base-200 p-4">
      <div class="flex items-center">
        <span class="truncate font-semibold">{props.username}</span>
      </div>

      <button class="btn btn-outline btn-sm" popovertarget="popover-1" style="anchor-name:--anchor-1">
        <ul class="dropdown menu w-52 rounded-box bg-base-100 shadow-sm"
          popover id="popover-1" style="position-anchor:--anchor-1">
          <li><a>🇫🇷 Français</a></li>
          <li><a>🇬🇧 English</a></li>
        </ul>
        <Switch fallback={<p>🌐 Unknown language</p>}>
          <Match when={props.locale === "fr-FR"}>
            <p>🇫🇷 Français</p>
          </Match>
          <Match when={props.locale === "en-US"}>
            <p>🇬🇧 English</p>
          </Match>
       </Switch>
      </button>

      <div class="divider my-0" />

      <ul class="menu w-full gap-1 p-0">
        <li>
          <A href="/" end activeClass="menu-active">
            Home
          </A>
        </li>
        <li>
          <A href="/collection" activeClass="menu-active">
            Collection
          </A>
        </li>
      </ul>
    </aside>
  );
};

export default Sidebar;