package handlers

import (
	"encoding/json"
	"fmt"
	"myapp/data"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/go-chi/chi/v5"
)

type msg struct {
	Message string `json:"message"`
	UserID  int    `json:"senderID"`
	MatchID int    `json:"receiverID"`
}

// GetMessages returns all sent messages between two users
func (h *Handlers) GetMessages(w http.ResponseWriter, r *http.Request) {
	// get the user ID from the url parameter
	userIDstr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		fmt.Println("Error converting userID to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the match IOD
	matchIDstr := chi.URLParam(r, "matchID")
	matchID, err := strconv.Atoi(matchIDstr)
	if err != nil {
		fmt.Println("Error converting matchID to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get a slice of all messages between two users
	messages, err := h.Models.Messages.GetAllForIDFromID(userID, matchID)
	if err != nil {
		fmt.Println("Error getting messages:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// send the messages as json to the browser
	err = json.NewEncoder(w).Encode(messages)
}

// GetThreads returns all message thread previews to the user, so that
// the user can select which of their linked matches they'd like to interact with.
func (h *Handlers) GetThreads(w http.ResponseWriter, r *http.Request) {
	// get the user ID from the url
	userIDstr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		fmt.Println("Error converting userID to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get the user's links (matches that were accepted by either person)
	links, err := h.Models.Links.GetAllForOneUser(userID)
	if err != nil {
		fmt.Println("Error getting threads:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// create a slice (dynamic array) that will hold message threads
	threads := []data.Thread{}
	for _, links := range links {
		var user *data.User
		var match *data.Match
		var userHasOpenedThread bool

		// get the match record from the match table that is linked by ID to the currently-iterated link
		match, err = h.Models.Matches.GetByBothUsers(links.User_A_ID, links.User_B_ID)
		if err != nil {
			fmt.Println("Error getting match:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// mark the thread as "viewed" if the match record has a true value for the current user
		if links.User_A_ID == userID {
			user, err = h.Models.Users.Get(links.User_B_ID)
			userHasOpenedThread = match.UserAViewedThread
		} else {
			user, err = h.Models.Users.Get(links.User_A_ID)
			userHasOpenedThread = match.UserBViewedThread
		}
		if err != nil {
			fmt.Println("Error getting user:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// get the ID of the other user (from the perspective of the current user)
		matchID := links.User_B_ID
		if matchID == userID {
			matchID = links.User_A_ID
		}

		// get display information for the thread previews from SQL query, and format it to show correctly
		latestMessagePreview,
			latestMessageTimeSent,
			otherUsersImage,
			timeSentISO,
			err := h.Models.RQ.ThreadPreviewQuery(userID, matchID)
		if err != nil {
			fmt.Println("Error in func ThreadPreviewQuery(), called in messages-handler.go.", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// If a thread has no messages sent on it, the SQL query above will return the latest-message-time-sent
		// as a zeroed-out ISO string. If this is the case, then set the latest-message-sent as the match-link's
		// creation time, so that the thread preview can still be organized on the the screen chronologically.
		var nullTime = "0001-01-01 00:00:00 +0000 UTC"
		if timeSentISO.String() == nullTime {
			timeSentISO = links.CreatedAt
		}

		// package the data we found into a thread struct
		tmp := data.Thread{
			UserID:                user.ID,
			MatchID:               matchID,
			MatchFirstName:        user.FirstName,
			LatestMessagePreview:  latestMessagePreview,
			LatestMessageTimeSent: latestMessageTimeSent,
			OtherUsersImage:       otherUsersImage,
			TimeSentISO:           timeSentISO,
			UserHasOpenedThread:   userHasOpenedThread,
		}

		// append that struct to a slice (a dynamic array)
		threads = append(threads, tmp)
	}

	// sort the slice of thread previews chronologically
	sort.Slice(threads, func(p, q int) bool {
		return threads[p].TimeSentISO.After(threads[q].TimeSentISO)
	})

	// send the thread previews to the view
	err = json.NewEncoder(w).Encode(threads)
	if err != nil {
		fmt.Println("Error encoding threads:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

// Thread displays the thread page. Actual messages will be called for display by JavaScript within the page
// hitting an endpoint calling GetMessages(), which is defined in this .go file.
func (h *Handlers) Thread(w http.ResponseWriter, r *http.Request) {
	// Check if user session exists. If not, then redirect back to login.
	if !h.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	// get the user ID from the session
	userID := h.App.Session.GetInt(r.Context(), "userID")

	// get the user ID of the other user
	otherUserIDstr := chi.URLParam(r, "fromUserID")
	otherUserID, err := strconv.Atoi(otherUserIDstr)
	if err != nil {
		fmt.Println("Error converting matchID to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get a data struct of the other user
	otherUser, err := h.Models.Users.Get(otherUserID)
	if err != nil {
		fmt.Println("Error getting match:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get a profile struct of the other user
	profile, err := h.Models.Profiles.GetByUserID(otherUserID)
	if err != nil {
		fmt.Println("Error getting profile ID", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// get settings of the user so we can pass the theme preference into the view
	settings, err := h.Models.Settings.GetByUserID(userID)
	if err != nil {
		h.App.ErrorLog.Println("error getting settings for user:", err)
	}

	// package data as json to send to the view
	vars := make(jet.VarMap)
	vars.Set("userID", userID)
	vars.Set("matchID", otherUserID)
	vars.Set("matchFirstName", otherUser.FirstName)
	vars.Set("matchProfileID", profile.ID)
	vars.Set("theme", settings.Theme)

	// save in match record that current user is viewing the thread
	match, err := h.Models.Matches.GetByBothUsers(userID, otherUserID)
	if match.User_A_ID == userID {
		match.UserAViewedThread = true
	} else if match.User_B_ID == userID {
		match.UserBViewedThread = true
	}
	err = h.Models.Matches.Update(*match)
	if err != nil {
		h.App.ErrorLog.Println("error updating matches table:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}

	// pass the json into the view and render the page
	err = h.App.Render.JetPage(w, r, "message_thread", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// CreateMessage is called by an endpoint that is hit when a user clicks 'send' in the threads page.
// The text that is entered in the message box is saved as a message in the database.
func (h *Handlers) CreateMessage(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1000000000)
	if err != nil {
		fmt.Println("Error parsing form:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	message := r.PostFormValue("message")
	senderIDstr := r.PostFormValue("senderID")
	receiverIDstr := r.PostFormValue("receiverID")

	if message == "" || senderIDstr == "" || receiverIDstr == "" {
		fmt.Println("Error: missing form value")
		http.Error(w, "missing form value", http.StatusBadRequest)
		return
	}

	senderID, err := strconv.Atoi(senderIDstr)
	if err != nil {
		fmt.Println("Error converting senderID to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	receiverID, err := strconv.Atoi(receiverIDstr)
	if err != nil {
		fmt.Println("Error converting receiverID to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	fmt.Println("Message:", message)
	fmt.Println("SenderID:", senderID)
	fmt.Println("ReceiverID:", receiverID)

	newMessage := data.Message{
		Content:   message,
		UserID:    senderID,
		MatchID:   receiverID,
		CreatedAt: time.Now(),
	}

	mID, err := h.Models.Messages.Insert(newMessage)
	if err != nil {
		fmt.Println("Error creating message:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	newMessage.ID = mID
}

// Messages renders the messages page. Thread previews on this page, however, are displayed via the page's
// JavaScript calling an endpoint that in turn, calls GetThreads(), which is also defined in this .go file.
func (h *Handlers) Messages(w http.ResponseWriter, r *http.Request) {
	// Check if user session exists. If not, then redirect back to login.
	if !h.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "users/login", http.StatusSeeOther)
		return
	}

	// get the current user's user ID
	userID := h.App.Session.GetInt(r.Context(), "userID")

	// get settings of the user so we can pass the theme preference into the view
	settings, err := h.Models.Settings.GetByUserID(userID)
	if err != nil {
		h.App.ErrorLog.Println("error getting settings for user:", err)
	}

	// package data as json to send to the view
	vars := make(jet.VarMap)
	vars.Set("userID", h.App.Session.GetInt(r.Context(), "userID"))
	vars.Set("theme", settings.Theme)

	// send data to the view and render the page
	err = h.App.Render.JetPage(w, r, "messages", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
