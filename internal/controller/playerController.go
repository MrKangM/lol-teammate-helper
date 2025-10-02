package controller

import (
	"encoding/json"
	"fmt"
	"lol-teammate-helper/internal/config"
	"lol-teammate-helper/internal/types"
	"lol-teammate-helper/internal/utils"
	"net/http"
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

	cfg, ok := config.Instance()
	if !ok {
		fmt.Println(logPrefix + " riot credentials are not initialised")
		return types.RankedStats{}
	}

	endpoint := fmt.Sprintf("/lol-ranked/v1/ranked-stats/%s", uuid)
	fmt.Println(logPrefix + " requesting " + endpoint)

	body, err := cfg.SendHttpRequest(endpoint, http.MethodGet)
	if err != nil {
		fmt.Printf("%s failed to request rank data: %v\n", logPrefix, err)
		return types.RankedStats{}
	}

	var rankInfo types.RankedStats
	if err := json.Unmarshal(body, &rankInfo); err != nil {
		fmt.Printf("%s failed to decode rank data: %v\n", logPrefix, err)
		return types.RankedStats{}
	}

	utils.ConvertRankDataToChinese(&rankInfo)
	fmt.Println(logPrefix + " ====================================")
	fmt.Printf("%s raw response: %#v\n", logPrefix, rankInfo)
	fmt.Println(logPrefix + " summary: " + utils.FormatRankInfo(rankInfo))

	return rankInfo
}
