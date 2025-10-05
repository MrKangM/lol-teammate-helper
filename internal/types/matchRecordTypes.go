package types

// MatchHistory represents the match history payload returned by the Riot client.
type MatchHistory struct {
	AccountID  int64             `json:"accountId"`
	Games      MatchHistoryGames `json:"games"`
	PlatformID string            `json:"platformId"`
}

// MatchHistoryGames wraps the metadata and match list for a history response.
type MatchHistoryGames struct {
	GameBeginDate  string `json:"gameBeginDate"`
	GameCount      int    `json:"gameCount"`
	GameEndDate    string `json:"gameEndDate"`
	GameIndexBegin int    `json:"gameIndexBegin"`
	GameIndexEnd   int    `json:"gameIndexEnd"`
	Games          []Game `json:"games"`
}

// Game captures the subset of match data required by the application.
type Game struct {
	EndOfGameResult string        `json:"endOfGameResult"`
	GameDuration    int           `json:"gameDuration"`
	QueueID         int           `json:"queueId"`
	GameMode        string        `json:"gameMode"`
	Participants    []Participant `json:"participants"`
}

// Participant stores the individual player stats within a match.
type Participant struct {
	ChampionID int   `json:"championId"`
	Stats      Stats `json:"stats"`
}

// Stats contains the combat summary used by the frontend.
type Stats struct {
	Win     bool `json:"win"`
	Kills   int  `json:"kills"`
	Deaths  int  `json:"deaths"`
	Assists int  `json:"assists"`
}
