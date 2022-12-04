package data

import (
	"testing"
)

func TestUserMusicProfile_Insert(t *testing.T) {
	u := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "testUMP@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error("Failed inserting new user", err)
		return
	}

	ump := &UserMusicProfile{
		UserID:           uID,
		Loudness:         10,
		Tempo:            10,
		TimeSig:          4,
		Acousticness:     0,
		Danceability:     0,
		Energy:           0,
		Instrumentalness: 0,
		Mode:             0,
		Speechiness:      0,
		Valence:          0,
	}

	_, err = ump.Insert(ump)
	if err != nil {
		t.Error("Failed inserting User Music Profile", err)
		return
	}
}

func TestUserMusicProfile_GetByUser(t *testing.T) {
	u := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "testUMPGetUser@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error("Failed inserting new user", err)
		return
	}

	u.ID = uID

	ump := &UserMusicProfile{
		UserID:           uID,
		Loudness:         10,
		Tempo:            10,
		TimeSig:          4,
		Acousticness:     0,
		Danceability:     0,
		Energy:           0,
		Instrumentalness: 0,
		Mode:             0,
		Speechiness:      0,
		Valence:          0,
	}

	umpID, err := ump.Insert(ump)
	if err != nil {
		t.Error("Failed inserting User Music Profile", err)
		return
	}

	ump.ID = umpID

	ump2, err := ump.GetByUser(u)
	if err != nil {
		t.Error("Failed getting User Music Profile", err)
		return
	}

	if ump2.ID != ump.ID {
		t.Error("User Music Profiles do not match")
	}
}

func TestUserMusicProfile_Insert_MissingData(t *testing.T) {
	u := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "testUMPMD@test.com",
		Active:    1,
	}

	_, err := u.Insert(u)
	if err != nil {
		t.Error("Failed inserting new user", err)
		return
	}

	ump := &UserMusicProfile{}

	_, err = ump.Insert(ump)
	if err == nil {
		t.Error("Should have failed insert with missing data")
	}
}

func TestUserMusicProfile_GetByUser_NoProfile(t *testing.T) {
	u := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "testUMPNP@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error("Failed inserting new user", err)
		return
	}

	u.ID = uID

	ump := UserMusicProfile{}

	_, err = ump.GetByUser(u)
	if err == nil {
		t.Error("No User Music Profile should have been returned")
	}
}

func TestUserMusicProfile_Get_BadID(t *testing.T) {
	u := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "testUMPBI@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error("Failed inserting new user", err)
		return
	}

	u.ID = uID

	ump := UserMusicProfile{
		UserID: u.ID,
	}

	umpID, err := ump.Insert(&ump)
	if err != nil {
		t.Error("Failed to insert")
	}

	newUMP, err := ump.GetByUserID(umpID + 10000)
	if err == nil {
		t.Error("Should not have found User Music Profile with this ID")
	}

	if newUMP != nil {
		t.Error("Object should be nil")
	}
}
