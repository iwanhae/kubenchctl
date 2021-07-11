package types

import (
	"encoding/json"
	"time"
)

type HTTPRequestReport struct {
	Status   int           `json:"status"`
	BodySize int64         `json:"body_size"`
	Duration time.Duration `json:"duration"`
}

func (m HTTPRequestReport) Print() {
	b, _ := json.Marshal(
		format{
			Kind: KindHTTPReport,
			Time: time.Now(),
			Data: m,
		})
	printByte(b)
}

type TCPRequestReport struct {
	DNSResolvingDuration        time.Duration `json:"dns_resolving_duration"`
	ConnectionEstablishDuration time.Duration `json:"conn_establish_duration"`
}

func (m TCPRequestReport) Print() {
	b, _ := json.Marshal(
		format{
			Kind: KindTCPReport,
			Time: time.Now(),
			Data: m,
		})
	printByte(b)
}
