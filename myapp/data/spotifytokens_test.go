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
