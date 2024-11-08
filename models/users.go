package models

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/aminGhafoory/daq/internal/database"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	DB *database.Queries
}

type User struct {
	ID           uuid.UUID
	Email        string
	PasswordHash string
}

func (us *UserService) CreateUser(email, password string) (*User, error) {
	email = strings.ToLower(email)
	HashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return &User{}, fmt.Errorf("create User: %w", err)
	}
	passwordHash := string(HashedBytes)

	newUser, err := us.DB.CreateUser(context.Background(), database.CreateUserParams{
		UserID:       uuid.New(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		Email:        email,
		PasswordHash: passwordHash,
	})

	if err != nil {
		return &User{}, fmt.Errorf("create User: %w", err)
	}

	return &User{
		ID:           newUser.UserID,
		Email:        newUser.Email,
		PasswordHash: newUser.PasswordHash,
	}, nil

}

func (us *UserService) Auth(email, passwordHash string) (*User, error) {
	user, err := us.DB.UserByEmail(context.Background(), email)
	if err != nil {
		return &User{}, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(passwordHash))
	if err != nil {
		return &User{}, fmt.Errorf("login unsuccessful %w", err)
	}

	return &User{
		ID:           user.UserID,
		Email:        email,
		PasswordHash: passwordHash,
	}, nil
}
