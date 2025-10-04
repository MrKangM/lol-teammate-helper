package types

type IPlayerBaseData struct {
	AccountId                   int     `json:"accountId,omitempty"`
	DisplayName                 string  `json:"displayName,omitempty"`
	GameName                    string  `json:"gameName,omitempty"`
	InternalName                string  `json:"internalName,omitempty"`
	NameChangeFlag              bool    `json:"nameChangeFlag,omitempty"`
	PercentCompleteForNextLevel int     `json:"percentCompleteForNextLevel,omitempty"` // 下一个等级百分比
	Privacy                     string  `json:"privacy,omitempty"`
	ProfileIconId               int     `json:"profileIconId,omitempty"` // 头像ID
	Puuid                       string  `json:"puuid,omitempty"`
	RerollPoints                IReroll `json:"rerollPoints,omitempty"`
	SummonerId                  int64   `json:"summonerId,omitempty"`
	SummonerLevel               int     `json:"summonerLevel,omitempty"` // 召唤师等级
	TagLine                     string  `json:"tagLine,omitempty"`
	Unnamed                     bool    `json:"unnamed,omitempty"`
	XpSinceLastLevel            int     `json:"xpSinceLastLevel,omitempty"` // 当前经验
	XpUntilNextLevel            int     `json:"xpUntilNextLevel,omitempty"` // 下一等级经验
	IconImgSrc                  string  `json:"iconImgSrc,omitempty"`       // 头像URL
	Region                      string  `json:"region,omitempty"`           // 服务器地区
}

type IReroll struct {
	CurrentPoints    int `json:"currentPoints"`
	MaxRolls         int `json:"maxRolls"`
	NumberOfRolls    int `json:"numberOfRolls"`
	PointsCostToRoll int `json:"pointsCostToRoll"`
	PointsToReroll   int `json:"pointsToReroll"`
}
