package gamestate

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/kyoukaya/go-lookup"
	"github.com/kyoukaya/rhine/proxy/gamestate/statestruct"
	"github.com/svyotov/mergo"
	"github.com/tidwall/gjson"
)

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func openAndRead(t *testing.T, filename string) []byte {
	b, err := ioutil.ReadFile(filename)
	check(t, err)
	return b
}

// TestHookWithPayload registers a hook which receives the value that has
// changed when an event is emitted.
func TestHookWithPayload(t *testing.T) {
	mod, _ := New(logShim{t})
	mod.stateMutex.Unlock()
	b := openAndRead(t, "testdata/syncdata.json")
	syncData, err := unmarshalSyncData(b)
	mod.state = &syncData.User
	check(t, err)
	b = openAndRead(t, "testdata/buildingsync.json")
	res := gjson.GetBytes(b, "playerDataDelta.modified")
	newUser, err := unmarshalUserData([]byte(res.Raw))
	check(t, err)
	err = mergo.Merge(mod.state, newUser, mergo.WithOverride)
	check(t, err)

	testChan := make(chan StateEvent, 1)
	go func() {
		for {
			evt := <-testChan
			fmt.Println(evt.Payload.(map[string]statestruct.EmptyStruct))
		}
	}()
	mod.Hook("building.rooms.ELEVATOR", "test", testChan, true)
	err = mod.WalkAndNotify([]byte(res.Raw))
	check(t, err)
}

type logShim struct{ t *testing.T }

func (l logShim) Flush()                              {}
func (l logShim) Println(i ...interface{})            { l.t.Log(i...) }
func (l logShim) Printf(s string, i ...interface{})   { l.t.Logf(s, i...) }
func (l logShim) Verboseln(i ...interface{})          { l.t.Log(i...) }
func (l logShim) Verbosef(s string, i ...interface{}) { l.t.Logf(s, i...) }
func (l logShim) Warnln(i ...interface{})             { l.t.Error(i...) }
func (l logShim) Warnf(s string, i ...interface{})    { l.t.Errorf(s, i...) }

func TestModificationHooks(t *testing.T) {
	const testPath = "dexNav.enemy.stage.camp_02"
	data := openAndRead(t, "testdata/buildingsync.json")
	b := []byte(gjson.GetBytes(data, "playerDataDelta.modified").Raw)
	testChan := make(chan StateEvent, 1)
	mod, _ := New(logShim{t})
	mod.stateMutex.Unlock()
	hook := mod.Hook(testPath, "test", testChan, false)
	mod.stateMutex.Lock()
	err := mod.WalkAndNotify(b)
	mod.stateMutex.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	var evt StateEvent
	select {
	case evt = <-testChan:
	default:
		t.Fatal()
	}
	if evt.Path != testPath {
		t.Fatal()
	}
	hook.Unhook()
	// Shouldn't be able to get path from testChan anymore after unhooking
	err = mod.WalkAndNotify(b)
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-testChan:
		t.Fatal()
	default:
	}
}

// map[string]interface{} did not add new entries in the map. Workaround with
// map[string]struct{} instead.
func TestMapInterfaceBug(t *testing.T) {
	b := openAndRead(t, "testdata/syncdata.json")
	syncData, err := unmarshalSyncData(b)
	user := syncData.User
	check(t, err)
	data := openAndRead(t, "testdata/buildingsync.json")
	b = []byte(gjson.GetBytes(data, "playerDataDelta.modified").Raw)
	newUser, err := unmarshalUserData(b)
	check(t, err)
	err = mergo.Merge(&user, newUser, mergo.WithOverride)
	check(t, err)
	if _, exists := user.Building.Rooms.Elevator["slot_9"]; !exists {
		t.Fatal()
	}
}

func TestStructAccess(t *testing.T) {
	b := openAndRead(t, "testdata/syncdata.json")
	syncData, err := unmarshalSyncData(b)
	check(t, err)
	user := syncData.User
	val, err := lookup.LookupString(user, "Building.Rooms", false)
	check(t, err)
	vali := val.Interface()
	unm, err := json.Marshal(vali)
	check(t, err)
	fmt.Printf("%s\n", unm)
}
