package usecase

import "github.com/Hiendang123/golang-server.git/internal/domain"

type UserUsecase struct {
	Repo domain.UserRepository
}

func NewUserUsecase(r domain.UserRepository) *UserUsecase {
	return &UserUsecase{Repo: r}
}

func (uc *UserUsecase) CreateUser(u *domain.User) error {
	return uc.Repo.Create(u)
}
