// Spotmeet - (Capstone Team E)
// 2022 Stephen Sumpter, John Neumeier,
// Zach Kods, Landon Wilson
package handlers

import (
	"encoding/base32"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"myapp/data"
	"net/http"
	"os"
	"time"

	"github.com/zmb3/spotify"

	"golang.org/x/oauth2"
)

var auth = spotify.Authenticator{}
var state string

// Set the scopes for our Spotify authorization. "Scopes" refers to how many permissions from the user we are asking for.
var spotScopes = []string{spotify.ScopeUserTopRead, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserFollowRead}

// SpotifyAuthorization formulates a URL that will be used as a redirect to Spotify's authorization endpoint.
func (h *Handlers) SpotifyAuthorization(w http.ResponseWriter, r *http.Request) {
	// define the url on our server where we'd like the user to come back to after authorizing with Spotify
	callback := os.Getenv("URL") + "/spotauth/callback"

	// create an "authenticator", which will formulate a url for the user to redirect to Spotify's auth endpoint
	auth = spotify.NewAuthenticator(callback, spotScopes...)
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
		return
	}

	// Define a "state", which is a complex code string that provides another layer of protection, ensuring
	// to our server that incoming requests from Spotify are legitimate.
	state = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	// Include the state string with the call to create a redirect url. This url will be unique only to this
	// user and redirect instance.
	url := auth.AuthURL(state)
	url = url + "&show_dialog=true"

	// redirect the user to Spotify for authorization
	http.Redirect(w, r, url, http.StatusSeeOther)
}

// SpotifyAuthorizationCallback is called when authorization information comes back from Spotify. This
// function updates/inserts that information into our database for the user.
func (h *Handlers) SpotifyAuthorizationCallback(w http.ResponseWriter, r *http.Request) {
	// get a token struct, formulated from the information received via http request from Spotify.
	// (The state code string is needed here in order for the information to be decoded.)
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
		fmt.Println(err)
		return
	}

	// Check if the state value recieved from Spotify is the same one that was sent out with the initial redirect
	// to Spotify for authorization. If it doesn't match, then redirect back to spotauth for re-authorization.
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
		fmt.Printf("State mismatch: %s != %s\n", st, state)
		return
	}

	// create a new client struct from the token we recieved from Spotify
	spotclient := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")

	// use that client struct with CurrentUser() in order to validate that the user (with Spotify) actually exists
	_, err = spotclient.CurrentUser()
	if err != nil {
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
		fmt.Println(err)
		return
	}

	// get the user ID from the session
	userID := h.App.Session.GetInt(r.Context(), "userID")
	if userID == 0 || userID == -1 {
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
		fmt.Println("The user_id of the current user could not be found in the session data.")
		return
	}

	// populate struct with the token information that we received
	spottoken := data.SpotifyToken{
		UserID:            userID,
		AccessToken:       tok.AccessToken,
		RefreshToken:      tok.RefreshToken,
		AccessTokenExpiry: tok.Expiry,
	}

	// update/insert the related information in the db, using the struct
	_, err = spottoken.Upsert(spottoken)
	if err != nil {
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
		fmt.Println(err)
	}

	// call the matches endpoint in order to begin rendering the matches page
	http.Redirect(w, r, "/matches", http.StatusSeeOther)
}

// NewAccessTokenRequest facilitates the acquisition of a new access token from Spotify,
// given that the user already has a refresh token to validate the new-token request.
func (h *Handlers) NewAccessTokenRequest(w http.ResponseWriter, r *http.Request) {
	callback := os.Getenv("URL") + "/newspotaccesstoken/callback"
	auth = spotify.NewAuthenticator(callback, spotScopes[0], spotScopes[1])
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		log.Fatal(err)
	}

	userID := h.App.Session.GetInt(r.Context(), "userID")
	userSpotTokens, _ := h.Models.SpotifyTokens.GetSpotifyTokenForUser(userID)
	oauth2SpotToken := oauth2.Token{
		AccessToken:  userSpotTokens.AccessToken,
		TokenType:    "bearer",
		RefreshToken: userSpotTokens.RefreshToken,
		Expiry:       userSpotTokens.AccessTokenExpiry,
	}

	client := auth.NewClient(&oauth2SpotToken)

	newtoken, _ := client.Token()
	userSpotTokens.AccessToken = newtoken.AccessToken
	err = userSpotTokens.UpdateAccessToken(newtoken.AccessToken, newtoken.Expiry)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Print("Spotify access token updated successfully!")
	}

	http.Redirect(w, r, "/matches", http.StatusSeeOther)
}

