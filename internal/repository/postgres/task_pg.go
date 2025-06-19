package postgres

import (
	"github.com/Hiendang123/golang-server.git/internal/domain"
	"gorm.io/gorm"
)

type TaskModel struct {
	ID    uint `gorm:"primaryKey"`
	Title string
	Done  bool
}

func toEntity(m *TaskModel) *domain.Task {
	return &domain.Task{ID: m.ID, Title: m.Title, Done: m.Done}
}

func toModel(e *domain.Task) *TaskModel {
	return &TaskModel{ID: e.ID, Title: e.Title, Done: e.Done}
}

type TaskPostgresRepo struct {
	DB *gorm.DB
}

func NewTaskPostgresRepo(db *gorm.DB) domain.TaskRepository {
	db.AutoMigrate(&TaskModel{})
	return &TaskPostgresRepo{DB: db}
}

func (r *TaskPostgresRepo) Create(t *domain.Task) error {
	model := toModel(t)
	if err := r.DB.Create(model).Error; err != nil {
		return err
	}
	t.ID = model.ID // Update the ID in the domain entity
	return nil
}

func (r *TaskPostgresRepo) GetAll(limit, offset int, filter domain.Task) ([]domain.Task, int64, error) {
	var models []TaskModel
	query := r.DB.Model(&TaskModel{})

	if filter.Title != "" {
		query = query.Where("title LIKE ?", "%"+filter.Title+"%")
	}
	if filter.Done {
		query = query.Where("done = ?", true)
	} else if !filter.Done {
		query = query.Where("done = ?", false)
	}

	var total int64
	query.Model(&TaskModel{}).Count(&total)
	if err := query.Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}

	tasks := []domain.Task{}
	for _, m := range models {
		tasks = append(tasks, *toEntity(&m))
	}
	return tasks, total, nil
}

func (r *TaskPostgresRepo) GetByID(id uint) (*domain.Task, error) {
	var m TaskModel
	if err := r.DB.First(&m, id).Error; err != nil {
		return nil, err
	}
	return toEntity(&m), nil
}

func (r *TaskPostgresRepo) Update(t *domain.Task) error {
	return r.DB.Save(toModel(t)).Error
}

// func (r *TaskPostgresRepo) Delete(id uint) error {
// 	return r.DB.Delete(&TaskModel{}, id).Error
// }

func (r *TaskPostgresRepo) DeleteAll() error {
	return r.DB.Unscoped().Delete(&TaskModel{}).Error
}

func (r *TaskPostgresRepo) Login(username, password string) (string, error) {
	var user domain.UserTask
	if err := r.DB.Where("email = ? AND password = ?", username, password).First(&user).Error; err != nil {
		return "", err
	}
	return "token", nil
}
