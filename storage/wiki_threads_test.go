package storage

import (
	"strings"
	"testing"
)

func TestWIKI(t *testing.T) {
	err := LoadWikiThreads()
	if err != nil {
		t.Error("Wiki threads must be loaded; ", err.Error())
	}
	testWikiSearch := [3]string{"GetPlayerName", "OnPlayerConnect", "Limits"}
	var found [3]bool
	for i := range GetWikiThread().Thread {
		for a := range testWikiSearch {
			if strings.Compare(GetWikiThread().Thread[i], testWikiSearch[a]) == 0 {
				found[a] = true
			}
		}
	}
	for tst := range found {
		if found[tst] != true {
			t.Error("\"wiki_threads.go\" not works correcly! ", err.Error())
		}
	}
}