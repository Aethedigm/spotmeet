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

func TestLikedArtst_Insert(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "likedartist1@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	a := Artist{
		SpotifyID: "likedtest",
		Name:      "likedtest",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}

	la := LikedArtist{
		UserID:     uID,
		ArtistID:   aID,
		LikedLevel: 100,
	}

	laID, err := la.Insert(la)
	if err != nil {
		t.Error(err)
	}

	if laID == 0 {
		t.Error("liked artist id should not be 0")
	}
}

func TestLikedArtist_Update(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "likedartist_update@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	a := Artist{
		SpotifyID: "likedtest_update",
		Name:      "likedtest_update",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}

	la := LikedArtist{
		UserID:     uID,
		ArtistID:   aID,
		LikedLevel: 100,
	}

	laID, err := la.Insert(la)
	if err != nil {
		t.Error(err)
	}

	if laID == 0 {
		t.Error("liked artist id should not be 0")
	}

	la.ID = laID

	la.LikedLevel = 50
	err = la.Update(la)
	if err != nil {
		t.Error(err)
	}

	la2, err := la.Get(laID)
	if err != nil {
		t.Error(err)
	}

	if la2.LikedLevel != 50 {
		t.Error("liked level should be 50")
	}
}

func TestLikedArtist_Delete(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "likedaritst_delete@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	a := Artist{
		SpotifyID: "likedtest_delete",
		Name:      "likedtest_delete",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}

	la := LikedArtist{
		UserID:     uID,
		ArtistID:   aID,
		LikedLevel: 100,
	}

	laID, err := la.Insert(la)
	if err != nil {
		t.Error(err)
	}

	if laID == 0 {
		t.Error("liked artist id should not be 0")
	}

	la.ID = laID

	err = la.Delete(laID)
	if err != nil {
		t.Error(err)
	}

	_, err = la.Get(laID)
	if err == nil {
		t.Error("liked artist should not exist")
	}
}
