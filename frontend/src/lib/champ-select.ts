import type { ChampSelectSnapshot } from "@/interface/champSelect"

export async function fetchChampSelectSnapshot(): Promise<ChampSelectSnapshot | null> {
  const api = (window as any)?.go?.main?.App?.GetCurrentChampSelectSnapshot
  if (typeof api !== "function") {
    throw new Error("GetCurrentChampSelectSnapshot bridge is unavailable")
  }

  const snapshot = await api()
  if (!snapshot || !Array.isArray(snapshot.team)) {
    return null
  }

  return snapshot as ChampSelectSnapshot
}