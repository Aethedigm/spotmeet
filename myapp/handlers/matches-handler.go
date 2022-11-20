package handlers

import (
	"encoding/json"
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
		h.App.ErrorLog.Println("Error converting matchID to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match, err := h.Models.Matches.Get(matchID)
	if err != nil {
		h.App.ErrorLog.Println("Error getting match:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match.Complete = true

	err = h.Models.Matches.Update(*match)
	if err != nil {
		h.App.ErrorLog.Println("Error updating match:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	link := data.Link{}
	link.User_A_ID = match.User_A_ID
	link.User_B_ID = match.User_B_ID
	link.PercentLink = 100
	link.ArtistID = match.ArtistID
	if link.ArtistID == 0 {
		link.ArtistID, err = h.Models.Artists.GetOneID()
		if err != nil {
			h.App.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	link.CreatedAt = time.Now()

	_, err = h.Models.Links.Insert(link)
	if err != nil {
		h.App.ErrorLog.Println("Error inserting link:", err)
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

	// get potential matches that qualify based on coordinates and looking-for
	users, err := h.Models.RQ.MatchQuery(*user, *settings)
	if err != nil {
		h.App.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// loop though the potential matches and confirm Match by checking against music profiles
	for i := range users {
		match := data.Match{}
		match.User_A_ID = userID
		match.User_B_ID = users[i]
		match.PercentMatch = 100
		match.CreatedAt = time.Now()
		match.ArtistID, err = h.Models.Artists.GetOneID()
		if err != nil {
			h.App.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		_, err := h.Models.Matches.Insert(match)
		if err != nil {
			h.App.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
	}

	matchesForDisplay, err := h.Models.RQ.MatchesDisplayQuery(userID)
	if err != nil {
		h.App.ErrorLog.Println("Error with MatchesDisplayQuery(), called in matches-handler.go:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	emptyJSON, err := json.Marshal("")

	if matchesForDisplay == nil {
		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(emptyJSON)
		if err != nil {
			h.App.ErrorLog.Println("error writing json")
			return
		}
	} else {
		matchesJSON, err := json.Marshal(matchesForDisplay)
		if err != nil {
			h.App.ErrorLog.Println("Error marshalling matchesForDisplay:", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(matchesJSON)
		if err != nil {
			h.App.ErrorLog.Println("error writing json")
			return
		}
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
		h.App.ErrorLog.Println("Error getting spotify token.", err)
		http.Redirect(w, r, "/users/spotauth", http.StatusSeeOther)
		return
	}

	// Commenting out because this is already called on login (auth-handlers.go, PostUserLogin())
	//err = h.SetSpotifyArtistsForUser(userID)
	//if err != nil {
	//	h.App.ErrorLog.Println("Error setting spotify artists for user.", err)
	//}

	expiry := userSpotTokens.AccessTokenExpiry.Unix() + 18000
	fiveMinutesFromNow := time.Now().Add(time.Minute * 5).Unix()
	if expiry < fiveMinutesFromNow {
		h.App.ErrorLog.Println(fiveMinutesFromNow)
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
