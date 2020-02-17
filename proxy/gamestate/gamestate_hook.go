package gamestate

type GameStateHook struct {
	path       string
	moduleName string
	listener   chan string
	gs         *GameState
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
	newHooks := make([]*GameStateHook, 0, len(oldHooks))
	for _, hook := range oldHooks {
		if hook == oldHook {
			continue
		}
		newHooks = append(newHooks, hook)
	}
	oldHook.gs.stateHooks[oldHook.path] = newHooks
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

func (mod *GameState) Hook(path, moduleName string, listener chan string) *GameStateHook {
	hook := &GameStateHook{
		path:       path,
		moduleName: moduleName,
		listener:   listener,
		gs:         mod,
	}
	select {
	case mod.hookQueue <- hook:
	default:
		mod.log.Warnf("gamestate: failed to add hook into hook queue %#v", hook)
		return nil
	}
	return hook
}
