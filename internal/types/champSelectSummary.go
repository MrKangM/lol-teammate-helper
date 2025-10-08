package types

import "time"

// RecentMatchSummary captures the essential details of a teammate\'s recent ranked match.
type RecentMatchSummary struct {
	ChampionID   int    `json:"championId"`
	ChampionName string `json:"championName"`
	ChampionIcon string `json:"championIcon"`
	Win          bool   `json:"win"`
	Kills        int    `json:"kills"`
	Deaths       int    `json:"deaths"`
	Assists      int    `json:"assists"`
	QueueID      int    `json:"queueId"`
	GameDuration int    `json:"gameDuration"`
}

// TeamMemberSummary aggregates lobby data and recent matches for a single teammate.
type TeamMemberSummary struct {
	Puuid            string               `json:"puuid"`
	GameName         string               `json:"gameName"`
	TagLine          string               `json:"tagLine"`
	AssignedPosition string               `json:"assignedPosition"`
	ChampionID       int                  `json:"championId"`
	ChampionName     string               `json:"championName"`
	ChampionIcon     string               `json:"championIcon"`
	CellID           int                  `json:"cellId"`
	RecentMatches    []RecentMatchSummary `json:"recentMatches"`
	SelectChampIcon  string               `json:"selectChampIcon"`
}

// ChampSelectSnapshot stores the latest champion select data for reuse in the UI.
type ChampSelectSnapshot struct {
	QueueID   int                 `json:"queueId"`
	GameID    int64               `json:"gameId"`
	UpdatedAt time.Time           `json:"updatedAt"`
	Team      []TeamMemberSummary `json:"team"`
}
