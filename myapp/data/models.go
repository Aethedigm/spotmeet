package data

import (
	"database/sql"
	"fmt"

	db2 "github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/postgresql"
)

var db *sql.DB
var upper db2.Session

// Models is the wrapper for all database models
type Models struct {
	// any models inserted here (and in the New function)
	// are easily accessible throughout the entire application
	Users             User
	Matches           Match
	Profiles          Profile
	SpotifyTokens     SpotifyToken
	Settings          Settings
	Messages          Message
	Links             Link
	RQ                RawQuery
	LikedArtists      LikedArtist
	Artists           Artist
	UserMusicProfiles UserMusicProfile
	Songs             Song
	LikedSongs        LikedSong
	RecoveryEmails    RecoveryEmail
}

// New initializes the models package for use
func New(databasePool *sql.DB) Models {
	db = databasePool
	upper, _ = postgresql.New(databasePool)

	return Models{
		Users: User{},
	}
}

// getInsertID returns the integer value of a newly inserted id (using upper)
func getInsertID(i db2.ID) int {
	idType := fmt.Sprintf("%T", i)
	if idType == "int64" {
		return int(i.(int64))
	}

	return i.(int)
}
