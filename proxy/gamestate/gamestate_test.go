package gamestate

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

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

// TestGameStateTraverse applies a delta to a base state and saves it to
// a file on the disk for manual checking.
func TestGameStateTraverse(t *testing.T) {
	b := openAndRead(t, "testdata/syncdata.json")
	syncData, err := unmarshalSyncData(b)
	user := syncData.User
	check(t, err)
	source, err := json.Marshal(user)
	check(t, err)
	f, err := os.Create("source.json")
	check(t, err)
	_, err = f.Write(source)
	check(t, err)
	f.Close()
	b = openAndRead(t, "testdata/buildingsync.json")
	res := gjson.GetBytes(b, "playerDataDelta.modified")
	newUser, err := unmarshalUserData([]byte(res.Raw))
	check(t, err)
	err = mergo.Merge(&user, newUser, mergo.WithOverride)
	check(t, err)
	merged, err := json.Marshal(user)
	check(t, err)
	f, err = os.Create("merged.json")
	check(t, err)
	_, err = f.Write(merged)
	check(t, err)
	f.Close()
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
	testChan := make(chan string, 1)
	mod, _ := New(logShim{t})
	mod.stateMutex.Unlock()
	hook := mod.Hook(testPath, "test", testChan)
	// WalkAndNotify should be called within the critical section
	mod.stateMutex.Lock()
	err := mod.WalkAndNotify(b)
	mod.stateMutex.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	var path string
	select {
	case path = <-testChan:
	default:
		t.Fatal()
	}
	if path != testPath {
		t.Fatal()
	}
	hook.Unhook()
	// Shouldn't be able to get path from testChan anymore after unhooking
	mod.stateMutex.Lock()
	err = mod.WalkAndNotify(b)
	mod.stateMutex.Unlock()
	if err != nil {
		t.Fatal(err)
	}
	select {
	case <-testChan:
		t.Fatal()
	default:
	}
}
