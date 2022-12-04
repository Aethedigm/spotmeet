package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) UpdateSettings(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	lookingFor := r.Form.Get("lookingFor")
	distance := r.Form.Get("distance")

	if settingsID := chi.URLParam(r, "settingsID"); settingsID != "" {
		sID, err := strconv.Atoi(settingsID)
		if err != nil {
			fmt.Println("Error converting settingsID to int:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		settings, err := h.Models.Settings.Get(sID)
		if err != nil {
			fmt.Println("Error getting settings:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		settings.LookingFor = lookingFor
		settings.Distance, err = strconv.Atoi(distance)
		if err != nil {
			fmt.Println("Error converting distance to int:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		u, err := h.Models.Users.Get(settings.UserID)
		if err != nil {
			fmt.Println("Error getting user:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Get u lat and long, save to settings with distance
		latMod := float64(settings.Distance) * 0.01492753623
		longMod := float64(settings.Distance) * 0.018315018315

		settings.LatMin = u.Latitude - latMod
		settings.LatMax = u.Latitude + latMod
		settings.LongMin = u.Longitude - longMod
		settings.LongMax = u.Longitude + longMod

		err = h.Models.Settings.Update(*settings)
	}

}

func (h *Handlers) Settings(w http.ResponseWriter, r *http.Request) {
	if !h.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "users/login", http.StatusSeeOther)
		return
	}

	userID := h.App.Session.GetInt(r.Context(), "userID")

	profile, err := h.Models.Profiles.GetByUserID(userID)
	if err != nil {
		fmt.Println("Error getting profile:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	settings, err := h.Models.Settings.GetByUserID(userID)
	if err != nil {
		fmt.Println("Error getting settings:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := make(jet.VarMap)
	vars.Set("userID", h.App.Session.GetInt(r.Context(), "userID"))
	vars.Set("profileID", profile.ID)
	vars.Set("usersProfileID", profile.UserID)
	vars.Set("distance", settings.Distance)
	vars.Set("settingsID", settings.ID)
	vars.Set("lookingFor", settings.LookingFor)
	vars.Set("matchSensitivity", settings.MatchSensitivityString())

	err = h.App.Render.JetPage(w, r, "settings", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

}
