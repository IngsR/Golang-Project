package main

import (
	"goproject/config"
	"goproject/database"
	"goproject/handlers"
	"goproject/middleware"
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

	// Halaman publik
	r.GET("/", handlers.HomeHandler)
	r.GET("/login", handlers.ShowLogin)
	r.POST("/login", handlers.Login)
	r.GET("/register", handlers.ShowRegister)
	r.POST("/register", handlers.Register)
	r.GET("/logout", handlers.Logout)

	// Halaman artikel (publik: lihat, login required: buat/edit/hapus)
	r.GET("/articles", handlers.ListArticles)
	r.GET("/articles/:id", handlers.ShowArticle)

	// Grup route yang memerlukan autentikasi
	auth := r.Group("/")
	auth.Use(middleware.AuthRequired())
	{
		auth.GET("/articles/create", handlers.ShowCreateForm)
		auth.POST("/articles/create", handlers.CreateArticle)
		auth.GET("/articles/:id/edit", handlers.ShowEditForm)
		auth.POST("/articles/:id/edit", handlers.UpdateArticle)
		auth.POST("/articles/:id/delete", handlers.DeleteArticle)
	}

	// Jalankan server
	log.Printf("🚀 %s berjalan di http://localhost%s\n", cfg.AppName, cfg.Port)
	if err := r.Run(cfg.Port); err != nil {
		log.Fatal("Gagal menjalankan server:", err)
	}
}
