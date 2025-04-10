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
    r.Get("/tag/{id}"       , s.ctrlr.TagPosts)
    r.Get("/search"         , s.ctrlr.Search)

    r.Get("/post/{id}"      , s.ctrlr.Post)

    r.Get("/profile"        , s.ctrlr.Profile)
    r.Put("/profile"        , s.ctrlr.ProfileUpdate)
    r.Get("/login"          , s.ctrlr.LoginPage)

    r.Post("/login"         , s.ctrlr.Login)
    r.Post("/logout"        , s.ctrlr.Logout)


    r.Mount("/admin", adminRoutes(s))

    return r
}


func adminRoutes(s *Server) chi.Router {
    r := chi.NewRouter()

    r.Use(jwtauth.Authenticator(s.ctrlr.TokenAuth))

    r.Get("/",              s.ctrlr.AdminDashboard)

    r.Route("/users", func(r chi.Router) {
        r.Get("/",          s.ctrlr.AdminUsers)
        r.Get("/{id}",      s.ctrlr.AdminUser)
        r.Get("/new",       s.ctrlr.AdminUserNew)
        r.Post("/",         s.ctrlr.AdminUserCreate)
        r.Put("/{id}",      s.ctrlr.AdminUserUpdate)
        r.Delete("/{id}",   s.ctrlr.AdminUserDelete)
    })

    r.Route("/categories", func(r chi.Router) {
        r.Get("/",          s.ctrlr.AdminCategories)
        r.Post("/",         s.ctrlr.AdminCategoryCreate)
        r.Put("/",          s.ctrlr.AdminCategoryUpdate)
        r.Put("/{id}/move", s.ctrlr.AdminCategoryMove)
        r.Delete("/{id}",   s.ctrlr.AdminCategoryDelete)
    })

    r.Route("/tags", func(r chi.Router) {
        r.Get("/",          s.ctrlr.AdminTags)
        r.Post("/",         s.ctrlr.AdminTagCreate)
        r.Put("/",          s.ctrlr.AdminTagUpdate)
        r.Delete("/{id}",   s.ctrlr.AdminTagDelete)
        r.Delete("/bulk",   s.ctrlr.AdminTagBulkDelete)
    })

    r.Route("/posts", func(r chi.Router) {
        r.Get("/",          s.ctrlr.AdminPosts)
        // r.Get("/new",       s.ctrlr.AdminPostNew)
        // r.Post("/",         s.ctrlr.AdminPostCreate)
        // r.Put("/{id}",      s.ctrlr.AdminPostUpdate)
        r.Delete("/bulk",   s.ctrlr.AdminPostBulkDelete)
    })

    r.Route("/images", func(r chi.Router) {
        r.Get("/",          s.ctrlr.AdminImages)
        r.Post("/",         s.ctrlr.AdminImageUpload)
        r.Delete("/bulk",   s.ctrlr.AdminImageBulkDelete)
    })

    // r.Get("/info",      s.ctrlr.AdminInfo)

    return r
}
