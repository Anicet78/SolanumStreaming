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
              <span class="text-2xl">{movie().title}</span>
              <Show when={movie().title !== movie().original_title}>
                <span class="text-lg">{movie().original_title}</span>
              </Show>

              <Capabilities/>
              <Show when={torrent()}>
                <video-player>
                  <video-skin style="display: block; width: 100%; aspect-ratio: 16 / 9;">
                    <hlsjs-video src={`http://localhost:8083/playlist.m3u8?torrent_link=${torrent()?.torrent_link}`} playsinline></hlsjs-video>
                  </video-skin>
                </video-player>

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