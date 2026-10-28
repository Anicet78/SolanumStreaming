import { A, action, useAction, useSubmission } from "@solidjs/router";
import { createSignal, For, Match, Show, Switch, type Accessor, type Component } from "solid-js";
import { usersApi, type Languages, type User } from "../api/auth";
import { useAuthContext } from "../contexts/AuthContext";

export const updateProfileAction = action(
  async (username: string, newLang: Languages) => {
    return usersApi.update({ new_username: username, new_language: newLang });
  }
);

const LANGS: Record<Languages, { flag: string; label: string }> = {
  "fr-FR": { flag: "🇫🇷", label: "Français" },
  "en-US": { flag: "🇬🇧", label: "English" },
};

const Sidebar = () => {
  const auth = useAuthContext();

  const updateProfile = useAction(updateProfileAction);
  const submission = useSubmission(updateProfileAction);

  const handleSetLang = async (lang: Languages) => {
    const user = auth.user();
    if (!user || user.lang === lang) return;
    await updateProfile(user.username, lang);
    auth.updateUser({ lang });
  };

  return (
    <aside class="flex min-h-full w-64 flex-col gap-4 bg-base-200 p-4">
      <Show when={auth.user()}>
        {(user) => (
          <>
            <div class="flex items-center">
              <span class="truncate pl-1 font-semibold">{user().username}</span>
            </div>

            <button
              class="btn btn-outline btn-sm"
              popovertarget="popover-1"
              style="anchor-name:--anchor-1"
              disabled={submission.pending}
            >
              <span>
                {LANGS[user().lang]?.flag ?? "🌐"}{" "}
                {LANGS[user().lang]?.label ?? "Unknown language"}
              </span>
            </button>

            <ul
              class="dropdown menu w-52 rounded-box bg-base-100 shadow-sm"
              popover
              id="popover-1"
              style="position-anchor:--anchor-1"
            >
              <For each={Object.entries(LANGS) as [Languages, (typeof LANGS)[Languages]][]}>
                {([code, { flag, label }]) => (
                  <li>
                    <button
                      type="button"
                      classList={{ "menu-active": user().lang === code }}
                      onClick={(e) => {
                        handleSetLang(code);
                        // e.currentTarget.closest("[popover]")?.hidePopover();
                      }}
                    >
                      {flag} {label}
                    </button>
                  </li>
                )}
              </For>
            </ul>
            <div class="divider my-0" />
          </>
        )}

      </Show>
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
