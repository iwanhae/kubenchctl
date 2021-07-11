package types

import (
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Kind string

const (
	KindMessage     = "message"
	KindServerStats = "server_stats"
	KindHTTPReport  = "cluster_http_report"
	KindTCPReport   = "cluster_tcp_report"
)

type format struct {
	Time time.Time   `json:"time"`
	Kind Kind        `json:"kind"`
	Data interface{} `json:"data"`
}

var (
	printWait, printMu sync.Mutex
)

var (
	buff   = make([][]byte, 1_000)
	offset int32
)

func printByte(b []byte) {
	o := atomic.AddInt32(&offset, 1)
	if o == 1 {
		go print()
	}
	o %= int32(len(buff))
	printMu.Lock()
	buff[o] = b
	printMu.Unlock()
}

func print() {
	var ptr, p int32
	printWait.Lock()
	for {
		for ptr == atomic.LoadInt32(&offset) {
			printWait.Unlock()
			time.Sleep(50 * time.Millisecond)
			printWait.Lock()
		}
		p = ptr % int32(len(buff))
		os.Stdout.Write(buff[p])
		os.Stdout.Write([]byte{'\n'})
		ptr += 1
	}
}

func WaitPrint() {
	printWait.Lock()
}
