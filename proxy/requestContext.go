package proxy

import (
	"fmt"
	"time"

	"github.com/elazarl/goproxy"
)

// RequestContext contains additional data to be stored in goproxy.ProxyCtx.UserData to
// provide request context to the response handler.
type RequestContext struct {
	// Was the request blocked by the proxy
	RequestIsBlocked bool
	RequestOp        string
	// Contains the request body. Nil if DispatchState is for a request.
	RequestData []byte
	// Start time for handling the request
	StartT time.Time
	// Contains the dispatch object for the corresponding user if this is
	// a response to a game request.
	dispatch *dispatch
}

// GetRequestContext returns the dispatch context for a goproxy.ProxyCtx, will panic if
// called on a goproxy.ProxyCtx not associated with game data. I.e., the request was not
// handled by proxy.HandleReq.
func GetRequestContext(ctx *goproxy.ProxyCtx) *RequestContext {
	reqCtx, ok := ctx.UserData.(*RequestContext)
	if reqCtx == nil || !ok {
		panic(fmt.Sprintf("Failed to cast ctx.UserData into *RequestContext"))
	}
	return reqCtx
}
