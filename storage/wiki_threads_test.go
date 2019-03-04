package storage

import (
	"testing"
	"fmt"
)

func TestWIKI(t *testing.T) {
	testWikiSearch := [3]string{"GetPlayerName", "onplayerconnect", "Db_get_field_assoc_float"}
	var found = [3]bool{false, false, false}
	for i := range testWikiSearch {
		_, found[i] = SearchThread(testWikiSearch[i])
	}
	for a := range found {
		if found[a] != true {
			t.Error("Something's gone wrong")
		}
	}
}