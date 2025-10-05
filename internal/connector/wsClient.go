package connector

import (
	"crypto/tls"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

const wsBaseURL = "wss://127.0.0.1:%d"

func Connection(port int, authHeader string) {
	url := fmt.Sprintf(wsBaseURL, port)

	dialer := websocket.Dialer{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	headers := http.Header{}
	headers.Set("Authorization", authHeader)
	headers.Set("Accept", "*/*")

	conn, _, err := dialer.Dial(url, headers)
	if err != nil {
		fmt.Printf("[connector.Connection] dial %s failed: %v\n", url, err)
		return
	}
	defer func() {
		_ = conn.Close()
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("[connector.Connection] read failed: %v\n", err)
			return
		}
		fmt.Printf("[connector.Connection] message: %s\n", message)
	}
}
