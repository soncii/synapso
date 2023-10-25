package model

import (
	"encoding/json"
	"fmt"
	"regexp"
	"synapso/repository"

	"golang.org/x/crypto/bcrypt"
)

// User defines struct of User data.
type User struct {
	ID           uint   `gorm:"primary_key" json:"id"`
	FirstName    string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email" gorm:"unique"`
	MobileNumber string `json:"mobile_number"`
	Gender       string `json:"gender"`
	DateOfBirth  string `json:"date_of_birth"`
	Password     string `json:"-"`
	Role         string `json:"role"`
}

type UserCreate struct {
	FirstName    string `json:"name"`
	Surname      string `json:"surname"`
	Email        string `json:"email"`
	MobileNumber string `json:"mobile_number"`
	Gender       string `json:"gender"`
	DateOfBirth  string `json:"date_of_birth"`
	Password     string `json:"password"`
	Role         string `json:"role"`
}

func (u UserCreate) ToModel() User {
	return User{
		FirstName:    u.FirstName,
		Surname:      u.Surname,
		Email:        u.Email,
		MobileNumber: u.MobileNumber,
		Gender:       u.Gender,
		DateOfBirth:  u.DateOfBirth,
		Password:     string(HashPassword(u.Password)),
		Role:         u.Role,
	}
}

func (u UserCreate) Validate() (bool, error) {
	if u.FirstName == "" {
		return false, fmt.Errorf("First name is required")
	}

	if u.Surname == "" {
		return false, fmt.Errorf("Surname is required")
	}

	if u.Email == "" {
		return false, fmt.Errorf("Email is required")
	}

	// Validate email format (basic regex pattern)
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(u.Email) {
		return false, fmt.Errorf("Invalid email format")
	}

	if u.MobileNumber == "" {
		return false, fmt.Errorf("Mobile number is required")
	}

	// You can add more specific validation rules for mobile numbers

	if u.Gender == "" {
		return false, fmt.Errorf("Gender is required")
	}

	dateOfBirthRegex := regexp.MustCompile(`^\d{2}\.\d{2}\.\d{4}$`)
	if !dateOfBirthRegex.MatchString(u.DateOfBirth) {
		return false, fmt.Errorf("Invalid date of birth format (DD.MM.YYYY)")
	}

	if u.Password == "" {
		return false, fmt.Errorf("Password is required")
	}

	if len(u.Password) < 6 {
		return false, fmt.Errorf("Password must be at least 6 characters")
	}

	if u.Role != "researcher" && u.Role != "subject" {
		return false, fmt.Errorf("Invalid role")
	}

	return true, nil
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
func (a *User) FindByEmail(rep *repository.Repository, email string) (*User, error) {
	var user User
	err := rep.Where("email = ?", email).First(&user).Error
	return &user, err
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
