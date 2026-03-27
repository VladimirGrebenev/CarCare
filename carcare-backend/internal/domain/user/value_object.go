package user

import (
	"errors"
	"regexp"
)

type Email string

func (e Email) Validate() error {
	if e == "" {
		return errors.New("email is required")
	}
	// Simple email regex
	re := regexp.MustCompile(`^[\w._%+-]+@[\w.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(string(e)) {
		return errors.New("invalid email format")
	}
	return nil
}

type Role string

const (
	RoleUser  Role = "user"
	RoleAdmin Role = "admin"
)

func (r Role) Validate() error {
	if r == "" {
		return errors.New("role is required")
	}
	switch r {
	case RoleUser, RoleAdmin:
		return nil
	default:
		return errors.New("invalid role")
	}
}
