package user

type UserRepository interface {
	GetUserById(id string) (*User, error)
	AddUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(id string) error
}
