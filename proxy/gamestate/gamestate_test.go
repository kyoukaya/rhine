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
