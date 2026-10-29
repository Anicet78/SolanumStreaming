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
    <main class="detail-page">
      <Show when={id()} fallback={<span>Not Found</span>}>
        <Show when={movieInfos()} fallback={<span>Loading...</span>}>
          {(movie) => (
            <div class="player-page">
              <div class="detail-backbar">
                <A class="detail-back" href="/collection">← Back to collection</A>
                <div class="detail-search"><SearchBar /></div>
              </div>
              <Show when={torrent.state === "ready"}>
                <div class="player-frame">
                  <video-player>
                    <video-skin>
                      <hlsjs-video src={`http://localhost:8083/playlist.m3u8?torrent_link=${encodeURIComponent(torrent()!.torrent_link)}`} playsinline></hlsjs-video>
                    </video-skin>
                  </video-player>
                </div>
              </Show>
              <div class="player-info">
                <div class="eyebrow">Now playing</div>
                <h1 class="detail-title">{movie().title}</h1>
                <Show when={movie().title !== movie().original_title}>
                  <div class="detail-original">Original title: {movie().original_title}</div>
                </Show>
                <div class="detail-stats">
                  <span class="detail-stat">{movie().release_date.slice(0, 4)}</span>
                  <span class="detail-stat">{Math.floor(movie().runtime / 60)}h {movie().runtime % 60}m</span>
                  <span class="detail-stat">★ {movie().vote_average.toFixed(1)} / 10</span>
                </div>
                <Show when={movie().tagline || movie().overview}>
                  <Show when={movie().tagline}><div class="detail-tagline">{movie().tagline}</div></Show>
                  <Show when={movie().overview}><p class="detail-overview">{movie().overview}</p></Show>
                </Show>
              </div>
              <Show when={torrent.error}>
                <div class="player-error">
                  No stream is available for this movie yet.
                </div>
              </Show>

              <script type="module" src="/src/player.ts"></script>
            </div>
          )}
        </Show>
      </Show>
    </main>
  )
}

export default CollectionMovie
