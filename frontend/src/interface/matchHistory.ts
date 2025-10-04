export interface MatchHistory {
  accountId?: number
  games?: MatchHistoryGames
  platformId?: string
}

export interface MatchHistoryGames {
  gameBeginDate?: string
  gameCount?: number
  gameEndDate?: string
  gameIndexBegin?: number
  gameIndexEnd?: number
  games?: Game[]
}

export interface Game {
  endOfGameResult?: string
  gameDuration?: number
  gameId?: number
  gameMode?: string
  queueId?: number
  participants?: Participant[]
}

export interface Participant {
  championId?: number
  highestAchievedSeasonTier?: string
  participantId?: number
  stats?: Stats
}

export interface Stats {
  assists?: number
  deaths?: number
  kills?: number
  win?: boolean
}

export interface IRecentMatchRecord {
  id: number
  champion: ChampionSummary
  isWin: boolean
  kills: number
  deaths: number
  assists: number
  queue: string
  duration: string
}

export interface ChampionSummary {
  name: string
  icon: string
}
