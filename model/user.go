package model

import (
	"encoding/json"

	"github.com/ybkuroki/go-webapp-project-template/repository"
	"golang.org/x/crypto/bcrypt"
)

// User defines struct of User data.
type User struct {
	ID           uint   `gorm:"primary_key" json:"id"`
	FirstName    string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobile_number"`
	Gender       string `json:"gender"`
	DateOfBirth  string `json:"date_of_birth"`
	Password     string `json:"-"`
	Role         string `json:"role"`
}

// TableName returns the table name of User struct and it is used by gorm.
func (User) TableName() string {
	return "user"
}

// NewUserWithPlainPassword is constructor. And it is encoded plain text password by using bcrypt.
func HashPassword(password string) []byte {
	hashed, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
	return hashed
}

// FindByName returns Users full matched given User name.
func (a *User) FindByEmail(rep *repository.Repository, email string) (user *User, err error) {
	err = rep.Preload("Authority").Where("email = ?", email).Find(user).Error
	return
}

// Create persists this User data.
func (a *User) Create(rep *repository.Repository) (*User, error) {
	if error := rep.Create(a).Error; error != nil {
		return nil, error
	}
	return a, nil
}

// ToString is return string of object
func (a *User) ToString() (string, error) {
	bytes, error := json.Marshal(a)
	return string(bytes), error
}
