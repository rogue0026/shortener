package domain

import "golang.org/x/crypto/bcrypt"

type User struct {
	UserID   string `json:"user_id,omitempty"`
	Login    string `json:"login"`
	Password string `json:"password"`
	Email    string `json:"email,omitempty"`
}

func (u *User) HashPassword() error {
	hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	if err != nil {
		return err
	}
	u.Password = string(hash)
	return nil
}

func (u *User) CheckPassword(inPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(inPassword))
}
