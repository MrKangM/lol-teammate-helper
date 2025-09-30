package utils

import (
	"fmt"
	"lol-teammate-helper/internal/types"
	"strings"
)

// ConvertRankDataToChinese 将段位数据转换为中文显示
func ConvertRankDataToChinese(rankInfo *types.RankedStats) {
	// 定义段位映射
	tierMap := map[string]string{
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

	// 定义队列名称映射
	queueMap := map[string]string{
		"RANKED_SOLO_5x5":      "单双排",
		"RANKED_FLEX_SR":       "灵活组排",
		"RANKED_TFT":           "云顶之弈",
		"RANKED_TFT_DOUBLE_UP": "云顶双人",
		"RANKED_TFT_TURBO":     "云顶快速",
	}

	// 转换最高段位字段
	if chineseTier, exists := tierMap[rankInfo.HighestCurrentSeasonReachedTierSR]; exists {
		rankInfo.HighestCurrentSeasonReachedTierSR = chineseTier
	}
	if chineseTier, exists := tierMap[rankInfo.HighestPreviousSeasonEndTier]; exists {
		rankInfo.HighestPreviousSeasonEndTier = chineseTier
	}

	// 转换队列映射中的段位数据
	for queueType, entry := range rankInfo.QueueMap {
		convertRankedEntry(&entry, tierMap, queueMap, queueType)
		rankInfo.QueueMap[queueType] = entry
	}

	// 转换队列列表中的段位数据
	for i := range rankInfo.Queues {
		convertRankedEntry(&rankInfo.Queues[i], tierMap, queueMap, rankInfo.Queues[i].QueueType)
	}

	// 转换最高排位条目
	convertRankedEntry(&rankInfo.HighestRankedEntry, tierMap, queueMap, rankInfo.HighestRankedEntry.QueueType)
	convertRankedEntry(&rankInfo.HighestRankedEntrySR, tierMap, queueMap, rankInfo.HighestRankedEntrySR.QueueType)
}

// convertRankedEntry 转换单个排位条目的段位信息
// 修复：参数类型应该是 *controller.RankedEntry 而不是 *controller.RankedStats
func convertRankedEntry(entry *types.RankedEntry, tierMap map[string]string, queueMap map[string]string, queueType string) {
	// 转换段位
	if chineseTier, exists := tierMap[entry.Tier]; exists {
		entry.Tier = chineseTier
	}
	if chineseTier, exists := tierMap[entry.HighestTier]; exists {
		entry.HighestTier = chineseTier
	}
	if chineseTier, exists := tierMap[entry.PreviousSeasonEndTier]; exists {
		entry.PreviousSeasonEndTier = chineseTier
	}
	if chineseTier, exists := tierMap[entry.PreviousSeasonHighestTier]; exists {
		entry.PreviousSeasonHighestTier = chineseTier
	}
	if chineseTier, exists := tierMap[entry.RatedTier]; exists {
		entry.RatedTier = chineseTier
	}

	// 转换队列类型显示名（可选，如果需要的话）
	if chineseName, exists := queueMap[queueType]; exists {
		// 这里可以选择是否修改QueueType字段，或者创建新的字段存储中文名
		// 为了避免破坏原始数据结构，这里只做转换显示，不修改原始QueueType
		_ = chineseName // 可以在打印时使用
	}
}

// 添加一个方法来格式化显示排位信息
func FormatRankInfo(rankInfo types.RankedStats) string {
	var result strings.Builder

	result.WriteString("=== 排位信息汇总 ===\n")

	// 整体信息
	result.WriteString(fmt.Sprintf("当前赛季最高段位: %s\n", rankInfo.HighestCurrentSeasonReachedTierSR))
	result.WriteString(fmt.Sprintf("上赛季最高段位: %s %s\n",
		rankInfo.HighestPreviousSeasonEndTier,
		rankInfo.HighestPreviousSeasonEndDivision))

	// 各队列信息
	for queueType, entry := range rankInfo.QueueMap {
		if entry.Tier != "" && entry.Tier != "未定级" {
			queueName := getQueueChineseName(queueType) // 修复：删除 pc. 前缀
			result.WriteString(fmt.Sprintf("\n【%s】\n", queueName))
			result.WriteString(fmt.Sprintf("当前段位: %s %s", entry.Tier, entry.Division))
			if entry.LeaguePoints > 0 {
				result.WriteString(fmt.Sprintf(" (%d胜点)", entry.LeaguePoints))
			}
			result.WriteString("\n")

			if entry.HighestTier != "" && entry.HighestTier != "未定级" {
				result.WriteString(fmt.Sprintf("历史最高: %s %s\n", entry.HighestTier, entry.HighestDivision))
			}

			if entry.Wins+entry.Losses > 0 {
				winRate := float64(entry.Wins) / float64(entry.Wins+entry.Losses) * 100
				result.WriteString(fmt.Sprintf("战绩: %d胜 %d负 (胜率%.1f%%)\n", entry.Wins, entry.Losses, winRate))
			}
		}
	}

	return result.String()
}

// getQueueChineseName 获取队列中文名称
func getQueueChineseName(queueType string) string {
	queueNames := map[string]string{
		"RANKED_SOLO_5x5":      "单双排",
		"RANKED_FLEX_SR":       "灵活组排",
		"RANKED_TFT":           "云顶之弈",
		"RANKED_TFT_DOUBLE_UP": "云顶双人",
		"RANKED_TFT_TURBO":     "云顶快速",
	}
	if name, exists := queueNames[queueType]; exists {
		return name
	}
	return queueType
}
