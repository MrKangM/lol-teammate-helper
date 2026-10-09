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
// The history list only contains the queried player as a participant; the
// single-game endpoint returns all ten participants and their identities.
type Game struct {
	GameID                int64                 `json:"gameId"`
	GameCreation          int64                 `json:"gameCreation"` // epoch milliseconds
	EndOfGameResult       string                `json:"endOfGameResult"`
	GameDuration          int                   `json:"gameDuration"`
	QueueID               int                   `json:"queueId"`
	GameMode              string                `json:"gameMode"`
	Participants          []Participant         `json:"participants"`
	ParticipantIdentities []ParticipantIdentity `json:"participantIdentities"`
	Teams                 []GameTeam            `json:"teams"`
}

// Participant stores the individual player stats within a match.
type Participant struct {
	ParticipantID int      `json:"participantId"`
	TeamID        int      `json:"teamId"`
	ChampionID    int      `json:"championId"`
	Spell1ID      int      `json:"spell1Id"`
	Spell2ID      int      `json:"spell2Id"`
	Stats         Stats    `json:"stats"`
	Timeline      Timeline `json:"timeline"`
}

// Timeline carries the lane assignment of a participant.
type Timeline struct {
	Lane string `json:"lane"`
	Role string `json:"role"`
}

// Stats contains the combat summary used by the frontend.
type Stats struct {
	Win                         bool `json:"win"`
	Kills                       int  `json:"kills"`
	Deaths                      int  `json:"deaths"`
	Assists                     int  `json:"assists"`
	ChampLevel                  int  `json:"champLevel"`
	GoldEarned                  int  `json:"goldEarned"`
	TotalDamageDealtToChampions int  `json:"totalDamageDealtToChampions"`
	TotalDamageTaken            int  `json:"totalDamageTaken"`
	TotalMinionsKilled          int  `json:"totalMinionsKilled"`
	NeutralMinionsKilled        int  `json:"neutralMinionsKilled"`
	VisionScore                 int  `json:"visionScore"`
	Item0                       int  `json:"item0"`
	Item1                       int  `json:"item1"`
	Item2                       int  `json:"item2"`
	Item3                       int  `json:"item3"`
	Item4                       int  `json:"item4"`
	Item5                       int  `json:"item5"`
	Item6                       int  `json:"item6"`
}

// Items returns the seven item slots (six items plus the trinket) in order.
func (s Stats) Items() [7]int {
	return [7]int{s.Item0, s.Item1, s.Item2, s.Item3, s.Item4, s.Item5, s.Item6}
}

// CS is the total creep score (lane minions plus jungle monsters).
func (s Stats) CS() int { return s.TotalMinionsKilled + s.NeutralMinionsKilled }

// ParticipantIdentity maps a participant slot to a player.
type ParticipantIdentity struct {
	ParticipantID int            `json:"participantId"`
	Player        IdentityPlayer `json:"player"`
}

// IdentityPlayer is the public identity of a participant.
type IdentityPlayer struct {
	Puuid        string `json:"puuid"`
	GameName     string `json:"gameName"`
	TagLine      string `json:"tagLine"`
	SummonerName string `json:"summonerName"`
	ProfileIcon  int    `json:"profileIcon"`
}

// GameTeam is the per-team result of a finished game.
type GameTeam struct {
	TeamID int    `json:"teamId"`
	Win    string `json:"win"` // "Win" or "Fail"
}
