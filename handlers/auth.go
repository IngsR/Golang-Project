package handlers

import (
	"goproject/database"
	"goproject/middleware"
	"goproject/models"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// isLoggedIn memeriksa apakah user sudah login dari context
func isLoggedIn(c *gin.Context) bool {
	loggedIn, exists := c.Get("logged_in")
	return exists && loggedIn == true
}

// ShowLogin menampilkan halaman login
func ShowLogin(c *gin.Context) {
	// Redirect jika sudah login
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
func Login(c *gin.Context) {
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")

	// Validasi input kosong
	if email == "" || password == "" {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title":     "Login",
			"error":     "Email dan password harus diisi",
			"email":     email,
			"logged_in": false,
		})
		return
	}

	// Validasi panjang minimum password
	if len(password) < 6 {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title":     "Login",
			"error":     "Password minimal 6 karakter",
			"email":     email,
			"logged_in": false,
		})
		return
	}

	// Cari user berdasarkan email
	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil || !user.CheckPassword(password) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title":     "Login",
			"error":     "Email atau password salah",
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
func ShowRegister(c *gin.Context) {
	// Redirect jika sudah login
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
func Register(c *gin.Context) {
	name := strings.TrimSpace(c.PostForm("name"))
	email := strings.TrimSpace(c.PostForm("email"))
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	// Validasi input kosong
	if name == "" || email == "" || password == "" {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title":     "Registrasi",
			"error":     "Semua field harus diisi",
			"name":      name,
			"email":     email,
			"logged_in": false,
		})
		return
	}

	// Validasi panjang minimum password
	if len(password) < 6 {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title":     "Registrasi",
			"error":     "Password minimal 6 karakter",
			"name":      name,
			"email":     email,
			"logged_in": false,
		})
		return
	}

	if password != confirmPassword {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title":     "Registrasi",
			"error":     "Password dan konfirmasi password tidak cocok",
			"name":      name,
			"email":     email,
			"logged_in": false,
		})
		return
	}

	// Cek apakah email sudah terdaftar
	var existingUser models.User
	if database.DB.Where("email = ?", email).First(&existingUser).Error == nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title":     "Registrasi",
			"error":     "Email sudah terdaftar",
			"name":      name,
			"email":     email,
			"logged_in": false,
		})
		return
	}

	// Buat user baru
	user := models.User{
		Name:     name,
		Email:    email,
		Password: password,
	}

	if err := user.HashPassword(); err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title":     "Registrasi",
			"error":     "Gagal memproses registrasi",
			"logged_in": false,
		})
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title":     "Registrasi",
			"error":     "Gagal menyimpan data pengguna",
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
func Logout(c *gin.Context) {
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
