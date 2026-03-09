package handlers

import (
	"fmt"
	"goproject/services"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	articleService services.ArticleService
}

func NewArticleHandler(service services.ArticleService) *ArticleHandler {
	return &ArticleHandler{articleService: service}
}

// ListArticles menampilkan daftar semua artikel dengan pagination
func (h *ArticleHandler) ListArticles(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit := 10 // Pagination limit
	articles, err := h.articleService.GetAllArticlesPaginated(page, limit)
	if err != nil {
		articles = nil
	}

	// Logic view aman mengatasi blank screen
	loggedIn, _ := c.Get("logged_in")
	userName, _ := c.Get("user_name")

	hasPrev := page > 1
	hasNext := len(articles) == limit

	c.HTML(http.StatusOK, "articles.html", gin.H{
		"title":       "Daftar Artikel",
		"articles":    articles,
		"logged_in":   loggedIn,
		"user_name":   userName,
		"page":        page,
		"has_prev":    hasPrev,
		"has_next":    hasNext,
		"prev_page":   page - 1,
		"next_page":   page + 1,
	})
}

// ShowArticle menampilkan detail satu artikel
func (h *ArticleHandler) ShowArticle(c *gin.Context) {
	id := c.Param("id")
	
	article, err := h.articleService.GetArticleByID(id)
	if err != nil {
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
func (h *ArticleHandler) ShowCreateForm(c *gin.Context) {
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
func (h *ArticleHandler) CreateArticle(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	userID, _ := c.Get("user_id")

	coverImage := ""
	file, err := c.FormFile("cover_image")
	if err == nil {
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		uploadPath := "./static/uploads/" + filename
		os.MkdirAll("./static/uploads", os.ModePerm)
		if errSave := c.SaveUploadedFile(file, uploadPath); errSave == nil {
			coverImage = "/static/uploads/" + filename
		}
	}

	article, err := h.articleService.CreateArticle(title, content, coverImage, userID.(uint))
	
	if err != nil {
		loggedIn, _ := c.Get("logged_in")
		userName, _ := c.Get("user_name")
		c.HTML(http.StatusOK, "article_form.html", gin.H{
			"title":     "Buat Artikel Baru",
			"error":     err.Error(),
			"logged_in": loggedIn,
			"user_name": userName,
			"is_edit":   false,
		})
		return
	}

	c.Redirect(http.StatusFound, "/articles/"+strconv.Itoa(int(article.ID)))
}

// ShowEditForm menampilkan form untuk mengedit artikel
func (h *ArticleHandler) ShowEditForm(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")
	
	article, err := h.articleService.GetArticleByID(id)
	if err != nil || article.AuthorID != userID.(uint) {
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
func (h *ArticleHandler) UpdateArticle(c *gin.Context) {
	id := c.Param("id")
	title := c.PostForm("title")
	content := c.PostForm("content")
	userID, _ := c.Get("user_id")

	coverImage := ""
	file, err := c.FormFile("cover_image")
	if err == nil {
		filename := fmt.Sprintf("%d_%s", time.Now().Unix(), filepath.Base(file.Filename))
		uploadPath := "./static/uploads/" + filename
		os.MkdirAll("./static/uploads", os.ModePerm)
		if errSave := c.SaveUploadedFile(file, uploadPath); errSave == nil {
			coverImage = "/static/uploads/" + filename
		}
	}

	err = h.articleService.UpdateArticle(id, title, content, coverImage, userID.(uint))
	
	if err != nil {
		loggedIn, _ := c.Get("logged_in")
		userName, _ := c.Get("user_name")
		
		// Re-fetch article for the form if validation failed (not not-found)
		article, _ := h.articleService.GetArticleByID(id)
	
		c.HTML(http.StatusOK, "article_form.html", gin.H{
			"title":     "Edit Artikel",
			"error":     err.Error(),
			"article":   article,
			"logged_in": loggedIn,
			"user_name": userName,
			"is_edit":   true,
		})
		return
	}

	c.Redirect(http.StatusFound, "/articles/"+id)
}

// DeleteArticle menghapus artikel
func (h *ArticleHandler) DeleteArticle(c *gin.Context) {
	id := c.Param("id")
	userID, _ := c.Get("user_id")

	err := h.articleService.DeleteArticle(id, userID.(uint))
	if err != nil {
		// Just redirect on error for now to keep it simple, or show alert
		c.Redirect(http.StatusFound, "/articles")
		return
	}

	c.Redirect(http.StatusFound, "/articles")
}
