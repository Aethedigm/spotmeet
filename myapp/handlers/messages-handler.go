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

func (h *Handlers) GetMessages(w http.ResponseWriter, r *http.Request) {
	userIDstr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		fmt.Println("Error converting userID to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	matchIDstr := chi.URLParam(r, "matchID")
	matchID, err := strconv.Atoi(matchIDstr)
	if err != nil {
		fmt.Println("Error converting matchID to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	messages, err := h.Models.Messages.GetAllForIDFromID(userID, matchID)
	if err != nil {
		fmt.Println("Error getting messages:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = json.NewEncoder(w).Encode(messages)
}

func (h *Handlers) GetThreads(w http.ResponseWriter, r *http.Request) {

	userIDstr := chi.URLParam(r, "userID")
	userID, err := strconv.Atoi(userIDstr)
	if err != nil {
		fmt.Println("Error converting userID to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	links, err := h.Models.Links.GetAllForOneUser(userID)
	if err != nil {
		fmt.Println("Error getting threads:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	threads := []data.Thread{}
	for _, links := range links {
		var user *data.User
		if links.User_A_ID == userID {
			user, err = h.Models.Users.Get(links.User_B_ID)
		} else {
			user, err = h.Models.Users.Get(links.User_A_ID)
		}

		if err != nil {
			fmt.Println("Error getting user:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		matchID := links.User_B_ID
		if matchID == userID {
			matchID = links.User_A_ID
		}

		if err != nil {
			fmt.Println("Error getting user:", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

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

		var nullTime = "0001-01-01 00:00:00 +0000 UTC"
		if timeSentISO.String() == nullTime {
			timeSentISO = links.CreatedAt
		}

		tmp := data.Thread{
			UserID:                user.ID,
			MatchID:               matchID,
			MatchFirstName:        user.FirstName,
			LatestMessagePreview:  latestMessagePreview,
			LatestMessageTimeSent: latestMessageTimeSent,
			OtherUsersImage:       otherUsersImage,
			TimeSentISO:           timeSentISO,
		}

		threads = append(threads, tmp)
	}

	sort.Slice(threads, func(p, q int) bool {
		return threads[p].TimeSentISO.After(threads[q].TimeSentISO)
	})

	err = json.NewEncoder(w).Encode(threads)
	if err != nil {
		fmt.Println("Error encoding threads:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func (h *Handlers) Thread(w http.ResponseWriter, r *http.Request) {
	if !h.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "/users/login", http.StatusSeeOther)
		return
	}

	matchIDstr := chi.URLParam(r, "fromUserID")
	matchID, err := strconv.Atoi(matchIDstr)
	if err != nil {
		fmt.Println("Error converting matchID to int:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	match, err := h.Models.Users.Get(matchID)
	if err != nil {
		fmt.Println("Error getting match:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	profile, err := h.Models.Profiles.GetByUserID(matchID)
	if err != nil {
		fmt.Println("Error getting profile ID", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	vars := make(jet.VarMap)
	vars.Set("userID", h.App.Session.GetInt(r.Context(), "userID"))
	vars.Set("matchID", matchID)
	vars.Set("matchFirstName", match.FirstName)
	vars.Set("matchProfileID", profile.ID)

	err = h.App.Render.JetPage(w, r, "message_thread", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

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

func (h *Handlers) Messages(w http.ResponseWriter, r *http.Request) {
	if !h.App.Session.Exists(r.Context(), "userID") {
		http.Redirect(w, r, "users/login", http.StatusSeeOther)
		return
	}

	vars := make(jet.VarMap)
	vars.Set("userID", h.App.Session.GetInt(r.Context(), "userID"))

	err := h.App.Render.JetPage(w, r, "messages", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println("error rendering:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}
