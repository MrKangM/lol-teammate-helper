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

type CredentialsProvider interface {
	RiotPort() int
	RiotToken() string
}

type PlayerController struct {
	provider CredentialsProvider
}

func NewPlayerController(provider CredentialsProvider) *PlayerController {
	return &PlayerController{provider: provider}
}

func (pc *PlayerController) GetPlayerRankData(uuid string) types.RankedStats {
	if uuid == "" {
		fmt.Println("player controller: empty uuid provided")
		return types.RankedStats{}
	}

	provider := pc.provider
	if provider == nil {
		fmt.Println("player controller: riot credentials provider is not configured")
		return types.RankedStats{}
	}

	port := provider.RiotPort()
	token := provider.RiotToken()
	if port <= 0 {
		fmt.Println("player controller: riot port unavailable")
		return types.RankedStats{}
	}
	if token == "" {
		fmt.Println("player controller: riot token unavailable")
		return types.RankedStats{}
	}

	url := fmt.Sprintf("https://127.0.0.1:%d/lol-ranked/v1/ranked-stats/%s", port, uuid)
	fmt.Println("閹烘帊缍呴弫鐗堝祦URL:" + url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		fmt.Printf("閸掓稑缂撻幒鎺嶇秴鐠囬攱鐪版径杈Е: %v\n", err)
		return types.RankedStats{}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	authString := "riot:" + token
	req.Header.Set("Authorization", "Basic "+base64.StdEncoding.EncodeToString([]byte(authString)))

	resp, err := reqLib.Instance().Do(req)
	if err != nil {
		fmt.Printf("鐠囬攱鐪伴幒鎺嶇秴閺佺増宓佹径杈Е: %v\n", err)
		return types.RankedStats{}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("ranked stats request returned status %d\n", resp.StatusCode)
		return types.RankedStats{}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("鐠囪褰囬幒鎺嶇秴閸濆秴绨叉径杈Е: %v\n", err)
		return types.RankedStats{}
	}

	var rankInfo types.RankedStats
	if err := json.Unmarshal(body, &rankInfo); err != nil {
		fmt.Printf("鐟欙絾鐎介幒鎺嶇秴閺佺増宓佹径杈Е: %v\n", err)
		return types.RankedStats{}
	}

	utils.ConvertRankDataToChinese(&rankInfo)
	fmt.Println("====================================")
	fmt.Println(rankInfo)
	fmt.Println("閹崵绮ㄦ穱鈩冧紖:" + utils.FormatRankInfo(rankInfo))

	return rankInfo
}
