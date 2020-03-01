// Package gamestate currently only marshals and stores the initial data sync, planned
// to mirror the client's state perfectly and enable other mods to hook onto it
// to query data or receive updates if values have changed.
package gamestate

import (
	"encoding/json"
	"sync"

	"github.com/kyoukaya/go-lookup"
	rhLog "github.com/kyoukaya/rhine/log"
	"github.com/kyoukaya/rhine/proxy/gamestate/statestruct"
	"github.com/svyotov/mergo"
	"github.com/tidwall/gjson"

	"github.com/elazarl/goproxy"
)

const hookQueueMax = 100

// GameState provides a handle in which users can obtain a reference to the
// gamestate struct.
type GameState struct {
	state      *statestruct.User
	stateMutex sync.Mutex
	log        rhLog.Logger
	loaded     bool
	// Hooks are first added into the hookQueue and then added into the
	// stateHooks map just before notifying listeners with parseHookQueue.
	hookQueue  chan *GameStateHook
	stateHooks map[string][]*GameStateHook
}

// New provides a newly instantiated GameState struct and a callback for the
// proxy to call on every game packet.
func New(log rhLog.Logger) (*GameState, func(string, []byte, *goproxy.ProxyCtx)) {
	gs := GameState{
		log:        log,
		stateHooks: make(map[string][]*GameStateHook),
		hookQueue:  make(chan *GameStateHook, hookQueueMax),
	}
	gs.stateMutex.Lock()
	return &gs, gs.handle
}

// IsLoaded checks if the initial sync packet has already been parsed and the
// gamestate instance is ready for use.
func (mod *GameState) IsLoaded() bool {
	return mod.loaded
}

// StateSync blocks until the gamestate is usable. The consistency provided by
// this function isn't strict, but considering the time between packets, and
// that that the gamestate is guaranteed to be at least as new as the seqnum
// of the packet context that StateSync was called on and any transitional
// states up to the seqnum of the latest packet context.
func (mod *GameState) StateSync() {
	mod.stateMutex.Lock()
	// No critical section necessary as we're just waiting for the parseDataDelta
	// routine to release the lock on the state so we know that any caller
	mod.stateMutex.Unlock() //nolint:staticcheck
}

// GetStateRef will block until the module finishes parsing S/account/syncData.
func (mod *GameState) GetStateRef() *statestruct.User {
	mod.stateMutex.Lock()
	defer mod.stateMutex.Unlock()
	return mod.state
}

// Get returns the value of the gamestate from the path specified. Path is a
// period separated string based on the JSON keys, see https://github.com/mcuadros/go-lookup
// for reference. Blocks until the state is ready.
func (mod *GameState) Get(path string) (interface{}, error) {
	mod.stateMutex.Lock()
	defer mod.stateMutex.Unlock()
	val, err := lookup.LookupString(mod.state, path, true)
	if err != nil {
		return nil, err
	}
	return val.Interface(), nil
}

func (mod *GameState) parseDataDelta(data []byte, op string) {
	defer mod.stateMutex.Unlock()
	res := gjson.GetBytes(data, "playerDataDelta.modified")
	if !res.Exists() {
		return
	}
	// Could do an unsafe string cast here for performance
	b := []byte(res.Raw)
	user, err := unmarshalUserData(b)
	if err != nil {
		mod.log.Warnf("Failed to unmarshal %s: %s", op, err.Error())
		return
	}
	err = mergo.Merge(mod.state, user, mergo.WithOverride)
	if err != nil {
		mod.log.Warnf("Failed to merge %s: %s", op, err.Error())
	}
	// Notify state listeners
	err = mod.WalkAndNotify(b)
	if err != nil {
		mod.log.Warnf("Error occurred while notifying game state listeners: %s",
			err.Error())
	}
}

func (mod *GameState) handle(op string, data []byte, pktCtx *goproxy.ProxyCtx) {
	if op == "S/account/syncData" {
		go mod.handleSyncData(data)
		mod.loaded = true
	} else if mod.loaded {
		mod.stateMutex.Lock()
		go mod.parseDataDelta(data, op)
	}
}

func (mod *GameState) handleSyncData(data []byte) []byte {
	defer mod.stateMutex.Unlock()
	syncData, err := unmarshalSyncData(data)
	if err != nil {
		mod.log.Warnln(err)
	}
	mod.state = &syncData.User
	if err != nil {
		mod.log.Warnln(err)
		return data
	}
	return data
}

func unmarshalSyncData(data []byte) (*syncData, error) {
	r := syncData{}
	err := json.Unmarshal(data, &r)
	return &r, err
}

func unmarshalUserData(data []byte) (*statestruct.User, error) {
	r := &statestruct.User{}
	err := json.Unmarshal(data, r)
	return r, err
}

type syncData struct {
	User statestruct.User `json:"user"`
	Ts   int64            `json:"ts"`
}
