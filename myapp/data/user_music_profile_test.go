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
		UserID:   uID,
		Loudness: 10,
		Tempo:    10,
		TimeSig:  4,
	}

	_, err = ump.Insert(ump)
	if err != nil {
		t.Error("Failed inserting User Music Profile", err)
		return
	}
}

func TestUserMusicProfile_Get(t *testing.T) {
	u := User{
		FirstName: "Test",
		LastName:  "User",
		Email:     "testUMPGet@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error("Failed inserting new user", err)
		return
	}

	ump := &UserMusicProfile{
		UserID:   uID,
		Loudness: 10,
		Tempo:    10,
		TimeSig:  4,
	}

	umpID, err := ump.Insert(ump)
	if err != nil {
		t.Error("Failed inserting User Music Profile", err)
		return
	}

	ump.ID = umpID

	ump2, err := ump.Get(umpID)
	if err != nil {
		t.Error("Failed getting User Music Profile", err)
		return
	}

	if ump2.ID != ump.ID {
		t.Error("User Music Profiles do not match")
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
		UserID:   uID,
		Loudness: 10,
		Tempo:    10,
		TimeSig:  4,
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
