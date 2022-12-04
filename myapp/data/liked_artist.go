package data

import (
	"time"

	up "github.com/upper/db/v4"
)

// LikedArtist is the type for a liked artist
type LikedArtist struct {
	ID         int       `db:"id,omitempty"`
	UserID     int       `db:"user_id" json:"user_id"`
	ArtistID   int       `db:"artist_id" json:"artist_id"`
	LikedLevel int       `db:"liked_level" json:"liked_level"`
	Expires    time.Time `db:"expiry" json:"expiry"`
}

// Table returns the table name associated with this model in the database
func (l *LikedArtist) Table() string {
	return "liked_artists"
}

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
