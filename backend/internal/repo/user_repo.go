package repo

import (
	"github.com/yayasan/cms/internal/domain"
	"gorm.io/gorm"
)

type UserRepo struct{ db *gorm.DB }

func NewUserRepo(db *gorm.DB) *UserRepo { return &UserRepo{db} }

func (r *UserRepo) FindByEmail(email string) (*domain.User, error) {
	var u domain.User
	err := r.db.Preload("Role").Where("email = ? AND deleted_at IS NULL", email).First(&u).Error
	return &u, err
}

func (r *UserRepo) FindByID(id uint) (*domain.User, error) {
	var u domain.User
	err := r.db.Preload("Role").First(&u, id).Error
	return &u, err
}

func (r *UserRepo) FindAll() ([]domain.User, error) {
	var users []domain.User
	err := r.db.Preload("Role").Where("deleted_at IS NULL").Find(&users).Error
	return users, err
}

func (r *UserRepo) Create(u *domain.User) error {
	return r.db.Create(u).Error
}

func (r *UserRepo) Update(u *domain.User) error {
	return r.db.Save(u).Error
}

func (r *UserRepo) Delete(id uint) error {
	return r.db.Delete(&domain.User{}, id).Error
}

func (r *UserRepo) FindAllRoles() ([]domain.Role, error) {
	var roles []domain.Role
	err := r.db.Find(&roles).Error
	return roles, err
}
