// go:build integration

// run tests with this command: go test . --tags integration --count=1
package data

import "testing"

func TestLikedSong_Table(t *testing.T) {
	ls := LikedSong{}
	table := ls.Table()
	if table != "liked_songs" {
		t.Error("table name should be liked_songs")
	}
}

func TestLikedSong_Insert(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "likedsong1@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	s := Song{
		SpotifyID:        "likedtest",
		Name:             "likedtest",
		ArtistName:       "artist_name",
		LoudnessAvg:      0.0,
		TempoAvg:         0.0,
		TimeSigAvg:       0,
		Acousticness:     0,
		Danceability:     0,
		Energy:           0,
		Instrumentalness: 0,
		Mode:             0,
		Speechiness:      0,
		Valence:          0,
	}

	sID, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}

	ls := LikedSong{
		UserID: uID,
		SongID: sID,
	}

	lsID, err := ls.Insert(ls)
	if err != nil {
		t.Error(err)
	}

	if lsID == 0 {
		t.Error("liked song id should not be 0")
	}
}

func TestLikedSong_GetAllByOneUser(t *testing.T) {
	u := User{
		FirstName: "temp",
		LastName:  "temp_last",
		Email:     "gabou_likedsong@test.com",
		Active:    1,
	}

	uID, err := u.Insert(u)
	if err != nil {
		t.Error(err)
	}

	s := Song{
		SpotifyID:        "likedtest_getallbyoneuser",
		Name:             "likedtest_getallbyoneuser",
		ArtistName:       "likedtest_getallbyoneuser",
		LoudnessAvg:      0.0,
		TempoAvg:         0.0,
		TimeSigAvg:       0,
		Acousticness:     0,
		Danceability:     0,
		Energy:           0,
		Instrumentalness: 0,
		Mode:             0,
		Speechiness:      0,
		Valence:          0,
	}

	sID, err := s.Insert(s)
	if err != nil {
		t.Error(err)
	}

	a2 := Song{
		SpotifyID:        "likedtest_getallbyoneuser2",
		Name:             "likedtest_getallbyoneuser2",
		ArtistName:       "likedtest_getallbyoneuser2",
		LoudnessAvg:      0.0,
		TempoAvg:         0.0,
		TimeSigAvg:       0,
		Acousticness:     0,
		Danceability:     0,
		Energy:           0,
		Instrumentalness: 0,
		Mode:             0,
		Speechiness:      0,
		Valence:          0,
	}

	sID2, err := a2.Insert(a2)
	if err != nil {
		t.Error(err)
	}

	ls := LikedSong{
		UserID: uID,
		SongID: sID,
	}

	laID, err := ls.Insert(ls)
	if err != nil {
		t.Error(err)
	}

	ls.ID = laID

	ls2 := LikedSong{
		UserID: uID,
		SongID: sID2,
	}

	lsID2, err := ls2.Insert(ls2)
	if err != nil {
		t.Error(err)
	}

	ls2.ID = lsID2

	laList, err := ls.GetAllByOneUser(uID)
	if err != nil {
		t.Error(err)
	}

	if !LikedSongInArray(ls, laList) {
		t.Error("liked song should be in list")
	}

	if !LikedSongInArray(ls2, laList) {
		t.Error("liked song should be in list")
	}
}
