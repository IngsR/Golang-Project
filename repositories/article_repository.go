package repositories

import (
	"goproject/database"
	"goproject/models"
)

type ArticleRepository interface {
	Create(article *models.Article) error
	Update(article *models.Article) error
	Delete(article *models.Article) error
	FindByID(id string) (*models.Article, error)
	GetPaginatedArticles(limit int, offset int) ([]models.Article, error)
	GetRecentArticles(limit int) ([]models.Article, error)
}

type articleRepository struct{}

func NewArticleRepository() ArticleRepository {
	return &articleRepository{}
}

func (r *articleRepository) Create(article *models.Article) error {
	return database.DB.Create(article).Error
}

func (r *articleRepository) Update(article *models.Article) error {
	return database.DB.Save(article).Error
}

func (r *articleRepository) Delete(article *models.Article) error {
	return database.DB.Delete(article).Error
}

func (r *articleRepository) FindByID(id string) (*models.Article, error) {
	var article models.Article
	err := database.DB.Preload("Author").First(&article, id).Error
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (r *articleRepository) GetPaginatedArticles(limit int, offset int) ([]models.Article, error) {
	var articles []models.Article
	// We use limit and offset for pagination so we don't load the entire database into memory.
	err := database.DB.Preload("Author").Order("created_at DESC").Limit(limit).Offset(offset).Find(&articles).Error
	return articles, err
}

func (r *articleRepository) GetRecentArticles(limit int) ([]models.Article, error) {
	var articles []models.Article
	err := database.DB.Preload("Author").Order("created_at DESC").Limit(limit).Find(&articles).Error
	return articles, err
}
