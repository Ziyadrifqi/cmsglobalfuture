package repo

import (
	"github.com/yayasan/cms/internal/domain"
	"gorm.io/gorm"
)

// ── Page ──────────────────────────────────────────────────────────────────────

type PageRepo struct{ db *gorm.DB }

func NewPageRepo(db *gorm.DB) *PageRepo { return &PageRepo{db} }

func (r *PageRepo) FindBySlug(slug string) (*domain.Page, error) {
	var p domain.Page
	err := r.db.Where("slug = ?", slug).First(&p).Error
	return &p, err
}

func (r *PageRepo) FindAll() ([]domain.Page, error) {
	var pages []domain.Page
	err := r.db.Find(&pages).Error
	return pages, err
}

func (r *PageRepo) FindByID(id uint) (*domain.Page, error) {
	var p domain.Page
	err := r.db.First(&p, id).Error
	return &p, err
}

func (r *PageRepo) Update(p *domain.Page) error {
	return r.db.Save(p).Error
}

// ── Banner ────────────────────────────────────────────────────────────────────

type BannerRepo struct{ db *gorm.DB }

func NewBannerRepo(db *gorm.DB) *BannerRepo { return &BannerRepo{db} }

func (r *BannerRepo) FindActive() ([]domain.Banner, error) {
	var list []domain.Banner
	err := r.db.Where("is_active = true").Order("order_num ASC").Find(&list).Error
	return list, err
}

func (r *BannerRepo) FindAll() ([]domain.Banner, error) {
	var list []domain.Banner
	err := r.db.Order("order_num ASC").Find(&list).Error
	return list, err
}

func (r *BannerRepo) FindByID(id uint) (*domain.Banner, error) {
	var b domain.Banner
	err := r.db.First(&b, id).Error
	return &b, err
}

func (r *BannerRepo) Create(b *domain.Banner) error { return r.db.Create(b).Error }
func (r *BannerRepo) Update(b *domain.Banner) error { return r.db.Save(b).Error }
func (r *BannerRepo) Delete(id uint) error          { return r.db.Delete(&domain.Banner{}, id).Error }
