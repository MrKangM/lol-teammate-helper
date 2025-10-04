export interface IPlayerBaseData {
  accountId?: number
  displayName?: string
  gameName?: string
  internalName?: string
  nameChangeFlag?: boolean
  percentCompleteForNextLevel?: number
  privacy?: string
  profileIconId?: number
  puuid?: string
  rerollPoints?: IReroll
  summonerId?: number
  summonerLevel?: number // Summoner level
  tagLine?: string
  unnamed?: boolean
  xpSinceLastLevel?: number // Experience already earned for the current level
  xpUntilNextLevel?: number // Remaining experience required for next level
  iconImgSrc?: string // Profile icon URL
  region?: string // Localised region label
}

interface IReroll {
  currentPoints: number
  maxRolls: number
  numberOfRolls: number
  pointsCostToRoll: number
  pointsToReroll: number
}
