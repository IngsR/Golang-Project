package handlers

import (
	"goproject/database"
	"goproject/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// HomeHandler menampilkan halaman beranda dengan artikel terbaru
func HomeHandler(c *gin.Context) {
	var articles []models.Article
	database.DB.Preload("Author").Order("created_at DESC").Limit(6).Find(&articles)

	loggedIn, _ := c.Get("logged_in")
	userName, _ := c.Get("user_name")

	c.HTML(http.StatusOK, "home.html", gin.H{
		"title":     "Beranda",
		"articles":  articles,
		"logged_in": loggedIn,
		"user_name": userName,
	})
}
