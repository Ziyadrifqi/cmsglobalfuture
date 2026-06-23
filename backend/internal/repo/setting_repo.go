package repo

import (
	"github.com/yayasan/cms/internal/domain"
	"gorm.io/gorm"
)

type SettingRepo struct{ db *gorm.DB }

func NewSettingRepo(db *gorm.DB) *SettingRepo { return &SettingRepo{db} }

// FindAll mengembalikan semua setting, terurut per section — dipakai untuk
// menampilkan form CMS.
func (r *SettingRepo) FindAll() ([]domain.SiteSetting, error) {
	var list []domain.SiteSetting
	err := r.db.Order("section ASC, order_num ASC, id ASC").Find(&list).Error
	return list, err
}

func (r *SettingRepo) FindByKey(key string) (*domain.SiteSetting, error) {
	var s domain.SiteSetting
	err := r.db.Where("key = ?", key).First(&s).Error
	return &s, err
}

func (r *SettingRepo) Update(s *domain.SiteSetting) error { return r.db.Save(s).Error }

// AsMap mengembalikan seluruh setting sebagai map[key]value — dipakai oleh
// endpoint publik GET /api/v1/settings yang dikonsumsi frontend React.
func (r *SettingRepo) AsMap() (map[string]string, error) {
	list, err := r.FindAll()
	if err != nil {
		return nil, err
	}
	m := make(map[string]string, len(list))
	for _, s := range list {
		m[s.Key] = s.Value
	}
	return m, nil
}
