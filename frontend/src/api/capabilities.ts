export type Support = { supported: boolean; smooth: boolean; powerEfficient: boolean };

export const TIERS = {
  sd:  { label: "480p",        width: 854,  height: 480 },
  hd:  { label: "720p",        width: 1280, height: 720 },
  fhd: { label: "1080p",       width: 1920, height: 1080 },
  qhd: { label: "1440p (QHD)", width: 2560, height: 1440 },
  uhd: { label: "4K UHD",      width: 3840, height: 2160 },
  dci: { label: "4K DCI",      width: 4096, height: 2160 },
} as const;

export type TierName = keyof typeof TIERS;
const TIER_ORDER = Object.keys(TIERS) as TierName[];

export type Detected = {
  tier: TierName;
  hdr: boolean;
  wideGamut: boolean;
};

export type DisplayOverride = {
  tier?: TierName;
  hdr?: boolean;
};

export type Capabilities = {
  version: number;
  detectedAt: number;
  detected: Detected;
  video: Record<string, Partial<Record<TierName, Support>>>;
  audio: Record<string, boolean>;
};

const VERSION = 2;
const MAX_AGE = 7 * 24 * 3600 * 1000;
const KEY_CAPS = "client-capabilities";
const KEY_OVERRIDE = "client-display-override";

const VIDEO_CODECS: Record<string, string> = {
  h264: "avc1.640033",
  hevc: "hvc1.2.4.L153.B0",
  av1: "av01.0.12M.10",
  vp9: "vp09.02.51.10",
};

const AUDIO_CODECS: Record<string, string> = {
  aac: "mp4a.40.2",
  ac3: "ac-3",
  eac3: "ec-3",
  opus: "opus",
  flac: "flac",
};

const TEST_BITRATE: Record<TierName, number> = {
  sd: 3_000_000,
  hd: 6_000_000,
  fhd: 10_000_000,
  qhd: 15_000_000,
  uhd: 20_000_000,
  dci: 22_000_000,
};

const NO: Support = { supported: false, smooth: false, powerEfficient: false };

function guessTier(): TierName {
  const dpr = window.devicePixelRatio || 1;
  const longSide = Math.max(screen.width, screen.height) * dpr;

  if (longSide >= 3968) return "dci";
  if (longSide >= 3500) return "uhd";
  if (longSide >= 2300) return "qhd";
  if (longSide >= 1700) return "fhd";
  if (longSide >= 1100) return "hd";
  return "sd";
}

async function checkVideo(codec: string, tier: TierName): Promise<Support> {
  const { width, height } = TIERS[tier];
  try {
    const r = await navigator.mediaCapabilities.decodingInfo({
      type: "media-source",
      video: {
        contentType: `video/mp4; codecs="${codec}"`,
        width,
        height,
        bitrate: TEST_BITRATE[tier],
        framerate: 24,
      },
    });
    return { supported: r.supported, smooth: r.smooth, powerEfficient: r.powerEfficient };
  } catch {
    return NO;
  }
}

async function checkAudio(codec: string): Promise<boolean> {
  try {
    const r = await navigator.mediaCapabilities.decodingInfo({
      type: "media-source",
      audio: { contentType: `audio/mp4; codecs="${codec}"`, channels: "2", bitrate: 192_000, samplerate: 48_000 },
    });
    return r.supported;
  } catch {
    return false;
  }
}

export async function detectCapabilities(): Promise<Capabilities> {
  const video: Capabilities["video"] = {};

  await Promise.all(
    Object.entries(VIDEO_CODECS).map(async ([name, codec]) => {
      const entries = await Promise.all(
        TIER_ORDER.map(async (t) => [t, await checkVideo(codec, t)] as const),
      );
      video[name] = Object.fromEntries(entries);
    }),
  );

  const audio: Capabilities["audio"] = {};
  await Promise.all(
    Object.entries(AUDIO_CODECS).map(async ([name, codec]) => {
      audio[name] = await checkAudio(codec);
    }),
  );

  return {
    version: VERSION,
    detectedAt: Date.now(),
    detected: {
      tier: guessTier(),
      hdr: matchMedia("(dynamic-range: high)").matches,
      wideGamut: matchMedia("(color-gamut: p3)").matches,
    },
    video,
    audio,
  };
}

function loadCaps(): Capabilities | null {
  try {
    const raw = localStorage.getItem(KEY_CAPS);
    if (!raw) return null;
    const c = JSON.parse(raw) as Capabilities;
    if (c.version !== VERSION || Date.now() - c.detectedAt > MAX_AGE) return null;
    return c;
  } catch {
    return null;
  }
}

function loadOverride(): DisplayOverride {
  try {
    const raw = localStorage.getItem(KEY_OVERRIDE);
    if (!raw) return {};
    const o = JSON.parse(raw) as DisplayOverride;
    if (o.tier && !(o.tier in TIERS)) delete o.tier;
    return o;
  } catch {
    return {};
  }
}

let caps: Capabilities | null = loadCaps();
let override: DisplayOverride = loadOverride();
const listeners = new Set<() => void>();
const notify = () => listeners.forEach((fn) => fn());

export function subscribeCapabilities(fn: () => void): () => void {
  listeners.add(fn);
  return () => listeners.delete(fn);
}

export function getCapabilities(): Capabilities | null {
  return caps;
}

export function getOverride(): DisplayOverride {
  return override;
}

export async function ensureCapabilities(force = false): Promise<Capabilities> {
  if (caps && !force) return caps;
  caps = await detectCapabilities();
  try {
    localStorage.setItem(KEY_CAPS, JSON.stringify(caps));
  } catch {}
  notify();
  return caps;
}

export function setOverride(o: DisplayOverride): void {
  override = { ...override, ...o };
  try {
    localStorage.setItem(KEY_OVERRIDE, JSON.stringify(override));
  } catch {}
  notify();
}

export function clearOverride(): void {
  override = {};
  try {
    localStorage.removeItem(KEY_OVERRIDE);
  } catch {}
  notify();
}

export function getDisplay(): { tier: TierName; width: number; height: number; hdr: boolean } | null {
  if (!caps) return null;
  const tier = override.tier ?? caps.detected.tier;
  return {
    tier,
    width: TIERS[tier].width,
    height: TIERS[tier].height,
    hdr: override.hdr ?? caps.detected.hdr,
  };
}

export type ClientProfile = {
  display: { tier: TierName; width: number; height: number; hdr: boolean };
  video: Record<string, { maxTier: TierName | null; smoothMaxTier: TierName | null; hw: boolean }>;
  audio: string[];
};

export function buildProfile(): ClientProfile | null {
  const display = getDisplay();
  if (!caps || !display) return null;

  const video: ClientProfile["video"] = {};
  for (const [codec, tiers] of Object.entries(caps.video)) {
    let maxTier: TierName | null = null;
    let smoothMaxTier: TierName | null = null;
    let hw = false;
    for (const t of TIER_ORDER) {
      const s = tiers[t];
      if (s?.supported) maxTier = t;
      if (s?.supported && s.smooth) {
        smoothMaxTier = t;
        hw = s.powerEfficient;
      }
    }
    video[codec] = { maxTier, smoothMaxTier, hw };
  }

  return {
    display,
    video,
    audio: Object.entries(caps.audio).filter(([, ok]) => ok).map(([n]) => n),
  };
}

export function profileToParam(p: ClientProfile): string {
  return encodeURIComponent(btoa(JSON.stringify(p)));
}
