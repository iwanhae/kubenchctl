package network_client_bandwidth

import (
	"context"
	"io"
	"net/http"
	"runtime"
	"sync/atomic"
	"time"

	"github.com/iwanhae/kubenchctl/pkg/types"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

var (
	url               *string
	connCount         *int
	reqNum            *int
	duration          *time.Duration
	timeout           *time.Duration
	disableKeepAlives *bool
)

var (
	runningCount int32
)

func NetworkClientCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "client-bandwidth",
		Aliases: []string{"cb"},
		RunE: func(cmd *cobra.Command, args []string) error {
			err := runClient()
			types.WaitPrint()
			return err
		},
	}
	f := cmd.Flags()
	url = f.String("url", "http://127.0.0.1:8080/", "Server URL to request")
	connCount = f.IntP("concurrency", "c", 10, "Number of multiple requests to make at a time")
	duration = f.DurationP("time", "t", 10*time.Second, "How long do you want to test?")
	timeout = f.Duration("timeout", 10*time.Second, "Timeout value")
	disableKeepAlives = f.BoolP("DisableKeepAlive", "k", false, "disable keep alive")
	return cmd
}
func runClient() error {
	runtime.GOMAXPROCS(-1)
	c := http.Client{
		Transport: &http.Transport{
			Proxy:                 http.ProxyFromEnvironment,
			DialContext:           DefaultDialerfunc,
			ForceAttemptHTTP2:     true,
			MaxIdleConns:          100,
			IdleConnTimeout:       90 * time.Second,
			TLSHandshakeTimeout:   10 * time.Second,
			ExpectContinueTimeout: 1 * time.Second,
			DisableKeepAlives:     *disableKeepAlives,
		},
		Timeout: *timeout,
	}

	// shared buff
	buff := make([]byte, 16*1024*1024)

	ctx, _ := context.WithTimeout(context.Background(), *duration)
	eg, ctx := errgroup.WithContext(ctx)
	for i := 0; i < *connCount; i++ {
		eg.Go(func() error {
			hr := types.HTTPRequestReport{}
			for {
				select {
				case <-ctx.Done():
					return nil
				default:
					atomic.AddInt32(&runningCount, 1)
					t := time.Now()
					///////////////////////////
					res, err := c.Get(*url)
					if err != nil {
						return err
					}
					for {
						_, err := res.Body.Read(buff)
						if err == io.EOF {
							break
						}
						if err != nil {
							return err
						}
					}
					hr.Duration = time.Since(t)
					///////////////////////////
					hr.Status = res.StatusCode
					hr.BodySize = res.ContentLength
					hr.Print()
				}
			}
		})
	}
	err := eg.Wait()
	c.CloseIdleConnections()
	return err
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
