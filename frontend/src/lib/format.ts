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

/** Text colour and badge style per rating label. */
export const RATING_STYLES: Record<string, { badge: string; strip: string }> = {
  大腿: { badge: "bg-horse-top/15 text-horse-top border-horse-top/50", strip: "bg-horse-top" },
  上等马: { badge: "bg-horse-good/15 text-horse-good border-horse-good/50", strip: "bg-horse-good" },
  中等马: { badge: "bg-horse-mid/15 text-horse-mid border-horse-mid/50", strip: "bg-horse-mid" },
  下等马: { badge: "bg-horse-low/15 text-horse-low border-horse-low/50", strip: "bg-horse-low" },
}

export const PHASE_LABELS: Record<string, string> = {
  None: "空闲",
  Lobby: "房间中",
  Matchmaking: "匹配中",
  ReadyCheck: "等待接受",
  ChampSelect: "选人中",
  GameStart: "加载中",
  InProgress: "游戏中",
  Reconnect: "重新连接",
  WaitingForStats: "结算中",
  PreEndOfGame: "结算中",
  EndOfGame: "对局结束",
}

export const hasBridge = () => typeof (window as any)?.go?.main?.App?.GetDiagnostics === "function"
