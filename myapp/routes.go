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
	})

	a.App.Routes.Route("/users", func(r chi.Router) {
		r.Get("/login", a.Handlers.UserLogin)
		r.Get("/logout", a.Handlers.Logout)
		r.Get("/register", a.Handlers.UserRegister)
		r.Get("/profile", a.Handlers.Profile)
		r.Get("/profile/{profileID:[0-9]+}", a.Handlers.ProfileByID)

		r.Post("/login", a.Handlers.PostUserLogin)
		r.Post("/create", a.Handlers.CreateUserAndProfile)

	})

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
