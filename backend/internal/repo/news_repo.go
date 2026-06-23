package repo

import (
	"github.com/yayasan/cms/internal/domain"
	"gorm.io/gorm"
)

type NewsRepo struct{ db *gorm.DB }

func NewNewsRepo(db *gorm.DB) *NewsRepo { return &NewsRepo{db} }

func (r *NewsRepo) FindAll(opts domain.ListOptions) ([]domain.News, int64, error) {
	var list []domain.News
	var total int64

	q := r.db.Model(&domain.News{}).Preload("Author").Preload("Category")
	if opts.Status != "" {
		q = q.Where("status = ?", opts.Status)
	}
	if opts.Search != "" {
		q = q.Where("title ILIKE ?", "%"+opts.Search+"%")
	}
	q.Count(&total)

	offset := (opts.Page - 1) * opts.Limit
	err := q.Order("created_at DESC").Limit(opts.Limit).Offset(offset).Find(&list).Error
	return list, total, err
}

func (r *NewsRepo) FindPublished(limit int) ([]domain.News, error) {
	var list []domain.News
	err := r.db.Preload("Category").
		Where("status = ?", domain.StatusPublished).
		Order("published_at DESC").Limit(limit).Find(&list).Error
	return list, err
}

func (r *NewsRepo) FindBySlug(slug string) (*domain.News, error) {
	var n domain.News
	err := r.db.Preload("Author").Preload("Category").
		Where("slug = ? AND status = ?", slug, domain.StatusPublished).First(&n).Error
	return &n, err
}

func (r *NewsRepo) FindByID(id uint) (*domain.News, error) {
	var n domain.News
	err := r.db.Preload("Author").Preload("Category").First(&n, id).Error
	return &n, err
}

func (r *NewsRepo) Create(n *domain.News) error {
	return r.db.Create(n).Error
}

func (r *NewsRepo) Update(n *domain.News) error {
	return r.db.Save(n).Error
}

func (r *NewsRepo) Delete(id uint) error {
	return r.db.Delete(&domain.News{}, id).Error
}

func (r *NewsRepo) IncrementView(id uint) {
	r.db.Model(&domain.News{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + 1"))
}

func (r *NewsRepo) CountByStatus() map[string]int64 {
	type result struct {
		Status string
		Count  int64
	}
	var results []result
	r.db.Model(&domain.News{}).Select("status, count(*) as count").Group("status").Scan(&results)
	m := map[string]int64{}
	for _, r := range results {
		m[r.Status] = r.Count
	}
	return m
}
