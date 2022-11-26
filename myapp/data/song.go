package data

import (
	up "github.com/upper/db/v4"
)

// Song is the type for a song
type Song struct {
	ID          int     `db:"id,omitempty"`
	SpotifyID   string  `db:"spotify_id" json:"spotify_id"`
	Name        string  `db:"song_name" json:"song_name"`
	ArtistName  string  `db:"artist_name" json:"artist_name"`
	LoudnessAvg float64 `db:"loudness" json:"loudness"`
	TempoAvg    float64 `db:"tempo" json:"tempo"`
	TimeSigAvg  int     `db:"time_sig" json:"time_sig"`
}

// Table returns the table name associated with this model in the database
func (a *Song) Table() string {
	return "songs"
}

func (a *Song) GetOneID() (int, error) {
	collection := upper.Collection(a.Table())

	var sng *Song

	res := collection.Find().Limit(1)
	err := res.One(&sng)
	if err != nil {
		return sng.ID, err
	}

	return sng.ID, nil
}

// GetAll returns a slice of all songs
func (a *Song) GetAll() ([]*Song, error) {
	collection := upper.Collection(a.Table())

	var all []*Song

	res := collection.Find().OrderBy("song_name")
	err := res.All(&all)
	if err != nil {
		return nil, err
	}

	return all, nil
}

// GetByName gets one song, by name
func (a *Song) GetByName(name string) (*Song, error) {
	var theSong Song
	collection := upper.Collection(a.Table())
	res := collection.Find(up.Cond{"song_name": name})
	err := res.One(&theSong)
	if err != nil {
		return nil, err
	}

	return &theSong, nil
}

// GetByNameAndArtist gets one song, by name
func (a *Song) GetByNameAndArtist(name string, artistName string) (*Song, error) {
	var theSong Song
	collection := upper.Collection(a.Table())
	res := collection.Find(up.Cond{"song_name": name, "artist_name": artistName})
	err := res.One(&theSong)
	if err != nil {
		return nil, err
	}

	return &theSong, nil
}

// Get gets one song by id
func (a *Song) Get(id int) (*Song, error) {
	var theSong Song
	collection := upper.Collection(a.Table())
	res := collection.Find(up.Cond{"id": id})

	err := res.One(&theSong)
	if err != nil {
		return nil, err
	}

	return &theSong, nil
}

// Update updates a song record in the database
func (a *Song) Update(theSong Song) error {
	collection := upper.Collection(a.Table())
	res := collection.Find(theSong.ID)
	err := res.Update(&theSong)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a song record by id
func (a *Song) Delete(id int) error {
	collection := upper.Collection(a.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// DeleteByName deletes a song record by name
func (a *Song) DeleteByName(songName string) error {
	collection := upper.Collection(a.Table())
	res := collection.Find(up.Cond{"song_name": songName})
	err := res.Delete()
	if err != nil {
		return err
	}
	return nil
}

// Insert inserts a new song, and returns the newly inserted id
func (a *Song) Insert(theSong Song) (int, error) {
	collection := upper.Collection(a.Table())

	// Make sure this song doesn't already exist
	song, err := a.GetByNameAndArtist(theSong.Name, theSong.ArtistName)
	if song != nil {
		return 0, nil
	}

	res, err := collection.Insert(&theSong)
	if err != nil {
		return 0, err
	}

	sID := getInsertID(res.ID())

	return sID, nil
}
