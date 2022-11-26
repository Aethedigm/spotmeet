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

var spotScopes = []string{spotify.ScopeUserTopRead, spotify.ScopeUserReadRecentlyPlayed, spotify.ScopeUserFollowRead}

func (h *Handlers) SpotifyAuthorization(w http.ResponseWriter, r *http.Request) {
	callback := os.Getenv("URL") + "/spotauth/callback"
	auth = spotify.NewAuthenticator(callback, spotScopes...)
	randomBytes := make([]byte, 16)
	_, err := rand.Read(randomBytes)
	if err != nil {
		fmt.Println(err)
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
		return
	}

	state = base32.StdEncoding.WithPadding(base32.NoPadding).EncodeToString(randomBytes)

	url := auth.AuthURL(state)
	url = url + "&show_dialog=true"

	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *Handlers) SpotifyAuthorizationCallback(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
		fmt.Println(err)
		return
	}

	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
		fmt.Printf("State mismatch: %s != %s\n", st, state)
		return
	}

	spotclient := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")

	_, err = spotclient.CurrentUser()
	if err != nil {
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
		fmt.Println(err)
		return
	}

	userID := h.App.Session.GetInt(r.Context(), "userID")
	if userID == 0 || userID == -1 {
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
		fmt.Println("The user_id of the current user could not be found in the session data.")
		return
	}

	spottoken := data.SpotifyToken{
		UserID:            userID,
		AccessToken:       tok.AccessToken,
		RefreshToken:      tok.RefreshToken,
		AccessTokenExpiry: tok.Expiry,
	}

	_, err = spottoken.Upsert(spottoken)
	if err != nil {
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
		fmt.Println(err)
	}

	http.Redirect(w, r, "/matches", http.StatusSeeOther)
}

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

func (h *Handlers) SetSpotifySongsForUser(userID int, songs spotify.FullTrackPage) error {
	SpotTok, err := h.Models.SpotifyTokens.GetSpotifyTokenForUser(userID)
	if err != nil {
		fmt.Println("No spotify token found for user")
		return err
	}

	spotifyToken := &oauth2.Token{
		AccessToken:  SpotTok.AccessToken,
		RefreshToken: SpotTok.RefreshToken,
	}

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
			return err
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
			trackAnalysis, err := client.GetAudioAnalysis(songs.Tracks[x].ID)
			if err != nil {
				fmt.Println("Error getting audio analysis for track", songs.Tracks[x].ID)
				// add faux data here for the update of the record
			} else {
				fmt.Println("ID:", songs.Tracks[x].ID, "| Name:", songs.Tracks[x].Name, "| Artist:", songs.Tracks[x].Artists[0].Name)
				loud, tempo, timeSig, err := h.BuildSectionAggregate(trackAnalysis.Sections)
				if err != nil {
					fmt.Println("Error building section aggregate")
					continue
				}

				// update the song in the songs table
				temp = data.Song{
					ID:          tID,
					SpotifyID:   id.String(),
					Name:        name,
					ArtistName:  artistName,
					LoudnessAvg: loud,
					TempoAvg:    tempo,
					TimeSigAvg:  timeSig,
				}

				err = temp.Update(temp)
				if err != nil {
					fmt.Println("Error updating song: ID", tID)
					return err
				}
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
	fmt.Println("Music profile last updated at: ", musicProfile.UpdatedAt)
	if musicProfile.UpdatedAt.Before(time.Now().Truncate(time.Hour * 24)) {
		updatedMusicProfile, err := h.GetTracksAnalysis(*user)
		_, err = h.Models.UserMusicProfiles.Update(*updatedMusicProfile)
		if err != nil {
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
		fmt.Println("There was an error setting the Spotify songs for the user: ", user.ID, ".")
		return &data.UserMusicProfile{}, err
	}

	var (
		loudnessAvg float64
		tempoAvg    float64
		timesigAvg  []int
		total       int
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
		total += 1
	}

	// find the averages of all the musical aspects from the user's liked songs
	loudnessAvg = loudnessAvg / float64(total)
	tempoAvg = tempoAvg / float64(total)
	timesigMode := h.Mode(timesigAvg, total)

	musicProfile := data.UserMusicProfile{
		UserID:   user.ID,
		Loudness: loudnessAvg,
		Tempo:    tempoAvg,
		TimeSig:  timesigMode,
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
