package types

type RankedStats struct {
	CurrentSeasonSplitPoints          int64                  `json:"currentSeasonSplitPoints"`
	EarnedRegaliaRewardIds            []string               `json:"earnedRegaliaRewardIds"`
	HighestCurrentSeasonReachedTierSR string                 `json:"highestCurrentSeasonReachedTierSR"`
	HighestPreviousSeasonEndDivision  string                 `json:"highestPreviousSeasonEndDivision"`
	HighestPreviousSeasonEndTier      string                 `json:"highestPreviousSeasonEndTier"`
	HighestRankedEntry                RankedEntry            `json:"highestRankedEntry"`
	HighestRankedEntrySR              RankedEntry            `json:"highestRankedEntrySR"`
	PreviousSeasonSplitPoints         int64                  `json:"previousSeasonSplitPoints"`
	QueueMap                          map[string]RankedEntry `json:"queueMap"`
	Queues                            []RankedEntry          `json:"queues"`
	RankedRegaliaLevel                int64                  `json:"rankedRegaliaLevel"`
	Seasons                           map[string]SeasonInfo  `json:"seasons"`
	SplitsProgress                    map[string]any         `json:"splitsProgress"`
}

type RankedEntry struct {
	CurrentSeasonWinsForRewards   int64  `json:"currentSeasonWinsForRewards"`
	Division                      string `json:"division"`
	HighestDivision               string `json:"highestDivision"`
	HighestTier                   string `json:"highestTier"`
	IsProvisional                 bool   `json:"isProvisional"`
	LeaguePoints                  int64  `json:"leaguePoints"`
	Losses                        int64  `json:"losses"`
	MiniSeriesProgress            string `json:"miniSeriesProgress"`
	PreviousSeasonEndDivision     string `json:"previousSeasonEndDivision"`
	PreviousSeasonEndTier         string `json:"previousSeasonEndTier"`
	PreviousSeasonHighestDivision string `json:"previousSeasonHighestDivision"`
	PreviousSeasonHighestTier     string `json:"previousSeasonHighestTier"`
	PreviousSeasonWinsForRewards  int64  `json:"previousSeasonWinsForRewards"`
	ProvisionalGameThreshold      int64  `json:"provisionalGameThreshold"`
	ProvisionalGamesRemaining     int64  `json:"provisionalGamesRemaining"`
	QueueType                     string `json:"queueType"`
	RatedRating                   int64  `json:"ratedRating"`
	RatedTier                     string `json:"ratedTier"`
	Tier                          string `json:"tier"`
	Warnings                      any    `json:"warnings"`
	Wins                          int64  `json:"wins"`
	QueueDisplayName              string `json:"-"`
}

type SeasonInfo struct {
	CurrentSeasonEnd int64 `json:"currentSeasonEnd"`
	CurrentSeasonID  int64 `json:"currentSeasonId"`
	NextSeasonStart  int64 `json:"nextSeasonStart"`
}
