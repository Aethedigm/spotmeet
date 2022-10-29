// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import (
	"testing"
	"time"
)

func TestSpotifyToken_Upsert(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "tpt_u@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	s := SpotifyToken{
		UserID:            uID,
		AccessToken:       "test",
		RefreshToken:      "test",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		AccessTokenExpiry: time.Now().Add(time.Hour),
	}

	sID, err := s.Upsert(s)
	if err != nil {
		t.Error(err)
	}

	s.ID = sID

	_, err = s.Get(sID)
	if err != nil {
		t.Error(err)
	}
}

func TestSpotifyToken_GetUserForRefreshToken(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "gufrt@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	u.ID = uID

	s := SpotifyToken{
		UserID:            uID,
		AccessToken:       "test",
		RefreshToken:      "refresh",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		AccessTokenExpiry: time.Now().Add(time.Hour),
	}

	sID, err := s.Upsert(s)
	if err != nil {
		t.Error(err)
	}

	s.ID = sID

	_, err = s.GetUserForRefreshToken("refresh")
	if err != nil {
		t.Error(err)
	}
}

func TestSpotifyToken_GetUserForAccessToken(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "gufat@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	u.ID = uID

	s := SpotifyToken{
		UserID:            uID,
		AccessToken:       "access",
		RefreshToken:      "refresh",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		AccessTokenExpiry: time.Now().Add(time.Hour),
	}

	sID, err := s.Upsert(s)
	if err != nil {
		t.Error(err)
	}

	s.ID = sID

	_, err = s.GetUserForAccessToken("access")
	if err != nil {
		t.Error(err)
	}
}

func TestSpotifyToken_GetSpotifyTokenForUser(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "gstfu@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	u.ID = uID

	s := SpotifyToken{
		UserID:            uID,
		AccessToken:       "access_new",
		RefreshToken:      "refresh_new",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		AccessTokenExpiry: time.Now().Add(time.Hour),
	}

	sID, err := s.Upsert(s)
	if err != nil {
		t.Error(err)
	}

	s.ID = sID

	_, err = s.GetSpotifyTokenForUser(uID)
	if err != nil {
		t.Error(err)
	}
}

func TestSpotifyToken_Delete(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "spottoken_delete@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	u.ID = uID

	s := SpotifyToken{
		UserID:            uID,
		AccessToken:       "access_delete",
		RefreshToken:      "refresh_delete",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		AccessTokenExpiry: time.Now().Add(time.Hour),
	}

	sID, err := s.Upsert(s)
	if err != nil {
		t.Error(err)
	}

	s.ID = sID

	err = s.Delete(sID)
	if err != nil {
		t.Error(err)
	}

	_, err = s.Get(sID)
	if err == nil {
		t.Error("should have returned error")
	}
}

func TestSpotifyToken_DeleteForUser(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "spottoken_deleteforuser@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	u.ID = uID

	s := SpotifyToken{
		UserID:            uID,
		AccessToken:       "access_deleteuser",
		RefreshToken:      "refresh_deleteuser",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
		AccessTokenExpiry: time.Now().Add(time.Hour),
	}

	sID, err := s.Upsert(s)
	if err != nil {
		t.Error(err)
	}

	s.ID = sID

	err = s.DeleteByUserID(uID)
	if err != nil {
		t.Error(err)
	}

	_, err = s.Get(sID)
	if err == nil {
		t.Error("should have returned error")
	}
}
