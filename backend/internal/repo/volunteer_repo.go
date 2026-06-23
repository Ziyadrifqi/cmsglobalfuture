package repo

import (
	"github.com/yayasan/cms/internal/domain"
	"gorm.io/gorm"
)

// ── VolunteerRepo ─────────────────────────────────────────────────────────────

type VolunteerRepo struct{ db *gorm.DB }

func NewVolunteerRepo(db *gorm.DB) *VolunteerRepo { return &VolunteerRepo{db} }

func (r *VolunteerRepo) FindAll(opts domain.ListOptions) ([]domain.Volunteer, int64, error) {
	var list []domain.Volunteer
	var total int64
	q := r.db.Model(&domain.Volunteer{})
	if opts.Status != "" {
		q = q.Where("status = ?", opts.Status)
	}
	if opts.Search != "" {
		q = q.Where("title ILIKE ? OR division ILIKE ?", "%"+opts.Search+"%", "%"+opts.Search+"%")
	}
	q.Count(&total)
	offset := (opts.Page - 1) * opts.Limit
	err := q.Order("created_at DESC").Limit(opts.Limit).Offset(offset).Find(&list).Error
	return list, total, err
}

func (r *VolunteerRepo) FindActive(limit int) ([]domain.Volunteer, error) {
	var list []domain.Volunteer
	err := r.db.Where("status = ?", domain.VolunteerActive).
		Order("created_at DESC").Limit(limit).Find(&list).Error
	return list, err
}

func (r *VolunteerRepo) FindBySlug(slug string) (*domain.Volunteer, error) {
	var v domain.Volunteer
	err := r.db.Where("slug = ? AND status = ?", slug, domain.VolunteerActive).First(&v).Error
	return &v, err
}

func (r *VolunteerRepo) FindBySlugAny(slug string) (*domain.Volunteer, error) {
	var v domain.Volunteer
	err := r.db.Where("slug = ?", slug).First(&v).Error
	return &v, err
}

func (r *VolunteerRepo) FindByID(id uint) (*domain.Volunteer, error) {
	var v domain.Volunteer
	err := r.db.First(&v, id).Error
	return &v, err
}

func (r *VolunteerRepo) Create(v *domain.Volunteer) error { return r.db.Create(v).Error }
func (r *VolunteerRepo) Update(v *domain.Volunteer) error { return r.db.Save(v).Error }
func (r *VolunteerRepo) Delete(id uint) error             { return r.db.Delete(&domain.Volunteer{}, id).Error }

// ── ApplicantRepo ─────────────────────────────────────────────────────────────

type ApplicantRepo struct{ db *gorm.DB }

func NewApplicantRepo(db *gorm.DB) *ApplicantRepo { return &ApplicantRepo{db} }

func (r *ApplicantRepo) FindByVolunteer(volunteerID uint) ([]domain.Applicant, error) {
	var list []domain.Applicant
	err := r.db.Where("volunteer_id = ?", volunteerID).Order("created_at DESC").Find(&list).Error
	return list, err
}

func (r *ApplicantRepo) FindByID(id uint) (*domain.Applicant, error) {
	var a domain.Applicant
	err := r.db.Preload("Volunteer").First(&a, id).Error
	return &a, err
}

func (r *ApplicantRepo) Create(a *domain.Applicant) error { return r.db.Create(a).Error }
func (r *ApplicantRepo) Update(a *domain.Applicant) error { return r.db.Save(a).Error }

func (r *ApplicantRepo) CountTotal() int64 {
	var count int64
	r.db.Model(&domain.Applicant{}).Count(&count)
	return count
}
