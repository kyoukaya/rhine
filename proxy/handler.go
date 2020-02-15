package proxy

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/kyoukaya/rhine/utils"
	"github.com/tidwall/gjson"

	"github.com/elazarl/goproxy"
)

var (
	gameHostMatcher = regexp.MustCompile(`^gs\.arknights\.(jp|global):8443$`)
)

// HandleReq processes an outgoing HTTP request, dispatching it if it's game traffic.
func (proxy *Proxy) HandleReq(req *http.Request, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	defer proxy.Flush()
	reqCtx := &RequestContext{}
	reqCtx.StartT = time.Now()
	ctx.UserData = reqCtx
	// Block telemetry requests
	if proxy.hostFilter != nil && proxy.hostFilter.MatchString(req.Host) {
		proxy.Verbosef("==== Rejecting %v", req.Host)
		// Use the UserData field as a flag to indicate to the response handler that the
		// request that generated the response was blocked.
		reqCtx.RequestIsBlocked = true
		return req, goproxy.NewResponse(req, goproxy.ContentTypeText, http.StatusOK, "")
	}
	// Return if not game traffic
	if !gameHostMatcher.MatchString(req.URL.Host) {
		return req, nil
	}
	body, err := ioutil.ReadAll(req.Body)
	utils.Check(err)
	req.Body = ioutil.NopCloser(bytes.NewReader(body))

	op := "C/" + strings.Trim(req.URL.Path, "/")
	uid := req.Header.Get("uid")
	region := regionMap[req.URL.Hostname()[13:]]
	var d *dispatch
	if uid == "" {
		if op != "C/account/login" {
			return req, nil
		}
		uid = gjson.GetBytes(body, "uid").String()
		d = proxy.addUser(uid, region)
	} else {
		d = proxy.getUser(uid, region)
	}
	if d == nil {
		return req, nil
	}
	reqCtx.dispatch = d
	reqCtx.RequestData = body
	reqCtx.RequestOp = op
	req, resp := d.dispatch(op, body, ctx)
	if proxy.options.Verbose {
		proxy.Verbosef(">>>> %s (%d)\n", op, time.Since(reqCtx.StartT).Milliseconds())
	}
	return req, resp
}

// HandleResp processes an incoming http(s) response.
func (proxy *Proxy) HandleResp(resp *http.Response, ctx *goproxy.ProxyCtx) *http.Response {
	defer proxy.Flush()
	reqCtx := GetRequestContext(ctx)
	// If request that generated response was blocked or response not OK.
	if reqCtx == nil || resp == nil || reqCtx.RequestIsBlocked {
		return resp
	}
	body, err := ioutil.ReadAll(resp.Body)
	utils.Check(err)
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	// Game traffic
	if reqCtx.dispatch != nil {
		recvT := time.Now()
		op := "S/" + strings.Trim(ctx.Req.URL.Path, "/")
		_, resp := reqCtx.dispatch.dispatch(op, body, ctx)
		proxy.Verbosef("<<<< %s (%d,%d)\n", op, recvT.Sub(reqCtx.StartT).Milliseconds(), time.Since(recvT).Milliseconds())
		return resp
	}
	return resp
}
