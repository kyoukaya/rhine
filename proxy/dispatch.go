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

func (d *Dispatch) dispatch(op string, data []byte, ctx *goproxy.ProxyCtx) (*http.Request, *http.Response) {
	// Run core handlers
	for _, hook := range d.coreHandlers {
		hook(op, data, ctx)
	}
	// Run wildcard hooks
	if hooks, ok := d.Hooks["*"]; ok {
		for _, hook := range hooks {
			data = d.hookWrapper(hook, op, data, ctx)
		}
	}
	// Run normal hooks for op
	if hooks, ok := d.Hooks[op]; ok {
		for _, hook := range hooks {
			data = d.hookWrapper(hook, op, data, ctx)
		}
	}
	return ctx.Req, ctx.Resp
}

// Wrap hook handlers in a recover so we don't crash the entire proxy if it a
// module throws a panic.
func (d *Dispatch) hookWrapper(hook *PacketHook, op string, data []byte, ctx *goproxy.ProxyCtx) []byte {
	defer func() {
		if err := recover(); err != nil {
			d.Warnf("Recovered from panic while executing %s:\n%+v", hook.name, err)
		}
	}()
	return hook.Handle(op, data, ctx)
}

// Dispatch contains all the state pertaining to an authenticated user connected with
// the game servers.
type Dispatch struct {
	mutex        *sync.Mutex
	UID          int
	Region       string
	Hooks        map[string][]*PacketHook
	coreHandlers []func(string, []byte, *goproxy.ProxyCtx)
	shutdownCBs  []ShutdownCb
	log.Logger
	State *gamestate.GameState
}

func (d *Dispatch) initMods(mods []*rhineModule) {
	startT := time.Now()
	// Load core modules
	gs, gsHandler := gamestate.New(d.Logger)
	d.State = gs
	d.coreHandlers = append(d.coreHandlers, gsHandler)
	// Load user modules
	for _, mod := range mods {
		hooks, cb := mod.initFunc(d)
		d.initHooks(hooks)
		d.shutdownCBs = append(d.shutdownCBs, cb)
		d.Infof("%s loaded.", mod.name)
	}
	d.sortHooks()
	d.Verbosef("Mods loaded in %dms", time.Since(startT).Milliseconds())
}

type byPriority []*PacketHook

func (a byPriority) Len() int           { return len(a) }
func (a byPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a byPriority) Less(i, j int) bool { return a[i].priority > a[j].priority }

// Sort hooks in descending order
func (d *Dispatch) sortHooks() {
	for _, v := range d.Hooks {
		sort.Sort(byPriority(v))
	}
}

func (d *Dispatch) initHooks(hooks []*PacketHook) {
	if d == nil {
		return
	}
	for _, hook := range hooks {
		destMap := d.Hooks
		if err := insertHook(destMap, hook); err != nil {
			d.Warnf("Error while loading %#v: %v", hook, err)
		}
	}
}
