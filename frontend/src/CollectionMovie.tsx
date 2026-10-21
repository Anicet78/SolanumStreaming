import { useParams } from "@solidjs/router";
import { createResource, ErrorBoundary, Show } from "solid-js";
import { moviesApi } from "./api/movies";
import { getProfileParam } from "./api/capabilities";
import Capabilities from "./Capabilities";
import '@videojs/html/video/player';
import '@videojs/html/video/skin';
import '@videojs/html/media/hlsjs-video';

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
    <ErrorBoundary fallback={(err, reset) => (
      <div>
        <p>Oups : {err.message}</p>
        <button onClick={reset}>Réessayer</button>
      </div>
    )}>
      <Show when={id()} fallback={<span>Not Found</span>}>
        <Show when={movieInfos()} fallback={<span>Loading...</span>}>
          {(movie) => (
            <div class="min-h-screen flex flex-col items-center p-4">
              <Show when={torrent()}>
                <video-player>
                  <video-skin style="display: block; width: 100%; aspect-ratio: 2 / 1;">
                    {/* <hlsjs-video src={`http://localhost:8083/playlist.m3u8?torrent_link=${encodeURIComponent(torrent()!.torrent_link)}`} playsinline></hlsjs-video> */}
                  </video-skin>
                </video-player>

                <div class="w-full m-4 pl-3">
                  <div class="flex items-center justify-evenly">
                    <div class="flex flex-col w-fit">
                      <span class="font-black text-2xl">{movie().title}</span>
                      <Show when={movie().title !== movie().original_title}>
                        <span class="self-center text-sm">{movie().original_title}</span>
                      </Show>
                    </div>
                    <span class="font-black text-2xl">📅 {movie().release_date.slice(0, 4)}</span>
                    <span class="font-black text-2xl">⏳ {Math.floor(movie().runtime / 60)}h {movie().runtime % 60}m</span>
                    <span class="font-black text-2xl">⭐ {movie().vote_average}</span>
                  </div>
                  <div class="mt-5 flex flex-col">
                    <span class=" font-black text-xl mb-3">Overview</span>
                    <span>{movie().overview}</span>
                  </div>
                </div>

                <script type="module" src="/src/player.ts"></script>
              </Show>
            </div>
          )}
        </Show>
      </Show>
    </ErrorBoundary>
  )
}

export default CollectionMovie
