package types

// 英雄选择会话数据
type ChampSelectData struct {
	//Actions [][]Action `json:"actions"` // 禁选和选择动作记录，按阶段分组
	//AllowBattleBoost                        bool           `json:"allowBattleBoost"`                        // 是否允许战斗加成
	//AllowDuplicatePicks                     bool           `json:"allowDuplicatePicks"`                     // 是否允许重复选择英雄
	//AllowLockedEvents                       bool           `json:"allowLockedEvents"`                       // 是否允许锁定事件
	//AllowPlayerPickSameChampion             bool           `json:"allowPlayerPickSameChampion"`             // 是否允许玩家选择相同英雄
	//AllowRerolling                          bool           `json:"allowRerolling"`                          // 是否允许重新随机
	//AllowSkinSelection                      bool           `json:"allowSkinSelection"`                      // 是否允许选择皮肤
	//AllowSubsetChampionPicks                bool           `json:"allowSubsetChampionPicks"`                // 是否允许子集英雄选择
	Bans Bans `json:"bans"` // 禁用信息
	//BenchChampions                          []interface{}  `json:"benchChampions"`                          // 板凳英雄列表
	//BenchEnabled                            bool           `json:"benchEnabled"`                            // 是否启用板凳功能
	//BoostableSkinCount                      int            `json:"boostableSkinCount"`                      // 可提升皮肤数量
	//ChatDetails ChatDetails `json:"chatDetails"` // 聊天详情信息
	//Counter int `json:"counter"` // 会话计数器
	//DisallowBanningTeammateHoveredChampions bool           `json:"disallowBanningTeammateHoveredChampions"` // 是否禁止禁用队友悬停的英雄
	GameID int64 `json:"gameId"` // 游戏ID
	//HasSimultaneousBans                     bool           `json:"hasSimultaneousBans"`                     // 是否同时进行禁用
	//HasSimultaneousPicks                    bool           `json:"hasSimultaneousPicks"`                    // 是否同时进行选择
	ID string `json:"id"` // 会话唯一标识
	//IsCustomGame                            bool           `json:"isCustomGame"`                            // 是否为自定义游戏
	//IsLegacyChampSelect                     bool           `json:"isLegacyChampSelect"`                     // 是否为旧版英雄选择
	//IsSpectating                            bool           `json:"isSpectating"`                            // 是否为观战模式
	//LocalPlayerCellID                       int            `json:"localPlayerCellId"`                       // 本地玩家单元格ID
	//LockedEventIndex                        int            `json:"lockedEventIndex"`                        // 锁定事件索引
	MyTeam []Player `json:"myTeam"` // 我方队伍玩家列表
	//PickOrderSwaps []interface{}  `json:"pickOrderSwaps"` // 选择顺序交换列表
	//PositionSwaps  []PositionSwap `json:"positionSwaps"`  // 位置交换信息
	QueueID int `json:"queueId"` // 队列ID
	//RerollsRemaining                        int            `json:"rerollsRemaining"`                        // 剩余重新随机次数
	//ShowQuitButton                          bool           `json:"showQuitButton"`                          // 是否显示退出按钮
	//SkipChampionSelect                      bool           `json:"skipChampionSelect"`                      // 是否跳过英雄选择
	//TheirTeam                               []Player       `json:"theirTeam"`                               // 敌方队伍玩家列表
	//Timer                                   Timer          `json:"timer"`                                   // 计时器信息
	//Trades                                  []Trade        `json:"trades"`                                  // 英雄交换信息
}

// 禁选/选择动作
type Action struct {
	ActorCellID  int    `json:"actorCellId"`  // 执行者单元格ID
	ChampionID   int    `json:"championId"`   // 英雄ID
	Completed    bool   `json:"completed"`    // 是否已完成
	Duration     int    `json:"duration"`     // 持续时间
	ID           int    `json:"id"`           // 动作ID
	IsAllyAction bool   `json:"isAllyAction"` // 是否为友方动作
	IsInProgress bool   `json:"isInProgress"` // 是否进行中
	PickTurn     int    `json:"pickTurn"`     // 选择回合
	Type         string `json:"type"`         // 动作类型: ban-禁用, pick-选择, ten_bans_reveal-显示10个禁用
}

