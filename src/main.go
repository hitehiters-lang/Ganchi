package main

import (
	"ganchi_app/additional"
	"ganchi_app/config"
	"ganchi_app/connection"
	_ "ganchi_app/docs"
	"ganchi_app/routers/blog"
	"ganchi_app/routers/contacts"
	"ganchi_app/routers/loaders"
	"ganchi_app/routers/manual_loaders"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"
)

// @title Gangchi API
// @version 0.2
// @description Документация к API.
// @BasePath /
func main() {
	additional.PrintMessage("Подключаемся к базе данных")
	db := connection.InitDatabase()
	if db == nil {
		os.Exit(1)
	}

	IPv4 := additional.GetLocalIP()

	r := chi.NewRouter()

	// CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/api/docs/*", httpSwagger.WrapHandler)

	r.Get("/api/loaders/getall", loaders.GetLoaders)
	r.Get("/api/loaders/getimage/{loader_id}", loaders.GetLoaderImage)

	r.Get("/api/manual_loaders/getall", manual_loaders.GetManualLoaders)
	r.Get("/api/manual_loaders/getimage/{loader_id}", manual_loaders.GetManualLoaderImage)

	r.Get("/api/blog/posts", blog.GetPosts)
	r.Get("/api/blog/posts/{slug}", blog.GetPostBySlug)
	r.Get("/api/blog/categories", blog.GetCategories)
	r.Get("/api/blog/posts/{slug}/cover", blog.GetBlogCoverImage)
	r.Get("/api/blog/posts/{slug}/comments", blog.GetComments)
	r.Post("/api/blog/posts/{slug}/comments", blog.CreateComment)

	r.Post("/api/contact/send_mail", contacts.ContactUs)

	additional.PrintSuccess("", "Стартуем сервер на "+IPv4)
	additional.PrintSuccess("/api/docs/", "Документация доступна по адресу")

	http.ListenAndServe(config.GetServerPort(), r)
}
