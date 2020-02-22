package proxy

import (
	"github.com/kyoukaya/rhine/proxy/gamestate"
	"github.com/kyoukaya/rhine/proxy/gamestate/statestruct"
)

// RhineModule provides modules with an interface to Rhine, allowing them to
// register hooks for events.
type RhineModule struct {
	Region string
	UID    int

	name        string
	initialized bool
	shutdownCB  ShutdownCb
	hooks       []*PacketHook
	gameState   *gamestate.GameState
	*dispatch
}

type initFunc struct {
	name string
	fun  ModuleInitFunc
}

var (
	modules []initFunc
)

// ShutdownCb will be called when the proxy is shutting down or when a user reconnects.
type ShutdownCb func(shuttingDown bool)

// ModuleInitFunc will be called when a user authenticates with the server to
// allow the module to setup its hooks and internal state for an individual user.
type ModuleInitFunc func(*RhineModule)

// Hooker is the generic Hook interface which exposes a single Unhook method which accepts
// and returns nothing. Receivers who implement this interface should fail silently
// if Unhook is called on a nil or already unhooked hook.
type Hooker interface {
	Unhook()
}

// Hook registers a new packet hook whose PacketHandler will be called back when
// the specified target packet is received by Rhine. A Hooker is returned, allowing
// the caller to Unhook the hook to stop receiving callbacks.
// The current implementation sorts the hooks to maintain priority ordering,
// while this isn't the most efficient, especially after the initial hooking is done
// when all the modules are initialized, doing a binary search and bisecting would
// result in a lot of expensive copying anyway.
func (m *RhineModule) Hook(target string, priority int, handler PacketHandler) Hooker {
	hook := &PacketHook{target, priority, handler, m}
	m.hooks = append(m.hooks, hook)
	m.dispatch.insertHook(hook)
	return hook
}

// StateHook registers a new game state hook whose listener chan will be notified
// when the specified game state has been modified. The StateEvent passed through
// the chan will include the new state at the path if the wantPayload bool is set
// to true.
func (m *RhineModule) StateHook(path string, listener chan gamestate.StateEvent, wantPayload bool) Hooker {
	return m.gameState.Hook(path, m.name, listener, wantPayload)
}

// OnShutdown registers a void function which accepts a boolean argument to be called
// back the program is killed with SIGINT or when an Arknights user reconnects.
// The boolean argument will be set to true if the callback is initiated because
// of a SIGINT event, and false if it's a user reconnecting.
func (m *RhineModule) OnShutdown(cb ShutdownCb) {
	m.shutdownCB = cb
}

// GetGameState will block until the gamestate module finishes parsing S/account/syncData.
func (m *RhineModule) GetGameState() *statestruct.User {
	return m.gameState.GetStateRef()
}
