package controller

import (
	"fmt"

	"lol-teammate-helper/internal/service"
	"lol-teammate-helper/internal/types"
	"lol-teammate-helper/internal/utils"
)

type PlayerController struct{}

func NewPlayerController() *PlayerController {
	return &PlayerController{}
}

func (pc *PlayerController) GetPlayerRankData(uuid string) types.RankedStats {
	const logPrefix = "[playerController.GetPlayerRankData]"

	if uuid == "" {
		fmt.Println(logPrefix + " empty uuid provided")
		return types.RankedStats{}
	}

	rankInfo, err := service.Shared().GetRankedStats(uuid)
	if err != nil {
		fmt.Printf("%s failed to load rank data: %v\n", logPrefix, err)
		return types.RankedStats{}
	}

	utils.ConvertRankDataToChinese(&rankInfo)
	return rankInfo
}
