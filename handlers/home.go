package handlers

import (
	"goproject/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type HomeHandler struct {
	articleService services.ArticleService
}

func NewHomeHandler(service services.ArticleService) *HomeHandler {
	return &HomeHandler{articleService: service}
}

// ShowHome menampilkan halaman beranda dengan artikel terbaru
func (h *HomeHandler) ShowHome(c *gin.Context) {
	articles, err := h.articleService.GetHomeArticles()
	if err != nil {
		articles = nil // Fallback safe
	}

	loggedIn, _ := c.Get("logged_in")
	userName, _ := c.Get("user_name")

	c.HTML(http.StatusOK, "home.html", gin.H{
		"title":     "Beranda",
		"articles":  articles,
		"logged_in": loggedIn,
		"user_name": userName,
	})
}
