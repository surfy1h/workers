package workers

import (
	"fmt"
	"io"
	"net/http"
	"syscall/js"

	"github.com/syumai/workers/internal/jsutil"
)

var httpHandler http.Handler

func init() {
	var handleRequestCallback js.Func
	handleRequestCallback = js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			panic(fmt.Errorf("too many args given to handleRequest: %d", len(args)))
		}
		var cb js.Func
		cb = js.FuncOf(func(_ js.Value, pArgs []js.Value) any {
			defer cb.Release()
			resolve := pArgs[0]
			go func() {
				res, err := handleRequest(args[0])
				if err != nil {
					panic(err)
				}
				resolve.Invoke(res)
			}()
			return js.Undefined()
		})
		return jsutil.NewPromise(cb)
	})
	jsutil.Global.Set("handleRequest", handleRequestCallback)
}

// handleRequest accepts a Request object and returns Response object.
func handleRequest(reqObj js.Value) (js.Value, error) {
	if httpHandler == nil {
		return js.Value{}, fmt.Errorf("Serve must be called before handleRequest.")
	}
	req, err := toRequest(reqObj)
	if err != nil {
		panic(err)
	}
	reader, writer := io.Pipe()
	w := &responseWriterBuffer{
		header:     http.Header{},
		statusCode: http.StatusOK,
		reader:     reader,
		writer:     writer,
		readyCh:    make(chan struct{}),
	}
	go func() {
		defer w.ready()
		defer writer.Close()
		httpHandler.ServeHTTP(w, req)
	}()
	return toJSResponse(w)
}

// Server serves http.Handler on Cloudflare Workers.
// if the given handler is nil, http.DefaultServeMux will be used.
func Serve(handler http.Handler) {
	if handler == nil {
		handler = http.DefaultServeMux
	}
	httpHandler = handler
	jsutil.Global.Call("ready")
	select {}
}
