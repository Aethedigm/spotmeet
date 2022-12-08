package handlers

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"myapp/data"
	"net/http"
	"os"
	"strconv"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi/v5"
)

func (h *Handlers) UpdateUserPicture(w http.ResponseWriter, r *http.Request) {
	// Get profile ID
	profileIDstr := chi.URLParam(r, "profileID")

	// Check for directory + create if missing
	if _, err := os.Stat("public/images/u/" + profileIDstr); os.IsNotExist(err) {
		h.App.InfoLog.Println("Directory does not exist", err)
		if err := os.MkdirAll("public/images/u/"+profileIDstr, os.ModePerm); err != nil {
			h.App.ErrorLog.Println("Error creating directory", err)
			http.Error(w, "Error creating image", http.StatusInternalServerError)
			return
		}
	}

	// Create file separately
	_, err := os.Create("public/images/u/" + profileIDstr + "/pfp.jpg")
	if err != nil {
		h.App.ErrorLog.Println("Error creating output file")
		http.Error(w, "Error creating image", http.StatusInternalServerError)
		return
	}

	// Read image from stream
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		h.App.ErrorLog.Println("Error", err)
		http.Error(w, "Error creating image", http.StatusInternalServerError)
		return
	}

	// Write image to file
	err = ioutil.WriteFile("public/images/u/"+profileIDstr+"/pfp.jpg", data, fs.ModePerm)
	if err != nil {
		h.App.ErrorLog.Println("Error writing file", err)
		http.Error(w, "Error creating image", http.StatusInternalServerError)
		return
	}

	// Get ID
	profileID, err := strconv.Atoi(profileIDstr)
	if err != nil {
		h.App.ErrorLog.Println("Error converting id to int", err)
		http.Error(w, "Error creating image", http.StatusInternalServerError)
		return
	}

	// get profile
	profile, err := h.Models.Profiles.Get(profileID)
	if err != nil {
		h.App.ErrorLog.Println("Error getting profile", err)
		http.Error(w, "Error creating image", http.StatusInternalServerError)
		return
	}

	// Update DB url
	profile.ImageURL = "/public/images/u/" + profileIDstr + "/pfp.jpg"

	err = h.Models.Profiles.Update(*profile)
	if err != nil {
		h.App.ErrorLog.Println("Error updating profile", err)
		http.Error(w, "Error creating image", http.StatusInternalServerError)
		return
	}
}

// UpdateProfile updates the user's profile information with the new information given on the settings page
func (h *Handlers) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	desc := r.Form.Get("description")

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

		profile.Description = desc

		err = h.Models.Profiles.Update(*profile)
		if err != nil {
			fmt.Println("Error updating profile:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success"))
	}

	http.Error(w, "Error updating profile", http.StatusBadRequest)
}

// EditProfile renders the edit profile view.
func (h *Handlers) EditProfile(w http.ResponseWriter, r *http.Request) {
	// see if a userID session variable exists. If not, send the user back to the login page.
	if !h.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	// From the url, get the profile ID of the profile record to edit.
	if profileID := chi.URLParam(r, "profileID"); profileID != "" {
		pID, err := strconv.Atoi(profileID)
		if err != nil {
			fmt.Println("Error converting profileID to int:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// get the profile record data as a struct
		profile, err := h.Models.Profiles.Get(pID)
		if err != nil {
			fmt.Println("Error getting profile:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// If the user ID that's stored as a foreign key in the profile record does not match the user ID
		// of the current user, then prevent the user from being able to edit the given profile.
		if profile.UserID != h.App.Session.GetInt(r.Context(), "userID") {
			http.Error(w, "You are not authorized to edit this profile", http.StatusForbidden)
			return
		} else {
			// get the current user information as a struct
			user, err := h.Models.Users.Get(profile.UserID)
			if err != nil {
				fmt.Println("Error getting user:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// get settings of the user so we can pass the theme preference into the view
			settings, err := h.Models.Settings.GetByUserID(profile.UserID)
			if err != nil {
				h.App.ErrorLog.Println("error getting settings for user:", err)
			}

			// package data as json
			vars := make(jet.VarMap)
			vars.Set("userID", h.App.Session.GetInt(r.Context(), "userID"))
			vars.Set("profileID", profile.ID)
			vars.Set("description", profile.Description)
			vars.Set("FirstName", user.FirstName)
			vars.Set("imgurl", profile.ImageURL)
			vars.Set("theme", settings.Theme)

			// send json to the view and render the page
			err = h.App.Render.JetPage(w, r, "editprofile", vars, nil)
			if err != nil {
				h.App.ErrorLog.Println("error rendering:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
			}
		}
	} else {
		http.Redirect(w, r, "matches", http.StatusSeeOther)
	}
}

// ProfileByID renders the profile page. The profile shown is based on the profile ID that is passed as a
// url parameter from the calling endpoint.
func (h *Handlers) ProfileByID(w http.ResponseWriter, r *http.Request) {
	// Ensure that a user is signed in. If not, go to the login page.
	if !h.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	// get the profile ID from the url
	if profileID := chi.URLParam(r, "profileID"); profileID != "" {

		// convert the profile ID into a string
		pID, err := strconv.Atoi(profileID)
		if err != nil {
			fmt.Println("Error converting profileID to int:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// get the profile record data as a struct
		profile, err := h.Models.Profiles.Get(pID)
		if err != nil {
			fmt.Println("Error getting profile:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// get the user data of the user ID
		user, err := h.Models.Users.Get(profile.UserID)
		if err != nil {
			fmt.Println("Error getting user:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// save the ID of the user that is accessing this particular profile page
		thisUser := h.App.Session.GetInt(r.Context(), "userID")

		// get the liked-songs of the user that is being accessed
		likedSongs, err := h.Models.LikedSongs.GetAllByOneUser(profile.UserID)
		if err != nil {
			fmt.Println("Error getting liked songs from user:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// convert the likedSongs slice into a slice of correlating Song-structs
		var songs []data.Song
		var numberOfSongNames = 0
		for i := range likedSongs {
			song, err := h.Models.Songs.Get(likedSongs[i].SongID)
			if err != nil {
				fmt.Println("Error getting song:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			songs = append(songs, *song)
			numberOfSongNames += 1
			if i == 4 {
				break
			}
		}

		// formulate full song names with artist names by iterating through the songs slice
		var fullSongNames []string
		for i := 0; i < 5; i++ {
			if i < numberOfSongNames {
				fullSongName := songs[i].Name + "  by  " + songs[i].ArtistName
				fullSongNames = append(fullSongNames, fullSongName)
			} else {
				fullSongNames = append(fullSongNames, "")
			}
		}

		// get settings of the user so we can pass the theme preference into the view
		settings, err := h.Models.Settings.GetByUserID(thisUser)
		if err != nil {
			h.App.ErrorLog.Println("error getting settings for user:", err)
		}

		// package variables as json
		vars := make(jet.VarMap)
		vars.Set("userID", thisUser)
		vars.Set("profileID", profile.ID)
		vars.Set("usersProfileID", profile.UserID)
		vars.Set("FirstName", user.FirstName)
		vars.Set("imgurl", profile.ImageURL)
		vars.Set("description", profile.Description)
		vars.Set("theme", settings.Theme)
		vars.Set("numberOfSongNames", numberOfSongNames)
		for i := range fullSongNames {
			varName := "song" + strconv.Itoa(i+1)
			vars.Set(varName, fullSongNames[i])
		}

		// send json to the view and render the page
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
		http.Redirect(w, r, "login", http.StatusSeeOther)
		return
	}
}
