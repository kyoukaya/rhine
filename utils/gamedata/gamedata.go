// Package gamedata provides Arknights gamedata lookup datastructures parsed from
// https://github.com/Kengxxiao/ArknightsGameData.
// The package will automatically query the ArknightsGameData github repository
// will update the local files if the files are different from the local files.
package gamedata

import (
	"errors"
	"fmt"
	"sync"

	"github.com/kyoukaya/rhine/log"
	"github.com/kyoukaya/rhine/utils/gamedata/itemtable"
	"github.com/kyoukaya/rhine/utils/gamedata/stagetable"
)

type gameDataState struct {
	stageTableMap map[string]*stagetable.StageTable
	itemTableMap  map[string]*itemtable.ItemTable
}

// GameData provides methods to get data structures that contain game related
// data.
type GameData struct {
	region string
}

var (
	// ErrInvalidRegion is returned if the specified region is invalid.
	ErrInvalidRegion = errors.New("Invalid region")
	state            *gameDataState
	stateMutex       sync.Mutex
	updated          bool
	regionMap        = map[string]string{
		"GL": "en_US",
		"JP": "ja_JP",
		"KR": "ko_KR",
	}
)

// New creates a new GameData struct, may return an error if an invalid region
// is provided. Refer to proxy.regionMap for valid region strings.
func New(region string, logger log.Logger) (*GameData, error) {
	if region != "GL" && region != "JP" && region != "KR" {
		return nil, ErrInvalidRegion
	}
	stateMutex.Lock()
	if state == nil {
		state = &gameDataState{
			stageTableMap: make(map[string]*stagetable.StageTable),
			itemTableMap:  make(map[string]*itemtable.ItemTable),
		}
	}
	stateMutex.Unlock()

	fileMutex.Lock()
	if !updated {
		go updateGameData(logger)
	} else {
		fileMutex.Unlock()
	}
	return &GameData{region}, nil
}

// GetStageInfo provides a reference to the StageTable struct which contains
// information about game stages. This call will block if the gamedata has not
// been loaded yet. Region is an optional argument, by default it will
// use the region associated with the GameData receiver.
func (d *GameData) GetStageInfo(region ...string) (*stagetable.StageTable, error) {
	var regionName string
	if len(region) > 0 {
		if _, exists := regionMap[region[0]]; !exists {
			return nil, ErrInvalidRegion
		}
		regionName = region[0]
	} else {
		regionName = d.region
	}
	stateMutex.Lock()
	defer stateMutex.Unlock()
	if _, exists := state.stageTableMap[regionName]; !exists {
		d.loadStageTable(regionName)
	}
	return state.stageTableMap[regionName], nil
}

// GetItemInfo provides a reference to the ItemTable struct which contains
// information about items. This call will block if the gamedata has not
// been loaded yet. Region is an optional argument, by default it will
// use the region associated with the GameData receiver.
func (d *GameData) GetItemInfo(region ...string) (*itemtable.ItemTable, error) {
	var regionName string
	if len(region) > 0 {
		if _, exists := regionMap[region[0]]; !exists {
			return nil, ErrInvalidRegion
		}
		regionName = region[0]
	} else {
		regionName = d.region
	}
	stateMutex.Lock()
	defer stateMutex.Unlock()
	if _, exists := state.itemTableMap[regionName]; !exists {
		d.loadItemTable(regionName)
	}
	return state.itemTableMap[regionName], nil
}

// ErrPathOutOfBounds is returned when the fileName specified breaks out the directory.
var ErrPathOutOfBounds = fmt.Errorf("Path specified breaks out of data directory")
