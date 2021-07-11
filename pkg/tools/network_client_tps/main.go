package network_client_tps

import (
	"runtime"
	"sync/atomic"
	"time"

	"github.com/iwanhae/kubenchctl/pkg/types"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"golang.org/x/net/context"
	"golang.org/x/sync/errgroup"
)

var (
	url       *string
	connCount *int
	duration  *time.Duration
	timeout   *time.Duration
)

var (
	runningCount int32
)

func NetworkClientCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "client-tps",
		Aliases: []string{"ct", "client"},
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
	// 16MB shared buff
	buff := make([]byte, 16*1024*1024)
	c.MaxResponseBodySize = len(buff)

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

					status, body, err := c.GetTimeout(buff, *url, *timeout)
					hr.Duration = time.Since(t)
					if err != nil {
						return err
					}
					hr.Status = status
					hr.BodySize = int64(len(body))
					hr.Print()
				}
			}
		})
	}
	err := eg.Wait()
	c.CloseIdleConnections()
	return err
}
