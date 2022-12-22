package models

import (
	"time"

	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int
	Email    string `validate:"required,email"`
	Password string `validate:"required,min=8,max=25"`
	Login    string `validate:"required"`
}

type Comments struct {
	Id         int
	News_id    int
	User_id    int
	Created_at time.Time
	Reply_to   int
	Content    string
	Status     bool
}

type News struct {
	Id         int
	Theme      string
	Title      string
	Content    string
	Posted_at  time.Time
	Moder_comm string
	Posted     bool
	User_id    int
	Moder_id   int
}

func (u *User) HashPassword() error {
	encpass, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	u.Password = string(encpass)
	return err
}

func (u *User) ValidateUser(v *validator.Validate) error {
	if err := v.Struct(u); err != nil {
		return err
	}
	return nil
}
