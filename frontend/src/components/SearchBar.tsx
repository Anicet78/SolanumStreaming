import { useNavigate } from "@solidjs/router";
import { createSignal } from "solid-js";

const SearchBar = () => {
  const navigate = useNavigate();
  const [query, setQuery] = createSignal("");

  const handleSubmit = (e: SubmitEvent) => {
    e.preventDefault();
    const trimmed = query().trim();
    if (!trimmed) return;

    const params = new URLSearchParams({ title: trimmed, page: "1" });
    navigate(`/search?${params.toString()}`);
  };

  return (
    <form onSubmit={handleSubmit}>
      <label class="input input-primary">
        <svg class="h-[1em] opacity-50" xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24">
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
          placeholder="Search"
          value={query()}
          onInput={(e) => setQuery(e.currentTarget.value)}
        />
      </label>
    </form>
  )
}

export default SearchBar