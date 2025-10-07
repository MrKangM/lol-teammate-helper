package connector

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"lol-teammate-helper/internal/dispatch"
	"lol-teammate-helper/internal/types"
	"net/http"
)

const (
	wsBaseURL   = "wss://127.0.0.1:%d"
	lobbyAPIURI = "/lol-lobby/v2/lobby"
)

var testJson = `{
  "data": {
    "bans": {
      "myTeamBans": [],
      "numBans": 0,
      "theirTeamBans": []
    },
    "gameId": 0,
    "id": "mock-session",
    "myTeam": [
      {
        "assignedPosition": "utility",
        "cellId": 5,
        "championId": 161,
        "championPickIntent": 0,
        "gameName": "爱听周杰伦的歌呀",
        "internalName": "",
        "isHumanoid": false,
        "nameVisibilityType": "VISIBLE",
        "obfuscatedPuuid": "",
        "obfuscatedSummonerId": 0,
        "pickMode": 0,
        "pickTurn": 0,
        "playerAlias": "",
        "playerType": "",
        "puuid": "aae61f13-a4aa-506a-961e-bcfa71a8dbe6",
        "selectedSkinId": 161020,
        "spell1Id": 14,
        "spell2Id": 4,
        "summonerId": 16617704853,
        "tagLine": "70511",
        "team": 2,
        "wardSkinId": -1
      },
      {
        "assignedPosition": "jungle",
        "cellId": 6,
        "championId": 64,
        "championPickIntent": 0,
        "gameName": "丿丿艹神话艹",
        "internalName": "",
        "isHumanoid": false,
        "nameVisibilityType": "VISIBLE",
        "obfuscatedPuuid": "",
        "obfuscatedSummonerId": 0,
        "pickMode": 0,
        "pickTurn": 0,
        "playerAlias": "",
        "playerType": "",
        "puuid": "29025398-9f23-5f50-a180-0dc84cc9668d",
        "selectedSkinId": 64001,
        "spell1Id": 4,
        "spell2Id": 11,
        "summonerId": 17528309496,
        "tagLine": "81215",
        "team": 2,
        "wardSkinId": -1
      },
      {
        "assignedPosition": "top",
        "cellId": 7,
        "championId": 897,
        "championPickIntent": 0,
        "gameName": "FY丶青铜组V",
        "internalName": "",
        "isHumanoid": false,
        "nameVisibilityType": "VISIBLE",
        "obfuscatedPuuid": "",
        "obfuscatedSummonerId": 0,
        "pickMode": 0,
        "pickTurn": 0,
        "playerAlias": "",
        "playerType": "",
        "puuid": "b82ef168-018a-533b-904d-b767609e4bb8",
        "selectedSkinId": 897000,
        "spell1Id": 4,
        "spell2Id": 12,
        "summonerId": 15737607533,
        "tagLine": "52566",
        "team": 2,
        "wardSkinId": -1
      },
      {
        "assignedPosition": "middle",
        "cellId": 8,
        "championId": 105,
        "championPickIntent": 0,
        "gameName": "我叫王俊凯",
        "internalName": "",
        "isHumanoid": false,
        "nameVisibilityType": "VISIBLE",
        "obfuscatedPuuid": "",
        "obfuscatedSummonerId": 0,
        "pickMode": 0,
        "pickTurn": 0,
        "playerAlias": "",
        "playerType": "",
        "puuid": "4f4ff871-cda6-58f0-b759-6964c3b144c0",
        "selectedSkinId": 105001,
        "spell1Id": 14,
        "spell2Id": 4,
        "summonerId": 16685912130,
        "tagLine": "82964",
        "team": 2,
        "wardSkinId": -1
      },
      {
        "assignedPosition": "bottom",
        "cellId": 9,
        "championId": 15,
        "championPickIntent": 0,
        "gameName": "这波让我拉扯",
        "internalName": "",
        "isHumanoid": false,
        "nameVisibilityType": "VISIBLE",
        "obfuscatedPuuid": "",
        "obfuscatedSummonerId": 0,
        "pickMode": 0,
        "pickTurn": 0,
        "playerAlias": "",
        "playerType": "",
        "puuid": "637f3866-b9ad-5f33-81b4-1d7c393cc770",
        "selectedSkinId": 15034,
        "spell1Id": 4,
        "spell2Id": 21,
        "summonerId": 17912364269,
        "tagLine": "82413",
        "team": 2,
        "wardSkinId": -1
      }
    ],
    "queueId": 420
  },
  "eventType": "Update",
  "uri": "/lol-champ-select/v1/session"
}`

func Connection(port int, authHeader string) {
	url := fmt.Sprintf(wsBaseURL, port)

	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	headers := http.Header{}
	headers.Set("Authorization", authHeader)
	headers.Set("Content-Type", "application/json")
	headers.Set("Accept", "*/*")

	conn, resp, err := dialer.Dial(url, headers)
	if err != nil {
		fmt.Printf("[connector.Connection] dial %s failed: %v\n", url, err)
		return
	}

	if resp != nil && resp.StatusCode != http.StatusSwitchingProtocols {
		fmt.Printf("[connector.Connection] dial %s failed: %v\n", url, resp.StatusCode)
		_ = conn.Close()
		return
	}

	payload := []interface{}{5, "OnJsonApiEvent", lobbyAPIURI}
	if err := conn.WriteJSON(payload); err != nil {
		fmt.Printf("subscribe %s failed: %v\n", lobbyAPIURI, err)
	}

	msgCh := make(chan types.WSMessageType, 10)
	defer func() {
		_ = conn.Close()
		close(msgCh)
	}()

	go msgReader(msgCh)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("[connector.Connection] read failed: %v\n", err)
			return
		}

		var rawPayload []json.RawMessage
		if err := json.Unmarshal(message, &rawPayload); err != nil {
			fmt.Printf("[connector.Connection] invalid payload: %v\n", err)
			continue
		}

		if len(rawPayload) < 3 {
			continue
		}

		var eventName string
		if err := json.Unmarshal(rawPayload[1], &eventName); err != nil {
			fmt.Printf("[connector.Connection] parse event name failed: %v\n", err)
			continue
		}

		if eventName != "OnJsonApiEvent" {
			continue
		}

		var msg types.WSMessageType
		//if err := json.Unmarshal(rawPayload[2], &msg); err != nil {
		//	fmt.Printf("[connector.Connection] failed to decode msg data: %v\n", err)
		//	if err := json.Unmarshal([]byte(testJson), &msg); err != nil {
		//		fmt.Printf("[connector.Connection] failed to decode fallback msg: %v\n", err)
		//		continue
		//	}
		//}

		if len(msg.Data) == 0 {
			if err := json.Unmarshal([]byte(testJson), &msg); err != nil {
				fmt.Printf("[connector.Connection] failed to decode fallback msg: %v\n", err)
				continue
			}
		}

		msgCh <- msg
		fmt.Println("写入数据")
	}
}

func msgReader(msgCh <-chan types.WSMessageType) {
	for msg := range msgCh {
		dispatch.EventHandler(msg.Uri, msg.Data)
		fmt.Printf("读取数据:%d\n", len(msgCh))
	}
}
