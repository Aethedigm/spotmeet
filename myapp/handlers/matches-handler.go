package handlers

import (
	"encoding/json"
	"fmt"
	"myapp/data"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/CloudyKit/jet/v6"
)

func (h *Handlers) RejectMatch(w http.ResponseWriter, r *http.Request) {
	matchIDstr := chi.URLParam(r, "matchID")

	matchID, err := strconv.Atoi(matchIDstr)
	if err != nil {
		h.App.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match, err := h.Models.Matches.Get(matchID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match.Complete = true

	err = h.Models.Matches.Update(*match)
	if err != nil {
		h.App.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": "true"}`))
}

func (h *Handlers) AcceptMatch(w http.ResponseWriter, r *http.Request) {
	matchIDstr := chi.URLParam(r, "matchID")

	matchID, err := strconv.Atoi(matchIDstr)
	if err != nil {
		fmt.Println("Error converting matchID to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match, err := h.Models.Matches.Get(matchID)
	if err != nil {
		fmt.Println("Error getting match:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match.Complete = true

	err = h.Models.Matches.Update(*match)
	if err != nil {
		fmt.Println("Error updating match:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	link := data.Link{}
	link.User_A_ID = match.User_A_ID
	link.User_B_ID = match.User_B_ID
	link.PercentLink = 100
	link.ArtistID = 1
	link.CreatedAt = time.Now()

	_, err = h.Models.Links.Insert(link)
	if err != nil {
		fmt.Println("Error inserting link:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": "true"}`))
}

func (h *Handlers) MyMatchResults(w http.ResponseWriter, r *http.Request) {
	userID := h.App.Session.GetInt(r.Context(), "userID")
	user, err := h.Models.Users.Get(userID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	settings, err := h.Models.Settings.GetByUserID(userID)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}

	users, err := h.Models.RQ.MatchQuery(*user, *settings)
	if err != nil {
		h.App.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	for i := range users {
		match := data.Match{}
		match.User_A_ID = userID
		match.User_B_ID = users[i]
		match.PercentMatch = 100
		match.CreatedAt = time.Now()

		_, err := h.Models.Matches.Insert(match)
		if err != nil {
			h.App.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	matches, err := h.Models.Matches.GetAllForOneUser(userID)
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

	userID := h.App.Session.GetInt(r.Context(), "userID")
	userSpotTokens, err := h.Models.SpotifyTokens.GetSpotifyTokenForUser(userID)
	if err != nil {
		fmt.Println("Error getting spotify token.", err)
		http.Redirect(w, r, "/users/login?spotConnFailed=true", http.StatusSeeOther)
		return
	} else {
		err = h.SetSpotifyArtistsForUser(userID)
		if err != nil {
			fmt.Println("Error setting spotify artists for user.", err)
		}
	}

	expiry := userSpotTokens.AccessTokenExpiry.Unix() + 18000
	fiveMinutesFromNow := time.Now().Add(time.Minute * 5).Unix()
	if expiry < fiveMinutesFromNow {
		fmt.Println(fiveMinutesFromNow)
		// go to users/newspotaccesstoken to get the new access token with the refresh token
		http.Redirect(w, r, "users/newspotaccesstoken", http.StatusSeeOther)
		return
	}

	vars := make(jet.VarMap)
	vars.Set("userID", userID)

	err = h.App.Render.Page(w, r, "matches", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}
