package usecase

import "github.com/Hiendang123/golang-server.git/internal/domain"

type TaskUsecase struct {
	Repo domain.TaskRepository
}

func NewTaskUsecase(r domain.TaskRepository) *TaskUsecase {
	return &TaskUsecase{Repo: r}
}

func (uc *TaskUsecase) CreateTask(t *domain.Task) error {
	return uc.Repo.Create(t)
}

func (uc *TaskUsecase) GetAll(limit, offset int, filter domain.Task) ([]domain.Task, int64, error) {
	return uc.Repo.GetAll(limit, offset, filter)
}

func (uc *TaskUsecase) GetTaskByID(id uint) (*domain.Task, error) {
	return uc.Repo.GetByID(id)
}

func (uc *TaskUsecase) UpdateTask(t *domain.Task) error {
	return uc.Repo.Update(t)
}

// func (uc *TaskUsecase) DeleteTask(id uint) error {
// 	return uc.Repo.Delete(id)
// }

func (uc *TaskUsecase) DeleteAll() error {
	return uc.Repo.DeleteAll()
}
