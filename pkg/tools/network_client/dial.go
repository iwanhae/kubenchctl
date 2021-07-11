package network_client

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"
	"time"

	"github.com/iwanhae/kubenchctl/pkg/types"
	"github.com/valyala/fasthttp"
)

var (
	defaultDialer       = &fasthttp.TCPDialer{Concurrency: 1000}
	dnsCount      int32 = 0
)

func DefaultDialer(addr string) (net.Conn, error) {

	tr := types.TCPRequestReport{}

	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return nil, err
	}
	if host == "" { // localhost
		host = "127.0.0.1"
	}

	ip := net.ParseIP(host)
	if ip == nil {
		t := time.Now()
		IPs, err := net.DefaultResolver.LookupIPAddr(context.TODO(), host)
		tr.DNSResolvingDuration = time.Since(t)
		if err != nil {
			return nil, err
		}
		if len(IPs) == 0 {
			return nil, fmt.Errorf("fail to resolve %q", host)
		}
		c := dnsCount
		atomic.AddInt32(&dnsCount, 1)
		ip = IPs[int(c)%len(IPs)].IP
	}

	addr = fmt.Sprintf("%s:%s", ip.String(), port)
	t := time.Now()
	net, err := defaultDialer.Dial(addr)
	tr.ConnectionEstablishDuration = time.Since(t)
	if err != nil {
		return nil, err
	}
	tr.Print()
	return net, nil
}
