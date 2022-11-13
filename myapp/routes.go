package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	_ "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func (a *application) routes() *chi.Mux {
	// middleware must come before any routes
	a.App.Routes.Use(render.SetContentType(render.ContentTypeJSON))
	a.App.Routes.Use(middleware.Timeout(60 * time.Second))

	// add routes here
	a.App.Routes.Get("/", a.Handlers.Home)

	a.App.Routes.Route("/matches", func(r chi.Router) {
		r.Get("/", a.Handlers.Matches)
		r.Get("/myresults", a.Handlers.MyMatchResults)
		r.Get("/accept/{matchID:[0-9]+}", a.Handlers.AcceptMatch)
		r.Get("/reject/{matchID:[0-9]+}", a.Handlers.RejectMatch)

		r.Post("/location", a.Handlers.Location)
	})

	a.App.Routes.Route("/users", func(r chi.Router) {
		r.Get("/login", a.Handlers.UserLogin)
		r.Get("/logout", a.Handlers.Logout)
		r.Get("/register", a.Handlers.UserRegister)
		r.Get("/spotauth", a.Handlers.SpotifyAuthorization)
		r.Get("/newspotaccesstoken", a.Handlers.NewAccessTokenRequest)
		r.Get("/profile", a.Handlers.Profile)
		r.Get("/profile/{profileID:[0-9]+}", a.Handlers.ProfileByID)
		r.Get("/edit-profile/{profileID:[0-9]+}", a.Handlers.EditProfile)

		r.Post("/login", a.Handlers.PostUserLogin)
		r.Post("/create", a.Handlers.CreateUserAndProfile)

		r.Put("/settings/{settingsID:[0-9]+}", a.Handlers.UpdateSettings)
		r.Put("/update-profile/{profileID:[0-9]+}", a.Handlers.UpdateProfile)
	})

	a.App.Routes.Get("/spotauth/callback", a.Handlers.SpotifyAuthorizationCallback)

	a.App.Routes.Route("/messages", func(r chi.Router) {
		r.Get("/", a.Handlers.Messages)
		r.Get("/getThreads/{userID:[0-9]+}", a.Handlers.GetThreads)
		r.Get("/getMessages/{fromUserID:[0-9]+}", a.Handlers.Thread)
		r.Get("/between/{userID:[0-9]+}/{matchID:[0-9]+}", a.Handlers.GetMessages)

		r.Post("/create", a.Handlers.CreateMessage)
	})

	a.App.Routes.Get("/settings", a.Handlers.Settings)

	a.App.Routes.Get("/about", a.Handlers.About)

	// static routes
	fileServer := http.FileServer(http.Dir("./public"))
	a.App.Routes.Handle("/public/*", http.StripPrefix("/public", fileServer))

	return a.App.Routes
}

type ErrResponse struct {
	Err            error  `json:"-"`
	HTTPStatusCode int    `json:"-"`
	StatusText     string `json:"status"`
	AppCode        int64  `json:"code,omitempty"`
	ErrorText      string `json:"error,omitempty"`
}
