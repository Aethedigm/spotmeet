package handlers

import (
	"fmt"
	"myapp/data"
	"net/http"
	"net/smtp"
	"os"
	"strconv"

	"github.com/CloudyKit/jet/v6"

	"github.com/go-chi/chi/v5"
)

func (h *Handlers) SendRecoveryEmail(w http.ResponseWriter, r *http.Request) {
	from := os.Getenv("MAIL_FROM")
	pass := os.Getenv("MAIL_PASSWORD")
	to := chi.URLParam(r, "email")

	var recEmail data.RecoveryEmail

	usr, err := h.Models.Users.GetByEmail(to)
	if err != nil {
		fmt.Println("Error sending recovery email", err)
		w.Write([]byte(`{"status":"error"}`))
		return
	}

	recEmail.UserID = usr.ID
	recEmail.ID, err = recEmail.Insert(recEmail)
	if err != nil {
		fmt.Println("Error inserting recovery email object", err)
		w.Write([]byte(`{"status":"error"}`))
		return
	}

	body := `<!DOCTYPE html PUBLIC “-//W3C//DTD XHTML 1.0 Transitional//EN” “https://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd”>
	<html xmlns=“https://www.w3.org/1999/xhtml”>
	<head>
	<meta http–equiv=“Content-Type” content=“text/html; charset=UTF-8” />
	<meta http–equiv=“X-UA-Compatible” content=“IE=edge” />
	<meta name=“viewport” content=“width=device-width, initial-scale=1.0 “ />
	</head>
	<body><h1>Reset password here: <a href="` +
		os.Getenv("MAIL_RECOVERY_URL") + `/` + strconv.Itoa(recEmail.ID) +
		`">Reset Password</a></h1></body></html>`

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: Recovery Email\n" +
		"MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n" +
		body

	err = smtp.SendMail(os.Getenv("MAIL_HOST")+":587", smtp.PlainAuth("", from, pass, os.Getenv("MAIL_HOST")),
		from, []string{to}, []byte(msg))

	if err != nil {
		fmt.Println("Error sending recovery email", err)
		w.Write([]byte(`{"status":"error"}`))
		return
	}

	fmt.Println("Sent recovery email")
	w.Write([]byte(`{"status":"ok"}`))
}

func (h *Handlers) RecoveryEmailAccepted(w http.ResponseWriter, r *http.Request) {
	recoverIDstr := chi.URLParam(r, "recoverID")
	recoverID, err := strconv.Atoi(recoverIDstr)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	recEmail, err := h.Models.RecoveryEmails.Get(recoverID)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}

	vars := make(jet.VarMap)
	vars.Set("userID", recEmail.UserID)

	err = h.App.Render.Page(w, r, "reset-password", vars, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
		return
	}
}

func (h *Handlers) ForgotPassword(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "forgot-password", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}

func (h *Handlers) PasswordResetSent(w http.ResponseWriter, r *http.Request) {
	err := h.App.Render.Page(w, r, "reset-complete", nil, nil)
	if err != nil {
		h.App.ErrorLog.Println(err)
	}
}
