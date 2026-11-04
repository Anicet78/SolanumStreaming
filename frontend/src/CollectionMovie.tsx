import { useParams } from "@solidjs/router";
import { createResource, Match, Show, Switch } from "solid-js";
import { moviesApi } from "./api/movies";
import { getProfileParam } from "./api/capabilities";
import '@videojs/html/video/player';
import '@videojs/html/video/skin';
import '@videojs/html/media/hlsjs-video';
import DetailPage from "./components/DetailPage";
import DetailBackbar from "./components/DetailBackbar";
import MovieDetails from "./components/MovieDetails";

const CollectionMovie = () => {
  const params = useParams()
  const id = () => Number(params.movie_id)

  const [profile] = createResource(getProfileParam)
  const [movieInfos] = createResource(id, moviesApi.getMovieInfos)
  const [torrent] = createResource(
    () => {
      const caps = profile()
      return id() && caps ? { id: id(), caps } : false
    },
    ({ id, caps }) => moviesApi.searchTorrent(id, caps)
  )

  return (
    <DetailPage>
      <Show when={id()} fallback={<span>Not Found</span>}>
        <Show when={movieInfos()} fallback={
          <Show when={torrent.loading}>
            <div class="flex justify-center font-bold">
              <span class="mr-2">Loading</span>
              <span class="loading loading-infinity loading-md"></span>
            </div>
          </Show>
        }>
          {(movie) => (
            <div class="mx-auto w-full 3xl:w-[min(100%,2200px)] 4xl:w-[min(100%,3000px)]">
              <DetailBackbar href="/collection" label="← Back to collection" />
              <Show when={torrent.state === "ready"}>
                <div class="overflow-hidden bg-[#05090d] shadow-[0_20px_50px_rgb(0_0_0/28%)] [&_hlsjs-video]:block [&_hlsjs-video]:h-full [&_hlsjs-video]:w-full [&_video-player]:block [&_video-player]:w-full [&_video-skin]:block [&_video-skin]:aspect-[2/1] [&_video-skin]:w-full 4xl:[&_video-skin]:aspect-video">
                  <video-player>
                    <video-skin>
                      <hlsjs-video src={`http://localhost:8083/playlist.m3u8?torrent_link=${encodeURIComponent(torrent()!.torrent_link)}`} playsinline></hlsjs-video>
                    </video-skin>
                  </video-player>
                </div>
              </Show>
              <div class="mt-6 rounded-2xl border border-line bg-surface p-[clamp(20px,4vw,34px)] 4xl:p-14.5">
                <MovieDetails movie={movie()} eyebrow="Now playing" compact />
              </div>
              <Show when={torrent.state === "errored"}>
                <div class="mt-4.5 rounded-[10px] border border-[#6b4c56] p-4 text-[#f2b8c5] 4xl:text-[1.55rem]">
                  <Switch fallback={<span>An error has occured</span>}>
                    <Match when={torrent.error.message === "No torrent found"}>
                      No stream is available for this movie yet.
                    </Match>
                    <Match when={torrent.error.message === "Movie not in collection"}>
                      This movie is not in your collection
                    </Match>
                  </Switch>
                </div>
              </Show>

              <script type="module" src="/src/player.ts"></script>
            </div>
          )}
        </Show>
      </Show>
    </DetailPage>
  )
}

export default CollectionMovie
