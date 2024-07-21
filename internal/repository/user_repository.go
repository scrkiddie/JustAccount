package repository

import (
	"awesomeProject12/internal/entity"
	"gorm.io/gorm"
)

type UserRepository struct {
}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

func (r *UserRepository) Create(db *gorm.DB, entity *entity.User) error {
	return db.Create(entity).Error
}

func (r *UserRepository) Update(db *gorm.DB, entity *entity.User) error {
	return db.Save(entity).Error
}

func (r *UserRepository) FindByUsername(db *gorm.DB, entity *entity.User, username string) error {
	return db.Where("username = ?", username).Take(entity).Error
}

func (r *UserRepository) FindById(db *gorm.DB, entity *entity.User, id int) error {
	return db.Where("id = ?", id).Take(entity).Error
}

func (r *UserRepository) CountByUsername(db *gorm.DB, username any) (int64, error) {
	var count int64
	if err := db.Model(&entity.User{}).Where("username = ?", username).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func (r *UserRepository) CountByEmail(db *gorm.DB, email any) (int64, error) {
	var count int64
	if err := db.Model(&entity.User{}).Where("email = ?", email).Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}
