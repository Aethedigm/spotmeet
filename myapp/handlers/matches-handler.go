package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"myapp/data"
	"myapp/middleware"
	"net/http"
	"time"
)

// MyMatchResults gets Matches that have already been created
func (h *Handlers) MyMatchResults(w http.ResponseWriter, r *http.Request) {
	matches, err := h.Models.Matches.GetAllForOneUser(h.App.Session.GetInt(r.Context(), "userID"))
	if err != nil {
		fmt.Println("Error getting matches:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	matchesJSON, err := json.Marshal(matches)
	if err != nil {
		fmt.Println("Error marshalling matches:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	_, err = w.Write(matchesJSON)
	if err != nil {
		h.App.ErrorLog.Println("error writing json")
		return
	}
}

// Matches calls for the rendering of the matches page, if a user exists for the browser session,
// and, if the current user has an access token that will not expire within 5 minutes.
func (h *Handlers) Matches(w http.ResponseWriter, r *http.Request) {
	if !h.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "users/login", http.StatusSeeOther)
		return
	}

	// Check the user's spotify access token expiry. If it's within 5 minutes of its
	// -- expiry, redirect to users/newspotaccesstoken to get a new one from spotify via the refresh token
	// Get user ID from the session
	userID := h.App.Session.GetInt(r.Context(), "userID")
	userSpotTokens, _ := h.Models.SpotifyTokens.GetSpotifyTokenForUser(userID)
	// convert to east-coat time zone, that's used by time.Time (add 4 hrs)
	expiry := userSpotTokens.AccessTokenExpiry.Unix() + 14400
	fiveMinutesFromNow := time.Now().Add(time.Minute * 5).Unix()
	if expiry < fiveMinutesFromNow {
		fmt.Println(fiveMinutesFromNow)
		// go to users/newspotaccesstoken to get the new access token with the refresh token
		http.Redirect(w, r, "users/newspotaccesstoken", http.StatusSeeOther)
		return
	}

	// find new matches for the user
	numberOfNewMatches, err := middleware.CreateMatches(userID)
	if err != nil {
		h.App.ErrorLog.Println("Error creating matches for user_id ", userID, ":", err)
	}

	// get variables ready to pass into the matches view
	u := data.User{}
	user, err := u.Get(userID)
	if err != nil {
		h.App.ErrorLog.Println("Error getting User struct for user_id ", userID, ":", err)
	}
	userName := user.FirstName
	vars := make(jet.VarMap)
	vars.Set("numberOfNewMatches", numberOfNewMatches)
	vars.Set("userName", userName)

	// render the page, and pass variables into it
	err = h.App.Render.Page(w, r, "matches", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}
