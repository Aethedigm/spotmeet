package data

import (
	up "github.com/upper/db/v4"
	"time"
)

// LikedArtist is the type for a liked artist
type LikedArtist struct {
	ID         int       `db:"id,omitempty"`
	UserID     int       `db:"user_id" json:"user_id"`
	ArtistID   int       `db:"artist_id" json:"artist_id"`
	LikedLevel string    `db:"liked_level" json:"liked_level"`
	Expires    time.Time `db:"expiry" json:"expiry"`
}

// Table returns the table name associated with this model in the database
func (l *LikedArtist) Table() string {
	return "liked_artists"
}

// GetAll returns a slice of all liked artists
func (l *LikedArtist) GetAll() ([]*LikedArtist, error) {
	collection := upper.Collection(l.Table())

	var all []*LikedArtist

	res := collection.Find().OrderBy("id")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// GetAllByOneUser returns a slice of all liked artists by a single user
func (l *LikedArtist) GetAllByOneUser(userID int) ([]*LikedArtist, error) {
	collection := upper.Collection(l.Table())

	var all []*LikedArtist

	res := collection.Find(up.Cond{"user_id =": userID}).OrderBy("id")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// Get gets one liked artist by id
func (l *LikedArtist) Get(id int) (*LikedArtist, error) {
	var theLikedArtist LikedArtist
	collection := upper.Collection(l.Table())
	res := collection.Find(up.Cond{"id =": id})

	err := res.One(&theLikedArtist)
	if err != nil {
		return nil, err
	}
	return &theLikedArtist, nil
}

// Update updates a liked artist record in the database
func (l *LikedArtist) Update(theLikedArtist LikedArtist) error {
	collection := upper.Collection(l.Table())
	res := collection.Find(theLikedArtist.ID)
	err := res.Update(&theLikedArtist)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a liked artist record by id
func (l *LikedArtist) Delete(id int) error {
	collection := upper.Collection(l.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a new liked artist, and returns the newly inserted id
func (l *LikedArtist) Insert(theLikedArtist LikedArtist) (int, error) {
	collection := upper.Collection(l.Table())

	theLikedArtist.Expires = time.Now().Add(time.Hour * 24 * 30) // expires after 30 days

	// make the insert
	res, err := collection.Insert(theLikedArtist)
	if err != nil {
		return 0, err
	}

	// get the id from the insert
	id := getInsertID(res.ID())

	return id, nil
}
