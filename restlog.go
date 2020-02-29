package zlog

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/caddyserver/caddy/caddyhttp/httpserver"
	"github.com/liuzl/filestore"
	"github.com/rs/zerolog"
)

var c Chain

func initZlog(dir, splitBy string, once sync.Once) {
	once.Do(func() {
		hostname, _ := os.Hostname()
		var out io.Writer
		f, err := filestore.NewFileStorePro(dir, splitBy)
		if err != nil {
			out = os.Stdout
			fmt.Fprintf(os.Stderr, "err: %+v, will zerolog to stdout\n", err)
		} else {
			out = f
		}
		log := zerolog.New(out).With().
			Timestamp().
			Str("service", filepath.Base(os.Args[0])).
			Str("host", hostname).
			Logger()

		c = NewChain()

		// Install the logger handler with default output on the console
		c = c.Append(NewHandler(log))

		c = c.Append(AccessHandler(func(r *http.Request,
			status, size int, duration time.Duration) {
			FromRequest(r).Debug().
				Str("method", r.Method).
				Str("url", r.URL.String()).
				Int("status", status).
				Int("size", size).
				Dur("duration", duration).
				Msg("")
		}))

		// Install some provided extra handler to set some request's context fields.
		// Thanks to those handler, all our logs will come with some pre-populated fields.
		c = c.Append(RemoteAddrHandler("server"))
		c = c.Append(HeaderHandler("X-Forwarded-For"))
		c = c.Append(HeaderHandler("User-Agent"))
		c = c.Append(HeaderHandler("Referer"))
		c = c.Append(RequestIDHandler("req_id", "Request-Id"))
		// keep in order
		c = c.Append(DelResponseHeaderHandler("Cost"))
		c = c.Append(ResponseHeaderHandler("Cost", "float"))
		c = c.Append(DumpResponseHandler("response"))
		c = c.Append(DumpRequestHandler("request"))
	})
}

func WithLog(h httpserver.Handler,
	dir, splitBy string, once sync.Once) httpserver.Handler {
	initZlog(dir, splitBy, once)
	return c.Then(h)
}
