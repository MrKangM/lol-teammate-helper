const POSITION_LABELS: Record<string, string> = {
  top: "上单",
  jungle: "打野",
  middle: "中单",
  bottom: "ADC",
  utility: "辅助",
}

export const positionLabel = (position?: string) => {
  if (!position) return "待定"
  return POSITION_LABELS[position.toLowerCase()] ?? position
}

const QUEUE_LABELS: Record<number, string> = {
  420: "单双排",
  440: "灵活组排",
  400: "匹配",
  430: "匹配",
  450: "大乱斗",
  700: "冠军杯赛",
}

export const queueLabel = (queueId?: number) => {
  if (!queueId) return "未知模式"
  return QUEUE_LABELS[queueId] ?? `队列 ${queueId}`
}

export const formatDuration = (seconds?: number) => {
  if (!seconds || seconds <= 0) return "--:--"
  const m = Math.floor(seconds / 60)
  const s = Math.floor(seconds % 60)
  return `${String(m).padStart(2, "0")}:${String(s).padStart(2, "0")}`
}

export const formatKda = (k?: number, d?: number, a?: number) => `${k ?? 0} / ${d ?? 0} / ${a ?? 0}`

/** (K+A)/D with deaths floored at 1, one decimal. */
export const kdaRatio = (k: number, d: number, a: number) => ((k + a) / Math.max(d, 1)).toFixed(1)

export const formatPercent = (value?: number) => `${Math.round(value ?? 0)}%`

/** 12345 -> "12.3k" */
export const compactNumber = (n?: number) => {
  if (!n) return "0"
  return n >= 1000 ? `${(n / 1000).toFixed(1)}k` : String(Math.round(n))
}

export const timeAgo = (epochMs?: number, now: number = Date.now()) => {
  if (!epochMs) return ""
  const minutes = Math.floor((now - epochMs) / 60000)
  if (minutes < 1) return "刚刚"
  if (minutes < 60) return `${minutes}分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days}天前`
  return `${Math.floor(days / 30)}个月前`
}

export const RATING_STYLES: Record<string, string> = {
  大腿: "bg-amber-500 text-white",
  上等马: "bg-emerald-500 text-white",
  中等马: "bg-sky-500/80 text-white",
  下等马: "bg-rose-500/80 text-white",
}

export const PHASE_LABELS: Record<string, string> = {
  ChampSelect: "选人中",
  GameStart: "加载中",
  InProgress: "游戏中",
  Reconnect: "游戏中",
}
