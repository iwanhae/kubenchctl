package network_server

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync/atomic"
	"time"

	"github.com/iwanhae/kubenchctl/pkg/types"
	"github.com/spf13/cobra"
	"github.com/valyala/fasthttp"
	"github.com/valyala/fasthttp/expvarhandler"
)

// Flags
var (
	addr *string
)

var (
	total, total_conns, active_conn, idle_conn int64
	// 16MB
	emptyBytes = make([]byte, 16*1024*1024)
)

func NetworkServerCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use: "server",
		RunE: func(cmd *cobra.Command, args []string) error {
			types.NewMessagef("Listening on %s", *addr).Print()
			go func() {
				var old int64
				for {
					t := atomic.LoadInt64(&total)
					active_c := atomic.LoadInt64(&active_conn)
					types.ServerStats{
						TotalRequests: t,
						NewRequests:   t - old,

						TotalConnections:  total_conns,
						ActiveConnections: active_c,
					}.Print()
					old = t
					time.Sleep(time.Second)
				}
			}()
			runServer()
			return nil
		},
	}
	f := cmd.Flags()
	addr = f.String("addr", ":8080", "TCP address to listen to")
	return cmd
}

func runServer() {
	s := fasthttp.Server{
		Handler: requestHandler,
		ConnState: func(c net.Conn, cs fasthttp.ConnState) {
			switch cs {
			case fasthttp.StateNew:
				atomic.AddInt64(&total_conns, 1)
				atomic.AddInt64(&active_conn, 1)
				// types.NewMessagef("Connection Established %q", c.RemoteAddr().String()).Print()
			case fasthttp.StateActive:
			case fasthttp.StateIdle:
			case fasthttp.StateClosed:
				atomic.AddInt64(&active_conn, -1)
				// types.NewMessagef("Connection Closed %q", c.RemoteAddr().String()).Print()
			}

		},
	}

	if err := s.ListenAndServe(*addr); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func requestHandler(ctx *fasthttp.RequestCtx) {
	atomic.AddInt64(&total, 1)
	switch string(ctx.Path()) {
	case "/stats":
		expvarhandler.ExpvarHandler(ctx)
	case "/":
		dummyHandler(ctx, 0)
	case "/1KB":
		dummyHandler(ctx, 1024)
	case "/10KB":
		dummyHandler(ctx, 1024*10)
	case "/100KB":
		dummyHandler(ctx, 1024*100)
	case "/1MB":
		dummyHandler(ctx, 1024*1024)
	case "/10MB":
		dummyHandler(ctx, 1024*1024*10)
	case "/100MB":
		dummyHandler(ctx, 1024*1024*100)
	case "/1GB":
		dummyHandler(ctx, 1024*1024*1024)
	case "/10GB":
		dummyHandler(ctx, 1024*1024*1024*10)
	default:
		echoHandler(ctx)
	}
}

func dummyHandler(ctx *fasthttp.RequestCtx, length int) {
	ctx.Response.SetBodyStream(&zeroReader{Length: length}, length)
}

func echoHandler(ctx *fasthttp.RequestCtx) {
	fmt.Fprintf(ctx, "Hello, world!\n\n")

	fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
	fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
	fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
	fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
	fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
	fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
	fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
	fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
	fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
	fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

	fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)

	ctx.SetContentType("text/plain; charset=utf8")

	// Set arbitrary headers
	ctx.Response.Header.Set("X-My-Header", "my-header-value")

	// Set cookies
	var c fasthttp.Cookie
	c.SetKey("cookie-name")
	c.SetValue("cookie-value")
	ctx.Response.Header.SetCookie(&c)
}

type zeroReader struct {
	Length int
	offset int
}

func (z *zeroReader) Read(p []byte) (n int, err error) {
	if z.Length == z.offset {
		return 0, io.EOF
	}
	size := z.Length - z.offset
	if len(emptyBytes) < size {
		size = len(emptyBytes)
	}
	n = copy(p, emptyBytes[:size])
	z.offset += n
	err = nil
	return
}
