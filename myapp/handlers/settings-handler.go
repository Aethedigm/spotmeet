// Spotmeet - (Capstone Team E)
// 2022 Stephen Sumpter, John Neumeier,
// Zach Kods, Landon Wilson
package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi/v5"
)

// UpdateSettings gets the information given by the user and updates it in the database
func (h *Handlers) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	// parse the information that was given by the user
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// save the information that was given by the user into variables
	lookingFor := r.Form.Get("lookingFor")
	distance := r.Form.Get("distance")
	sensitivity := r.Form.Get("sensitivity")
	theme := r.Form.Get("theme")

	// If a settings ID was given in the url, continue. If not, do not continue.
	if settingsID := chi.URLParam(r, "settingsID"); settingsID != "" {
		// convert settings ID to a string
		sID, err := strconv.Atoi(settingsID)
		if err != nil {
			fmt.Println("Error converting settingsID to int:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// get the before-changes settings record as a struct
		settings, err := h.Models.Settings.Get(sID)
		if err != nil {
			fmt.Println("Error getting settings:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// change LookingFor value in the settings struct
		settings.LookingFor = lookingFor

		// change distance value in the settings struct
		settings.Distance, err = strconv.Atoi(distance)
		if err != nil {
			fmt.Println("Error converting distance to int:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// change MatchSensitivity value in the settings struct
		settings.MatchSensitivity, err = strconv.Atoi(sensitivity)
		if err != nil {
			fmt.Println("Error converting sensitivity to int", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// change Theme value in the settings struct
		settings.Theme = theme

		// get the user for the settings we've been updating
		u, err := h.Models.Users.Get(settings.UserID)
		if err != nil {
			fmt.Println("Error getting user:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get u lat and long, save to settings with distance
		latMod := float64(settings.Distance) * 0.01492753623
		longMod := float64(settings.Distance) * 0.018315018315

		// change max/min values in settings struct, so matching on location data can be more accurate
		settings.LatMin = u.Latitude - latMod
		settings.LatMax = u.Latitude + latMod
		settings.LongMin = u.Longitude - longMod
		settings.LongMax = u.Longitude + longMod

		// pass the updated struct into Update so the settings for the user can be updated
		err = h.Models.Settings.Update(*settings)
	}
}

// Settings renders the settings page.
func (h *Handlers) Settings(w http.ResponseWriter, r *http.Request) {
	// ensure that a userID variable exists within the browser session
	if !h.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "users/login", http.StatusSeeOther)
		return
	}

	// get the user ID of the user
	userID := h.App.Session.GetInt(r.Context(), "userID")

	// get the profile struct of the user
	profile, err := h.Models.Profiles.GetByUserID(userID)
	if err != nil {
		fmt.Println("Error getting profile:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the settings struct of the user
	settings, err := h.Models.Settings.GetByUserID(userID)
	if err != nil {
		fmt.Println("Error getting settings:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// package data as json for the view
	vars := make(jet.VarMap)
	vars.Set("userID", h.App.Session.GetInt(r.Context(), "userID"))
	vars.Set("profileID", profile.ID)
	vars.Set("usersProfileID", profile.UserID)
	vars.Set("distance", settings.Distance)
	vars.Set("settingsID", settings.ID)
	vars.Set("lookingFor", settings.LookingFor)
	vars.Set("matchSensitivity", settings.MatchSensitivity)
	vars.Set("matchSensitivityString", settings.MatchSensitivityString())
	vars.Set("theme", settings.Theme)

	// send json to the view and render the page
	err = h.App.Render.JetPage(w, r, "settings", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}
