package types

import "time"

// RecentMatchSummary captures the essential details of a player's recent ranked match.
type RecentMatchSummary struct {
	GameID       int64  `json:"gameId"`
	GameCreation int64  `json:"gameCreation"` // epoch milliseconds
	ChampionID   int    `json:"championId"`
	ChampionName string `json:"championName"`
	ChampionIcon string `json:"championIcon"`
	Position     string `json:"position"` // top / jungle / middle / bottom / utility, or empty
	Win          bool   `json:"win"`
	Kills        int    `json:"kills"`
	Deaths       int    `json:"deaths"`
	Assists      int    `json:"assists"`
	CS           int    `json:"cs"`
	Damage       int    `json:"damage"`
	Gold         int    `json:"gold"`
	VisionScore  int    `json:"visionScore"`
	QueueID      int    `json:"queueId"`
	GameDuration int    `json:"gameDuration"`
}

// RankSummary is one ranked queue entry of a player.
type RankSummary struct {
	QueueType    string  `json:"queueType"`
	QueueName    string  `json:"queueName"`
	TierKey      string  `json:"tierKey"` // e.g. "GOLD"; empty when unranked
	Tier         string  `json:"tier"`    // localised display name
	Division     string  `json:"division"`
	LeaguePoints int     `json:"leaguePoints"`
	Wins         int     `json:"wins"`
	Losses       int     `json:"losses"`
	WinRate      float64 `json:"winRate"` // percent, 0-100
}

// PlayerStats aggregates a player's recent ranked games.
type PlayerStats struct {
	Games        int     `json:"games"`
	Wins         int     `json:"wins"`
	WinRate      float64 `json:"winRate"` // percent
	AvgKills     float64 `json:"avgKills"`
	AvgDeaths    float64 `json:"avgDeaths"`
	AvgAssists   float64 `json:"avgAssists"`
	KDA          float64 `json:"kda"`
	AvgCS        float64 `json:"avgCs"`
	AvgDamage    float64 `json:"avgDamage"`
	AvgVision    float64 `json:"avgVision"`
	Streak       int     `json:"streak"` // >0 consecutive wins, <0 consecutive losses
	ChampGames   int     `json:"champGames"`
	ChampWins    int     `json:"champWins"`
	ChampWinRate float64 `json:"champWinRate"`
	PosGames     int     `json:"posGames"` // games recently played in the assigned position
}

// Rating is the composite strength estimate used for the 大腿/上等马/... label.
type Rating struct {
	Score float64 `json:"score"` // 0-100; meaningful only when Valid
	Valid bool    `json:"valid"`
	Label string  `json:"label"` // 大腿 / 上等马 / 中等马 / 下等马, empty when not enough data
}

// TeamMemberSummary aggregates lobby data, rank and recent matches for a single player.
type TeamMemberSummary struct {
	Puuid            string `json:"puuid"`
	GameName         string `json:"gameName"`
	TagLine          string `json:"tagLine"`
	AssignedPosition string `json:"assignedPosition"`
	ChampionID       int    `json:"championId"`
	ChampionName     string `json:"championName"`
	ChampionIcon     string `json:"championIcon"`
	CellID           int    `json:"cellId"`

	SummonerLevel int      `json:"summonerLevel"`
	Spells        []string `json:"spells"`
	MasteryLevel  int      `json:"masteryLevel"`
	MasteryPoints int      `json:"masteryPoints"`

	Solo          RankSummary          `json:"solo"`
	Flex          RankSummary          `json:"flex"`
	Stats         PlayerStats          `json:"stats"`
	Rating        Rating               `json:"rating"`
	Tags          []string             `json:"tags"`
	RecentMatches []RecentMatchSummary `json:"recentMatches"`
}

// ChampSelectSnapshot stores the latest champion select data for reuse in the UI.
type ChampSelectSnapshot struct {
	QueueID   int                 `json:"queueId"`
	GameID    int64               `json:"gameId"`
	Phase     string              `json:"phase"` // ChampSelect / InProgress ...
	UpdatedAt time.Time           `json:"updatedAt"`
	Team      []TeamMemberSummary `json:"team"`
	Enemy     []TeamMemberSummary `json:"enemy"`
}
