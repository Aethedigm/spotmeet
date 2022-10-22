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

	a.App.Routes.Route("/users", func(r chi.Router) {
		r.Get("/login", a.Handlers.UserLogin)
		r.Post("/login", a.Handlers.PostUserLogin)
		r.Get("/logout", a.Handlers.Logout)

		r.Post("/create", a.Handlers.CreateUser)
	})

	// a.App.Routes.Get("/get-user/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	id, _ := strconv.Atoi(chi.URLParam(r, "id"))

	// 	u, err := a.Models.Users.Get(id)
	// 	if err != nil {
	// 		a.App.ErrorLog.Println(err)
	// 		return
	// 	}

	// 	fmt.Fprintf(w, "%s %s %s", u.FirstName, u.LastName, u.Email)
	// })

	// a.App.Routes.Get("/update-user/{id}", func(w http.ResponseWriter, r *http.Request) {
	// 	id, _ := strconv.Atoi(chi.URLParam(r, "id"))
	// 	u, err := a.Models.Users.Get(id)
	// 	if err != nil {
	// 		a.App.ErrorLog.Println(err)
	// 		return
	// 	}

	// 	u.LastName = a.App.RandomString(10)
	// 	err = u.Update(*u)
	// 	if err != nil {
	// 		a.App.ErrorLog.Println(err)
	// 		return
	// 	}

	// 	fmt.Fprintf(w, "updated last name to %s", u.LastName)

	// })

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
