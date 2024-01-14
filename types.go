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
		PermissionName: "All",
		RawName:        "all.all_permissions",
	}
	userAdminGroup := &Group{

		GroupName:   "Admin",
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

func CreateAccount(username, password, email, firstName, lastName string) (*User, error) {
	user, err := createUser(username, password, email)
	if err != nil {
		return nil, err
	}

	profile := createProfile(firstName, lastName)

	return &User{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
		Profil:   *profile,
		IsAdmin:  false,
		IsActive: true,
		Groups:   make([]Group, 0),
	}, nil
}

func createUser(username, password, email string) (*UserCreate, error) {
	hashedPassword, err := HashPassword(password)
	if err != nil {
		return nil, err
	}
	return &UserCreate{Username: username, Password: hashedPassword, Email: email}, nil
}

func createProfile(firstName, lastName string) *Profil {
	return &Profil{FirstName: firstName, LastName: lastName, ImageProfile: "images/default/avatar.jpg"}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
