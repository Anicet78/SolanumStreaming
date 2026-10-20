import { JSX } from "solid-js";

declare module "solid-js" {
  namespace JSX {
    interface IntrinsicElements {
      "video-player": JSX.HTMLAttributes<HTMLElement> & {
        src?: string;
        poster?: string;
      };
      "video-skin": JSX.HTMLAttributes<HTMLElement>;
      "hlsjs-video": JSX.HTMLAttributes<HTMLElement>  & {
        src?: string;
        playsinline: boolean;
      };
    }
  }
}
