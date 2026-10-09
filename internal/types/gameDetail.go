package types

// GameDetailPlayer is one row of the end-of-game scoreboard.
type GameDetailPlayer struct {
	Puuid        string      `json:"puuid"`
	Name         string      `json:"name"`
	ChampionID   int         `json:"championId"`
	ChampionName string      `json:"championName"`
	ChampionIcon string      `json:"championIcon"`
	Level        int         `json:"level"`
	Kills        int         `json:"kills"`
	Deaths       int         `json:"deaths"`
	Assists      int         `json:"assists"`
	CS           int         `json:"cs"`
	Gold         int         `json:"gold"`
	Damage       int         `json:"damage"`
	DamageTaken  int         `json:"damageTaken"`
	VisionScore  int         `json:"visionScore"`
	Spells       []SpellInfo `json:"spells"`
	Items        []string    `json:"items"` // data URIs; empty string for an empty slot
	IsTarget     bool        `json:"isTarget"`
}

// GameDetailTeam is one side of the scoreboard.
type GameDetailTeam struct {
	TeamID  int                `json:"teamId"`
	Win     bool               `json:"win"`
	Kills   int                `json:"kills"`
	Gold    int                `json:"gold"`
	Players []GameDetailPlayer `json:"players"`
}

// GameDetail is the full end-of-game picture of a single match.
type GameDetail struct {
	GameID       int64            `json:"gameId"`
	GameCreation int64            `json:"gameCreation"`
	GameDuration int              `json:"gameDuration"`
	QueueID      int              `json:"queueId"`
	Teams        []GameDetailTeam `json:"teams"`
}
