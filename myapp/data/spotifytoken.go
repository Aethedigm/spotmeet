package data

import (
	"errors"
	"time"

	up "github.com/upper/db/v4"
)

// SpotifyToken stores token information that is retrieved from Spotify.
// RefreshToken is constant, and therefore CreatedAt reflect when the RefreshToken
// -- was created.
// AccessToken is derived from the RefreshToken upon request from Spotify, and therefore
// -- the Updated at field reflects the time the latest AccessToken was created.
// AccessTokenExpiry holds the time that the AccessToken will expire. The system is
// -- designed to hand the RefreshToken to Spotify at an ordinate amount of time before the
// -- AccessToken is set to expire, get a new AccessToken, and set it here. This keeps the user
// -- logged in constantly, without the need to reenter their Spotify credentials.
type SpotifyToken struct {
	ID                int       `db:"id,omitempty" json:"id"`
	UserID            int       `db:"user_id" json:"user_id"`
	AccessToken       string    `db:"access_token" json:"access_token"`
	RefreshToken      string    `db:"refresh_token" json:"refresh_token"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
	AccessTokenExpiry time.Time `db:"access_expiry" json:"access_expiry"`
}

func (t *SpotifyToken) Table() string {
	return "spotify_token"
}

// GetUserForRefreshToken returns a User for the given refresh_token string.
func (t *SpotifyToken) GetUserForRefreshToken(refreshtoken string) (*User, error) {
	var u User
	var theSpotifyToken SpotifyToken

	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"refresh_token": refreshtoken})
	err := res.One(&theSpotifyToken)
	if err != nil {
		return nil, err
	}

	collection = upper.Collection("users")
	res = collection.Find(up.Cond{"id": theSpotifyToken.UserID})
	err = res.One(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetUserForAccessToken returns a User for the given access_token string.
func (t *SpotifyToken) GetUserForAccessToken(accesstoken string) (*User, error) {
	var u User
	var theSpotifyToken SpotifyToken

	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"access_token": accesstoken})
	err := res.One(&theSpotifyToken)
	if err != nil {
		return nil, err
	}

	collection = upper.Collection("users")
	res = collection.Find(up.Cond{"id": theSpotifyToken.UserID})
	err = res.One(&u)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// GetSpotifyTokenForUser returns the SpotifyToken for the given user_id. This function
// assumes that there is only one SpotifyToken (one refresh and one access token) for a
// single SpotMeet user at one time.
func (t *SpotifyToken) GetSpotifyTokenForUser(id int) (*SpotifyToken, error) {
	var spotifytoken *SpotifyToken
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"user_id": id}) // removed "access_expiry >": time.Now() because SpotifyToken struct
	// contains both access and refresh tokens. Need to be able to access this struct even when the access token expiry
	// is before the current time.
	err := res.One(&spotifytoken)
	if err != nil {
		return nil, err
	}

	return spotifytoken, nil
}

// Get SpotifyToken by id
func (t *SpotifyToken) Get(id int) (*SpotifyToken, error) {
	var spotifytoken *SpotifyToken
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"id": id})
	err := res.One(&spotifytoken)
	if err != nil {
		return nil, err
	}

	return spotifytoken, nil
}

// Delete deletes a SpotifyToken by id
func (t *SpotifyToken) Delete(id int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(id)
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

// DeleteByUserID deletes a SpotifyToken by user_id
func (t *SpotifyToken) DeleteByUserID(userID int) error {
	collection := upper.Collection(t.Table())
	res := collection.Find(up.Cond{"user_id": userID})
	err := res.Delete()
	if err != nil {
		return err
	}

	return nil
}

// Upsert inserts a new spotifytoken record into the spotifytokens table.
// This returns the ID of the new spotifytoken record, and type error.
func (t *SpotifyToken) Upsert(spotifytoken SpotifyToken) (int, error) {
	collection := upper.Collection(t.Table())

	// make sure AccessToken and RefreshToken are set
	if spotifytoken.AccessToken == "" || spotifytoken.RefreshToken == "" {
		return 0, errors.New("AccessToken and RefreshToken must be set")
	}

	// delete existing spotifytoken record, if one with the same user_id exists
	res := collection.Find(up.Cond{"user_id": spotifytoken.UserID})
	err := res.Delete()
	if err != nil {
		return 0, err
	}

	spotifytoken.CreatedAt = time.Now()
	spotifytoken.UpdatedAt = time.Now()

	res2, err := collection.Insert(spotifytoken)

	if err != nil {
		return 0, err
	}

	id := getInsertID(res2.ID())

	return id, nil
}

// UpdateAccessToken updates the access_token, updated_at, and access_expiry fields of
// a user's spotifytoken record
func (t *SpotifyToken) UpdateAccessToken(accessToken string, accessExpiry time.Time) error {
	collection := upper.Collection(t.Table())

	// update the calling struct's values
	t.AccessToken = accessToken
	t.UpdatedAt = time.Now()
	t.AccessTokenExpiry = accessExpiry

	// use UpdateReturning to automatically identify the record in the db our calling struct is referring to
	err := collection.UpdateReturning(t)
	if err != nil {
		return err
	}

	return nil
}
