package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

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

	err := h.App.Render.Page(w, r, "matches", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}
