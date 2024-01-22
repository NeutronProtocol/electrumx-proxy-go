package ws

import (
	"electrumx-proxy-go/common"
	"electrumx-proxy-go/config"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"sync"
	"time"
)

var Callbacks = make(map[uint32]chan *common.JsonRpcResponse)

var WsTx *websocket.Conn

var CallbacksLock = sync.RWMutex{}

func InitWebSocket(urlStr string) {
	for {
		err := connectWebSocket(urlStr)
		if err != nil {
			log.Printf("Error connecting to WebSocket, will retry in 5 seconds: %v", err)
			time.Sleep(5 * time.Second)
			continue
		}
		break
	}
}
func connectWebSocket(urlStr string) error {
	u, err := url.Parse(urlStr)
	if err != nil {
		return fmt.Errorf("error parsing URL: %v", err)
	}
	dialer := websocket.DefaultDialer
	conn, _, err := dialer.Dial(u.String(), nil)
	if err != nil {
		return fmt.Errorf("error connecting to WebSocket: %v", err)
	}
	WsTx = conn
	go readMessages(conn)
	return nil
}
func readMessages(conn *websocket.Conn) {
	defer conn.Close()
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket closed unexpectedly: %v", err)
				InitWebSocket(config.Conf.ElectrumxServer)
			} else {
				log.Printf("Error reading websocket message: %v", err)
			}
			break
		}
		var response common.JsonRpcResponse
		if err := json.Unmarshal(message, &response); err != nil {
			log.Printf("Error unmarshalling websocket message: %v", err)
			continue
		}
		CallbacksLock.Lock()
		if ch, ok := Callbacks[response.ID]; ok {
			ch <- &response
			delete(Callbacks, response.ID)
		}
		CallbacksLock.Unlock()
	}
}
