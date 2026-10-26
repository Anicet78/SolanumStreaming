import { A, useParams } from "@solidjs/router";
import { createResource, Show } from "solid-js";
import { moviesApi } from "./api/movies";
import { getProfileParam } from "./api/capabilities";
import '@videojs/html/video/player';
import '@videojs/html/video/skin';
import '@videojs/html/media/hlsjs-video';
import SearchBar from "./components/SearchBar";

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
    <>
      <div class="flex items-center justify-between ml-4 mt-4">
        <A href={`/collection`}>
          <div class="btn btn-dash btn-primary w-xs">⮐ Back to Collection</div>
        </A>
        <SearchBar/>
        <div class="w-xs"></div>
      </div>

      <Show when={id()} fallback={<span>Not Found</span>}>
        <Show when={movieInfos()} fallback={<span>Loading...</span>}>
          {(movie) => (
            <div class="min-h-screen p-4 flex flex-col items-center">
              <Show when={torrent.state === "ready"}>
                <video-player>
                  <video-skin style="display: block; width: 100%; aspect-ratio: 2 / 1;">
                    <hlsjs-video src={`http://localhost:8083/playlist.m3u8?torrent_link=${encodeURIComponent(torrent()!.torrent_link)}`} playsinline></hlsjs-video>
                  </video-skin>
                </video-player>
              </Show>

              <div class="w-full p-4 pl-3">
                <div class="flex items-center justify-evenly">
                  <div class="flex flex-col w-fit">
                    <span class="font-black text-2xl">{movie().title}</span>
                    <Show when={movie().title !== movie().original_title}>
                      <span class="self-center text-sm text-gray-300">{movie().original_title}</span>
                    </Show>
                  </div>
                  <div class="flex flex-col items-center p-2">
                    <span class="text-2xl">🗓️ </span>
                    <span class="font-black text-2xl">{movie().release_date.slice(0, 4)}</span>
                  </div>
                  <div class="flex flex-col items-center p-2">
                    <span class="text-2xl">⏳ </span>
                    <span class="font-black text-2xl">{Math.floor(movie().runtime / 60)}h {movie().runtime % 60}m</span>
                  </div>
                  <div class="flex flex-col items-center p-2">
                    <span class="text-2xl">⭐ </span>
                    <span class="font-black text-2xl">{movie().vote_average}</span>
                  </div>
                </div>
                <Show when={movie().overview !== "" || movie().tagline !== ""}>
                  <div class="mt-5 flex flex-col">
                    <span class="text-gray-400">{movie().tagline.toUpperCase()}</span>
                    <Show when={movie().overview !== "" && movie().tagline !== ""}><div class="mb-1"/></Show>
                    <span>{movie().overview}</span>
                  </div>
                </Show>
              </div>
              <Show when={torrent.error}>
                <div class="mt-10 flex justify-center items-center h-full">
                  <span class="text-2xl font-serif">
                    Oups, ce film n'a pas de torrent disponible
                  </span>
                </div>
              </Show>

              <script type="module" src="/src/player.ts"></script>
            </div>
          )}
        </Show>
      </Show>
    </>
  )
}

export default CollectionMovie