// NewAccessTokenAssign is a depricated version of NewAccessTokenRequest(), defined above.
func (h *Handlers) NewAccessTokenAssign(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		fmt.Println(err)
		return
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		fmt.Printf("State mismatch: %s != %s\n", st, state)
		return
	}

	spotclient := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")

	_, err = spotclient.CurrentUser()
	if err != nil {
		fmt.Println(err)
	}

	userID := h.App.Session.GetInt(r.Context(), "userID")
	if userID == 0 || userID == -1 {
		fmt.Println("The user_id of the current user could not be found in the session data.")
		return
	}

	spottoken := data.SpotifyToken{
		UserID:            userID,
		AccessToken:       tok.AccessToken,
		RefreshToken:      tok.RefreshToken,
		AccessTokenExpiry: tok.Expiry,
	}

	spottoken.Upsert(spottoken)

	http.Redirect(w, r, "/matches", http.StatusSeeOther)
}

// SetSpotifySongsForUser
func (h *Handlers) SetSpotifySongsForUser(userID int, songs spotify.FullTrackPage) error {
	// get the Spotify tokens for the current user
	SpotTok, err := h.Models.SpotifyTokens.GetSpotifyTokenForUser(userID)
	if err != nil {
		fmt.Println("No spotify token found for user")
		return err
	}

	// create an oauth2 token with our Spotify tokens
	spotifyToken := &oauth2.Token{
		AccessToken:  SpotTok.AccessToken,
		RefreshToken: SpotTok.RefreshToken,
	}

	// Use the oauth2 token to see if the current user has that token information in their browser
	// in order to confirm them as a current user.
	client := auth.NewClient(spotifyToken)
	_, err = client.CurrentUser()
	if err != nil {
		fmt.Println("Current user is nil")
		return err
	}

	// delete user's liked songs (currently the top 20) before inserting new tracks
	err = h.Models.LikedSongs.DeleteAllForUser(userID)
	if err != nil {
		fmt.Println("Error deleting all liked songs for user ", userID, ". ")
		return err
	}

	// save all Spotify song IDs from the songs that were passed in
	var songIDs []spotify.ID
	for i := range songs.Tracks {
		songIDs = append(songIDs, songs.Tracks[i].ID)
	}

	// get audio features of the songs
	var allTracksFeatures []*spotify.AudioFeatures
	if songIDs != nil {
		allTracksFeatures, err = client.GetAudioFeatures(songIDs...)
		if err != nil {
			fmt.Println("Error getting audio features from Spotify.")
			return err
		}
	}

	// insert new songs and liked songs for the user (currently their top 20)
	for x := range songs.Tracks {
		id := songs.Tracks[x].ID
		name := songs.Tracks[x].Name
		artistName := songs.Tracks[x].Album.Artists[0].Name

		temp := data.Song{
			SpotifyID:  id.String(),
			Name:       name,
			ArtistName: artistName,
		}

		// Insert the song into the db, and get the table's song ID. tID will == 0 if song already
		// exists in the songs table
		tID, err := temp.Insert(temp)
		if err != nil {
			fmt.Println("Error inserting song: ID", tID)
		}

		// If song already exists in songs table, get the song ID by name, and no need to get song aspects.
		// Else, use the current tID, get song aspects, and update the song records with musical aspects
		if tID == 0 {
			song, err := h.Models.Songs.GetByName(name)
			if err != nil {
				fmt.Println("Error getting song that has already been inserted into songs table", tID)
				return err
			}
			tID = song.ID
		} else {
			// trackFeatures, err := client.GetAudioFeatures(songs.Tracks[x].ID)
			if err != nil {
				fmt.Println("Error getting audio features for track", songs.Tracks[x].ID)
				return err
			}
			fmt.Println("ID:", songs.Tracks[x].ID, "| Name:", songs.Tracks[x].Name, "| Artist:", songs.Tracks[x].Artists[0].Name)

			// update the song in the songs table
			temp = data.Song{
				ID:               tID,
				SpotifyID:        id.String(),
				Name:             name,
				ArtistName:       artistName,
				LoudnessAvg:      float64(allTracksFeatures[x].Loudness),
				TempoAvg:         float64(allTracksFeatures[x].Tempo),
				TimeSigAvg:       allTracksFeatures[x].TimeSignature,
				Acousticness:     allTracksFeatures[x].Acousticness,
				Danceability:     allTracksFeatures[x].Danceability,
				Energy:           allTracksFeatures[x].Energy,
				Instrumentalness: allTracksFeatures[x].Instrumentalness,
				Mode:             allTracksFeatures[x].Mode,
				Speechiness:      allTracksFeatures[x].Speechiness,
				Valence:          allTracksFeatures[x].Valence,
			}

			err = temp.Update(temp)
			if err != nil {
				fmt.Println("Error updating song: ID", tID)
				return err
			}

		}

		tempLsng := data.LikedSong{
			UserID: userID,
			SongID: tID,
		}

		_, err = h.Models.LikedSongs.Insert(tempLsng)
		if err != nil {
			fmt.Println("Error inserting song ", temp.Name, " with id ", tID, " into liked_songs table")
			return err
		}
	}
	return nil
}

