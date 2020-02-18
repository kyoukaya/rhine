package proxy

import (
	"github.com/elazarl/goproxy"
)

// PacketHook contains information about the hook and allows for execution of the underlying
// Handle and Shutdown methods in the PacketHandler.
type PacketHook struct {
	target   string
	priority int
	handler  PacketHandler
	mod      *RhineModule
}

// handle calls the handle method of the underlying PacketHandler.
func (hook *PacketHook) handle(op string, data []byte, pktCtx *goproxy.ProxyCtx) []byte {
	return hook.handler(op, data, pktCtx)
}

// Unhook will unhook the receiving PacketHook if it's hooked. Otherwise, it
// will fail silently.
func (hook *PacketHook) Unhook() {
	if hook == nil {
		// Fail silently
		return
	}
	hook.mod.dispatch.removeHook(hook)
}

// PacketHandler represents handler functions exposed by a module.
type PacketHandler func(op string, data []byte, pktCtx *goproxy.ProxyCtx) []byte

func (d *dispatch) insertHook(hook *PacketHook) {
	var hookSlice []*PacketHook

	hookSlice, ok := d.hooks[hook.target]
	if !ok {
		hookSlice = make([]*PacketHook, 0)
	}
	hookSlice = append(hookSlice, hook)
	d.hooks[hook.target] = hookSlice
}
