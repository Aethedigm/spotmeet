package handlers

import (
	"fmt"
	"io/fs"
	"io/ioutil"
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

func (h *Handlers) EditProfile(w http.ResponseWriter, r *http.Request) {
	if !h.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
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

		if profile.UserID != h.App.Session.GetInt(r.Context(), "userID") {
			http.Error(w, "You are not authorized to edit this profile", http.StatusForbidden)
			return
		} else {

			user, err := h.Models.Users.Get(profile.UserID)
			if err != nil {
				fmt.Println("Error getting user:", err)
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			vars := make(jet.VarMap)
			vars.Set("userID", h.App.Session.GetInt(r.Context(), "userID"))
			vars.Set("profileID", profile.ID)
			vars.Set("description", profile.Description)
			vars.Set("FirstName", user.FirstName)

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

		thisUser := h.App.Session.GetInt(r.Context(), "userID")

		vars := make(jet.VarMap)
		vars.Set("userID", thisUser)
		vars.Set("profileID", profile.ID)
		vars.Set("usersProfileID", profile.UserID)
		vars.Set("FirstName", user.FirstName)
		vars.Set("imgurl", profile.ImageURL)
		vars.Set("description", profile.Description)

		// GET TOP 3 ARTISTS

		Artists := []string{"Artist #1", "Artist #2", "Artist #3"}

		lart, err := h.Models.LikedArtists.GetAllByOneUser(profile.UserID)
		if err != nil {
			fmt.Println("Error getting liked artists for user", profile.UserID, err)
		}

		maxArt := 3
		if len(lart) < 3 {
			maxArt = len(lart)
		}

		for i := 0; i < maxArt; i++ {
			tmpArt, err := h.Models.Artists.Get(lart[i].ArtistID)
			if err != nil {
				fmt.Println("Error getting artist from liked artists")
				break
			}
			Artists[i] = tmpArt.Name
		}

		vars.Set("Artist1", Artists[0])
		vars.Set("Artist2", Artists[1])
		vars.Set("Artist3", Artists[2])

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
