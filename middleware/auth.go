package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/sessions"
)

// Store adalah session store global
var Store *sessions.CookieStore

// InitSession menginisialisasi session store dengan secret key
func InitSession(secretKey string) {
	Store = sessions.NewCookieStore([]byte(secretKey))
}

// AuthRequired adalah middleware yang memastikan user sudah login
func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := Store.Get(c.Request, "session")
		if err != nil || session.Values["user_id"] == nil {
			c.Redirect(http.StatusFound, "/login")
			c.Abort()
			return
		}
		// Simpan data user ke context agar bisa diakses di handler
		c.Set("user_id", session.Values["user_id"])
		c.Set("user_name", session.Values["user_name"])
		c.Set("user_email", session.Values["user_email"])
		c.Next()
	}
}

// SetUserContext adalah middleware yang menyimpan info user ke context (opsional, tidak redirect)
func SetUserContext() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := Store.Get(c.Request, "session")
		if err == nil && session.Values["user_id"] != nil {
			c.Set("user_id", session.Values["user_id"])
			c.Set("user_name", session.Values["user_name"])
			c.Set("user_email", session.Values["user_email"])
			c.Set("logged_in", true)
		} else {
			c.Set("logged_in", false)
		}
		c.Next()
	}
}
