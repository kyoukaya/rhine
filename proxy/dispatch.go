package proxy

import (
	"net/http"
	"sort"
	"sync"
	"time"

	"github.com/elazarl/goproxy"
	"github.com/kyoukaya/rhine/log"
	"github.com/kyoukaya/rhine/proxy/gamestate"
)

// Dispatch contains all the state pertaining to an authenticated user connected with
// the game servers.
type dispatch struct {
	// Exported fields
	log.Logger

	// Private fields
	mutex        *sync.Mutex
	uid          int
	region       string
	hooks        map[string][]*PacketHook
	coreHandlers []func(string, []byte, *goproxy.ProxyCtx)
	modules      []*RhineModule
	intialized   bool

	// Core modules
	state *gamestate.GameState
}

func (d *dispatch) dispatch(op string, data []byte, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// Run core handlers
	for _, hook := range d.coreHandlers {
		hook(op, data, ctx)
	}
	// Run wildcard hooks
	if hooks, ok := d.hooks["*"]; ok {
		for _, hook := range hooks {
			data = d.hookWrapper(hook, op, data, ctx)
		}
	}
	// Run normal hooks for op
	if hooks, ok := d.hooks[op]; ok {
		for _, hook := range hooks {
			data = d.hookWrapper(hook, op, data, ctx)
		}
	}
	return ctx.Req, ctx.Resp
}

// Wrap hook handlers in a recover so we don't crash the entire proxy if it a
// module throws a panic.
func (d *dispatch) hookWrapper(hook *PacketHook, op string, data []byte, ctx *goproxy.ProxyCtx) []byte {
	defer func() {
		if err := recover(); err != nil {
			d.Warnf("Recovered from panic while executing %s:\n%+v", hook.mod.name, err)
		}
	}()
	return hook.handle(op, data, ctx)
}

func (d *dispatch) initMods(mods []initFunc) {
	startT := time.Now()
	// Load core modules
	gs, gsHandler := gamestate.New(d.Logger)
	d.state = gs
	d.coreHandlers = append(d.coreHandlers, gsHandler)
	// Load user modules
	for _, mod := range mods {
		newMod := &RhineModule{
			name:      mod.name,
			Region:    d.region,
			UID:       d.uid,
			gameState: gs,
			dispatch:  d,
		}
		mod.fun(newMod)
		d.modules = append(d.modules, newMod)
		newMod.initialized = true
		d.Printf("%s loaded.", mod.name)
	}
	d.sortHooks()
	d.Verbosef("Mods loaded in %dms", time.Since(startT).Milliseconds())
	d.intialized = true
}

func (d *dispatch) removeHook(oldHook *PacketHook) {
	hooks, ok := d.hooks[oldHook.target]
	if !ok {
		d.Warnf("Tried to remove hook that doesn't exist: %#v", oldHook)
		return
	}
	newHooks := make([]*PacketHook, 0)
	for _, hook := range hooks {
		if hook == oldHook {
			continue
		}
		newHooks = append(newHooks, hook)
	}
	d.hooks[oldHook.target] = newHooks
}

type byPriority []*PacketHook

func (a byPriority) Len() int           { return len(a) }
func (a byPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPriority) Less(i, j int) bool { return a[i].priority > a[j].priority }

// Sort hooks in descending order
func (d *dispatch) sortHooks() {
	for _, v := range d.hooks {
		sortHookSl(byPriority(v))
	}
}

func sortHookSl(hooks []*PacketHook) {
	sort.Sort(byPriority(hooks))
}
