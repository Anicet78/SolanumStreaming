import { useNavigate, useSearchParams } from "@solidjs/router";
import { createSignal } from "solid-js";

const SearchBar = () => {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const [query, setQuery] = createSignal("");

  const handleSubmit = (e: SubmitEvent) => {
    e.preventDefault();
    const trimmed = query().trim();
    if (!trimmed) return;

    const params = new URLSearchParams({ title: trimmed, page: "1" });
    navigate(`/search?${params.toString()}`);
  };

  return (
    <form
      class="flex items-center gap-2.5 rounded-xl border border-line-strong bg-surface py-[0.3125rem] pr-[0.4375rem] pl-4 shadow-[inset_0_1px_0_rgb(255_255_255/4%)] focus-within:border-brand focus-within:shadow-[0_0_0_0.1875rem_rgb(103_211_196/12%)] pointer-coarse:min-h-11"
      onSubmit={handleSubmit}
    >
        <svg class="w-4.5 flex-none text-brand" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" aria-hidden="true">
          <g
          stroke-linejoin="round"
          stroke-linecap="round"
          stroke-width="2.5"
          fill="none"
          stroke="currentColor"
          >
          <circle cx="11" cy="11" r="8"></circle>
          <path d="m21 21-4.3-4.3"></path>
          </g>
        </svg>
        <input
          type="search"
          required
          class="w-full min-w-0 self-stretch border-0 bg-transparent text-snow outline-none placeholder:text-[#71859a]"
          placeholder="Search for a movie"
          aria-label="Search for a movie"
          value={searchParams.title ?? query()}
          onInput={(e) => setQuery(e.currentTarget.value)}
        />
        <button class="cursor-pointer rounded-lg border-0 bg-brand px-3.5 py-[0.5625rem] text-[0.82rem] font-bold text-ink mobile:hidden" type="submit">Search</button>
    </form>
  )
}

export default SearchBar
