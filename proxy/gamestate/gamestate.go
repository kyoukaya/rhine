// Package gamestate currently only marshals and stores the initial data sync, planned
// to mirror the client's state perfectly and enable other mods to hook onto it
// to query data or receive updates if values have changed.
package gamestate

import (
	"encoding/json"
	"sync"

	rhLog "github.com/kyoukaya/rhine/log"
	"github.com/kyoukaya/rhine/proxy/gamestate/statestruct"

	"github.com/elazarl/goproxy"
)

// GameState provides a handle in which users can obtain a reference to the
// gamestate struct.
type GameState struct {
	state      *statestruct.User
	stateMutex sync.Mutex
	log        rhLog.Logger
}

func (mod *GameState) handle(op string, data []byte, pktCtx *goproxy.ProxyCtx) {
	if op == "S/account/syncData" {
		go mod.handleSyncData(data)
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

// GetStateRef will block until the module finishes parsing S/account/syncData.
func (mod *GameState) GetStateRef() *statestruct.User {
	mod.stateMutex.Lock()
	defer mod.stateMutex.Unlock()
	return mod.state
}

// New provides a newly instantiated GameState struct and a callback for the
// proxy to call on every game packet.
func New(log rhLog.Logger) (*GameState, func(string, []byte, *goproxy.ProxyCtx)) {
	gs := GameState{log: log}
	gs.stateMutex.Lock()
	return &gs, gs.handle
}

func unmarshalSyncData(data []byte) (syncData, error) {
	r := syncData{}
	err := json.Unmarshal(data, &r)
	return r, err
}

type syncData struct {
	User statestruct.User `json:"user"`
	Ts   int64            `json:"ts"`
}
