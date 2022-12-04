package data

import (
	up "github.com/upper/db/v4"
)

// Artist is the type for an artist
type Artist struct {
	ID        int    `db:"id,omitempty"`
	SpotifyID string `db:"spotify_id" json:"spotify_id"`
	Name      string `db:"artist_name" json:"artist_name"`
}

// Table returns the table name associated with this model in the database
func (a *Artist) Table() string {
	return "artists"
}

func (a *Artist) GetOneID() (int, error) {
	collection := upper.Collection(a.Table())

	var art *Artist

	res := collection.Find().Limit(1)
	err := res.One(&art)
	if err != nil {
		return art.ID, err
	}

	return art.ID, nil
}

func (a *Artist) GetByName(name string) (*Artist, error) {
	var theArtist Artist
	collection := upper.Collection(a.Table())
	res := collection.Find(up.Cond{"artist_name": name})
	err := res.One(&theArtist)
	if err != nil {
		return nil, err
	}

	return &theArtist, nil
}

// Insert inserts a new artist, and returns the newly inserted id
func (a *Artist) Insert(theArtist Artist) (int, error) {
	collection := upper.Collection(a.Table())

	// Make sure this artist doesn't already exist
	_, err := a.GetByName(theArtist.Name)
	if err == nil {
		return 0, nil
	}

	res, err := collection.Insert(&theArtist)
	if err != nil {
		return 0, err
	}

	aID := getInsertID(res.ID())

	return aID, nil
}
