package repo

import (
	"github.com/yayasan/cms/internal/domain"
	"gorm.io/gorm"
)

// ── GalleryCategoryRepo ───────────────────────────────────────────────────────

type GalleryCategoryRepo struct{ db *gorm.DB }

func NewGalleryCategoryRepo(db *gorm.DB) *GalleryCategoryRepo {
	return &GalleryCategoryRepo{db}
}

func (r *GalleryCategoryRepo) FindAll() ([]domain.GalleryCategory, error) {
	var list []domain.GalleryCategory
	err := r.db.Order("order_num ASC, name ASC").Find(&list).Error
	return list, err
}

func (r *GalleryCategoryRepo) FindByID(id uint) (*domain.GalleryCategory, error) {
	var c domain.GalleryCategory
	err := r.db.First(&c, id).Error
	return &c, err
}

func (r *GalleryCategoryRepo) FindBySlug(slug string) (*domain.GalleryCategory, error) {
	var c domain.GalleryCategory
	err := r.db.Where("slug = ?", slug).First(&c).Error
	return &c, err
}

func (r *GalleryCategoryRepo) Create(c *domain.GalleryCategory) error { return r.db.Create(c).Error }
func (r *GalleryCategoryRepo) Update(c *domain.GalleryCategory) error { return r.db.Save(c).Error }
func (r *GalleryCategoryRepo) Delete(id uint) error {
	return r.db.Delete(&domain.GalleryCategory{}, id).Error
}

// ── GalleryRepo ───────────────────────────────────────────────────────────────

type GalleryRepo struct{ db *gorm.DB }

func NewGalleryRepo(db *gorm.DB) *GalleryRepo { return &GalleryRepo{db} }

func (r *GalleryRepo) FindAll(opts domain.ListOptions) ([]domain.GalleryItem, int64, error) {
	var list []domain.GalleryItem
	var total int64
	q := r.db.Model(&domain.GalleryItem{}).Preload("Category")
	if opts.Category != "" {
		q = q.Joins("JOIN gallery_categories gc ON gc.id = gallery_items.category_id").
			Where("gc.slug = ?", opts.Category)
	}
	q.Count(&total)
	offset := (opts.Page - 1) * opts.Limit
	err := q.Order("gallery_items.order_num ASC, gallery_items.created_at DESC").
		Limit(opts.Limit).Offset(offset).Find(&list).Error
	return list, total, err
}

func (r *GalleryRepo) FindActive(opts domain.ListOptions) ([]domain.GalleryItem, int64, error) {
	var list []domain.GalleryItem
	var total int64
	q := r.db.Model(&domain.GalleryItem{}).Preload("Category").Where("is_active = true")
	if opts.Category != "" {
		q = q.Joins("JOIN gallery_categories gc ON gc.id = gallery_items.category_id").
			Where("gc.slug = ?", opts.Category)
	}
	q.Count(&total)
	offset := (opts.Page - 1) * opts.Limit
	err := q.Order("gallery_items.order_num ASC, gallery_items.created_at DESC").
		Limit(opts.Limit).Offset(offset).Find(&list).Error
	return list, total, err
}

func (r *GalleryRepo) FindByID(id uint) (*domain.GalleryItem, error) {
	var item domain.GalleryItem
	err := r.db.Preload("Category").First(&item, id).Error
	return &item, err
}

func (r *GalleryRepo) Create(item *domain.GalleryItem) error { return r.db.Create(item).Error }
func (r *GalleryRepo) Update(item *domain.GalleryItem) error { return r.db.Save(item).Error }
func (r *GalleryRepo) Delete(id uint) error {
	return r.db.Delete(&domain.GalleryItem{}, id).Error
}
