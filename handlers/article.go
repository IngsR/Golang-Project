package handlers

import (
	"goproject/database"
	"goproject/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// ListArticles menampilkan daftar semua artikel
func ListArticles(c *gin.Context) {
	var articles []models.Article
	database.DB.Preload("Author").Order("created_at DESC").Find(&articles)

	loggedIn, _ := c.Get("logged_in")
	userName, _ := c.Get("user_name")

	c.HTML(http.StatusOK, "articles.html", gin.H{
		"title":     "Daftar Artikel",
		"articles":  articles,
		"logged_in": loggedIn,
		"user_name": userName,
	})
}

// ShowArticle menampilkan detail satu artikel
func ShowArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article
	result := database.DB.Preload("Author").First(&article, id)
	if result.Error != nil {
		c.HTML(http.StatusNotFound, "home.html", gin.H{
			"title": "Tidak Ditemukan",
			"error": "Artikel tidak ditemukan",
		})
		return
	}

	loggedIn, _ := c.Get("logged_in")
	userName, _ := c.Get("user_name")
	userID, _ := c.Get("user_id")

	c.HTML(http.StatusOK, "article_detail.html", gin.H{
		"title":     article.Title,
		"article":   article,
		"logged_in": loggedIn,
		"user_name": userName,
		"user_id":   userID,
	})
}

// ShowCreateForm menampilkan form untuk membuat artikel baru
func ShowCreateForm(c *gin.Context) {
	loggedIn, _ := c.Get("logged_in")
	userName, _ := c.Get("user_name")

	c.HTML(http.StatusOK, "article_form.html", gin.H{
		"title":     "Buat Artikel Baru",
		"logged_in": loggedIn,
		"user_name": userName,
		"is_edit":   false,
	})
}

// CreateArticle memproses pembuatan artikel baru
func CreateArticle(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	userID, _ := c.Get("user_id")

	if title == "" || content == "" {
		loggedIn, _ := c.Get("logged_in")
		userName, _ := c.Get("user_name")
		c.HTML(http.StatusOK, "article_form.html", gin.H{
			"title":     "Buat Artikel Baru",
			"error":     "Judul dan konten harus diisi",
			"logged_in": loggedIn,
			"user_name": userName,
			"is_edit":   false,
		})
		return
	}

	article := models.Article{
		Title:    title,
		Content:  content,
		AuthorID: userID.(uint),
	}

	if err := database.DB.Create(&article).Error; err != nil {
		loggedIn, _ := c.Get("logged_in")
		userName, _ := c.Get("user_name")
		c.HTML(http.StatusOK, "article_form.html", gin.H{
			"title":     "Buat Artikel Baru",
			"error":     "Gagal menyimpan artikel",
			"logged_in": loggedIn,
			"user_name": userName,
			"is_edit":   false,
		})
		return
	}

	c.Redirect(http.StatusFound, "/articles/"+strconv.Itoa(int(article.ID)))
}

// ShowEditForm menampilkan form untuk mengedit artikel
func ShowEditForm(c *gin.Context) {
	id := c.Param("id")
	var article models.Article
	result := database.DB.First(&article, id)
	if result.Error != nil {
		c.Redirect(http.StatusFound, "/articles")
		return
	}

	// Pastikan hanya pemilik artikel yang bisa edit
	userID, _ := c.Get("user_id")
	if article.AuthorID != userID.(uint) {
		c.Redirect(http.StatusFound, "/articles")
		return
	}

	loggedIn, _ := c.Get("logged_in")
	userName, _ := c.Get("user_name")

	c.HTML(http.StatusOK, "article_form.html", gin.H{
		"title":     "Edit Artikel",
		"article":   article,
		"logged_in": loggedIn,
		"user_name": userName,
		"is_edit":   true,
	})
}

// UpdateArticle memproses pembaruan artikel
func UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article
	result := database.DB.First(&article, id)
	if result.Error != nil {
		c.Redirect(http.StatusFound, "/articles")
		return
	}

	// Pastikan hanya pemilik artikel yang bisa update
	userID, _ := c.Get("user_id")
	if article.AuthorID != userID.(uint) {
		c.Redirect(http.StatusFound, "/articles")
		return
	}

	title := c.PostForm("title")
	content := c.PostForm("content")

	if title == "" || content == "" {
		loggedIn, _ := c.Get("logged_in")
		userName, _ := c.Get("user_name")
		c.HTML(http.StatusOK, "article_form.html", gin.H{
			"title":     "Edit Artikel",
			"error":     "Judul dan konten harus diisi",
			"article":   article,
			"logged_in": loggedIn,
			"user_name": userName,
			"is_edit":   true,
		})
		return
	}

	database.DB.Model(&article).Updates(models.Article{
		Title:   title,
		Content: content,
	})

	c.Redirect(http.StatusFound, "/articles/"+id)
}

// DeleteArticle menghapus artikel
func DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	var article models.Article
	result := database.DB.First(&article, id)
	if result.Error != nil {
		c.Redirect(http.StatusFound, "/articles")
		return
	}

	// Pastikan hanya pemilik artikel yang bisa hapus
	userID, _ := c.Get("user_id")
	if article.AuthorID != userID.(uint) {
		c.Redirect(http.StatusFound, "/articles")
		return
	}

	database.DB.Delete(&article)
	c.Redirect(http.StatusFound, "/articles")
}
