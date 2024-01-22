package router

import (
	"electrumx-proxy-go/common"
	"electrumx-proxy-go/ws"
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"time"
)

func handleRequestToWs(method string, params []interface{}) (*common.JsonRpcResponse, error) {
	id := GetUniqueID()
	request := common.JsonRpcRequest{
		ID:     id,
		Method: method,
		Params: params,
	}

	requestBytes, err := json.Marshal(request)
	if err != nil {
		return nil, err
	}
	responseChan := make(chan *common.JsonRpcResponse, 1)
	ws.CallbacksLock.Lock()
	ws.Callbacks[id] = responseChan
	ws.CallbacksLock.Unlock()

	if err := ws.WsTx.WriteMessage(websocket.TextMessage, requestBytes); err != nil {
		ws.CallbacksLock.Lock()
		delete(ws.Callbacks, id)
		ws.CallbacksLock.Unlock()
		return nil, err
	}

	select {
	case response := <-responseChan:
		ws.CallbacksLock.Lock()
		delete(ws.Callbacks, id)
		ws.CallbacksLock.Unlock()

		return response, nil
	case <-time.After(30 * time.Second):
		ws.CallbacksLock.Lock()
		delete(ws.Callbacks, id)
		ws.CallbacksLock.Unlock()

		return nil, errors.New("request timeout")
	}
}
func HandleProxyGet(c *gin.Context) {
	method := c.Param("method")
	queryParams := c.Request.URL.Query()

	var params []interface{}
	if paramsStr := queryParams.Get("params"); paramsStr != "" {
		var paramsInterface interface{}
		if err := json.Unmarshal([]byte(paramsStr), &paramsInterface); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		params, _ = paramsInterface.([]interface{})
	}

	response, err := handleRequestToWs(method, params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newResponse := map[string]interface{}{
		"success":  true,
		"response": response.Result,
	}
	c.JSON(http.StatusOK, newResponse)
}

func HandleProxyPost(c *gin.Context) {
	method := c.Param("method")
	var body struct {
		Params []interface{} `json:"params"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := handleRequestToWs(method, body.Params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newResponse := map[string]interface{}{
		"success":  true,
		"response": response.Result,
	}
	c.JSON(http.StatusOK, newResponse)
}
