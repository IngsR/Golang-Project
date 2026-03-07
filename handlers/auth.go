package handlers

import (
	"goproject/middleware"
	"goproject/services"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler(service services.AuthService) *AuthHandler {
	return &AuthHandler{authService: service}
}

// isLoggedIn memeriksa apakah user sudah login dari context
func isLoggedIn(c *gin.Context) bool {
	loggedIn, exists := c.Get("logged_in")
	return exists && loggedIn == true
}

// ShowLogin menampilkan halaman login
func (h *AuthHandler) ShowLogin(c *gin.Context) {
	if isLoggedIn(c) {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":     "Login",
		"logged_in": false,
	})
}

// Login memproses form login
func (h *AuthHandler) Login(c *gin.Context) {
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")

	user, err := h.authService.Login(email, password)
	if err != nil {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title":     "Login",
			"error":     err.Error(),
			"email":     email,
			"logged_in": false,
		})
		return
	}

	// Simpan data user ke session
	session, err := middleware.Store.Get(c.Request, "session")
	if err != nil {
		log.Printf("Error mendapatkan session: %v", err)
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"title":     "Login",
			"error":     "Terjadi kesalahan pada server, silakan coba lagi",
			"email":     email,
			"logged_in": false,
		})
		return
	}

	session.Values["user_id"] = user.ID
	session.Values["user_name"] = user.Name
	session.Values["user_email"] = user.Email

	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Printf("Error menyimpan session: %v", err)
		c.HTML(http.StatusInternalServerError, "login.html", gin.H{
			"title":     "Login",
			"error":     "Gagal menyimpan sesi login, silakan coba lagi",
			"email":     email,
			"logged_in": false,
		})
		return
	}

	c.Redirect(http.StatusFound, "/")
}

// ShowRegister menampilkan halaman registrasi
func (h *AuthHandler) ShowRegister(c *gin.Context) {
	if isLoggedIn(c) {
		c.Redirect(http.StatusFound, "/")
		return
	}

	c.HTML(http.StatusOK, "register.html", gin.H{
		"title":     "Registrasi",
		"logged_in": false,
	})
}

// Register memproses form registrasi
func (h *AuthHandler) Register(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	user, err := h.authService.Register(name, email, password, confirmPassword)
	if err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title":     "Registrasi",
			"error":     err.Error(),
			"name":      name,
			"email":     email,
			"logged_in": false,
		})
		return
	}

	// Auto-login setelah register
	session, err := middleware.Store.Get(c.Request, "session")
	if err != nil {
		log.Printf("Error mendapatkan session: %v", err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	session.Values["user_id"] = user.ID
	session.Values["user_name"] = user.Name
	session.Values["user_email"] = user.Email

	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Printf("Error menyimpan session: %v", err)
		c.Redirect(http.StatusFound, "/login")
		return
	}

	c.Redirect(http.StatusFound, "/")
}

// Logout menghapus session user
func (h *AuthHandler) Logout(c *gin.Context) {
	session, err := middleware.Store.Get(c.Request, "session")
	if err != nil {
		c.Redirect(http.StatusFound, "/")
		return
	}

	session.Options.MaxAge = -1
	if err := session.Save(c.Request, c.Writer); err != nil {
		log.Printf("Error menghapus session: %v", err)
	}

	c.Redirect(http.StatusFound, "/")
}
