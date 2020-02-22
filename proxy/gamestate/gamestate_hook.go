package gamestate

type GameStateHook struct {
	path        string
	moduleName  string
	listener    chan StateEvent
	gs          *GameState
	wantPayload bool
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
	oldHooks := oldHook.gs.stateHooks[oldHook.path]
	i := 0
	for _, hook := range oldHooks {
		if hook == oldHook {
			break
		}
		i++
	}
	oldHook.gs.stateHooks[oldHook.path] = append(oldHooks[:i], oldHooks[i+1:]...)
}

// Caller is assumed to be holding the mutex.
func (mod *GameState) parseHookQueue() {
	for {
		var hook *GameStateHook
		select {
		case hook = <-mod.hookQueue:
			mod.stateHooks[hook.path] = append(mod.stateHooks[hook.path], hook)
		default:
			return
		}
	}
}

func (mod *GameState) Hook(path, moduleName string, listener chan StateEvent, wantPayload bool) *GameStateHook {
	hook := &GameStateHook{
		path:        path,
		moduleName:  moduleName,
		listener:    listener,
		gs:          mod,
		wantPayload: wantPayload,
	}
	select {
	case mod.hookQueue <- hook:
	default:
		mod.log.Warnf("gamestate: failed to add hook into hook queue %#v", hook)
		return nil
	}
	return hook
}
