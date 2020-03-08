package gamestate

type GameStateHook struct {
	target     string
	moduleName string
	listener   chan StateEvent
	gs         *GameState
	event      bool
}

type StateEvent struct {
	Path    string
	Payload interface{}
}

// Unhook unhooks the hook from the gamestate.
func (oldHook *GameStateHook) Unhook() {
	if oldHook == nil {
		oldHook.gs.log.Warnf("GameStateHook.Unhook called on nil receiver")
		return
	}
	oldHook.gs.stateMutex.Lock()
	defer oldHook.gs.stateMutex.Unlock()
	oldHooks := oldHook.gs.stateHooks[oldHook.target]
	i := 0
	for _, hook := range oldHooks {
		if hook == oldHook {
			break
		}
		i++
	}
	if i == len(oldHooks) {
		// Hook not found
		return
	}
	oldHook.gs.stateHooks[oldHook.target] = append(oldHooks[:i], oldHooks[i+1:]...)
}

func (mod *GameState) parseHookQueue() {
	for hook := range mod.hookQueue {
		mod.stateMutex.Lock()
		mod.stateHooks[hook.target] = append(mod.stateHooks[hook.target], hook)
		mod.stateMutex.Unlock()
	}
}

// Hook creates a GameStateHook and attaches it as soon as possible. Notably, users
// should not expect the hook to be attached when the function returns as the attaching
// is done in a separate goroutine, allowing users to hook without blocking when
// the module is initialized on account/login, i.e., before game state is initialized
// from the SyncData packet.
func (mod *GameState) Hook(target, moduleName string, listener chan StateEvent, event bool) *GameStateHook {
	hook := &GameStateHook{
		target:     target,
		moduleName: moduleName,
		listener:   listener,
		gs:         mod,
		event:      event,
	}
	select {
	case mod.hookQueue <- hook:
	default:
		mod.log.Warnf("gamestate: failed to add hook into hook queue %#v", hook)
		return nil
	}
	return hook
}
