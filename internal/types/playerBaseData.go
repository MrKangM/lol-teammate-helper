package types

type IPlayerBaseData struct {
	AccountId                   int     `json:"accountId,omitempty"`
	DisplayName                 string  `json:"displayName,omitempty"`
	GameName                    string  `json:"gameName,omitempty"`
	InternalName                string  `json:"internalName,omitempty"`
	NameChangeFlag              bool    `json:"nameChangeFlag,omitempty"`
	PercentCompleteForNextLevel int     `json:"percentCompleteForNextLevel,omitempty"` // Percent progress towards next level
	Privacy                     string  `json:"privacy,omitempty"`
	ProfileIconId               int     `json:"profileIconId,omitempty"`
	Puuid                       string  `json:"puuid,omitempty"`
	RerollPoints                IReroll `json:"rerollPoints,omitempty"`
	SummonerId                  int64   `json:"summonerId,omitempty"`
	SummonerLevel               int     `json:"summonerLevel,omitempty"` // Summoner level
	TagLine                     string  `json:"tagLine,omitempty"`
	Unnamed                     bool    `json:"unnamed,omitempty"`
	XpSinceLastLevel            int     `json:"xpSinceLastLevel,omitempty"`
	XpUntilNextLevel            int     `json:"xpUntilNextLevel,omitempty"`
	IconImgSrc                  string  `json:"iconImgSrc,omitempty"`
	Region                      string  `json:"region,omitempty"`
}

type IReroll struct {
	CurrentPoints    int `json:"currentPoints"`
	MaxRolls         int `json:"maxRolls"`
	NumberOfRolls    int `json:"numberOfRolls"`
	PointsCostToRoll int `json:"pointsCostToRoll"`
	PointsToReroll   int `json:"pointsToReroll"`
}
