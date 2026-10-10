package types

// GameflowSession is the subset of /lol-gameflow/v1/session we need to find
// the opposing team once the game is loading or running.
type GameflowSession struct {
	Phase    string       `json:"phase"`
	GameData GameflowData `json:"gameData"`
}

// GameflowData holds the teams of the current game.
type GameflowData struct {
	GameID int64 `json:"gameId"`
	Queue  struct {
		ID int `json:"id"`
	} `json:"queue"`
	TeamOne []GameflowPlayer `json:"teamOne"`
	TeamTwo []GameflowPlayer `json:"teamTwo"`
}

// GameflowPlayer is one player entry in the gameflow teams.
type GameflowPlayer struct {
	Puuid            string `json:"puuid"`
	GameName         string `json:"gameName"`
	TagLine          string `json:"tagLine"`
	SummonerName     string `json:"summonerName"`
	ChampionID       int    `json:"championId"`
	SelectedPosition string `json:"selectedPosition"`
	Spell1ID         int    `json:"spell1Id"`
	Spell2ID         int    `json:"spell2Id"`
}