// 禁用信息
type Bans struct {
	MyTeamBans    []int `json:"myTeamBans"`    // 我方队伍禁用英雄ID列表
	NumBans       int   `json:"numBans"`       // 禁用数量
	TheirTeamBans []int `json:"theirTeamBans"` // 敌方队伍禁用英雄ID列表
}

// 聊天详情
type ChatDetails struct {
	MucJwtDto             MucJwtDto `json:"mucJwtDto"`             // MUC JWT数据传输对象
	MultiUserChatID       string    `json:"multiUserChatId"`       // 多人聊天ID
	MultiUserChatPassword string    `json:"multiUserChatPassword"` // 多人聊天密码
}

// MUC JWT数据传输对象
type MucJwtDto struct {
	ChannelClaim string `json:"channelClaim"` // 频道声明
	Domain       string `json:"domain"`       // 域名
	Jwt          string `json:"jwt"`          // JWT令牌
	TargetRegion string `json:"targetRegion"` // 目标区域
}

// 玩家信息
type Player struct {
	AssignedPosition     string `json:"assignedPosition"`     // 分配位置: top-上单, jungle-打野, middle-中单, bottom-ADC, utility-辅助
	CellID               int    `json:"cellId"`               // 单元格ID
	ChampionID           int    `json:"championId"`           // 选择的英雄ID
	ChampionPickIntent   int    `json:"championPickIntent"`   // 英雄选择意图
	GameName             string `json:"gameName"`             // 游戏昵称
	InternalName         string `json:"internalName"`         // 内部名称
	IsHumanoid           bool   `json:"isHumanoid"`           // 是否为人形（是否为真实玩家）
	NameVisibilityType   string `json:"nameVisibilityType"`   // 名称可见性类型: VISIBLE-可见, HIDDEN-隐藏
	ObfuscatedPuuid      string `json:"obfuscatedPuuid"`      // 混淆的玩家唯一标识
	ObfuscatedSummonerID int    `json:"obfuscatedSummonerId"` // 混淆的召唤师ID
	PickMode             int    `json:"pickMode"`             // 选择模式
	PickTurn             int    `json:"pickTurn"`             // 选择回合
	PlayerAlias          string `json:"playerAlias"`          // 玩家别名
	PlayerType           string `json:"playerType"`           // 玩家类型
	Puuid                string `json:"puuid"`                // 玩家唯一标识
	SelectedSkinID       int    `json:"selectedSkinId"`       // 选择的皮肤ID
	Spell1ID             int    `json:"spell1Id"`             // 召唤师技能1ID
	Spell2ID             int    `json:"spell2Id"`             // 召唤师技能2ID
	SummonerID           int64  `json:"summonerId"`           // 召唤师ID
	TagLine              string `json:"tagLine"`              // 标签行
	Team                 int    `json:"team"`                 // 队伍: 1-蓝队, 2-红队
	WardSkinID           int    `json:"wardSkinId"`           // 眼皮肤ID
}

// 位置交换信息
type PositionSwap struct {
	CellID int    `json:"cellId"` // 单元格ID
	ID     int    `json:"id"`     // 交换ID
	State  string `json:"state"`  // 状态: AVAILABLE-可用, INVALID-无效
}

// 计时器信息
type Timer struct {
	AdjustedTimeLeftInPhase int    `json:"adjustedTimeLeftInPhase"` // 调整后的阶段剩余时间(毫秒)
	InternalNowInEpochMs    int64  `json:"internalNowInEpochMs"`    // 内部当前时间(纪元毫秒)
	IsInfinite              bool   `json:"isInfinite"`              // 是否无限时间
	Phase                   string `json:"phase"`                   // 当前阶段: FINALIZATION-最终确认阶段
	TotalTimeInPhase        int    `json:"totalTimeInPhase"`        // 阶段总时间(毫秒)
}

// 英雄交换信息
type Trade struct {
	CellID int    `json:"cellId"` // 单元格ID
	ID     int    `json:"id"`     // 交换ID
	State  string `json:"state"`  // 状态: AVAILABLE-可用, INVALID-无效
}
