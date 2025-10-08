export interface RecentMatchSummary {
  championId: number
  championName?: string
  championIcon?: string
  win: boolean
  kills: number
  deaths: number
  assists: number
  queueId: number
  gameDuration: number
}

export interface TeamMemberSummary {
  puuid: string
  gameName: string
  tagLine: string
  assignedPosition: string
  championId: number
  cellId: number
  recentMatches: RecentMatchSummary[]
}

export interface ChampSelectSnapshot {
  queueId: number
  gameId: number
  updatedAt: string
  team: TeamMemberSummary[]
}