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
	ArtistID    int       `db:"artist_id" json:"artist_id"`
	ThreadID    int       `db:"thread_id" json:"thread_id"`
	CreatedAt   time.Time `db:"created_at"`
}

// Table returns the table name associated with this model in the database
func (m *Link) Table() string {
	return "links"
}

// GetAll returns a slice of all links.
func (m *Link) GetAll() ([]*Link, error) {
	collection := upper.Collection(m.Table())

	var all []*Link

	res := collection.Find().OrderBy("id")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
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

// Get gets one link by id
func (m *Link) Get(id int) (*Link, error) {
	var thelink Link
	collection := upper.Collection(m.Table())
	res := collection.Find(up.Cond{"id =": id})

	err := res.One(&thelink)
	if err != nil {
		return nil, err
	}

	return &thelink, nil
}

// Delete deletes a link by id
func (m *Link) Delete(id int) error {
	collection := upper.Collection(m.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil

}

// Insert inserts a new link, and returns the newly inserted link's id
func (m *Link) Insert(thelink Link) (int, error) {
	if thelink.User_A_ID == 0 || thelink.User_B_ID == 0 {
		return 0, errors.New("User_A_ID and User_B_ID must be set")
	}

	if thelink.User_A_ID == thelink.User_B_ID {
		return 0, errors.New("User_A_ID and User_B_ID cannot be the same")
	}

	if thelink.ThreadID == 0 {
		return 0, errors.New("ThreadID must be set")
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
