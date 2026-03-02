package handlers

import (
	"goproject/database"
	"goproject/middleware"
	"goproject/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ShowLogin menampilkan halaman login
func ShowLogin(c *gin.Context) {
	loggedIn, _ := c.Get("logged_in")
	c.HTML(http.StatusOK, "login.html", gin.H{
		"title":     "Login",
		"logged_in": loggedIn,
	})
}

// Login memproses form login
func Login(c *gin.Context) {
	email := c.PostForm("email")
	password := c.PostForm("password")

	var user models.User
	result := database.DB.Where("email = ?", email).First(&user)
	if result.Error != nil || !user.CheckPassword(password) {
		c.HTML(http.StatusOK, "login.html", gin.H{
			"title": "Login",
			"error": "Email atau password salah",
		})
		return
	}

	// Simpan data user ke session
	session, _ := middleware.Store.Get(c.Request, "session")
	session.Values["user_id"] = user.ID
	session.Values["user_name"] = user.Name
	session.Values["user_email"] = user.Email
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusFound, "/")
}

// ShowRegister menampilkan halaman registrasi
func ShowRegister(c *gin.Context) {
	loggedIn, _ := c.Get("logged_in")
	c.HTML(http.StatusOK, "register.html", gin.H{
		"title":     "Registrasi",
		"logged_in": loggedIn,
	})
}

// Register memproses form registrasi
func Register(c *gin.Context) {
	name := c.PostForm("name")
	email := c.PostForm("email")
	password := c.PostForm("password")
	confirmPassword := c.PostForm("confirm_password")

	// Validasi input
	if name == "" || email == "" || password == "" {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "Registrasi",
			"error": "Semua field harus diisi",
		})
		return
	}

	if password != confirmPassword {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "Registrasi",
			"error": "Password dan konfirmasi password tidak cocok",
		})
		return
	}

	// Cek apakah email sudah terdaftar
	var existingUser models.User
	if database.DB.Where("email = ?", email).First(&existingUser).Error == nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "Registrasi",
			"error": "Email sudah terdaftar",
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
			"title": "Registrasi",
			"error": "Gagal memproses registrasi",
		})
		return
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.HTML(http.StatusOK, "register.html", gin.H{
			"title": "Registrasi",
			"error": "Gagal menyimpan data pengguna",
		})
		return
	}

	// Auto-login setelah register
	session, _ := middleware.Store.Get(c.Request, "session")
	session.Values["user_id"] = user.ID
	session.Values["user_name"] = user.Name
	session.Values["user_email"] = user.Email
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusFound, "/")
}

// Logout menghapus session user
func Logout(c *gin.Context) {
	session, _ := middleware.Store.Get(c.Request, "session")
	session.Options.MaxAge = -1
	session.Save(c.Request, c.Writer)

	c.Redirect(http.StatusFound, "/")
}
