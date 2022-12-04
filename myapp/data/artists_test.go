// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import "testing"

func TestArtist_GetByName(t *testing.T) {
	a := Artist{
		Name:      "test3",
		SpotifyID: "123",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}
	a.ID = aID

	artist, err := a.GetByName(a.Name)
	if err != nil {
		t.Error(err)
	}

	if artist.ID != aID {
		t.Error("incorrect artist returned")
	}
}
