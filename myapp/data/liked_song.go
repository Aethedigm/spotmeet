// Spotmeet - (Capstone Team E)
// 2022 Stephen Sumpter, John Neumeier,
// Zach Kods, Landon Wilson
package data

import (
	up "github.com/upper/db/v4"
)

// LikedSong is the type for a liked song
type LikedSong struct {
	ID     int `db:"id,omitempty"`
	UserID int `db:"user_id" json:"user_id"`
	SongID int `db:"song_id" json:"song_id"`
}

// Table returns the table name associated with this model in the database
func (l *LikedSong) Table() string {
	return "liked_songs"
}

func (l *LikedSong) GetAllByOneUser(userID int) ([]*LikedSong, error) {
	collection := upper.Collection(l.Table())

	var all []*LikedSong

	res := collection.Find(up.Cond{"user_id =": userID}).OrderBy("id")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// DeleteAllForUser deletes all liked songs for a user
func (l *LikedSong) DeleteAllForUser(userID int) error {
	collection := upper.Collection(l.Table())
	res := collection.Find(up.Cond{"user_id =": userID})
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a new liked song, and returns the newly inserted id
func (l *LikedSong) Insert(theLikedSong LikedSong) (int, error) {
	collection := upper.Collection(l.Table())

	// make the insert
	res, err := collection.Insert(theLikedSong)
	if err != nil {
		return 0, err
	}

	// get the id from the insert
	id := getInsertID(res.ID())

	return id, nil
}
