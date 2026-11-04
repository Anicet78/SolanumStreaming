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

const navLink =
  "flex items-center gap-3 rounded-[9px] px-3 py-[11px] text-[0.91rem] transition-[background,color] duration-200 ease-[ease] hover:bg-surface-hover hover:text-snow 3xl:text-[1rem] 4xl:text-[1.55rem]";
const navInactive = "text-[#91a4b7]";
const navActive = "bg-surface-hover text-snow";
const navIcon = "w-4.5 text-center text-[1.1rem] text-brand";

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
    <aside class="relative z-1 min-h-screen w-62 border-r border-line-soft bg-[rgb(10_17_24/98%)] px-4.5 py-7 max-lg:w-[min(84vw,300px)] 3xl:w-75">
      <Show when={auth.user()}>
        {(user) => (
          <>
            <div class="flex items-center gap-[11px] px-2.5 pb-6">
              <span class="grid size-8.5 place-items-center rounded-[10px] bg-brand font-heading font-bold text-ink">S</span>
              <div>
                <strong class="font-heading text-[1.05rem] tracking-[0.14em] text-snow">SOLANUM</strong>
                <small class="mt-0.5 block text-[0.56rem] font-bold tracking-[0.22em] text-brand">STREAMING</small>
              </div>
            </div>
            <div class="mb-3 flex items-center gap-2 rounded-[9px] border border-line px-2.5 py-[9px] text-[0.83rem] text-fog 3xl:text-[0.95rem] 4xl:text-[1.2rem]">
              <span class="size-[7px] rounded-[50%] bg-brand shadow-[0_0_10px_var(--color-brand)]" />{user().username}
            </div>

            <button
              class="w-full cursor-pointer rounded-[9px] border border-line-strong bg-transparent px-2.5 py-[9px] text-left text-[0.82rem] text-fog hover:border-brand 3xl:text-[0.95rem] 4xl:text-[1.2rem]"
              disabled={submission.pending}
              aria-expanded={languageMenuOpen()}
              onClick={() => setLanguageMenuOpen((open) => !open)}
            >
              <span class="block">
                {LANGS[user().lang]?.flag ?? "🌐"}{" "}
                {LANGS[user().lang]?.label ?? "Unknown language"}
              </span>
            </button>

            <Show when={languageMenuOpen()}>
              <ul class="relative z-2 mt-1.5 grid gap-[3px] rounded-[9px] border border-line-strong bg-surface p-[5px] shadow-[0_12px_24px_rgb(0_0_0/28%)]">
                <For each={Object.entries(LANGS) as [Languages, (typeof LANGS)[Languages]][]}>
                  {([code, { flag, label }]) => (
                    <li>
                      <button
                        type="button"
                        class="w-full cursor-pointer rounded-md border-0 px-2.5 py-[9px] text-left hover:bg-surface-hover hover:text-brand"
                        classList={{
                          "bg-transparent text-fog": user().lang !== code,
                          "bg-surface-hover text-brand": user().lang === code,
                        }}
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
            <div class="mx-1.5 my-5.5 h-px bg-line-soft" />
          </>
        )}

      </Show>
      <ul class="grid gap-[5px]">
        <li>
          <A href="/" end class={navLink} inactiveClass={navInactive} activeClass={navActive}>
            <span class={navIcon}>⌂</span> Home
          </A>
        </li>
        <li>
          <A href="/collection" class={navLink} inactiveClass={navInactive} activeClass={navActive}>
            <span class={navIcon}>▣</span> Collection
          </A>
        </li>
      </ul>
    </aside>
  );
};

export default Sidebar;
