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

	// Create user's music profile preferences, if user has no music profile record yet (the check
	// is inside of the function).
	// The Update version of this is called from the view via fetch() for behind-the-scenes updating.
	// The Create version is called here when we have a new user, in order to prioritize gathering this
	// info, so that Matches can be displayed upon first matches.jet page load.
	err = h.CreateUserMusicProfile(userID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	settings, err := h.Models.Settings.GetByUserID(userID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	// get potential matches that qualify based on coordinates and looking-for
	users, err := h.Models.RQ.MatchQuery(*user, *settings)
	if err != nil {
		h.App.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// if any potential matches were found
	if users != nil {
		// get the musical preference profile of the current user
		currentUserMusicProfile, err := h.Models.UserMusicProfiles.GetByUser(*user)
		if err != nil {
			h.App.ErrorLog.Println(err)
			return
		}

		// loop though the potential matches and confirm Match by checking against music profiles
		for i := range users {
			// get the musical preference profile of the other user
			otherUserMusicProfile, err := h.Models.UserMusicProfiles.GetByUserID(users[i])
			if err != nil {
				h.App.ErrorLog.Println(err)
			}

			// get settings of the other user
			otherUserSettings, err := h.Models.Settings.GetByUserID(users[i])
			if err != nil {
				h.App.ErrorLog.Println(err)
			}

			// check if the two users' musical preference profiles are compatible
			itsAMatch, matchPercentage := h.CompareUserMusicProfiles(*currentUserMusicProfile,
				*otherUserMusicProfile, settings.MatchSensitivity, otherUserSettings.MatchSensitivity)

			if itsAMatch == true {

				// if true, then create the match
				match := data.Match{}
				match.User_A_ID = userID
				match.User_B_ID = users[i]
				match.PercentMatch = float32(matchPercentage)
				match.CreatedAt = time.Now()
				// match.ArtistID, err = h.Models.Artists.GetOneID()
				match.ArtistID = 0
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

	var isFirstLogin bool
	musicProfile, _ := h.Models.UserMusicProfiles.GetByUserID(userID)
	fmt.Println("THIS IS THE CURRENT USER ID: ", userID)
	if musicProfile == nil {
		isFirstLogin = true
	} else {
		isFirstLogin = false
	}

	vars := make(jet.VarMap)
	vars.Set("userID", userID)
	vars.Set("isFirstLogin", isFirstLogin)

	err = h.App.Render.Page(w, r, "matches", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) CompareUserMusicProfiles(profileA data.UserMusicProfile,
	profileB data.UserMusicProfile, matchSensitivityUserA int, matchSensitivityUserB int) (bool, int) {

	// Here are the values of the user_music_profile table that we are checking (for now)
	// Loudness  float64   `db:"loudness" json:"loudness"`
	// Tempo     float64   `db:"tempo" json:"tempo"`
	// TimeSig   int       `db:"time_sig" json:"time_sig"`
	var similarCount int
	var highestMatchSensitivity int
	var matchOnProfiles bool // return value
	var matchPercentage int  // return value
	const aspectCount = 3    // Loudness, Tempo, TimeSig (more to come)
	const loudnessRange = float64(4.0)
	const tempoRange = float64(12.0)

	// lower sensitivity means more chances of matching (looser aspect count restrictions)
	// higher sensitivity means less chances of matching (stricter aspect count restrictions)
	const lowSensitivityRange = 2
	const mediumSensitivityRange = 1
	const highSensitivityRange = 0

	// using the highest match sensitivity given
	if matchSensitivityUserA >= matchSensitivityUserB {
		highestMatchSensitivity = matchSensitivityUserA
	} else {
		highestMatchSensitivity = matchSensitivityUserB
	}

	// compare loudness averages
	if (profileA.Loudness <= profileB.Loudness+loudnessRange) &&
		(profileA.Loudness >= profileB.Loudness-loudnessRange) {
		similarCount += 1
	}

	// compare tempo averages
	if (profileA.Tempo <= profileB.Tempo+tempoRange) &&
		(profileA.Tempo >= profileB.Tempo-tempoRange) {
		similarCount += 1
	}

	// Compare time signature averages. See if users prefer signatures that are
	// multiples of 2, eg. 4/4, 2/4, multiples of 3, eg. 3/4, 6/4, 12/8,
	// or if both users prefer signatures that do not fit into these two standard,
	// widely-common-in-every-genre time signatures.
	// (Time signatures that are only divisible by 2 and not 3 warrant a much different
	// sound than time signatures that are divisible by 2, but that are also divisible by 3.)
	if ((profileA.TimeSig%2 == 0 && profileA.TimeSig%3 != 0) &&
		(profileB.TimeSig%2 == 0 && profileB.TimeSig%3 != 0)) ||
		(profileA.TimeSig%3 == 0 && profileB.TimeSig%3 == 0) ||
		((profileA.TimeSig%2 != 0 && profileA.TimeSig%3 != 0) &&
			(profileB.TimeSig%2 != 0 && profileB.TimeSig%3 != 0)) {
		similarCount += 1
	}

	// Decide if a match, based on the aspect count, the match sensitivity, and
	// the number of musical aspects that were marked as similar.
	switch {
	case highestMatchSensitivity <= 33:
		if similarCount <= aspectCount+lowSensitivityRange &&
			similarCount >= aspectCount-lowSensitivityRange {
			matchOnProfiles = true
		} else {
			matchOnProfiles = false
		}
	case highestMatchSensitivity >= 34 && matchSensitivityUserA <= 67:
		if similarCount <= aspectCount+mediumSensitivityRange &&
			similarCount >= aspectCount-mediumSensitivityRange {
			matchOnProfiles = true
		} else {
			matchOnProfiles = false
		}
	case highestMatchSensitivity >= 68:
		if similarCount <= aspectCount+highSensitivityRange &&
			similarCount >= aspectCount-highSensitivityRange {
			matchOnProfiles = true
		} else {
			matchOnProfiles = false
		}
	}

	matchPercentage = (similarCount / aspectCount) * 100

	return matchOnProfiles, matchPercentage
}
