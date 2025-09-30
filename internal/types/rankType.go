package types

// RankedStats 排位统计数据
type RankedStats struct {
	CurrentSeasonSplitPoints          int64                  `json:"currentSeasonSplitPoints"`          // 当前赛季分段点数
	EarnedRegaliaRewardIds            []string               `json:"earnedRegaliaRewardIds"`            // 已获得的荣誉奖励ID列表
	HighestCurrentSeasonReachedTierSR string                 `json:"highestCurrentSeasonReachedTierSR"` // 当前赛季达到的最高段位（单双排）
	HighestPreviousSeasonEndDivision  string                 `json:"highestPreviousSeasonEndDivision"`  // 上赛季结束时的最高段位分级
	HighestPreviousSeasonEndTier      string                 `json:"highestPreviousSeasonEndTier"`      // 上赛季结束时的最高段位
	HighestRankedEntry                RankedEntry            `json:"highestRankedEntry"`                // 最高排位条目
	HighestRankedEntrySR              RankedEntry            `json:"highestRankedEntrySR"`              // 最高排位条目（单双排）
	PreviousSeasonSplitPoints         int64                  `json:"previousSeasonSplitPoints"`         // 上赛季分段点数
	QueueMap                          map[string]RankedEntry `json:"queueMap"`                          // 各队列排位数据映射
	Queues                            []RankedEntry          `json:"queues"`                            // 排位队列列表
	RankedRegaliaLevel                int64                  `json:"rankedRegaliaLevel"`                // 排位荣誉等级
	Seasons                           map[string]SeasonInfo  `json:"seasons"`                           // 各队列赛季信息
	SplitsProgress                    map[string]interface{} `json:"splitsProgress"`                    // 分段进度
}

// RankedEntry 排位条目信息
type RankedEntry struct {
	CurrentSeasonWinsForRewards   int64       `json:"currentSeasonWinsForRewards"`   // 当前赛季奖励所需胜场数
	Division                      string      `json:"division"`                      // 当前段位分级（如I、II、III、IV）
	HighestDivision               string      `json:"highestDivision"`               // 历史最高段位分级
	HighestTier                   string      `json:"highestTier"`                   // 历史最高段位
	IsProvisional                 bool        `json:"isProvisional"`                 // 是否为定级赛状态
	LeaguePoints                  int64       `json:"leaguePoints"`                  // 胜点
	Losses                        int64       `json:"losses"`                        // 失败场次
	MiniSeriesProgress            string      `json:"miniSeriesProgress"`            // 晋级赛进度
	PreviousSeasonEndDivision     string      `json:"previousSeasonEndDivision"`     // 上赛季结束时的段位分级
	PreviousSeasonEndTier         string      `json:"previousSeasonEndTier"`         // 上赛季结束时的段位
	PreviousSeasonHighestDivision string      `json:"previousSeasonHighestDivision"` // 上赛季最高段位分级
	PreviousSeasonHighestTier     string      `json:"previousSeasonHighestTier"`     // 上赛季最高段位
	PreviousSeasonWinsForRewards  int64       `json:"previousSeasonWinsForRewards"`  // 上赛季奖励所需胜场数
	ProvisionalGameThreshold      int64       `json:"provisionalGameThreshold"`      // 定级赛阈值
	ProvisionalGamesRemaining     int64       `json:"provisionalGamesRemaining"`     // 剩余定级赛场次
	QueueType                     string      `json:"queueType"`                     // 队列类型
	RatedRating                   int64       `json:"ratedRating"`                   // 评分
	RatedTier                     string      `json:"ratedTier"`                     // 评分段位
	Tier                          string      `json:"tier"`                          // 当前段位
	Warnings                      interface{} `json:"warnings"`                      // 警告信息
	Wins                          int64       `json:"wins"`                          // 胜利场次
}

// SeasonInfo 赛季信息
type SeasonInfo struct {
	CurrentSeasonEnd int64 `json:"currentSeasonEnd"` // 当前赛季结束时间戳
	CurrentSeasonId  int64 `json:"currentSeasonId"`  // 当前赛季ID
	NextSeasonStart  int64 `json:"nextSeasonStart"`  // 下赛季开始时间戳
}
