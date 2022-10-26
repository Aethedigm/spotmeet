package data

import (
	up "github.com/upper/db/v4"
)

// definitive matching-motives for a user to select in their Settings (see Insert() method)
type relation string

// Settings is a struct that holds a user's settings configuration
//
// MatchSensitivity represents an integer from 1 to 100, and will be the minimum value of potential matches'
// liked_artists' liked_levels that the app will move forward with for matching.
//
// LikedArtistSensitivity (integer 1 - 100, representing percentage)
// This will control how sensitive the app is in saving artists a user has listened to as actual liked_artists. Without
// this, for every user the app would save every single listened-to artist as a liked_artist. For the moment, this will
// be hidden from the user, and every user's liked_artist_sensitivity will be set to a default of 5 (5%). This way,
// when SpotMeet pulls in all songs the user has listened to in the past month (each time they log in), it can filter
// out the listened-to artists whose songs accounted for less than 5% of the total songs listened to by the user in the
// past month.

type Settings struct {
	ID                     int      `db:"id,omitempty"`
	UserID                 int      `db:"user_id" json:"user_id"`
	Distance               int      `db:"distance" json:"distance"`
	LookingFor             relation `db:"looking_for" json:"looking_for"`
	MatchSensitivity       int      `db:"match_sensitivity" json:"match_sensitivity"`
	LikedArtistSensitivity int      `db:"liked_artist_sensitivity" json:"liked_artist_sensitivity"`
}

// Table returns the table name associated with this model in the database
func (s *Settings) Table() string {
	return "settings"
}

// GetAll returns a slice of all settings configurations
func (s *Settings) GetAll() ([]*Settings, error) {
	collection := upper.Collection(s.Table())

	var all []*Settings

	res := collection.Find().OrderBy("id")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// Get gets one settings record by settings id
func (s *Settings) Get(id int) (*Settings, error) {
	var settings Settings
	collection := upper.Collection(s.Table())
	res := collection.Find(up.Cond{"id =": id})

	err := res.One(&settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

// GetByUserID gets one settings record by user id
func (s *Settings) GetByUserID(user_id int) (*Settings, error) {
	var settings Settings
	collection := upper.Collection(s.Table())
	res := collection.Find(up.Cond{"user_id =": user_id})

	err := res.One(&settings)
	if err != nil {
		return nil, err
	}

	return &settings, nil
}

// Delete deletes a settings entry by settings id
func (s *Settings) Delete(id int) error {
	collection := upper.Collection(s.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil

}

// DeleteByUserID deletes a settings entry by a user's id
func (s *Settings) DeleteByUserID(user_id int) error {
	collection := upper.Collection(s.Table())
	res := collection.Find(up.Cond{"user_id =": user_id})
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil

}

// Insert inserts a new record for a user's settings, and returns the newly inserted record's id
func (s *Settings) Insert(settings Settings) (int, error) {

	// set possible selections for a user's desired relation-type by using the app
	relations := [6]relation{"friendship", "dating", "workout partner", "musicians", "concert-goers", "not sure yet"}

	// Default settings for all new users
	settings.MatchSensitivity = 50
	settings.LikedArtistSensitivity = 5
	settings.Distance = 50 // in miles
	settings.LookingFor = relations[0]

	collection := upper.Collection(s.Table())
	res, err := collection.Insert(settings)
	if err != nil {
		return 0, err
	}

	id := getInsertID(res.ID())

	return id, nil
}
