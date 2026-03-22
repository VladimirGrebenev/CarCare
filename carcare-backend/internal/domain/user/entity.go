package user

import (
	"errors"
)

type User struct {
	ID    string
	Email Email
	Name  string
	Role  Role
}

// Validate checks business/domain rules for User
func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("name is required")
	}
	if err := u.Email.Validate(); err != nil {
		return err
	}
	if err := u.Role.Validate(); err != nil {
		return err
	}
	return nil
}

func (u *User) IsAdmin() bool {
	return u.Role == RoleAdmin
}

func (u *User) IsActive() bool {
	// For extensibility: добавить поле Active, если потребуется
	return true
}
