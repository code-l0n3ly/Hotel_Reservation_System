package model

import (
	"errors"
	"net/mail"
	"regexp"
)

type User struct {
	UserID     string `json:"userID"`
	Name       string `json:"name"`
	Email      string `json:"email,omitempty"`
	Password   string `json:"password"`
	CreateTime string `json:"createTime,omitempty"`
	UserRole   string `json:"userRole"`
}

func (u *User) Validate() error {
	if u.UserID == "" {
		return errors.New("UserID is required")
	}
	if u.Name == "" {
		return errors.New("name is required")
	}
	if u.Email == "" {
		return errors.New("email is required")
	}
	if _, err := mail.ParseAddress(u.Email); err != nil {
		return errors.New("email is not valid")
	}
	if u.Password == "" {
		return errors.New("password is required")
	}
	if len(u.Password) < 8 {
		return errors.New("password must be at least 8 characters long")
	}
	if u.UserRole == "" {
		return errors.New("UserRole is required")
	}
	return nil
}

func (u *User) IsEmailValid() bool {
	_, err := mail.ParseAddress(u.Email)
	return err == nil
}

func (u *User) IsPasswordStrong() bool {
	var (
		hasMinLen    = false
		hasUppercase = false
		hasSpecial   = false
	)

	var specialCharPattern = regexp.MustCompile(`[!@#\$%\^&\*]`)

	for _, char := range u.Password {
		switch {
		case !hasMinLen:
			hasMinLen = len(u.Password) >= 8
		case !hasUppercase:
			hasUppercase = char >= 'A' && char <= 'Z'
		case !hasSpecial:
			hasSpecial = specialCharPattern.MatchString(u.Password)
		}
	}

	return hasMinLen && hasUppercase && hasSpecial
}
