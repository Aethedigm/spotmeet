// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import (
	"testing"
	"time"
)

func TestLink_Table(t *testing.T) {
	s := Link{}
	table := s.Table()
	if table != "links" {
		t.Error("table name should be links")
	}
}

func TestLink_Insert(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "link_insert@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	u2 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "link_insert2@test.com",
		Active:    1,
	}

	uID2, err := u2.Insert(u2)
	if err != nil {
		t.Error(err)
	}

	s := Song{
		SpotifyID:   "test1333",
		Name:        "test1333",
		ArtistName:  "test1333",
		LoudnessAvg: 0.0,
		TempoAvg:    0.0,
		TimeSigAvg:  0,
	}

	sID, err := s.Insert(s)
	if err != nil || sID == 0 {
		t.Error(err)
	}

	song, err := s.Get(sID)
	if err != nil {
		t.Error(err)
	}

	l := Link{
		User_A_ID:   uID,
		User_B_ID:   uID2,
		SongID:      song.ID,
		PercentLink: 100,
		CreatedAt:   time.Now(),
	}

	lID, err := l.Insert(l)
	if err != nil {
		t.Error(err)
	}

	l.ID = lID

	if l.ID != lID {
		t.Error("link not inserted")
	}
}

func TestLink_Get_Bad_User_A_ID(t *testing.T) {
	ub := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "testub@test.com",
		Active:    1,
	}

	ubID, err := ub.Insert(ub)
	if err != nil {
		t.Error(err)
	}

	ub.ID = ubID

	s := Song{
		SpotifyID:   "testA",
		Name:        "testA",
		ArtistName:  "testA",
		LoudnessAvg: 0.0,
		TempoAvg:    0.0,
		TimeSigAvg:  0,
	}

	sID, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}

	l := Link{
		User_A_ID: 0,
		User_B_ID: ubID,
		SongID:    sID,
	}

	_, err = l.Insert(l)
	if err == nil {
		t.Error("User A ID should not be 0")
	}
}

func TestLink_Get_Bad_User_B_ID(t *testing.T) {
	ua := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "testua@test.com",
		Active:    1,
	}

	usID, err := ua.Insert(ua)
	if err != nil {
		t.Error(err)
	}

	ua.ID = usID

	s := Song{
		SpotifyID:   "testF",
		Name:        "testF",
		ArtistName:  "testF",
		LoudnessAvg: 0.0,
		TempoAvg:    0.0,
		TimeSigAvg:  0,
	}

	sID, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}

	l := Link{
		User_A_ID: usID,
		User_B_ID: 0,
		SongID:    sID,
	}

	_, err = l.Insert(l)
	if err == nil {
		t.Error("User B ID should not be 0")
	}
}

func TestLink_GetAllForOneUser(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "gafou_link@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	u2 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "gafou_link2@test.com",
		Active:    1,
	}

	uID2, err := u2.Insert(u2)
	if err != nil {
		t.Error(err)
	}

	s := Song{
		SpotifyID:   "test_gafou",
		Name:        "test_gafou",
		ArtistName:  "test_gafou",
		LoudnessAvg: 0.0,
		TempoAvg:    0.0,
		TimeSigAvg:  0,
	}

	sID, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}

	l := Link{
		User_A_ID: uID,
		User_B_ID: uID2,
		SongID:    sID,
	}

	lID, err := l.Insert(l)
	if err != nil {
		t.Error(err)
	}

	l.ID = lID

	links, err := l.GetAllForOneUser(uID)
	if err != nil {
		t.Error(err)
	}

	if len(links) == 0 {
		t.Error("no links found")
	}
}
