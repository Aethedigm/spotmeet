package handlers

import (
	"encoding/json"
	"fmt"
	"myapp/data"
	"net/http"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/aethedigm/celeritas"
)

type Handlers struct {
	App    *celeritas.Celeritas
	Models data.Models
}

func (h *Handlers) MyMatchResults(w http.ResponseWriter, r *http.Request) {
	var matches []data.Match

	// Temporarily provide fake data
	m1 := data.Match{
		ID:           1,
		User_A_ID:    h.App.Session.GetInt(r.Context(), "user_id"),
		User_B_ID:    2,
		PercentMatch: 100,
		ArtistID:     1,
		CreatedAt:    time.Now(),
		Expires:      time.Now().Add(1 * time.Hour),
	}

	matches = append(matches, m1)

	m2 := data.Match{
		ID:           1,
		User_A_ID:    h.App.Session.GetInt(r.Context(), "user_id"),
		User_B_ID:    3,
		PercentMatch: 90,
		ArtistID:     3,
		CreatedAt:    time.Now(),
		Expires:      time.Now().Add(1 * time.Hour),
	}

	matches = append(matches, m2)

	m3 := data.Match{
		ID:           1,
		User_A_ID:    h.App.Session.GetInt(r.Context(), "user_id"),
		User_B_ID:    4,
		PercentMatch: 50,
		ArtistID:     2,
		CreatedAt:    time.Now(),
		Expires:      time.Now().Add(24 * 5 * time.Hour),
	}

	matches = append(matches, m3)

	js, err := json.Marshal(matches)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (h *Handlers) Matches(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "matches", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
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
