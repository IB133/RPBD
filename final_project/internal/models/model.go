package models

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"-"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Login    string `json:"login"`
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

func (u *User) Validate() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 25)),
		validation.Field(&u.Login, validation.Required))
}

func (u *User) ClearPassword() {
	u.Password = ""
}
