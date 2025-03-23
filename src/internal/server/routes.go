package server

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/jwtauth/v5"

	"sgublogsite/src/internal/controller"
)

func registerHandlers() http.Handler {
    r := chi.NewRouter()
    useMiddleware(r)
    registerRoutes(r)
    return r
}

func useMiddleware(r *chi.Mux) {
    r.Use(middleware.Logger)
    r.Use(middleware.Recoverer)
    r.Use(cors.Handler(cors.Options{
        AllowedOrigins:   []string{"https://*", "http://*"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        AllowCredentials: true,
        MaxAge:           300,
    }))

    r.Use(jwtauth.Verifier(controller.TokenAuth))
}

func registerRoutes(r *chi.Mux) http.Handler {
    r.Handle("/assets/*",  http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

    r.Get("/"               , controller.Home)
    r.Get("/about"          , controller.About)
    r.Get("/contact"        , controller.Contact)

    r.Get("/news"           , controller.News)
    r.Get("/announcements"  , controller.Announcements)
    r.Get("/categories"     , controller.Categories)
    r.Get("/category/{id}"  , controller.CategoryPosts)
    r.Get("/search"         , controller.Search)

    r.Get("/post/{id}"      , controller.Post)

    r.Get("/profile"        , controller.Profile)
    r.Get("/login"          , controller.LoginPage)

    r.Post("/register"      , controller.Register)
    r.Post("/login"         , controller.Login)
    r.Post("/logout"        , controller.Logout)

    return r
}
