package gamestate

import (
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/kyoukaya/rhine/proxy/gamestate/statestruct"
)

type logShim struct{ t *testing.T }

func (l logShim) Flush()                              {}
func (l logShim) Println(i ...interface{})            { l.t.Log(i...) }
func (l logShim) Printf(s string, i ...interface{})   { l.t.Logf(s, i...) }
func (l logShim) Verboseln(i ...interface{})          { l.t.Log(i...) }
func (l logShim) Verbosef(s string, i ...interface{}) { l.t.Logf(s, i...) }
func (l logShim) Warnln(i ...interface{})             { l.t.Error(i...) }
func (l logShim) Warnf(s string, i ...interface{})    { l.t.Errorf(s, i...) }

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
	mod, _ := New(logShim{t}, true)
	mod.handle("S/account/syncData", openAndRead(t, "testdata/syncdata.json"), nil)
	testChan := make(chan StateEvent, 1)
	done := make(chan error)
	go func() {
		evt := <-testChan
		payload := evt.Payload.(map[string]statestruct.EmptyStruct)
		if len(payload) != 12 {
			done <- fmt.Errorf("Expected a map with 12 entries in the payload")
		}
		close(done)
	}()
	mod.Hook("building.rooms.ELEVATOR", "test", testChan, false)
	mod.handle("S/building/sync", openAndRead(t, "testdata/buildingsync.json"), nil)
	check(t, <-done)
}

func TestGamestate(t *testing.T) {
	mod, _ := New(logShim{t}, true)
	mod.handle("S/account/syncData", openAndRead(t, "testdata/syncdata.json"), nil)
	mod.StateSync()
	mod.handle("S/building/sync", openAndRead(t, "testdata/buildingsync.json"), nil)
	mod.StateSync()
}

func TestModificationHooks(t *testing.T) {
	const testPath = "dexNav.enemy.stage.camp_02"
	testChan := make(chan StateEvent, 1)
	mod, _ := New(logShim{t}, true)
	hook := mod.Hook(testPath, "test", testChan, false)
	mod.handle("S/account/syncData", openAndRead(t, "testdata/syncdata.json"), nil)
	b := openAndRead(t, "testdata/buildingsync.json")
	mod.handle("S/building/sync", openAndRead(t, "testdata/buildingsync.json"), nil)
	// Block until the handle is done
	mod.StateSync()
	var evt StateEvent
	select {
	case evt = <-testChan:
	default:
		t.Fatal()
	}
	if evt.Path != testPath {
		t.Fatal("evt.Path != testPath")
	}
	// Shouldn't be able to get path from testChan anymore after unhooking
	hook.Unhook()
	mod.handle("S/building/sync", b, nil)
	mod.StateSync()
	select {
	case <-testChan:
		t.Fatal()
	default:
	}
}

// map[string]interface{} did not add new entries in the map. Workaround with
// map[string]struct{} instead.
func TestMapInterfaceBug(t *testing.T) {
	mod, _ := New(logShim{t}, true)
	mod.handle("S/account/syncData", openAndRead(t, "testdata/syncdata.json"), nil)
	mod.handle("S/building/sync", openAndRead(t, "testdata/buildingsync.json"), nil)
	mod.StateSync()
	if _, exists := mod.state.Building.Rooms.Elevator["slot_9"]; !exists {
		t.Fatal()
	}
}

func TestStructAccess(t *testing.T) {
	mod, _ := New(logShim{t}, true)
	mod.handle("S/account/syncData", openAndRead(t, "testdata/syncdata.json"), nil)
	val, err := mod.Get("building.rooms")
	check(t, err)
	rooms := val.(*statestruct.Rooms)
	// Don't really want to go through the trouble of testing deep equality, so
	// this will do for now.
	if rooms.Dormitory["slot_20"].Comfort != 3000 || len(rooms.Corridor) != 8 {
		t.Fatal("Rooms output not as expected.")
	}
}

// mod.handle has a latency of about 150-200ms, there's a lot of room for improvement
// here but it's usually not noticeable since it's typically faster than the game client.
func BenchmarkHandleSync(b *testing.B) {
	mod, _ := New(nil, true)
	data := openAndRead(nil, "testdata/syncdata.json")
	for i := 0; i < b.N; i++ {
		mod.handle("S/account/syncData", data, nil)
		mod.stateMutex.Lock()
	}
}
