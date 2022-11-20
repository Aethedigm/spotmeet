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

func (h *Handlers) SetSpotifyArtistsForUser(userID int) error {
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

	artists, err := client.CurrentUsersFollowedArtistsOpt(50, "0")
	if err != nil {
		fmt.Println("Error getting artists", err)
		return err
	}

	if len(artists.Artists) < 1 {
		fmt.Println("No artists returned")
		return err
	}

	for i := range artists.Artists {
		temp := data.Artist{
			SpotifyID: artists.Artists[i].ID.String(),
			Name:      artists.Artists[i].Name,
		}

		tID, err := temp.Insert(temp)
		if err != nil {
			fmt.Println("Error inserting artist ID", tID)
			return err
		}

		tempLart := data.LikedArtist{
			UserID:     userID,
			ArtistID:   tID,
			LikedLevel: 100,
		}

		_, err = h.Models.LikedArtists.Insert(tempLart)
		if err != nil {
			fmt.Println("Error inserting liked artist", tID, temp.Name)
			return err
		}
	}
	return nil
}

func (h *Handlers) Tracks(w http.ResponseWriter, r *http.Request) {
	userID := h.App.Session.Get(r.Context(), "userID")

	user, err := h.Models.Users.Get(userID.(int))
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get music profile for current user. If doesn't exist, make an empty one with current date/time.
	musicProfile, err := h.Models.UserMusicProfiles.GetByUser(*user)
	if musicProfile == nil {
		newMusicProfile, err := h.GetUserMusicProfile(*user) // move this inside if statement, if errs go away
		_, err = h.Models.UserMusicProfiles.Insert(newMusicProfile)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else if err != nil {
		fmt.Println(err)
		return
	} else {
		// if the user's music profile was updated more than a day ago, update the music profile
		fmt.Println("Music profile last updated at: ", musicProfile.UpdatedAt)
		if musicProfile.UpdatedAt.Before(time.Now().Truncate(time.Hour * 24)) {
			updatedMusicProfile, err := h.GetUserMusicProfile(*user)
			_, err = h.Models.UserMusicProfiles.Update(updatedMusicProfile)
			if err != nil {
				fmt.Println(err)
				return
			}

			w.Write([]byte("Success"))
		}
	}
}

func (h *Handlers) GetUserMusicProfile(user data.User) (*data.UserMusicProfile, error) {
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

	var (
		loudnessAvg float64
		tempoAvg    float64
		timesigAvg  int
		total       int
	)

	for x := range songs.Tracks {
		trackAnalysis, err := client.GetAudioAnalysis(songs.Tracks[x].ID)
		if err != nil {
			fmt.Println("Error getting audio analysis for track", songs.Tracks[x].ID)
		} else {
			fmt.Println("ID:", songs.Tracks[x].ID, "| Name:", songs.Tracks[x].Name, "| Artist:", songs.Tracks[x].Artists[0].Name)
			loud, temp, timeSig, err := BuildSectionAggregate(trackAnalysis.Sections)
			if err != nil {
				fmt.Println("Error building section aggregate")
				continue
			}

			loudnessAvg += loud
			tempoAvg += temp
			timesigAvg += timeSig
			total++
		}
	}

	loudnessAvg = loudnessAvg / float64(total)
	tempoAvg = tempoAvg / float64(total)
	timesigAvg = timesigAvg / total

	musicProfile := data.UserMusicProfile{
		UserID:   user.ID,
		Loudness: loudnessAvg,
		Tempo:    tempoAvg,
		TimeSig:  timesigAvg,
	}

	return &musicProfile, nil
}

func BuildSectionAggregate(sections []spotify.Section) (float64, float64, int, error) {
	var (
		loudness float64
		tempo    float64
		timesig  int
	)

	if len(sections) < 1 {
		return 0, 0, 0, errors.New("No sections provided")
	}

	for x := range sections {
		loudness += sections[x].Loudness
		tempo += sections[x].Tempo
		timesig += sections[x].TimeSignature
	}

	loudness = loudness / float64(len(sections))
	tempo = tempo / float64(len(sections))
	timesig = timesig / len(sections)

	return loudness, tempo, timesig, nil
}
