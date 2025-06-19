package postgres

import (
	"github.com/Hiendang123/golang-server.git/internal/domain"
	"gorm.io/gorm"
)

type UserModel struct {
	ID       uint   `gorm:"primaryKey"`
	Email    string `gorm:"uniqueIndex"`
	Password string
}

func toUserEntity(m *UserModel) *domain.User {
	return &domain.User{ID: m.ID, Email: m.Email, Password: m.Password}
}

func toUserModel(e *domain.User) *UserModel {
	return &UserModel{ID: e.ID, Email: e.Email, Password: e.Password}
}

type UserPostgresRepo struct {
	DB *gorm.DB
}

func NewUserPostgresRepo(db *gorm.DB) domain.UserRepository {
	db.AutoMigrate(&UserModel{})
	return &UserPostgresRepo{DB: db}
}

func (r *UserPostgresRepo) Create(u *domain.User) error {
	model := toUserModel(u)
	if err := r.DB.Create(model).Error; err != nil {
		return err
	}
	u.ID = model.ID // Update the ID after creation
	return nil
}
