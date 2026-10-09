import { GetPositionIcon, GetRankEmblem } from "../../wailsjs/go/main/App"
import { hasBridge } from "@/lib/format"

// Images come from the League client itself. Each request is cached for the
// whole session so every row asking for the same emblem shares one lookup.
const cache = new Map<string, Promise<string>>()

function cached(key: string, load: () => Promise<string>): Promise<string> {
  let hit = cache.get(key)
  if (!hit) {
    hit = hasBridge() ? load().catch(() => "") : Promise.resolve("")
    cache.set(key, hit)
  }
  return hit
}

export const loadRankEmblem = (tierKey: string) => cached(`rank:${tierKey}`, () => GetRankEmblem(tierKey))
export const loadPositionIcon = (position: string) => cached(`pos:${position}`, () => GetPositionIcon(position))

/** Forget misses so images are retried after the client (re)connects. */
export const resetAssetCache = () => cache.clear()
