import { A, action, useAction, useSubmission } from "@solidjs/router";
import { createSignal, For, Show } from "solid-js";
import { usersApi, type Languages } from "../api/auth";
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
  const [languageMenuOpen, setLanguageMenuOpen] = createSignal(false);

  const handleSetLang = async (lang: Languages) => {
    const user = auth.user();
    if (!user || user.lang === lang) return;
    await updateProfile(user.username, lang);
    auth.updateUser({ lang });
  };

  return (
    <aside class="app-sidebar">
      <Show when={auth.user()}>
        {(user) => (
          <>
            <div class="brand-lockup">
              <span class="brand-symbol">S</span>
              <div><strong>SOLANUM</strong><small>STREAMING</small></div>
            </div>
            <div class="user-chip"><span class="user-dot" />{user().username}</div>

            <button
              class="language-button"
              disabled={submission.pending}
              aria-expanded={languageMenuOpen()}
              onClick={() => setLanguageMenuOpen((open) => !open)}
            >
              <span class="language-label">
                {LANGS[user().lang]?.flag ?? "🌐"}{" "}
                {LANGS[user().lang]?.label ?? "Unknown language"}
              </span>
            </button>

            <Show when={languageMenuOpen()}>
              <ul class="language-menu">
                <For each={Object.entries(LANGS) as [Languages, (typeof LANGS)[Languages]][]}>
                  {([code, { flag, label }]) => (
                    <li>
                      <button
                        type="button"
                        classList={{ "menu-active": user().lang === code }}
                        onClick={() => {
                          void handleSetLang(code);
                          setLanguageMenuOpen(false);
                        }}
                      >
                        {flag} {label}
                      </button>
                    </li>
                  )}
                </For>
              </ul>
            </Show>
            <div class="sidebar-divider" />
          </>
        )}

      </Show>
      <ul class="sidebar-nav">
        <li>
          <A href="/" end activeClass="menu-active">
            <span>⌂</span> Home
          </A>
        </li>
        <li>
          <A href="/collection" activeClass="menu-active">
            <span>▣</span> Collection
          </A>
        </li>
      </ul>
    </aside>
  );
};

export default Sidebar;
