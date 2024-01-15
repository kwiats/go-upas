package models

import "github.com/kwiats/go-upas/utils"

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

type UserLogin struct {
	ID       string `json:"id,omitempty"`
	Username string `json:"username"`
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

func (uc *UserCreate) CreateAccount() (*UserCreate, error) {
	hashedPassword, err := utils.HashPassword(uc.Password)
	if err != nil {
		return nil, err // Handle the error in the caller
	}
	profile := uc.createProfile()
	return &UserCreate{
		Username:     uc.Username,
		Password:     hashedPassword,
		Email:        uc.Email,
		FirstName:    profile.FirstName,
		LastName:     profile.LastName,
		ImageProfile: profile.ImageProfile,
	}, nil
}

func (uc *UserCreate) createProfile() *Profil {
	if uc.ImageProfile == "" {
		uc.ImageProfile = "images/default/avatar.png"
	}
	return &Profil{
		FirstName:    uc.FirstName,
		LastName:     uc.LastName,
		ImageProfile: uc.ImageProfile,
	}
}
