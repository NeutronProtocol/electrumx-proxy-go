package router

import (
	"electrumx-proxy-go/common"
	"electrumx-proxy-go/ws"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	cors "github.com/itsjamie/gin-cors"
	"log"
	"net/http"
	"sync/atomic"
	"time"
)

var UniqueID uint32

func GetUniqueID() uint32 {
	return atomic.AddUint32(&UniqueID, 1)
}
func InitMasterRouter() *gin.Engine {
	r := gin.Default()
	gin.SetMode(gin.ReleaseMode)
	r.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Content-Type",
		MaxAge:          50 * time.Second,
		Credentials:     false,
		ValidateHeaders: false,
	}))
	r.GET("/", HandlerHello)
	r.GET("/proxy", HandleProxy)
	r.GET("/proxy/health", HandleHealth)
	r.POST("/proxy/health", HandleHealth)
	r.GET("/proxy/:method", HandleProxyGet)
	r.POST("/proxy/:method", HandleProxyPost)
	return r
}

func HandlerHello(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "Hello from Atomical, wish you have a pleasant user experience. --Atomical Neutron Protocol",
	})
}

func HandleProxy(c *gin.Context) {
	info := common.Info{
		Note: "Atomicals Neutron ElectrumX Digital Object Proxy Online",
		UsageInfo: common.UsageInfo{
			Note: "The service offers both POST and GET requests for proxying requests to ElectrumX. To handle larger broadcast transaction payloads use the POST method instead of GET.",
			POST: "POST /proxy/:method with string encoded array in the field \"params\" in the request body.",
			GET:  "GET /proxy/:method?params=[\"value1\"] with string encoded array in the query argument \"params\" in the URL.",
		},
		HealthCheck: "GET /proxy/health",
		Github:      "https://www.atomicalneutron.com",
		License:     "MIT",
	}

	response := common.Result{
		Success: true,
		Info:    info,
	}
	c.JSON(http.StatusOK, response)
}

type ResponseHealth struct {
	Success bool `json:"success"`
	Health  bool `json:"health,omitempty"`
}

func HandleHealth(c *gin.Context) {
	id := GetUniqueID()
	responseCh := make(chan *common.JsonRpcResponse)
	ws.CallbacksLock.Lock()
	ws.Callbacks[id] = responseCh
	ws.CallbacksLock.Unlock()
	request := common.JsonRpcRequest{
		ID:     id,
		Method: "blockchain.atomicals.get_global",
		Params: []interface{}{},
	}
	requestText, err := json.Marshal(request)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to marshal JSON-RPC request"})
		return
	}

	err = ws.WsTx.WriteMessage(websocket.TextMessage, requestText)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send message over WebSocket"})
		return
	}

	select {
	case rep := <-responseCh:
		ws.CallbacksLock.Lock()
		delete(ws.Callbacks, id)
		ws.CallbacksLock.Unlock()
		if rep.Result != nil {
			c.JSON(http.StatusOK, ResponseHealth{Success: true, Health: true})
		} else {
			c.JSON(http.StatusOK, ResponseHealth{Success: true, Health: false})
		}
	case <-time.After(5 * time.Second):
		log.Printf("<= id: %d, check health timeout, no response received after 5 seconds", id)
		ws.CallbacksLock.Lock()
		delete(ws.Callbacks, id)
		ws.CallbacksLock.Unlock()
		c.JSON(http.StatusOK, ResponseHealth{Success: true, Health: false})
	}
}
