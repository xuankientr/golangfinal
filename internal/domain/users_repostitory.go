package domain

type UserRepository interface {
	Create(user *User) error
}
