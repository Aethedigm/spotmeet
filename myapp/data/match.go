// Spotmeet - (Capstone Team E)
// 2022 Stephen Sumpter, John Neumeier,
// Zach Kods, Landon Wilson
package data

import (
	"errors"
	"time"

	up "github.com/upper/db/v4"
)

// Match is the type for a match
type Match struct {
	ID                int       `db:"id,omitempty"`
	User_A_ID         int       `db:"user_a_id" json:"user_A_id"`
	User_B_ID         int       `db:"user_b_id" json:"user_B_id"`
	PercentMatch      float32   `db:"percent_match" json:"percent_match"`
	SongID            int       `db:"song_id" json:"song_id"`
	CreatedAt         time.Time `db:"created_at"`
	Expires           time.Time `db:"expiry" json:"expiry"`
	Complete          bool      `db:"complete" json:"complete"`
	UserAViewedThread bool      `db:"user_a_viewed" json:"user_a_viewed"`
	UserBViewedThread bool      `db:"user_b_viewed" json:"user_b_viewed"`
}

// Table returns the table name associated with this model in the database
func (m *Match) Table() string {
	return "matches"
}

func (m *Match) Update(match Match) error {
	collection := upper.Collection(match.Table())
	res := collection.Find(match.ID)
	err := res.Update(match)
	if err != nil {
		return err
	}
	return nil
}

// Get gets one match by id
func (m *Match) Get(id int) (*Match, error) {
	var thematch Match
	collection := upper.Collection(m.Table())
	res := collection.Find(up.Cond{"id": id})

	err := res.One(&thematch)
	if err != nil {
		return nil, err
	}

	return &thematch, nil
}

// Get gets one match by id
func (m *Match) GetByBothUsers(id int, id2 int) (*Match, error) {
	var thematch Match
	collection := upper.Collection(m.Table())
	res1 := collection.Find(up.Cond{"user_a_id": id, "user_b_id": id2})
	res2 := collection.Find(up.Cond{"user_a_id": id2, "user_b_id": id})

	err := res1.One(&thematch)
	if err != nil {
		err := res2.One(&thematch)
		if err != nil {
			return nil, err
		}
	}

	return &thematch, nil
}

// Insert inserts a new match, and returns the newly inserted match's id
func (m *Match) Insert(thematch Match) (int, error) {

	if thematch.User_A_ID == 0 || thematch.User_B_ID == 0 {
		return 0, errors.New("User_A_ID and User_B_ID must be set")
	}

	if thematch.User_A_ID == thematch.User_B_ID {
		return 0, errors.New("User_A_ID and User_B_ID cannot be the same")
	}

	// thematch.MessageStatus = 0
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
