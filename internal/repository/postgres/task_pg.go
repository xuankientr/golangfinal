package postgres

import (
	"time"

	"github.com/Hiendang123/golang-server.git/internal/domain"
	"gorm.io/gorm"
)

type TaskModel struct {
	ID        uint `gorm:"primaryKey"`
	Title     string
	Done      bool
	CreatedAt time.Time `gorm:"autoCreateTime"`
	CreatedBy uint
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	UpdatedBy uint
}

func toEntity(m *TaskModel) *domain.Task {
	return &domain.Task{
		ID:        m.ID,
		Title:     m.Title,
		Done:      m.Done,
		CreatedAt: m.CreatedAt.Format(time.RFC3339),
		CreatedBy: m.CreatedBy,
		UpdatedAt: m.UpdatedAt.Format(time.RFC3339),
		UpdatedBy: m.UpdatedBy,
	}
}

func toModel(e *domain.Task) *TaskModel {
	return &TaskModel{
		ID:        e.ID,
		Title:     e.Title,
		Done:      e.Done,
		CreatedBy: e.CreatedBy,
		UpdatedBy: e.UpdatedBy,
	}
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
	t.ID = model.ID // Update the ID after creation
	return nil
}

func (r *TaskPostgresRepo) GetAll(limit, offset int, filters domain.Task) ([]domain.Task, int64, error) {
	var models []TaskModel
	query := r.DB.Model(&TaskModel{})

	if filters.Title != "" {
		query = query.Where("title LIKE ?", "%"+filters.Title+"%")
	}
	if filters.Done {
		query = query.Where("done = ?", filters.Done)
	} else if !filters.Done {
		query = query.Where("done = ?", false)
	}
	if filters.CreatedBy != 0 {
		query = query.Where("created_by = ?", filters.CreatedBy)
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

func (r *TaskPostgresRepo) Delete(id uint) error {
	return r.DB.Delete(&TaskModel{}, id).Error
}

// DeleteAll deletes all tasks from the database.
func (r *TaskPostgresRepo) DeleteAll() error {
	return r.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&TaskModel{}).Error
}
