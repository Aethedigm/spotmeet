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

// GetAll returns a slice of all artists
func (a *Artist) GetAll() ([]*Artist, error) {
	collection := upper.Collection(a.Table())

	var all []*Artist

	res := collection.Find().OrderBy("artist_name")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// GetByName gets one artist, by name
func (a *Artist) GetByName(name string) (*Artist, error) {
	var theArtist Artist
	collection := upper.Collection(a.Table())
	res := collection.Find(up.Cond{"artist_name =": name})
	err := res.One(&theArtist)
	if err != nil {
		return nil, err
	}

	return &theArtist, nil
}

// Get gets one artist by id
func (a *Artist) Get(id int) (*Artist, error) {
	var theArtist Artist
	collection := upper.Collection(a.Table())
	res := collection.Find(up.Cond{"id =": id})

	err := res.One(&theArtist)
	if err != nil {
		return nil, err
	}

	return &theArtist, nil
}

// Update updates an artist record in the database
func (a *Artist) Update(theArtist Artist) error {
	collection := upper.Collection(a.Table())
	res := collection.Find(theArtist.ID)
	err := res.Update(&theArtist)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes an artist record by id
func (a *Artist) Delete(id int) error {
	collection := upper.Collection(a.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteByName deletes an artist record by name
func (a *Artist) DeleteByName(artistName int) error {
	collection := upper.Collection(a.Table())
	res := collection.Find(up.Cond{"artist_name =": artistName})
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a new artist, and returns the newly inserted id
func (a *Artist) Insert(theArtist Artist) (int, error) {
	collection := upper.Collection(a.Table())

	// NEED TO CHECK IF LINES 100-110 WORK

	// grab the artist if it already exists in the database
	checkIfExists, err1 := a.GetByName(theArtist.Name)
	if err1 != nil {
		return 0, err1
	}

	// if empty
	blankArtist := Artist{}
	if *checkIfExists != blankArtist {
		return 0, nil
	}

	// make the insert
	res, err2 := collection.Insert(theArtist)
	if err2 != nil {
		return 0, err2
	}

	// get the id from the insert
	id := getInsertID(res.ID())

	return id, nil
}
