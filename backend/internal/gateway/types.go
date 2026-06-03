package gateway

import "encoding/json"

type WbRequest struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

type WbResponse struct {
	Type    string `json:"type"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

func toJSON(v any) []byte {
	data, _ := json.Marshal(v)
	return data
}
