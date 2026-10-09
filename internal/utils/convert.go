package utils

import (
	"fmt"
	"lol-teammate-helper/internal/types"
	"strings"
)

// tierNames maps LCU tier identifiers to display names.
var tierNames = map[string]string{
	"IRON":        "黑铁",
	"BRONZE":      "青铜",
	"SILVER":      "白银",
	"GOLD":        "黄金",
	"PLATINUM":    "白金",
	"EMERALD":     "翡翠",
	"DIAMOND":     "钻石",
	"MASTER":      "大师",
	"GRANDMASTER": "宗师",
	"CHALLENGER":  "王者",
	"":            "未定级",
	"NA":          "未定级",
}

// queueNames maps LCU queue types to display names.
var queueNames = map[string]string{
	"RANKED_SOLO_5x5":      "单双排",
	"RANKED_FLEX_SR":       "灵活组排",
	"RANKED_TFT":           "云顶之弈",
	"RANKED_TFT_DOUBLE_UP": "云顶双人",
	"RANKED_TFT_TURBO":     "云顶极速",
}

func ConvertRankDataToChinese(rankInfo *types.RankedStats) {
	if chineseTier, ok := tierNames[rankInfo.HighestCurrentSeasonReachedTierSR]; ok {
		rankInfo.HighestCurrentSeasonReachedTierSR = chineseTier
	}
	if chineseTier, ok := tierNames[rankInfo.HighestPreviousSeasonEndTier]; ok {
		rankInfo.HighestPreviousSeasonEndTier = chineseTier
	}

	for queueType, entry := range rankInfo.QueueMap {
		convertRankedEntry(&entry, queueType)
		rankInfo.QueueMap[queueType] = entry
	}

	for i := range rankInfo.Queues {
		convertRankedEntry(&rankInfo.Queues[i], rankInfo.Queues[i].QueueType)
	}

	convertRankedEntry(&rankInfo.HighestRankedEntry, rankInfo.HighestRankedEntry.QueueType)
	convertRankedEntry(&rankInfo.HighestRankedEntrySR, rankInfo.HighestRankedEntrySR.QueueType)
}

func convertRankedEntry(entry *types.RankedEntry, queueType string) {
	if chineseTier, ok := tierNames[entry.Tier]; ok {
		entry.Tier = chineseTier
	}
	if chineseTier, ok := tierNames[entry.HighestTier]; ok {
		entry.HighestTier = chineseTier
	}
	if chineseTier, ok := tierNames[entry.PreviousSeasonEndTier]; ok {
		entry.PreviousSeasonEndTier = chineseTier
	}
	if chineseTier, ok := tierNames[entry.PreviousSeasonHighestTier]; ok {
		entry.PreviousSeasonHighestTier = chineseTier
	}
	if chineseTier, ok := tierNames[entry.RatedTier]; ok {
		entry.RatedTier = chineseTier
	}

	if chineseName, ok := queueNames[queueType]; ok {
		entry.QueueDisplayName = chineseName
	}
}

func FormatRankInfo(rankInfo types.RankedStats) string {
	var result strings.Builder

	result.WriteString("=== 排位信息汇总 ===\n")
	result.WriteString(fmt.Sprintf("当前赛季最高段位: %s\n", rankInfo.HighestCurrentSeasonReachedTierSR))
	result.WriteString(fmt.Sprintf("上赛季最高段位: %s %s\n",
		rankInfo.HighestPreviousSeasonEndTier,
		rankInfo.HighestPreviousSeasonEndDivision))

	for queueType, entry := range rankInfo.QueueMap {
		if entry.Tier == "" || entry.Tier == "未定级" {
			continue
		}

		queueName := entry.QueueDisplayName
		if queueName == "" {
			queueName = getQueueChineseName(queueType)
		}

		result.WriteString(fmt.Sprintf("\n[%s]\n", queueName))
		result.WriteString(fmt.Sprintf("当前段位: %s %s", entry.Tier, entry.Division))
		if entry.LeaguePoints > 0 {
			result.WriteString(fmt.Sprintf(" (%d胜点)", entry.LeaguePoints))
		}
		result.WriteString("\n")

		if entry.HighestTier != "" && entry.HighestTier != "未定级" {
			result.WriteString(fmt.Sprintf("历史最高: %s %s\n", entry.HighestTier, entry.HighestDivision))
		}

		if total := entry.Wins + entry.Losses; total > 0 {
			winRate := float64(entry.Wins) / float64(total) * 100
			result.WriteString(fmt.Sprintf("战绩: %d胜 %d负 (胜率%.1f%%)\n", entry.Wins, entry.Losses, winRate))
		}
	}

	return result.String()
}

func getQueueChineseName(queueType string) string {
	if name, ok := queueNames[queueType]; ok {
		return name
	}
	return queueType
}

// serverNames maps platform IDs to Chinese server names.
var serverNames = map[string]string{
	"HN1":   "艾欧尼亚",
	"HN10":  "黑色玫瑰",
	"TJ100": "联盟四区",
	"TJ101": "联盟五区",
	"NJ100": "联盟一区",
	"GZ100": "联盟二区",
	"CQ100": "联盟三区",
	"BGP2":  "峡谷之巅",
	"PBE":   "体验服",
}

func GetServerChineseName(serverName string) string {
	if name, ok := serverNames[serverName]; ok {
		return name
	}
	return serverName
}

// TierName returns the Chinese display name for an LCU tier key such as "GOLD".
func TierName(tierKey string) string {
	if name, ok := tierNames[tierKey]; ok {
		return name
	}
	return tierKey
}

// QueueName returns the Chinese display name of a ranked queue type.
func QueueName(queueType string) string {
	return getQueueChineseName(queueType)
}
