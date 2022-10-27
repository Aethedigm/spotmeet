package handlers

import (
	"encoding/json"
	"fmt"
	"myapp/data"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/CloudyKit/jet/v6"
	"github.com/aethedigm/celeritas"
)

type Handlers struct {
	App    *celeritas.Celeritas
	Models data.Models
}

func (h *Handlers) ProfileByID(w http.ResponseWriter, r *http.Request) {
	if !h.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "users/login", http.StatusSeeOther)
		return
	}

	if profileID := chi.URLParam(r, "profileID"); profileID != "" {

		pID, err := strconv.Atoi(profileID)
		if err != nil {
			fmt.Println("Error converting profileID to int:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		profile, err := h.Models.Profiles.Get(pID)
		if err != nil {
			fmt.Println("Error getting profile:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err := h.Models.Users.Get(profile.UserID)
		if err != nil {
			fmt.Println("Error getting user:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		vars := make(jet.VarMap)
		vars.Set("FirstName", user.FirstName)
		vars.Set("imgurl", profile.ImageURL)
		vars.Set("description", profile.Description)

		// GET TOP 3 ARTISTS
		vars.Set("Artist1", "Artist#1")
		vars.Set("Artist2", "Artist#2")
		vars.Set("Artist3", "Artist#3")

		err = h.App.Render.JetPage(w, r, "profile", vars, nil)
		if err != nil {
			h.App.ErrorLog.Println("error rendering:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
	} else {
		http.Error(w, "No ID provided", http.StatusBadRequest)
	}
}

func (h *Handlers) Profile(w http.ResponseWriter, r *http.Request) {
	if userID := h.App.Session.GetInt(r.Context(), "userID"); userID != 0 {
		profile, err := h.Models.Profiles.GetByUserID(userID)
		if err != nil {
			fmt.Println("Error getting profile:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		http.Redirect(w, r, "/users/profile/"+fmt.Sprint(profile.ID), http.StatusSeeOther)
	} else {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}
}

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
	w.Write(matchesJSON)
}

func (h *Handlers) Matches(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "matches", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	if h.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "/matches", http.StatusSeeOther)
		return
	}

	err := h.App.Render.Page(w, r, "home", nil, nil)

	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) GoPage(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.GoPage(w, r, "home", nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) JetPage(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.JetPage(w, r, "jet-template", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) SessionTest(w http.ResponseWriter, r *http.Request) {
	myData := "bar"

	h.App.Session.Put(r.Context(), "foo", myData)

	myValue := h.App.Session.GetString(r.Context(), "foo")

	vars := make(jet.VarMap)
	vars.Set("foo", myValue)

	err := h.App.Render.JetPage(w, r, "sessions", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) CreateUserAndProfile(w http.ResponseWriter, r *http.Request) {

	// create empty user struct
	u := &data.User{}

	// convert the body of the http request to json, and decode it into the empty user struct
	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		fmt.Println("Error decoding json:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// insert the new user information into the users table
	userID, err := u.Insert(*u)
	if err != nil {
		fmt.Println("Error inserting user:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Record the user-insert success message into the outgoing http response
	fmt.Fprintf(w, "User created with id: %d", userID)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%d", userID)))

	// create profile struct with the new user id, and default information
	p := &data.Profile{
		UserID:      userID,
		Description: "Hello, I'm new to SpotMeet! This is a default message.",
		ImageURL:    "/public/images/default-profile-pic.png",
	}

	// insert the new profile into the database
	profileID, err := p.Insert(*p)
	if err != nil {
		fmt.Println("Error inserting profile:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Record the profile-insert success message into the outgoing http response
	fmt.Fprintf(w, "Profile created with id: %d", profileID)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%d", profileID)))
}
