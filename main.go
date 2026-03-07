package main

import (
	"goproject/config"
	"goproject/database"
	"goproject/handlers"
	"goproject/middleware"
	"goproject/repositories"
	"goproject/services"
	"html/template"
	"log"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func main() {
	// Muat konfigurasi
	cfg := config.LoadConfig()

	// Inisialisasi session store
	middleware.InitSession(cfg.SecretKey)

	// Inisialisasi database
	database.InitDB(cfg.DBPath)

	// Setup Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// Inisialisasi Repositories (Data Layer)
	userRepo := repositories.NewUserRepository()
	articleRepo := repositories.NewArticleRepository()

	// Inisialisasi Services (Business Logic Layer)
	authService := services.NewAuthService(userRepo)
	articleService := services.NewArticleService(articleRepo)

	// Inisialisasi Handlers (HTTP/Controller Layer)
	homeHandler := handlers.NewHomeHandler(articleService)
	authHandler := handlers.NewAuthHandler(authService)
	articleHandler := handlers.NewArticleHandler(articleService)

	// Custom template functions
	r.SetFuncMap(template.FuncMap{
		"truncate": func(s string, length int) string {
			if len(s) <= length {
				return s
			}
			return s[:length] + "..."
		},
		"formatDate": func(t time.Time) string {
			return t.Format("02 Jan 2006, 15:04")
		},
		"nl2br": func(s string) template.HTML {
			return template.HTML(strings.ReplaceAll(template.HTMLEscapeString(s), "\n", "<br>"))
		},
	})

	// Load templates
	r.LoadHTMLGlob("templates/**/*")

	// Serve static files
	r.Static("/static", "./static")

	// Middleware global
	r.Use(middleware.SetUserContext())

	// === ROUTES ===

	// Halaman publik (Guest/Public)
	r.GET("/", homeHandler.ShowHome)
	r.GET("/login", authHandler.ShowLogin)
	r.POST("/login", authHandler.Login)
	r.GET("/register", authHandler.ShowRegister)
	r.POST("/register", authHandler.Register)

	r.GET("/articles", articleHandler.ListArticles)
	r.GET("/articles/:id", articleHandler.ShowArticle)

	// Grup route yang memerlukan autentikasi (Wajib Login)
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		auth.GET("/logout", authHandler.Logout)

		auth.GET("/articles/create", articleHandler.ShowCreateForm)
		auth.POST("/articles/create", articleHandler.CreateArticle)
		auth.GET("/articles/:id/edit", articleHandler.ShowEditForm)
		auth.POST("/articles/:id/edit", articleHandler.UpdateArticle)
		auth.POST("/articles/:id/delete", articleHandler.DeleteArticle)
	}

	// Jalankan server
	log.Printf("🚀 %s berjalan di http://localhost%s\n", cfg.AppName, cfg.Port)
	if err := r.Run(cfg.Port); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
