package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"lol-teammate-helper/internal/reqLib"
	"lol-teammate-helper/internal/types"
	"lol-teammate-helper/internal/utils"
	"net/http"
)

type PlayerController struct {
}

func (pc *PlayerController) GetPlayerRankData(uuid string) types.RankedStats {
	const (
		defaultPort   = 56737
		defaultToken  = "iOkyLnA-zyiJNOB9cQ1svA"
		fallbackPUUID = "637f3866-b9ad-5f33-81b4-1d7c393cc770"
	)

	identifier := uuid
	if len(identifier) == 0 {
		identifier = fallbackPUUID
	}

	url := fmt.Sprintf("https://127.0.0.1:%d/lol-ranked/v1/ranked-stats/%s", defaultPort, identifier)
	fmt.Println("排位数据URL:" + url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Printf("创建排位请求失败: %v\n", err)
		return types.RankedStats{}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	authString := "riot:" + defaultToken
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(authString)))

	resp, err := reqLib.Instance().Do(req)
	if err != nil {
		fmt.Printf("请求排位数据失败: %v\n", err)
		return types.RankedStats{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("排位请求返回非200状态: %d\n", resp.StatusCode)
		return types.RankedStats{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("读取排位响应失败: %v\n", err)
		return types.RankedStats{}
	}

	var rankInfo types.RankedStats
	if err := json.Unmarshal(body, &rankInfo); err != nil {
		fmt.Printf("解析排位数据失败: %v\n", err)
		return types.RankedStats{}
	}

	utils.ConvertRankDataToChinese(&rankInfo)
	fmt.Println("====================================")
	fmt.Println(rankInfo)
	fmt.Println("总结信息:" + utils.FormatRankInfo(rankInfo))

	return rankInfo
}
