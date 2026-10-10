import { ref } from "vue"
import { GetCurrentChampSelectSnapshot, GetCurrentSummoner, GetDiagnostics, GetImgSrc, GetMyCareer } from "../../wailsjs/go/main/App"
import { diag, types } from "../../wailsjs/go/models"
import { EventsOn } from "../../wailsjs/runtime"
import { resetAssetCache } from "@/lib/assets"
import { hasBridge } from "@/lib/format"

export const status = ref<diag.Snapshot | null>(null)
export const snapshot = ref<types.ChampSelectSnapshot | null>(null)
export const summoner = ref<types.IPlayerBaseData | null>(null)
export const summonerIcon = ref("")
export const career = ref<types.TeamMemberSummary | null>(null)
export const careerLoading = ref(false)
export const careerError = ref("")

let started = false

async function refreshStatus() {
  try {
    const next = await GetDiagnostics()
    const wasConnected = status.value?.connected ?? false
    status.value = next
    if (next.connected && !wasConnected) {
      resetAssetCache()
      await loadSummoner()
    }
    if (!next.connected && wasConnected) {
      summoner.value = null
    }
  } catch (error) {
    console.warn("[store] failed to read diagnostics", error)
  }
}

async function loadSummoner() {
  try {
    const me = await GetCurrentSummoner()
    if (!me?.puuid) return
    summoner.value = me
    summonerIcon.value = (await GetImgSrc(me.profileIconId ?? 0)) || ""
  } catch (error) {
    console.warn("[store] failed to load summoner", error)
  }
}

export async function loadCareer() {
  if (careerLoading.value) return
  careerLoading.value = true
  careerError.value = ""
  try {
    const me = await GetMyCareer()
    if (!me?.puuid) {
      careerError.value = "还没有取到生涯数据，请确认客户端已登录。"
      return
    }
    career.value = me
  } catch (error) {
    console.error("[store] failed to load career", error)
    careerError.value = "生涯数据加载失败，详情见“诊断”页或日志文件。"
  } finally {
    careerLoading.value = false
  }
}

/** Starts listening for backend events and polling connection state (idempotent). */
export function startStore() {
  if (started || !hasBridge()) return
  started = true

  EventsOn("champ-select:snapshot", (payload: unknown) => {
    snapshot.value = types.ChampSelectSnapshot.createFrom(payload)
  })
  EventsOn("champ-select:ended", () => {
    snapshot.value = null
  })

  GetCurrentChampSelectSnapshot()
    .then((current) => {
      if (current?.team?.length && !snapshot.value) snapshot.value = current
    })
    .catch(() => undefined)

  refreshStatus()
  window.setInterval(refreshStatus, 2500)
}
