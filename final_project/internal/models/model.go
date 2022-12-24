package models

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id       int    `json:"-" db:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Login    string `json:"login"`
}

type Comments struct {
	Id         int    `json:"-"`
	News_id    int    `json:"news_id"`
	User_id    int    `json:"user_id"`
	Created_at string `json:"created_at"`
	Reply_to   int    `json:"reply_to"`
	Content    string `json:"content"`
	Status     bool   `json:"-"`
}

type News struct {
	Id         int    `json:"-"`
	Theme      string `json:"theme"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	Posted_at  string `json:"posted_at"`
	Moder_comm string `json:"moder_comm,omitempty"`
	Posted     bool   `json:"posted,omitempty"`
	User_id    int    `json:"-"`
	Moder_id   int    `json:"-"`
}

func (u *User) HashPassword() error {
	encpass, err := bcrypt.GenerateFromPassword([]byte(u.Password), 14)
	u.Password = string(encpass)
	return err
}

func (u *User) ValidateRegistration() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 25)),
		validation.Field(&u.Login, validation.Required))
}

func (u *User) ValidateAuth() error {
	return validation.ValidateStruct(
		u,
		validation.Field(&u.Email, validation.Required, is.Email),
		validation.Field(&u.Password, validation.Required, validation.Length(8, 25)),
	)
}

func (n *News) ValidateNews() error {
	return validation.ValidateStruct(
		n,
		validation.Field(&n.Title, validation.Required),
		validation.Field(&n.Theme, validation.Required),
		validation.Field(&n.Content, validation.Required),
	)
}

func (n *News) ValidateModerComm() error {
	return validation.ValidateStruct(
		n,
		validation.Field(&n.Moder_comm, validation.Required),
	)
}

func (u *User) ClearPassword() {
	u.Password = ""
}
