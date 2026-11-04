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
      class="flex items-center gap-2.5 rounded-xl border border-line-strong bg-surface py-[5px] pr-[7px] pl-4 shadow-[inset_0_1px_0_rgb(255_255_255/4%)] focus-within:border-brand focus-within:shadow-[0_0_0_3px_rgb(103_211_196/12%)] 4xl:rounded-2xl 4xl:py-2.5 4xl:pr-3 4xl:pl-5.5"
      onSubmit={handleSubmit}
    >
        <svg class="w-4.5 flex-none text-brand" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
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
          class="w-full min-w-0 border-0 bg-transparent text-snow outline-none placeholder:text-[#71859a] 4xl:text-[1.2rem]"
          placeholder="Search for a movie"
          value={searchParams.title ?? query()}
          onInput={(e) => setQuery(e.currentTarget.value)}
        />
        <button class="cursor-pointer rounded-lg border-0 bg-brand px-3.5 py-[9px] text-[0.82rem] font-bold text-ink mobile:hidden 4xl:px-5.5 4xl:py-3.5 4xl:text-[1rem]" type="submit">Search</button>
    </form>
  )
}

export default SearchBar
