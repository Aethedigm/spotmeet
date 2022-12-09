// Spotmeet - (Capstone Team E)
// 2022 Stephen Sumpter, John Neumeier,
// Zach Kods, Landon Wilson
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

// About displays the About page
func (h *Handlers) About(w http.ResponseWriter, r *http.Request) {
	// check if user ID exists in current browser session
	if !h.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "users/login", http.StatusSeeOther)
		return
	}

	// get userID
	userID := h.App.Session.GetInt(r.Context(), "userID")

	// get settings for the current user
	settings, err := h.Models.Settings.GetByUserID(userID)
	if err != nil {
		fmt.Println("Error getting settings:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the desired theme to pass into the view
	vars := make(jet.VarMap)
	vars.Set("theme", settings.Theme)

	err = h.App.Render.Page(w, r, "about", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
	}
}

// Location is an endpoint that gets the location from the browser and saves it to the db
// in relation to the current user.
func (h *Handlers) Location(w http.ResponseWriter, r *http.Request) {
	// get info from the page that is hitting this endpoint
	err := r.ParseForm()
	if err != nil {
		h.App.ErrorLog.Println("error parsing form:", err)
		return
	}

	// extract the userID from the page
	userIDstr := r.Form.Get("userID")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		h.App.ErrorLog.Println("error converting userID to int:", err)
		return
	}

	// get currentUser struct from the userID
	currentUser, err := h.Models.Users.Get(userID)
	if err != nil {
		h.App.ErrorLog.Println("error getting user object:", err)
		return
	}

	// get settings of current user
	currentUserSettings, err := h.Models.Settings.GetByUserID(userID)
	if err != nil {
		h.App.ErrorLog.Println("error getting settings object:", err)
		return
	}

	// get location data from the page
	lat := r.Form.Get("lat")
	lng := r.Form.Get("long")

	// convert location data to correct type for calculations
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

	// If new user, perform calculations on user location data that was just acquired
	// in order to provide settings parameters for matching based on location.
	if currentUser.Latitude == 0 && currentUser.Longitude == 0 {
		defaultDistance := 50
		latMod := float64(defaultDistance) * 0.01492753623
		longMod := float64(defaultDistance) * 0.018315018315

		currentUserSettings.LatMin = latFloat - latMod
		currentUserSettings.LatMax = latFloat + latMod
		currentUserSettings.LongMin = lngFloat - longMod
		currentUserSettings.LongMax = lngFloat + longMod

		err = currentUserSettings.Update(*currentUserSettings)
		if err != nil {
			h.App.ErrorLog.Println("error updating user ", userID, " settings", err)
			return
		}
	}

	// set the lat and long that was gathered from the browser, to the current user struct
	currentUser.Latitude = latFloat
	currentUser.Longitude = lngFloat

	// update the db entry for the user with the current user struct
	err = currentUser.Update(*currentUser)
	if err != nil {
		h.App.ErrorLog.Println("error updating user ", userID, err)
		return
	}

	// write success message back to browser
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%d", userID)))
}

// Home renders the home page
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

// CreateUserAndProfile ensures that a new user has a record created within the user and profile tables
func (h *Handlers) CreateUserAndProfile(w http.ResponseWriter, r *http.Request) {

	u := &data.User{}

	err := json.NewDecoder(r.Body).Decode(u)
	if err != nil {
		fmt.Println("Error decoding json:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	userID, err := u.Insert(*u)
	if err != nil {
		fmt.Println("Error inserting user:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Fprintf(w, "User created with id: %d", userID)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%d", userID)))

	p := &data.Profile{
		UserID:      userID,
		Description: "Hello, I'm new to SpotMeet! This is a default message.",
		ImageURL:    "/public/images/default-profile-pic.png",
	}

	profileID, err := p.Insert(*p)
	if err != nil {
		fmt.Println("Error inserting profile:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s := &data.Settings{
		UserID:   userID,
		Distance: 50,
		Theme:    "light",
	}

	_, err = s.Insert(*s)
	if err != nil {
		fmt.Println("Error inserting settings:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("%d", profileID)))
}
