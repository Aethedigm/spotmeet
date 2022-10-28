package handlers

import (
	"fmt"
	"github.com/zmb3/spotify"
	"log"
	"myapp/data"
	"net/http"
	"os"
)

var auth = spotify.Authenticator{}
var state = "abc123"

func (h *Handlers) UserRegister(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "register", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "login", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Redirect(w, r, "/users/login?loginFailed=true", http.StatusSeeOther)
		return
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	user, err := h.Models.Users.GetByEmail(email)
	if err != nil {
		http.Redirect(w, r, "/users/login?loginFailed=true", http.StatusSeeOther)
		return
	}

	matches, err := user.PasswordMatches(password)
	if err != nil {
		http.Redirect(w, r, "/users/login?loginFailed=true", http.StatusSeeOther)
		return
	}

	if !matches {
		http.Redirect(w, r, "/users/login?loginFailed=true", http.StatusSeeOther)
		return
	}

	h.App.Session.Put(r.Context(), "userID", user.ID)

	// Need to get the specific Spotify redirect and access tokens for the user_id we just found.
	// If we do not do this, then if the browser signs in as a new or other user, it keeps the tokens from the
	// last user who was logged in.
	// Also, we need to wipe these spotify tokens from the session data once the app user purposefully
	// logs out of SpotMeet.

	http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
	// http.Redirect(w, r, "/", http.StatusSeeOther)

}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	h.App.Session.RenewToken(r.Context())
	h.App.Session.Remove(r.Context(), "userID")
	http.Redirect(w, r, "/users/login", http.StatusSeeOther)
}

func (h *Handlers) SpotifyAuthorization(w http.ResponseWriter, r *http.Request) {
	callback := os.Getenv("LOCALHOST_URL") + "/spotauth/callback"
	auth = spotify.NewAuthenticator(
		callback,
		// these scopes may need to be changed out for this app
		spotify.ScopeUserTopRead,
		spotify.ScopeUserReadRecentlyPlayed)

	url := auth.AuthURL(state)
	//fmt.Println("Log in to Spotify by visiting this page:", url)
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func (h *Handlers) SpotifyAuthorizationCallback(w http.ResponseWriter, r *http.Request) {
	tok, err := auth.Token(state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	// use the token to get an authenticated client
	client := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")
	//ch <- &client

	user, err := client.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)

	// get the user_id of the user who has gained access to spotify
	userID := h.App.Session.GetInt(r.Context(), "userID")
	if userID == 0 {
		log.Fatal("The user_id of the current user could not be found in the session data.")
	}

	// Write the new refresh token and new access token to a new
	// spotifytoken struct, and assign that to the current User's SpotifyToken variable.
	spottoken := data.SpotifyToken{
		UserID:            userID,
		AccessToken:       tok.AccessToken,
		RefreshToken:      tok.RefreshToken,
		AccessTokenExpiry: tok.Expiry,
		// TokenType: tok.TokenType, // I don't think we need this b/c a user will always be a 'bearer'
	}

	// insert the new SpotifyToken into the database
	spottoken.Insert(spottoken) // need error-handling here once migrations have been made

	http.Redirect(w, r, "/matches", http.StatusSeeOther)
}
