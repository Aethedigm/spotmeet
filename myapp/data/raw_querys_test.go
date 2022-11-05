package data

import "testing"

func TestRawQuery(t *testing.T) {
	u := User{
		FirstName: "trq1",
		LastName:  "trq2",
		Email:     "trq3@test.com",
		Active:    1,
		Latitude:  100,
		Longitude: 100,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error("Error inserting user")
	}

	u.ID = uID

	u1Setting := Settings{
		UserID:                 uID,
		Distance:               100,
		LookingFor:             "friends",
		MatchSensitivity:       100,
		LikedArtistSensitivity: 100,
		LatMin:                 100,
		LatMax:                 102,
		LongMin:                100,
		LongMax:                102,
	}

	u1SID, err := u1Setting.Insert(u1Setting)
	if err != nil {
		t.Error("Error inserting user1's settings")
	}

	u1Setting.ID = u1SID

	u2 := User{
		FirstName: "trq4",
		LastName:  "trq5",
		Email:     "trq6@test.com",
		Active:    1,
		Latitude:  101,
		Longitude: 101,
	}

	u2ID, err := u2.Insert(u2)
	if err != nil {
		t.Error("Error inserting user 2")
	}

	u2.ID = u2ID

	u2Setting := Settings{
		UserID:                 u2ID,
		Distance:               100,
		LookingFor:             "friends",
		MatchSensitivity:       100,
		LikedArtistSensitivity: 100,
		LatMin:                 100,
		LatMax:                 102,
		LongMin:                100,
		LongMax:                102,
	}

	u2SID, err := u2Setting.Insert(u2Setting)
	if err != nil {
		t.Error("Error inserting user 2's settings")
	}

	u2Setting.ID = u2SID

	users, err := models.RQ.MatchQuery(u, u1Setting)
	if err != nil {
		t.Error("Error getting users matches")
	}

	for i := 0; i < len(users); i++ {
		if users[i] == u2.ID {
			return
		}
	}

	t.Error("User 2's ID was not found in matched users")
}
