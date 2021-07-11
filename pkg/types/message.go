package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type Message struct {
	Message string `json:"message"`
}

func (m Message) Print() {
	b, _ := json.Marshal(
		format{
			Kind: KindMessage,
			Time: time.Now(),
			Data: m,
		})
	printByte(b)
}

func NewMessagef(msg string, args ...interface{}) Message {
	return Message{
		Message: fmt.Sprintf(msg, args...),
	}
}
