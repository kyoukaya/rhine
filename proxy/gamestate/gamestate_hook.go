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

// Unhook unhooks the hook from the gamestate. Will fail silently if the hook
// is already unhooked.
func (oldHook *GameStateHook) Unhook() {
	if oldHook == nil {
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
	oldHook.gs.stateHooks[oldHook.target] = append(oldHooks[:i], oldHooks[i+1:]...)
}

// Caller is assumed to be holding the mutex.
func (mod *GameState) parseHookQueue() {
	for {
		var hook *GameStateHook
		select {
		case hook = <-mod.hookQueue:
			mod.stateHooks[hook.target] = append(mod.stateHooks[hook.target], hook)
		default:
			return
		}
	}
}

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
