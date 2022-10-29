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

func TestLikedArtist_GetAll(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "likedartist_getall@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	a := Artist{
		SpotifyID: "likedtest_getall",
		Name:      "likedtest_getall",
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

	la2 := LikedArtist{
		UserID:     uID,
		ArtistID:   aID,
		LikedLevel: 50,
	}

	laID2, err := la2.Insert(la2)
	if err != nil {
		t.Error(err)
	}

	if laID2 == 0 {
		t.Error("liked artist id should not be 0")
	}

	la2.ID = laID2

	laList, err := la.GetAll()
	if err != nil {
		t.Error(err)
	}

	if !LikedArtistInArray(la, laList) {
		t.Error("liked artist should be in list")
	}
}

func TestLikedArtist_GetAllByOneUser(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "gabou_likedartist@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	a := Artist{
		SpotifyID: "likedtest_getallbyoneuser",
		Name:      "likedtest_getallbyoneuser",
	}

	aID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}

	a2 := Artist{
		SpotifyID: "likedtest_getallbyoneuser2",
		Name:      "likedtest_getallbyoneuser2",
	}

	aID2, err := a2.Insert(a2)
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

	la.ID = laID

	la2 := LikedArtist{
		UserID:     uID,
		ArtistID:   aID2,
		LikedLevel: 50,
	}

	laID2, err := la2.Insert(la2)
	if err != nil {
		t.Error(err)
	}

	la2.ID = laID2

	laList, err := la.GetAllByOneUser(uID)
	if err != nil {
		t.Error(err)
	}

	if !LikedArtistInArray(la, laList) {
		t.Error("liked artist should be in list")
	}

	if !LikedArtistInArray(la2, laList) {
		t.Error("liked artist should be in list")
	}
}