func (h *Handlers) CreateUserMusicProfile(userID int) error {
	user, err := h.Models.Users.Get(userID)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// If user's music profile doesn't exist, make an empty one with current date/time.
	musicProfile, _ := h.Models.UserMusicProfiles.GetByUser(*user)
	if musicProfile == nil {
		newMusicProfile, err := h.GetTracksAnalysis(*user)
		_, err = h.Models.UserMusicProfiles.Insert(newMusicProfile)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (h *Handlers) UpdateUserMusicProfile(w http.ResponseWriter, r *http.Request) {
	userID := h.App.Session.GetInt(r.Context(), "userID")
	user, err := h.Models.Users.Get(userID)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get music profile for current user.
	musicProfile, err := h.Models.UserMusicProfiles.GetByUser(*user)
	if err != nil {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"No user music profile created for user yet": "true"}`))
		return
	}

	// if the user's music profile was updated more than a day ago, update the music profile
	// fmt.Println("Music profile last updated at: ", musicProfile.UpdatedAt)
	resetTime := time.Now().Add(-time.Hour * 24)
	if musicProfile.UpdatedAt.After(resetTime) {
		updatedMusicProfile, err := h.GetTracksAnalysis(*user)
		if err != nil {
			h.App.ErrorLog.Println(err)
			return
		}

		_, err = h.Models.UserMusicProfiles.Update(*updatedMusicProfile)
		if err != nil {
			h.App.ErrorLog.Println(err)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Musical preference profile updated": "true"}`))
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"Musical preference profile updated": "false"}`))
	}
	return
}

func (h *Handlers) GetTracksAnalysis(user data.User) (*data.UserMusicProfile, error) {
	// By default will only collect the top 20 songs for this user
	SpotTok, err := h.Models.SpotifyTokens.GetSpotifyTokenForUser(user.ID)
	if err != nil {
		fmt.Println("No spotify token found for user")
		return &data.UserMusicProfile{}, err
	}

	spotifyToken := &oauth2.Token{
		AccessToken:  SpotTok.AccessToken,
		RefreshToken: SpotTok.RefreshToken,
	}

	client := auth.NewClient(spotifyToken)
	_, err = client.CurrentUser()
	if err != nil {
		fmt.Println("Current user is nil")
		return &data.UserMusicProfile{}, err
	}

	songs, err := client.CurrentUsersTopTracks()
	if err != nil {
		fmt.Println(err)
		return &data.UserMusicProfile{}, err
	}

	// Analyze songs' musical aspects, add them to songs table, and to current user's liked_songs
	err = h.SetSpotifySongsForUser(user.ID, *songs)
	if err != nil {
		fmt.Println("There was an error setting the Spotify songs for the user: ", user.ID)
		return &data.UserMusicProfile{}, err
	}

	var (
		loudnessAvg         float64
		tempoAvg            float64
		timesigAvg          []int
		acousticnessAvg     float32
		danceabilityAvg     float32
		energyAvg           float32
		instrumentalnessAvg float32
		modeMostCommon      []int
		speechinessAvg      float32
		valenceAvg          float32
		total               int
	)

	likedSongs, err := h.Models.LikedSongs.GetAllByOneUser(user.ID)
	if err != nil {
		fmt.Println("There was an error getting all liked songs for user ", user.ID, ".")
		return &data.UserMusicProfile{}, err
	}

	// add up all the musical aspects from the user's liked songs
	for x := range likedSongs {
		song, err := h.Models.Songs.Get(likedSongs[x].SongID)
		if err != nil {
			fmt.Println("There was an error getting a song from the db. Song ID: ", likedSongs[x].SongID, ".")
			return &data.UserMusicProfile{}, err
		}

		loudnessAvg += song.LoudnessAvg
		tempoAvg += song.TempoAvg
		timesigAvg = append(timesigAvg, song.TimeSigAvg)
		acousticnessAvg += song.Acousticness
		danceabilityAvg += song.Danceability
		energyAvg += song.Energy
		instrumentalnessAvg += song.Instrumentalness
		modeMostCommon = append(modeMostCommon, song.Mode)
		speechinessAvg += song.Speechiness
		valenceAvg += song.Valence
		total += 1
	}

	// find the averages of all the musical aspects from the user's liked songs
	loudnessAvg = loudnessAvg / float64(total)
	tempoAvg = tempoAvg / float64(total)
	timesigMode := h.Mode(timesigAvg, total)
	acousticnessAvg = acousticnessAvg / float32(total)
	danceabilityAvg = danceabilityAvg / float32(total)
	energyAvg = energyAvg / float32(total)
	instrumentalnessAvg = instrumentalnessAvg / float32(total)
	modeMode := h.Mode(modeMostCommon, total)
	speechinessAvg = speechinessAvg / float32(total)
	valenceAvg = valenceAvg / float32(total)

	musicProfile := data.UserMusicProfile{
		UserID:           user.ID,
		Loudness:         loudnessAvg,
		Tempo:            tempoAvg,
		TimeSig:          timesigMode,
		Acousticness:     acousticnessAvg,
		Danceability:     danceabilityAvg,
		Energy:           energyAvg,
		Instrumentalness: instrumentalnessAvg,
		Mode:             modeMode,
		Speechiness:      speechinessAvg,
		Valence:          valenceAvg,
	}

	return &musicProfile, nil
}

func (h *Handlers) BuildSectionAggregate(sections []spotify.Section) (float64, float64, int, error) {
	var (
		loudness float64
		tempo    float64
		timesig  []int
	)

	if len(sections) < 1 {
		return 0, 0, 0, errors.New("No sections provided")
	}

	for x := range sections {
		loudness += sections[x].Loudness
		tempo += sections[x].Tempo
		timesig = append(timesig, sections[x].TimeSignature)
	}

	loudness = loudness / float64(len(sections))
	tempo = tempo / float64(len(sections))
	timesigMode := h.Mode(timesig, len(sections))

	return loudness, tempo, timesigMode, nil
}
