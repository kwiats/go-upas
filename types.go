package main

import (
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
)

type Base struct {
	ID        string    `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
	UpdatedAt time.Time `json:"updated_at,omitempty"`
}

type User struct {
	Base
	UserCreate
	IsAdmin     bool         `json:"is_admin,omitempty"`
	IsActive    bool         `json:"is_active,omitempty"`
	Groups      []Group      `json:"groups,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
	Profil      Profil       `json:"profil,omitempty"`
}

type UserCreate struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UserDTO struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}
type Profil struct {
	Base
	ProfilDTO
}

type ProfilDTO struct {
	ImageProfile string `json:"image_profile,omitempty"`
	FirstName    string `json:"first_name,omitempty"`
	LastName     string `json:"last_name,omitempty"`
}

type Group struct {
	Base
	GroupDTO
}

type GroupDTO struct {
	Name        string       `json:"name,omitempty"`
	Permissions []Permission `json:"permissions,omitempty"`
}

type Permission struct {
	Base
	PermissionDTO
}

type PermissionDTO struct {
	Name    string `json:"name,omitempty"`
	RawName string `json:"raw_name,omitempty"`
}

func NewTestUser(username, email, password string) *User {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		log.Fatalln("Cannot create new account")
	}
	userProfile := &Profil{
		ImageProfile: "images/id/avatar.jpg",
		FirstName:    "Admin",
		LastName:     "Admin",
	}

	userAdminPermission := &Permission{
		Name:    "All",
		RawName: "all.all_permissions",
	}
	userAdminGroup := &Group{
		Name:        "Admin",
		Permissions: []Permission{*userAdminPermission},
	}
	return &User{
		Username: username,
		Email:    email,
		Password: hashedPassword,
		Profil:   *userProfile,
		IsAdmin:  true,
		IsActive: true,
		Groups:   []Group{*userAdminGroup},
	}

}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
