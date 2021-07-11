package types

import (
	"encoding/json"
	"time"
)

type ServerStats struct {
	TotalRequests int64 `json:"total_requests"`
	NewRequests   int64 `json:"new_requests"`

	TotalConnections  int64 `json:"total_conns"`
	ActiveConnections int64 `json:"active_conn"`
}

func (m ServerStats) Print() {
	b, _ := json.Marshal(
		format{
			Kind: KindMessage,
			Time: time.Now(),
			Data: m,
		})
	printByte(b)
}
