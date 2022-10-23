package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// Match is the type for a match
type Match struct {
	ID           int       `db:"id,omitempty"`
	User_A_ID    int       `db:"user_A_id" json:"user_A_id"`
	User_B_ID    int       `db:"user_B_id" json:"user_B_id"`
	PercentMatch int       `db:"percent_match" json:"percent_match"`
	Artist_ID    int       `db:"artist_id" json:"artist_id"`
	CreatedAt    time.Time `db:"created_at"`
	Expires      time.Time `db:"expiry" json:"expiry"`
}

// Table returns the table name associated with this model in the database
func (m *Match) Table() string {
	return "match"
}

// GetAll returns a slice of all matches.
func (m *Match) GetAll() ([]Match, error) {
	collection := upper.Collection(m.Table())

	var all []Match

	res := collection.Find().OrderBy("user_A_id")
	err := res.All(all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// GetByUserID returns a slice of all matches for a user.
func (m *Match) GetAllForOneUser(userID int) ([]*Match, error) {

	var all []*Match
	collection := upper.Collection(m.Table())
	res := collection.Find(up.Cond{"user_A_id =": userID, "user_B_id =": userID})
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// Get gets one match by id
func (m *Match) Get(id int) (*Match, error) {
	var thematch Match
	collection := upper.Collection(m.Table())
	res := collection.Find(up.Cond{"id =": id})

	err := res.One(&thematch)
	if err != nil {
		return nil, err
	}

	return &thematch, nil
}

// Delete deletes a match by id
func (m *Match) Delete(id int) error {
	collection := upper.Collection(m.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil

}

// Insert inserts a new match, and returns the newly inserted match's id
func (m *Match) Insert(thematch Match) (int, error) {

	// ------------------
	// create something here that ensures thematch has all fields filled except for CreatedAt and
	// Expires
	// ------------------

	thematch.CreatedAt = time.Now()
	thematch.Expires = time.Now().Add(time.Hour * 24 * 5) // expires after 5 days

	collection := upper.Collection(m.Table())
	res, err := collection.Insert(thematch)
	if err != nil {
		return 0, err
	}

	id := getInsertID(res.ID())

	return id, nil
}
