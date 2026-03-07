package services

import (
	"errors"
	"goproject/models"
	"goproject/repositories"
)

type ArticleService interface {
	CreateArticle(title, content string, authorID uint) (*models.Article, error)
	UpdateArticle(id string, title, content string, authorID uint) error
	DeleteArticle(id string, authorID uint) error
	GetArticleByID(id string) (*models.Article, error)
	GetAllArticlesPaginated(page, limit int) ([]models.Article, error)
	GetHomeArticles() ([]models.Article, error)
}

type articleService struct {
	articleRepo repositories.ArticleRepository
}

func NewArticleService(repo repositories.ArticleRepository) ArticleService {
	return &articleService{articleRepo: repo}
}

func (s *articleService) CreateArticle(title, content string, authorID uint) (*models.Article, error) {
	if title == "" || content == "" {
		return nil, errors.New("Judul dan konten harus diisi")
	}

	article := &models.Article{
		Title:    title,
		Content:  content,
		AuthorID: authorID,
	}

	err := s.articleRepo.Create(article)
	if err != nil {
		return nil, errors.New("Gagal menyimpan artikel")
	}

	return article, nil
}

func (s *articleService) UpdateArticle(id string, title, content string, authorID uint) error {
	if title == "" || content == "" {
		return errors.New("Judul dan konten harus diisi")
	}

	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return errors.New("Artikel tidak ditemukan")
	}

	if article.AuthorID != authorID {
		return errors.New("Anda tidak memiliki izin untuk mengedit artikel ini")
	}

	article.Title = title
	article.Content = content

	return s.articleRepo.Update(article)
}

func (s *articleService) DeleteArticle(id string, authorID uint) error {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return errors.New("Artikel tidak ditemukan")
	}

	if article.AuthorID != authorID {
		return errors.New("Anda tidak memiliki izin untuk menghapus artikel ini")
	}

	return s.articleRepo.Delete(article)
}

func (s *articleService) GetArticleByID(id string) (*models.Article, error) {
	article, err := s.articleRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("Artikel tidak ditemukan")
	}
	return article, nil
}

func (s *articleService) GetAllArticlesPaginated(page, limit int) ([]models.Article, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10 // default 10 per page
	}
	offset := (page - 1) * limit

	return s.articleRepo.GetPaginatedArticles(limit, offset)
}

func (s *articleService) GetHomeArticles() ([]models.Article, error) {
	return s.articleRepo.GetRecentArticles(6) // Home displays 6 articles 
}
