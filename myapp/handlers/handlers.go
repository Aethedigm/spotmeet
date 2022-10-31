package handlers

import (
	"encoding/json"
	"fmt"
	"myapp/data"
	"net/http"
	"strconv"

	"github.com/CloudyKit/jet/v6"
	"github.com/aethedigm/celeritas"
)

type Handlers struct {
	App    *celeritas.Celeritas
	Models data.Models
}

func (h *Handlers) Location(w http.ResponseWriter, r *http.Request) {
	if !h.App.Session.Exists(r.Context(), "userID") {
		h.App.ErrorLog.Println("error retrieving user_id from session")
		return
	}
	userID := h.App.Session.GetInt(r.Context(), "userID")

	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println("error parsing form:", err)
		return
	}

	// get the current User from the database as a struct
	currentUser, err := h.Models.Users.Get(userID)
	if err != nil {
		h.App.ErrorLog.Println("error getting user object:", err)
		return
	}

	lat := r.Form.Get("lat")
	lng := r.Form.Get("long")

	// convert the lat and long to float64
	latFloat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		h.App.ErrorLog.Println("error converting lat to float64:", err)
		return
	}

	lngFloat, err := strconv.ParseFloat(lng, 64)
	if err != nil {
		h.App.ErrorLog.Println("error converting long to float64:", err)
		return
	}

	currentUser.Latitude = latFloat
	currentUser.Longitude = lngFloat
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

	// Create settings
	s := &data.Settings{
		UserID:   userID,
		Distance: 50,
	}

	// insert the new settings into the database
	_, err = s.Insert(*s)
	if err != nil {
		fmt.Println("Error inserting settings:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Record the profile-insert success message into the outgoing http response
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%d", profileID)))
}
