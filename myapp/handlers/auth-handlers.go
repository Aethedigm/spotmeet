package handlers

import (
	"fmt"
	"log"
	"myapp/data"
	"net/http"
	"os"

	"github.com/zmb3/spotify"
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

	_, err = h.Models.SpotifyTokens.GetSpotifyTokenForUser(user.ID)
	if err != nil {
		// User does not have current token, so redirect to Spotify auth
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
	}

	// User has current token, so redirect to matches page
	http.Redirect(w, r, "/matches", http.StatusSeeOther)
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

	spotclient := auth.NewClient(tok)
	fmt.Fprintf(w, "Login Completed!")

	_, err = spotclient.CurrentUser()
	if err != nil {
		log.Fatal(err)
	}

	userID := h.App.Session.GetInt(r.Context(), "userID")
	if userID == 0 || userID == -1 {
		log.Fatal("The user_id of the current user could not be found in the session data.")
	}

	// Write the new refresh token and new access token to a new
	// spotifytoken struct, and assign that to the current User's SpotifyToken variable.
	spottoken := data.SpotifyToken{
		UserID:            userID,
		AccessToken:       tok.AccessToken,
		RefreshToken:      tok.RefreshToken,
		AccessTokenExpiry: tok.Expiry,
	}

	spottoken.Upsert(spottoken)

	http.Redirect(w, r, "/matches", http.StatusSeeOther)
}
