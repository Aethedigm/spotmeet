package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

type PasswordRequest struct {
	Password string `json:"password"`
}

// ChangePassword handles the request for a different password to be implemented
func (h *Handlers) ChangePassword(w http.ResponseWriter, r *http.Request) {
	var p PasswordRequest

	userIDstr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		h.App.ErrorLog.Println("Unable to convert to int", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = json.NewDecoder(r.Body).Decode(&p)
	if err != nil {
		h.App.ErrorLog.Println("Unable to deserialize", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	err = h.Models.Users.ResetPassword(userID, p.Password)
	if err != nil {
		h.App.ErrorLog.Println("Unable to change password", err)
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	w.Write([]byte(`{"status":"ok"}`))

	// Delete all requests for this user
	err = h.Models.RecoveryEmails.DeleteAllForUser(userID)
	if err != nil {
		h.App.ErrorLog.Println("Error clearing password reset requests for user", userID)
		return
	}
}

// UserRegister displays the registration view
func (h *Handlers) UserRegister(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "register", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

// UserLogin displays the login view
func (h *Handlers) UserLogin(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "login", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

// PostUserLogin handles login form information by validating it, t
func (h *Handlers) PostUserLogin(w http.ResponseWriter, r *http.Request) {
	// get user input from the form
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println("Error parsing form", err)
		http.Redirect(w, r, "/users/login?loginFailed=true", http.StatusSeeOther)
		return
	}

	// break out user input from the form
	email := r.Form.Get("email")
	password := r.Form.Get("password")

	// get User's information while simultaneously checking if they exist in db
	user, err := h.Models.Users.GetByEmail(email)
	if err != nil {
		h.App.ErrorLog.Println("Error getting user by email", err)
		http.Redirect(w, r, "/users/login?loginFailed=true", http.StatusSeeOther)
		return
	}

	// verify if the password is correct
	matches, err := user.PasswordMatches(password)
	if err != nil {
		h.App.ErrorLog.Println("Error checking password", err)
		http.Redirect(w, r, "/users/login?loginFailed=true", http.StatusSeeOther)
		return
	}

	if !matches {
		h.App.ErrorLog.Println("Password does not match")
		http.Redirect(w, r, "/users/login?loginFailed=true", http.StatusSeeOther)
		return
	}

	// put the user's id into the browser session for later use
	h.App.Session.Put(r.Context(), "userID", user.ID)

	// verify if the user already has been authorized with Spotify. If not, send them to Spotify
	// to get authorized.
	_, err = h.Models.SpotifyTokens.GetSpotifyTokenForUser(user.ID)
	if err != nil {
		h.App.ErrorLog.Println("Spotify token for user does not exist. Going to /users/spotauth to get one.", err)
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
		return
	}

	//err = h.SetSpotifyArtistsForUser(user.ID)
	//if err != nil {
	//	http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
	//	return
	//}

	// go to matches page
	http.Redirect(w, r, "/matches", http.StatusSeeOther)
}

func (h *Handlers) Logout(w http.ResponseWriter, r *http.Request) {
	h.App.Session.Remove(r.Context(), "userID")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
