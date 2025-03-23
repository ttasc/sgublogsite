package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"
)

func (s *Server) registerHandlers() http.Handler {
    r := chi.NewRouter()
    s.useMiddleware(r)
    s.registerRoutes(r)
    return r
}

func (s *Server) useMiddleware(r *chi.Mux) {
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"https://*", "http://*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    r.Use(jwtauth.Verifier(s.ctrlr.TokenAuth))
}

func (s *Server) registerRoutes(r *chi.Mux) http.Handler {
    r.Handle("/assets/*",  http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

    r.Get("/"               , s.ctrlr.Home)
    r.Get("/about"          , s.ctrlr.About)
    r.Get("/contact"        , s.ctrlr.Contact)

    r.Get("/news"           , s.ctrlr.News)
    r.Get("/announcements"  , s.ctrlr.Announcements)
    r.Get("/categories"     , s.ctrlr.Categories)
    r.Get("/category/{id}"  , s.ctrlr.CategoryPosts)
    r.Get("/search"         , s.ctrlr.Search)

    r.Get("/post/{id}"      , s.ctrlr.Post)

    r.Get("/profile"        , s.ctrlr.Profile)
    r.Get("/login"          , s.ctrlr.LoginPage)

    r.Post("/register"      , s.ctrlr.Register)
    r.Post("/login"         , s.ctrlr.Login)
    r.Post("/logout"        , s.ctrlr.Logout)

    return r
}
