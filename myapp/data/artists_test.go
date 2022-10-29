// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import "testing"

func TestArtist_GetAll(t *testing.T) {
	a := Artist{
		Name:      "test",
		SpotifyID: "123",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}
	a.ID = aID

	b := Artist{
		Name:      "test2",
		SpotifyID: "1234",
	}

	bID, err := b.Insert(b)
	if err != nil {
		t.Error(err)
	}
	b.ID = bID

	artists, err := a.GetAll()

	if ArtistInArray(a, artists) == false {
		t.Error("failed to return artist")
	}

	if ArtistInArray(b, artists) == false {
		t.Error("failed to return artist")
	}
}

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

func TestArtist_Get(t *testing.T) {
	a := Artist{
		Name:      "test4",
		SpotifyID: "123",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}
	a.ID = aID

	artist, err := a.Get(aID)
	if err != nil {
		t.Error(err)
	}

	if artist.ID != aID {
		t.Error("incorrect artist returned")
	}
}

func TestArtist_Update(t *testing.T) {
	a := Artist{
		Name:      "test5",
		SpotifyID: "123",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}
	a.ID = aID

	a.Name = "test6"

	err = a.Update(a)
	if err != nil {
		t.Error(err)
	}

	artist, err := a.Get(aID)
	if err != nil {
		t.Error(err)
	}

	if artist.Name != "test6" {
		t.Error("incorrect artist returned")
	}
}

func TestArtist_Delete(t *testing.T) {
	a := Artist{
		Name:      "test7",
		SpotifyID: "123",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}
	a.ID = aID

	err = a.Delete(aID)
	if err != nil {
		t.Error(err)
	}
}

func TestArtist_DeleteByName(t *testing.T) {
	a := Artist{
		Name:      "test8",
		SpotifyID: "123",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}
	a.ID = aID

	err = a.DeleteByName(a.Name)
	if err != nil {
		t.Error(err)
	}
}
