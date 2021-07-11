package network_client

import (
	"runtime"
	"sync/atomic"
	"time"

	"github.com/iwanhae/kubenchctl/pkg/types"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"golang.org/x/sync/errgroup"
)

var (
	url        *string
	connCount  *int
	reqNum     *int
	timeout    *time.Duration
	keepAlived *bool
)

var (
	runningCount int32
)

func NetworkClientCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use: "client",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runClient()
			types.WaitPrint()
			return err
		},
	}
	f := cmd.Flags()
	url = f.String("url", "http://127.0.0.1:8080/", "Server URL to request")
	connCount = f.IntP("concurrency", "c", 10, "Number of multiple requests to make at a time")
	reqNum = f.IntP("requests", "r", 1000, "Number of requests to perform (-1 is unlimited)")
	timeout = f.DurationP("timeout", "t", 10*time.Second, "Timeout value")
	keepAlived = f.Bool("keep-alive", true, "Use established connection. (or it will close and reopen connection everytime)")
	return cmd
}
func runClient() error {
	runtime.GOMAXPROCS(-1)
	c := fasthttp.Client{
		Name:            "netbench-client",
		MaxConnsPerHost: *connCount,
		ReadTimeout:     *timeout,
		WriteTimeout:    *timeout,
		Dial:            DefaultDialer,
	}
	if !*keepAlived {
		c.MaxConnDuration = time.Nanosecond
	}

	// shared buff
	buff := make([]byte, 1024*1024*1024)

	eg := errgroup.Group{}
	for i := 0; i < *connCount; i++ {
		eg.Go(func() error {
			hr := types.HTTPRequestReport{}
			for !IsFinished() {
				atomic.AddInt32(&runningCount, 1)
				t := time.Now()

				status, body, err := c.GetTimeout(buff, *url, *timeout)
				hr.Duration = time.Since(t)
				if err != nil {
					return err
				}
				hr.Status = status
				hr.BodySize = len(body)
				hr.Print()
			}
			return nil
		})
	}
	return eg.Wait()
}

func IsFinished() bool {
	if *reqNum < 0 {
		return false
	}
	if int(atomic.LoadInt32(&runningCount)) < *reqNum {
		return false
	}
	return true
}
