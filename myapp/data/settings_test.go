// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import "testing"

func TestSettings_Table(t *testing.T) {
	s := Settings{}
	table := s.Table()
	if table != "settings" {
		t.Error("table name should be settings")
	}
}

func TestSettings_GetAll(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "getall_settings@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	s := Settings{
		UserID:                 uID,
		Distance:               10,
		LookingFor:             "friend",
		MatchSensitivity:       10,
		LikedArtistSensitivity: 10,
	}

	sID, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}

	s.ID = sID

	u2 := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "settings_getall2@test.com",
		Active:    1,
	}

	uID2, err := u2.Insert(u2)
	if err != nil {
		t.Error(err)
	}

	s2 := Settings{
		UserID:                 uID2,
		Distance:               10,
		LookingFor:             "friend",
		MatchSensitivity:       10,
		LikedArtistSensitivity: 10,
	}

	sID2, err := s2.Insert(s2)
	if err != nil {
		t.Error(err)
	}

	s2.ID = sID2

	settings, err := s.GetAll()
	if err != nil {
		t.Error(err)
	}

	if !SettingsInArray(s, settings) {
		t.Error("settings not in array")
	}

	if !SettingsInArray(s2, settings) {
		t.Error("settings not in array")
	}
}

func TestSettings_Delete(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "delsettings@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	s := Settings{
		UserID:                 uID,
		Distance:               10,
		LookingFor:             "friend",
		MatchSensitivity:       10,
		LikedArtistSensitivity: 10,
	}

	sID, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}

	s.ID = sID

	err = s.Delete(s.ID)
	if err != nil {
		t.Error(err)
	}

	_, err = s.GetByUserID(uID)
	if err == nil {
		t.Error("settings should not exist")
	}
}

func TestSettings_DeleteByUserID(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "dbui_settings@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	s := Settings{
		UserID:                 uID,
		Distance:               10,
		LookingFor:             "friend",
		MatchSensitivity:       10,
		LikedArtistSensitivity: 10,
	}

	sID, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}

	s.ID = sID

	err = s.DeleteByUserID(uID)
	if err != nil {
		t.Error(err)
	}

	_, err = s.GetByUserID(uID)
	if err == nil {
		t.Error("settings should not exist")
	}
}

func TestSettings_Get(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "settings_get@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	s := Settings{
		UserID:                 uID,
		Distance:               10,
		LookingFor:             "friend",
		MatchSensitivity:       10,
		LikedArtistSensitivity: 10,
	}

	sID, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}

	s.ID = sID

	s2, err := s.Get(s.ID)
	if err != nil {
		t.Error(err)
	}

	if s.ID != s2.ID {
		t.Error("settings not equal")
	}
}
