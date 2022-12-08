package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"myapp/data"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
)

// matches-handler.go contains logic relating specifically to the creation and removal of Matches

// RejectMatch processes a user's decision to decline a match
func (h *Handlers) RejectMatch(w http.ResponseWriter, r *http.Request) {
	// get the match ID from the url parameter
	matchIDstr := chi.URLParam(r, "matchID")
	matchID, err := strconv.Atoi(matchIDstr)
	if err != nil {
		h.App.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the match's struct
	match, err := h.Models.Matches.Get(matchID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// set it as complete, but refrain from creating a "Link" due to this match being declined
	match.Complete = true

	// update the match in the db by passing it the updated struct
	err = h.Models.Matches.Update(*match)
	if err != nil {
		h.App.ErrorLog.Println(err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send success message to browser
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": "true"}`))
}

// AcceptMatch processes a user's decision to accept a match
func (h *Handlers) AcceptMatch(w http.ResponseWriter, r *http.Request) {
	// get the match ID from the url parameter
	matchIDstr := chi.URLParam(r, "matchID")
	matchID, err := strconv.Atoi(matchIDstr)
	if err != nil {
		h.App.ErrorLog.Println("Error converting matchID to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the match's struct
	match, err := h.Models.Matches.Get(matchID)
	if err != nil {
		h.App.ErrorLog.Println("Error getting match:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// set it as complete
	match.Complete = true

	// update the match in the db by passing it the updated struct
	err = h.Models.Matches.Update(*match)
	if err != nil {
		h.App.ErrorLog.Println("Error updating match:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// since the match was accepted by the user, create a "Link" as a result.
	link := data.Link{}
	link.User_A_ID = match.User_A_ID
	link.User_B_ID = match.User_B_ID
	link.PercentLink = 100
	link.SongID = match.SongID
	if link.SongID == 0 {
		link.SongID, err = h.Models.Artists.GetOneID()
		if err != nil {
			h.App.ErrorLog.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	link.CreatedAt = time.Now()

	// insert the new link into the db
	_, err = h.Models.Links.Insert(link)
	if err != nil {
		h.App.ErrorLog.Println("Error inserting link:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// write success message to the browser
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"success": "true"}`))
}

// MyMatchResults creates new matches for a user based on both user's location and preferences
// designated in their settings. Then, all newly-made matches and previously-existing matches are gathered in the
// proper way for display to the user.
func (h *Handlers) MyMatchResults(w http.ResponseWriter, r *http.Request) {
	// get the current user's user-table data as a struct
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

	// get a struct of the current user's settings
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
			if otherUserMusicProfile == nil {
				continue
			}

			// get settings of the other user
			otherUserSettings, err := h.Models.Settings.GetByUserID(users[i])
			if err != nil {
				h.App.ErrorLog.Println(err)
			}

			// check if the two users' musical preference profiles are compatible
			itsAMatch, matchPercentage, songIDMatchedOn := h.CompareUserMusicProfiles(*currentUserMusicProfile,
				*otherUserMusicProfile, settings.MatchSensitivity, otherUserSettings.MatchSensitivity)

			if itsAMatch == true {

				// if true, then create the match
				match := data.Match{}
				match.User_A_ID = userID
				match.User_B_ID = users[i]
				match.PercentMatch = float32(matchPercentage)
				match.CreatedAt = time.Now()
				// match.ArtistID, err = h.Models.Artists.GetOneID()
				match.SongID = songIDMatchedOn
				match.UserAViewedThread = false
				match.UserBViewedThread = false
				if err != nil {
					h.App.ErrorLog.Println(err)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}

				// insert the new match into the db
				_, err := h.Models.Matches.Insert(match)
				if err != nil {
					h.App.ErrorLog.Println(err)
					http.Error(w, err.Error(), http.StatusBadRequest)
					return
				}
			}
		}
	}

	// gather the necessary match-information for displaying matches to the user
	matchesForDisplay, err := h.Models.RQ.MatchesDisplayQuery(userID)
	if err != nil {
		h.App.ErrorLog.Println("Error with MatchesDisplayQuery(), called in matches-handler.go:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// initialize json container
	emptyJSON, err := json.Marshal("")

	// if no matches, send empty container to user
	// else, convert match information into json and send to the user
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

// Matches renders the matches page. Actual match information is called from the JS within the page
// that hits an endpoint, calling MyMatchResults() within this .go file.
func (h *Handlers) Matches(w http.ResponseWriter, r *http.Request) {
	// only continue if current user session exists
	if !h.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "users/login", http.StatusSeeOther)
		return
	}

	// get current user's ID
	userID := h.App.Session.GetInt(r.Context(), "userID")

	// Get current user's Spotify authorization credentials. If they don't exist,
	// go to Spotify to get them.
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

	// Send the user to renew their Spotify credentials if they expire within the next 5 minutes
	expiry := userSpotTokens.AccessTokenExpiry.Unix() + 18000
	fiveMinutesFromNow := time.Now().Add(time.Minute * 5).Unix()
	if expiry < fiveMinutesFromNow {
		h.App.ErrorLog.Println(fiveMinutesFromNow)
		http.Redirect(w, r, "users/newspotaccesstoken", http.StatusSeeOther)
		return
	}

	// Create flags that we will pass into the view, so that the JavaScript can call endpoints that
	// may be needed to be hit.
	var isFirstLogin bool
	var locationUpdateNeeded bool
	musicProfile, _ := h.Models.UserMusicProfiles.GetByUserID(userID)

	// If the user's music profile doesn't exist yet, we know that the user is new, and location data
	// is still needed. Else, indicate the user is not new and further evaluate whether or not the user
	// needs to have their location re-gathered.
	if musicProfile == nil {
		isFirstLogin = true
		locationUpdateNeeded = true
	} else {
		isFirstLogin = false
		user, err := h.Models.Users.Get(userID)
		if err != nil {
			h.App.ErrorLog.Println("error getting struct for user ", userID, ". ")
		}

		// Only get user's location if the user was updated more than 10 minutes prior to the current time.
		if user.UpdatedAt.Before(time.Now().Truncate(time.Minute * 10)) {
			locationUpdateNeeded = true
		} else {
			locationUpdateNeeded = false
		}
	}

	// get settings of the user so we can pass the theme preference into the view
	settings, err := h.Models.Settings.GetByUserID(userID)
	if err != nil {
		h.App.ErrorLog.Println("error getting settings for user:", err)
	}

	// prepare data needed by the view
	vars := make(jet.VarMap)
	vars.Set("userID", userID)
	vars.Set("isFirstLogin", isFirstLogin)
	vars.Set("locationUpdateNeeded", locationUpdateNeeded)
	vars.Set("theme", settings.Theme)

	// send data to the matches view and render the page
	err = h.App.Render.Page(w, r, "matches", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

// CompareUserMusicProfiles is the last check in the matching algorithm. This function is called
// with potential matches (profileA, profileB) that have already passed through location and
// settings-preferences checks.
// This function returns a boolean indicating whether or not the users are a successful match (matchOnProfiles),
// the "percentage" or "strength" the match (matchPercentage), and songIDMatchedOn, the ID of the song between both users' liked-songs
// that a new musical profile (created from the averaging of both user's musical preferences) matches closest to.
// the matchOnProfiles, matchPercentage, songIDMatchedOn
func (h *Handlers) CompareUserMusicProfiles(profileA data.UserMusicProfile,
	profileB data.UserMusicProfile, matchSensitivityUserA int, matchSensitivityUserB int) (bool, int, int) {

	// Here are the values of the user_music_profile table that we are checking:
	// Loudness  		float64   `db:"loudness" json:"loudness"`
	// Tempo     		float64   `db:"tempo" json:"tempo"`
	// TimeSig   		int       `db:"time_sig" json:"time_sig"`
	// Acousticness     float32   `db:"acousticness" json:"acousticness"`
	// Danceability     float32   `db:"danceability" json:"danceability"`
	// Energy           float32   `db:"energy" json:"energy"`
	// Instrumentalness float32   `db:"instrumentalness" json:"instrumentalness"`
	// Mode             int       `db:"mode" json:"mode"`
	// Speechiness      float32   `db:"speechiness" json:"speechiness"`
	// Valence          float32   `db:"valence" json:"valence"`
	var similarCount int
	var highestMatchSensitivity int
	var matchOnProfiles bool // return value
	var matchPercentage int  // return value
	var songIDMatchedOn int  // return value
	const aspectCount = 10
	const loudnessRange = float64(4.0)
	const tempoRange = float64(12.0)
	const acousticnessRange = float32(0.1)
	const danceabilityRange = float32(0.1)
	const energyRange = float32(0.1)
	const instrumentalnessRange = float32(0.1)
	const speechinessRange = float32(0.1)
	const valenceRange = float32(0.1)

	// lower sensitivity means more chances of matching (looser aspect count restrictions)
	// higher sensitivity means less chances of matching (stricter aspect count restrictions)
	const lowSensitivityRange = 8
	const mediumSensitivityRange = 5
	const highSensitivityRange = 2

	userALikedSongs, err := h.Models.LikedSongs.GetAllByOneUser(profileA.UserID)
	if err != nil {
		fmt.Println("Error getting liked songs from user ", profileA.UserID, ". ")
		return false, 0, 0
	}

	userBLikedSongs, err := h.Models.LikedSongs.GetAllByOneUser(profileB.UserID)
	if err != nil {
		fmt.Println("Error getting liked songs from user ", profileB.UserID, ". ")
		return false, 0, 0
	}

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

	// compare acousticness averages
	if (profileA.Acousticness <= profileB.Acousticness+acousticnessRange) &&
		(profileA.Acousticness >= profileB.Acousticness-acousticnessRange) {
		similarCount += 1
	}

	// compare danceability averages
	if (profileA.Danceability <= profileB.Danceability+danceabilityRange) &&
		(profileA.Danceability >= profileB.Danceability-danceabilityRange) {
		similarCount += 1
	}

	// compare energy averages
	if (profileA.Energy <= profileB.Energy+energyRange) &&
		(profileA.Energy >= profileB.Energy-energyRange) {
		similarCount += 1
	}

	// compare instrumentalness averages
	if (profileA.Instrumentalness <= profileB.Instrumentalness+instrumentalnessRange) &&
		(profileA.Instrumentalness >= profileB.Instrumentalness-instrumentalnessRange) {
		similarCount += 1
	}

	// compare mode (major/minor key) preferences
	if profileA.Mode == profileB.Mode {
		similarCount += 1
	}

	// compare speechiness averages
	if (profileA.Speechiness <= profileB.Speechiness+speechinessRange) &&
		(profileA.Speechiness >= profileB.Speechiness-speechinessRange) {
		similarCount += 1
	}

	// compare valence averages
	if (profileA.Valence <= profileB.Valence+valenceRange) &&
		(profileA.Valence >= profileB.Valence-valenceRange) {
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
		if similarCount >= aspectCount-lowSensitivityRange {
			matchOnProfiles = true
		} else {
			matchOnProfiles = false
		}
	case highestMatchSensitivity >= 34 && highestMatchSensitivity <= 67:
		if similarCount >= aspectCount-mediumSensitivityRange {
			matchOnProfiles = true
		} else {
			matchOnProfiles = false
		}
	case highestMatchSensitivity >= 68:
		if similarCount >= aspectCount-highSensitivityRange {
			matchOnProfiles = true
		} else {
			matchOnProfiles = false
		}
	}

	matchPercentageFloat := (float64(similarCount) / float64(aspectCount)) * 100.00
	matchPercentage = int(matchPercentageFloat)

	// If a match on musical profiles, find the song ID of the liked song of both users that is closest
	// to the aggregate of both user's musical preferences.
	if matchOnProfiles == true {
		// Get the averages of musical preferences of both users
		aggregateMusicProfile := data.UserMusicProfile{
			Loudness: (profileA.Loudness + profileB.Loudness) / 2,
			Tempo:    (profileA.Tempo + profileB.Tempo) / 2,
			TimeSig:  profileA.TimeSig, // Averaging timesigs does not make sense, so just pick one since
			// we found previously that both for each user of the match are either the same, or are similar.
			Acousticness:     (profileA.Acousticness + profileB.Acousticness) / 2,
			Danceability:     (profileA.Danceability + profileB.Danceability) / 2,
			Energy:           (profileA.Energy + profileB.Energy) / 2,
			Instrumentalness: (profileA.Instrumentalness + profileB.Instrumentalness) / 2,
			Mode:             profileA.Mode,
			Speechiness:      (profileA.Speechiness + profileB.Speechiness) / 2,
			Valence:          (profileA.Valence + profileB.Valence) / 2,
		}

		// put all songs in a slice
		var songs []data.Song
		for x := range userALikedSongs {
			song, err := h.Models.Songs.Get(userALikedSongs[x].SongID)
			if err != nil {
				fmt.Println("Error getting a song from userALikedSongs. Song ID: ", userALikedSongs[x].SongID, ".")
			}
			songs = append(songs, *song)
		}
		for x := range userBLikedSongs {
			song, err := h.Models.Songs.Get(userBLikedSongs[x].SongID)
			if err != nil {
				fmt.Println("Error getting a song from userBLikedSongs. Song ID: ", userBLikedSongs[x].SongID, ".")
			}
			songs = append(songs, *song)
		}

		// Initialize variables that will hold song IDs of songs that--given their designated musical aspect--
		// match closest to the aggregate music profile.
		var loudnessClosestSongID int
		var loudnessClosestSongID2 int
		// var loudnessClosestSongID3 int

		var tempoClosestSongID int
		var tempoClosestSongID2 int
		// var tempoClosestSongID3 int

		var acousticnessClosestSongID int
		var acousticnessClosestSongID2 int
		// var acousticnessClosestSongID3 int

		var danceabilityClosestSongID int
		var danceabilityClosestSongID2 int
		// var danceabilityClosestSongID3 int

		var energyClosestSongID int
		var energyClosestSongID2 int
		// var energyClosestSongID3 int

		var instrumentalnessClosestSongID int
		var instrumentalnessClosestSongID2 int
		// var instrumentalnessClosestSongID3 int

		var speechinessClosestSongID int
		var speechinessClosestSongID2 int
		// var speechinessClosestSongID3 int

		var valenceClosestSongID int
		var valenceClosestSongID2 int
		// var valenceClosestSongID3 int

		var timeSigsMatchingSongIDs []int
		var modesMatchingSongIDs []int
		var differences []float64

		// find the song that is closest to the aggregate music profile.

		// *** Commenting within Loudness should be considered for blocks of code after it that
		// are similar. ***

		// get songs with closest Loudness
		// Create a map (key-value pairs) to store--for each song iterated through--the difference between
		// the aggregate music profile's given song aspect and the currently-iterated song's given aspect.
		var diffMap = make(map[float64]int)
		// iterate through each liked-song of both users
		for i := range songs {
			// get the difference in the aggregate profile's preference and the actual song's level of the aspect
			// (Lower difference == higher chance the song will be considered a closest-song within the given aspect)
			difference := aggregateMusicProfile.Loudness - songs[i].LoudnessAvg
			if difference < 0 {
				difference *= -1
			}
			// add the difference as a key and the song ID as a value to the map
			diffMap[difference] = songs[i].ID
			// append just the difference to a slice (dynamic array) so we can sort them later to find the smallest
			differences = append(differences, difference)
		}
		// Sort the slice and find its length. Sorts by highest to lowest.
		sort.Float64s(differences)
		length := len(differences)
		// use the tail-end 2 values of the differences slice as keys to the map in order to find the song IDs
		// of songs that are the closest to the aggregate musical profile in the given musical aspect.
		loudnessClosestSongID = diffMap[differences[length-1]]
		loudnessClosestSongID2 = diffMap[differences[length-2]]
		// loudnessClosestSongID3 = diffMap[differences[length-3]]

		// get songs with closest Tempo
		// for k := range diffMap { delete(diffMap, k) }
		diffMap = make(map[float64]int)
		differences = nil
		for i := range songs {
			difference := aggregateMusicProfile.Tempo - songs[i].TempoAvg
			if difference < 0 {
				difference *= -1
			}
			diffMap[difference] = songs[i].ID
			differences = append(differences, difference)
		}
		sort.Float64s(differences)
		length = len(differences)
		tempoClosestSongID = diffMap[differences[length-1]]
		tempoClosestSongID2 = diffMap[differences[length-2]]
		// tempoClosestSongID3 = diffMap[differences[length-3]]

		// get songs with closest Acousticness
		diffMap = make(map[float64]int)
		differences = nil
		for i := range songs {
			difference := float64(aggregateMusicProfile.Acousticness - songs[i].Acousticness)
			if difference < 0 {
				difference *= -1
			}
			diffMap[difference] = songs[i].ID
			differences = append(differences, difference)
		}
		sort.Float64s(differences)
		length = len(differences)
		acousticnessClosestSongID = diffMap[differences[length-1]]
		acousticnessClosestSongID2 = diffMap[differences[length-2]]
		// acousticnessClosestSongID3 = diffMap[differences[length-3]]

		// get songs with closest Danceability
		diffMap = make(map[float64]int)
		differences = nil
		for i := range songs {
			difference := float64(aggregateMusicProfile.Danceability - songs[i].Danceability)
			if difference < 0 {
				difference *= -1
			}
			diffMap[difference] = songs[i].ID
			differences = append(differences, difference)
		}
		sort.Float64s(differences)
		length = len(differences)
		danceabilityClosestSongID = diffMap[differences[length-1]]
		danceabilityClosestSongID2 = diffMap[differences[length-2]]
		// danceabilityClosestSongID3 = diffMap[differences[length-3]]

		// get songs with closest Energy
		diffMap = make(map[float64]int)
		differences = nil
		for i := range songs {
			difference := float64(aggregateMusicProfile.Energy - songs[i].Energy)
			if difference < 0 {
				difference *= -1
			}
			diffMap[difference] = songs[i].ID
			differences = append(differences, difference)
		}
		sort.Float64s(differences)
		length = len(differences)
		energyClosestSongID = diffMap[differences[length-1]]
		energyClosestSongID2 = diffMap[differences[length-2]]
		// energyClosestSongID3 = diffMap[differences[length-3]]

		// get songs with closest Instrumentalness
		diffMap = make(map[float64]int)
		differences = nil
		for i := range songs {
			difference := float64(aggregateMusicProfile.Instrumentalness - songs[i].Instrumentalness)
			if difference < 0 {
				difference *= -1
			}
			diffMap[difference] = songs[i].ID
			differences = append(differences, difference)
		}
		sort.Float64s(differences)
		length = len(differences)
		instrumentalnessClosestSongID = diffMap[differences[length-1]]
		instrumentalnessClosestSongID2 = diffMap[differences[length-2]]
		// instrumentalnessClosestSongID3 = diffMap[differences[length-3]]

		// get songs with closest Speechiness
		diffMap = make(map[float64]int)
		differences = nil
		for i := range songs {
			difference := float64(aggregateMusicProfile.Speechiness - songs[i].Speechiness)
			if difference < 0 {
				difference *= -1
			}
			diffMap[difference] = songs[i].ID
			differences = append(differences, difference)
		}
		sort.Float64s(differences)
		length = len(differences)
		speechinessClosestSongID = diffMap[differences[length-1]]
		speechinessClosestSongID2 = diffMap[differences[length-2]]
		// speechinessClosestSongID3 = diffMap[differences[length-3]]

		// get songs with closest Valence
		diffMap = make(map[float64]int)
		differences = nil
		for i := range songs {
			difference := float64(aggregateMusicProfile.Valence - songs[i].Valence)
			if difference < 0 {
				difference *= -1
			}
			diffMap[difference] = songs[i].ID
			differences = append(differences, difference)
		}
		sort.Float64s(differences)
		length = len(differences)
		valenceClosestSongID = diffMap[differences[length-1]]
		valenceClosestSongID2 = diffMap[differences[length-2]]
		// valenceClosestSongID3 = diffMap[differences[length-3]]

		// get songs with identical TimeSigs and Modes
		for i := range songs {
			if aggregateMusicProfile.TimeSig == songs[i].TimeSigAvg {
				timeSigsMatchingSongIDs = append(timeSigsMatchingSongIDs, songs[i].ID)
			}
			if aggregateMusicProfile.Mode == songs[i].Mode {
				modesMatchingSongIDs = append(modesMatchingSongIDs, songs[i].ID)
			}
		}

		// Get the mode (highest occurring) of a slice of all of the IDs found to be the closest in each
		// aspect area (except timesig, just yet).
		IDs := []int{loudnessClosestSongID, loudnessClosestSongID2, // loudnessClosestSongID3,
			tempoClosestSongID, tempoClosestSongID2, // tempoClosestSongID3,
			acousticnessClosestSongID, acousticnessClosestSongID2, // acousticnessClosestSongID3,
			danceabilityClosestSongID, danceabilityClosestSongID2, // danceabilityClosestSongID3,
			energyClosestSongID, energyClosestSongID2, // energyClosestSongID3,
			instrumentalnessClosestSongID, instrumentalnessClosestSongID2, // instrumentalnessClosestSongID3,
			speechinessClosestSongID, speechinessClosestSongID2, // speechinessClosestSongID3,
			valenceClosestSongID, valenceClosestSongID2, // valenceClosestSongID3,
		}
		FinalIDs := IDs

		// Now, timesigs...
		// If there were any songs with matching timesigs to the aggregate music profile, AND
		// those songs were already in the IDs slice, add them to the new slice of IDs.
		if timeSigsMatchingSongIDs != nil {
			for t := range timeSigsMatchingSongIDs {
				idExists := false
				for x := range IDs {
					if IDs[x] == timeSigsMatchingSongIDs[t] {
						idExists = true
					}
				}
				if idExists == true {
					FinalIDs = append(FinalIDs, timeSigsMatchingSongIDs[t])
				}
			}
		}

		// Same thing for modes...
		if modesMatchingSongIDs != nil {
			for t := range modesMatchingSongIDs {
				idExists := false
				for x := range IDs {
					if IDs[x] == modesMatchingSongIDs[t] {
						idExists = true
					}
				}
				if idExists == true {
					FinalIDs = append(FinalIDs, modesMatchingSongIDs[t])
				}
			}
		}

		// Find the most common-occurring song ID in our slice that was made from the song IDs that are musically
		// closest to our music profile that was made from the aggregate of both users in the confirmed match.
		songIDMatchedOn = h.Mode(FinalIDs, aspectCount)

	} else {
		songIDMatchedOn = 0
	}

	return matchOnProfiles, matchPercentage, songIDMatchedOn
}

// Mode finds the highest-occurring integer within a slice (dynamic array)
func (h *Handlers) Mode(IDs []int, size int) int {
	if size == 0 {
		return 0
	}
	mp := make(map[int]int)
	for _, i := range IDs {
		mp[i]++
	}
	maxCount := 0
	mode := 0

	for value, count := range mp {
		if count > maxCount {
			maxCount = count
			mode = value
		}
	}

	return mode
}
