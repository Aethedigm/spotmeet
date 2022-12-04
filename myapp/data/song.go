package data

import (
	"errors"

	up "github.com/upper/db/v4"
)

// Song is the type for a song
type Song struct {
	ID               int     `db:"id,omitempty"`
	SpotifyID        string  `db:"spotify_id" json:"spotify_id"`
	Name             string  `db:"song_name" json:"song_name"`
	ArtistName       string  `db:"artist_name" json:"artist_name"`
	LoudnessAvg      float64 `db:"loudness" json:"loudness"`
	TempoAvg         float64 `db:"tempo" json:"tempo"`
	TimeSigAvg       int     `db:"time_sig" json:"time_sig"`
	Acousticness     float32 `db:"acousticness" json:"acousticness"`
	Danceability     float32 `db:"danceability" json:"danceability"`
	Energy           float32 `db:"energy" json:"energy"`
	Instrumentalness float32 `db:"instrumentalness" json:"instrumentalness"`
	Mode             int     `db:"mode" json:"mode"`
	Speechiness      float32 `db:"speechiness" json:"speechiness"`
	Valence          float32 `db:"valence" json:"valence"`
}

// Table returns the table name associated with this model in the database
func (a *Song) Table() string {
	return "songs"
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

// GetBySpotifyID gets one song, by Spotify ID
func (a *Song) GetBySpotifyID(spotifyID string) (*Song, error) {
	var theSong Song
	collection := upper.Collection(a.Table())
	res := collection.Find(up.Cond{"spotify_id": spotifyID})
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
	song, err := a.GetBySpotifyID(theSong.SpotifyID)
	if err != nil {
		if err.Error() != "upper: no more rows in this result set" {
			return 0, err
		}
	}

	if song != nil {
		return 0, errors.New("song id already exists")
	}

	res, err := collection.Insert(&theSong)
	if err != nil {
		return 0, err
	}

	sID := getInsertID(res.ID())

	return sID, nil
}
