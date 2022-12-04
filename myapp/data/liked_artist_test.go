// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import "testing"

func TestLikedArtist_Table(t *testing.T) {
	la := LikedArtist{}
	table := la.Table()
	if table != "liked_artists" {
		t.Error("table name should be liked_artists")
	}
}
