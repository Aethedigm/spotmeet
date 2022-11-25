// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import "testing"

func TestSong_GetAll(t *testing.T) {
	a := Song{
		Name:      "test",
		SpotifyID: "123",
	}

	sID, err := a.Insert(a)
	if err != nil {
		t.Error(err)
	}
	a.ID = sID

	s := Song{
		SpotifyID:   "test2",
		Name:        "test2",
		ArtistName:  "test2",
		LoudnessAvg: 0.0,
		TempoAvg:    0.0,
		TimeSigAvg:  0,
	}

	sID2, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}
	s.ID = sID2

	songs, err := a.GetAll()

	if SongInArray(a, songs) == false {
		t.Error("failed to return song")
	}

	if SongInArray(s, songs) == false {
		t.Error("failed to return song")
	}
}

func TestSong_GetByName(t *testing.T) {
	s := Song{
		SpotifyID:   "test3",
		Name:        "test3",
		ArtistName:  "test3",
		LoudnessAvg: 0.0,
		TempoAvg:    0.0,
		TimeSigAvg:  0,
	}

	sID, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}
	s.ID = sID

	song, err := s.GetByName(s.Name)
	if err != nil {
		t.Error(err)
	}

	if song.ID != sID {
		t.Error("incorrect song returned")
	}
}

func TestSong_Get(t *testing.T) {
	s := Song{
		SpotifyID:   "test4",
		Name:        "test4",
		ArtistName:  "test4",
		LoudnessAvg: 0.0,
		TempoAvg:    0.0,
		TimeSigAvg:  0,
	}

	sID, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}
	s.ID = sID

	song, err := s.Get(sID)
	if err != nil {
		t.Error(err)
	}

	if song.ID != sID {
		t.Error("incorrect song returned")
	}
}

func TestSong_Update(t *testing.T) {
	s := Song{
		SpotifyID:   "test5",
		Name:        "test5",
		ArtistName:  "test5",
		LoudnessAvg: 0.0,
		TempoAvg:    0.0,
		TimeSigAvg:  0,
	}

	sID, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}
	s.ID = sID

	s.Name = "test6"

	err = s.Update(s)
	if err != nil {
		t.Error(err)
	}

	song, err := s.Get(sID)
	if err != nil {
		t.Error(err)
	}

	if song.Name != "test6" {
		t.Error("incorrect song returned")
	}
}

func TestSong_Delete(t *testing.T) {
	s := Song{
		SpotifyID:   "test7",
		Name:        "test7",
		ArtistName:  "test7",
		LoudnessAvg: 0.0,
		TempoAvg:    0.0,
		TimeSigAvg:  0,
	}

	sID, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}
	s.ID = sID

	err = s.Delete(sID)
	if err != nil {
		t.Error(err)
	}
}

func TestSong_DeleteByName(t *testing.T) {
	s := Song{
		SpotifyID:   "test8",
		Name:        "test8",
		ArtistName:  "test8",
		LoudnessAvg: 0.0,
		TempoAvg:    0.0,
		TimeSigAvg:  0,
	}

	sID, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}
	s.ID = sID

	err = s.DeleteByName(s.Name)
	if err != nil {
		t.Error(err)
	}
}
