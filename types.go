package main

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Base struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	Base
	Username    string       `json:"username"`
	Email       string       `json:"email"`
	Password    string       `json:"password"`
	IsAdmin     bool         `json:"is_admin,omitempty"`
	IsActive    bool         `json:"is_active,omitempty"`
	Groups      []Group      `json:"groups,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
	Profil      Profil       `json:"profil,omitempty"`
}

type UserCreate struct {
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
	ImageProfile string `json:"image_profile,omitempty"`
}

type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
type Profil struct {
	Base
	ImageProfile string `json:"image_profile,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
}

type Group struct {
	Base
	GroupName   string       `json:"group_name,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
}

type Permission struct {
	Base
	PermissionName string `json:"permission_name,omitempty"`
	RawName        string `json:"raw_name,omitempty"`
}

func CreateBaseObject() *Base {
	currentTime := time.Now().UTC()
	return &Base{
		CreatedAt: currentTime,
		UpdatedAt: currentTime,
	}
}
func UpdatedBaseObject() *Base {
	currentTime := time.Now().UTC()
	return &Base{
		UpdatedAt: currentTime,
	}
}

func (user *UserCreate) CreateAccount() (*UserCreate, error) {
	user.HashPassword()
	profile := user.createProfile()
	return &UserCreate{
		Username:     user.Username,
		Password:     user.Password,
		Email:        user.Email,
		FirstName:    profile.FirstName,
		LastName:     profile.LastName,
		ImageProfile: profile.ImageProfile,
	}, nil
}

func (user *UserCreate) createProfile() *Profil {
	if user.ImageProfile == "" {
		user.ImageProfile = "images/default/avatar.png"
	}
	return &Profil{
		FirstName:    user.FirstName,
		LastName:     user.LastName,
		ImageProfile: user.ImageProfile,
	}
}

func (user *UserCreate) HashPassword() {
	hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		log.Fatalln("Cannot hash password")
		return
	}
	user.Password = string(hashed)
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
