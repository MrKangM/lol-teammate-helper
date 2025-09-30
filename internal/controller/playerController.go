package controller

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"lol-teammate-helper/internal/reqLib"
	"lol-teammate-helper/internal/types"
	"lol-teammate-helper/internal/utils"
	"net/http"
)

type PlayerController struct {
}

func (pc *PlayerController) GetPlayerRankData(uuid string) types.RankedStats {
	var port int = 56737
	var token string = "JaffTt5fcKeFztB5KHO4Hw"
	puuid := "637f3866-b9ad-5f33-81b4-1d7c393cc770"
	baseDataUrl := fmt.Sprintf("/lol-ranked/v1/ranked-stats/%s", puuid)
	url := fmt.Sprintf("https://127.0.0.1:%d%s", port, baseDataUrl)
	fmt.Println("排位数据URL:" + url)
	client := reqLib.Instance().Client()
	// 创建请求
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Sprintf("创建请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// 创建Authorization头：Basic + base64("riot:token")
	authString := "riot:" + token
	encodedAuth := base64.StdEncoding.EncodeToString([]byte(authString))
	req.Header.Set("Authorization", "Basic "+encodedAuth)

	// 发送请求
	resp, err := client.Do(req)
	if err != nil {
		fmt.Sprintf("请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Sprintf("读取响应失败: %v", err)
	}

	// 返回响应结果
	var rankInfo types.RankedStats
	result := fmt.Sprintf("排位!!!状态码: %d\n响应体: %s", resp.StatusCode, string(body))
	fmt.Println("排位数据:" + result)
	err = json.Unmarshal(body, &rankInfo)
	if err != nil {
		fmt.Printf("解析排位数据失败: %v\n", err)
		//return // 或者处理错误
	}
	utils.ConvertRankDataToChinese(&rankInfo)
	fmt.Println("====================================")
	fmt.Println(rankInfo)
	fmt.Println("总结信息：" + utils.FormatRankInfo(rankInfo))
	return rankInfo
}
