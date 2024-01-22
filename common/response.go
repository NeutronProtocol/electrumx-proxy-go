package common

import (
	"encoding/json"
)

// Result is the response structure for JSON output.
type Result struct {
	Success bool        `json:"success"`
	Info    interface{} `json:"info"`
}

// Info contains detailed information to be returned in the proxy response.
type Info struct {
	Note        string    `json:"note"`
	UsageInfo   UsageInfo `json:"usageInfo"`
	HealthCheck string    `json:"healthCheck"`
	Github      string    `json:"github"`
	License     string    `json:"license"`
}

// UsageInfo provides usage details for the proxy service.
type UsageInfo struct {
	Note string `json:"note"`
	POST string `json:"POST"`
	GET  string `json:"GET"`
}

type JsonRpcRequest struct {
	ID     uint32        `json:"id"`
	Method string        `json:"method"`
	Params []interface{} `json:"params"`
}

type JsonRpcResponse struct {
	Result *json.RawMessage `json:"result"`
	Error  *json.RawMessage `json:"error"`
	ID     uint32           `json:"id"`
}
