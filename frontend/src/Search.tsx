import { createResource, For, Show } from "solid-js"
import MovieCard from "./components/MovieCard"
import MovieGrid from "./components/MovieGrid"
import Page from "./components/Page"
import PageHeading from "./components/PageHeading"
import { moviesApi } from "./api/movies"
import { A, useNavigate, useSearchParams } from "@solidjs/router";

function toSingleString(v: string | string[] | undefined): string | undefined {
  return Array.isArray(v) ? v[0] : v;
}

const pageLink =
  "grid size-9 place-items-center rounded-[9px] border text-[0.86rem] transition-[background,border-color,color] duration-200 ease-[ease]";
const pageLinkIdle = "border-line-strong text-fog hover:border-brand hover:bg-brand hover:text-ink";
const pageLinkCurrent = "border-brand bg-brand text-ink";

const Search = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();

  if (!searchParams.title || searchParams.title == "")
    navigate("/");

  const [results] = createResource(
    () => ({
      title: toSingleString(searchParams.title),
      page: toSingleString(searchParams.page) ? Number(toSingleString(searchParams.page)) : undefined,
    }),
    moviesApi.search
  );

  return (
    <Page>
        <PageHeading eyebrow="Library search" title="Search results" subtitle="Explore titles" />
        <Show when={results.loading}>
          <MovieGrid loading>
            <For each={Array.from({ length: 18 })}>
              {() => <div class="aspect-2/3 animate-shimmer rounded-[14px] bg-[linear-gradient(110deg,#16232e_30%,#223543_45%,#16232e_60%)] bg-size-[200%_100%]"></div>}
            </For>
          </MovieGrid>
        </Show>

        <Show when={results.error || (results() && results()?.total_results == 0)}>
          <div class="rounded-2xl border border-dashed border-line-strong px-6 py-18 text-center text-dim">
            <strong>Nothing found</strong><br />
            <span>{results.error?.message ?? "No result has been found"}</span>
          </div>
        </Show>

        <Show when={!results.loading && !results.error}>
          <MovieGrid>
            <For each={results()?.results}>
              {(movie) =>
                <A href={`/${movie.id}`}>
                  <MovieCard {...movie} />
                </A>
              }
            </For>
          </MovieGrid>
          <Show when={results()?.total_pages! > 1}>
            <div class="mt-8.5 flex flex-wrap justify-center gap-[7px]">
              <For each={Array.from({ length: results()?.total_pages!})}>
                {(_, i) =>
                  <Show when={i() == 0 || i() + 1 == results()?.total_pages || (results()?.page! - i() < 4 && results()?.page! - i() >= 0) || (i() - results()?.page! < 4 && i() - results()?.page! >= 0)}>
                    <A
                      class={pageLink}
                      classList={{
                        [pageLinkIdle]: i() + 1 !== results()?.page,
                        [pageLinkCurrent]: i() + 1 === results()?.page,
                      }}
                      href={`/search?title=${searchParams.title}&page=${i() + 1}`}
                    >
                      {i() + 1}
                    </A>
                  </Show>
                }
              </For>
            </div>
          </Show>
        </Show>
    </Page>
  );
};

export default Search
