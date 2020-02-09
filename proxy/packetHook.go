package proxy

import (
	"github.com/elazarl/goproxy"
)

// PacketHook contains information about the hook and allows for execution of the underlying
// Handle and Shutdown methods in the PacketHandler.
type PacketHook struct {
	name     string
	target   string
	priority int
	handler  PacketHandler
}

// Handle calls the Handle method of the underlying PacketHandler.
func (hook *PacketHook) Handle(op string, data []byte, pktCtx *goproxy.ProxyCtx) []byte {
	return hook.handler(op, data, pktCtx)
}

// PacketHandler represents handler functions exposed by a module.
type PacketHandler func(op string, data []byte, pktCtx *goproxy.ProxyCtx) []byte

// NewPacketHook returns an initialized PacketHook struct.
func NewPacketHook(name, target string, priority int, handler PacketHandler) *PacketHook {
	return &PacketHook{name, target, priority, handler}
}

func insertHook(hookMap map[string][]*PacketHook, hook *PacketHook) error {
	var hookSlice []*PacketHook
	hookSlice, ok := hookMap[hook.target]
	if !ok {
		hookSlice = make([]*PacketHook, 0)
	}
	hookSlice = append(hookSlice, hook)
	hookMap[hook.target] = hookSlice
	return nil
}
