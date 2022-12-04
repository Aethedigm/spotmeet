package data

import (
	"errors"
	"time"

	up "github.com/upper/db/v4"
)

// Link is the type for a link
type Link struct {
	ID          int       `db:"id,omitempty"`
	User_A_ID   int       `db:"user_a_id" json:"user_A_id"`
	User_B_ID   int       `db:"user_b_id" json:"user_B_id"`
	PercentLink int       `db:"percent_match" json:"percent_match"`
	SongID      int       `db:"song_id" json:"song_id"`
	CreatedAt   time.Time `db:"created_at"`
}

// Table returns the table name associated with this model in the database
func (m *Link) Table() string {
	return "links"
}

// GetAllForOneUser returns a slice of all links for a user.
func (m *Link) GetAllForOneUser(userID int) ([]Link, error) {

	var all []Link

	var a1 []Link
	var a2 []Link

	collection := upper.Collection(m.Table())
	res1 := collection.Find(up.Cond{"user_a_id =": userID})
	res2 := collection.Find(up.Cond{"user_b_id =": userID})

	err := res1.All(&a1)
	if err != nil {
		return nil, err
	}

	err = res2.All(&a2)
	if err != nil {
		return nil, err
	}

	all = append(a1, a2...)

	return all, nil
}

// Insert inserts a new link, and returns the newly inserted link's id
func (m *Link) Insert(thelink Link) (int, error) {
	if thelink.User_A_ID == 0 || thelink.User_B_ID == 0 {
		return 0, errors.New("User_A_ID and User_B_ID must be set")
	}

	if thelink.User_A_ID == thelink.User_B_ID {
		return 0, errors.New("User_A_ID and User_B_ID cannot be the same")
	}

	thelink.CreatedAt = time.Now()

	collection := upper.Collection(m.Table())
	res, err := collection.Insert(thelink)
	if err != nil {
		return 0, err
	}

	id := getInsertID(res.ID())

	return id, nil
}
