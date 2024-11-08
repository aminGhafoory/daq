package models

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/aminGhafoory/daq/internal/database"
	"github.com/aminGhafoory/daq/rand"
	"github.com/google/uuid"
)

const MinBytesPerToken = 32

type Session struct {
	ID        uuid.UUID
	Token     string
	TokenHash string
}

type SessionService struct {
	DB            *database.Queries
	BytesPerToken int
}

func (ss *SessionService) Create(ID uuid.UUID) (*Session, error) {

	bytesPerToken := ss.BytesPerToken
	if bytesPerToken < MinBytesPerToken {
		bytesPerToken = MinBytesPerToken
	}
	sessionToken, err := rand.String(bytesPerToken)
	if err != nil {
		return nil, err
	}
	tokenhash := ss.hash(sessionToken)
	s := Session{
		ID:        ID,
		Token:     sessionToken,
		TokenHash: tokenhash,
	}

	_, err = ss.DB.CreateSession(context.Background(), database.CreateSessionParams{
		UserID:    s.ID,
		TokenHash: s.TokenHash,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, fmt.Errorf("create Session %w", err)
	}
	return &s, nil
}

func (ss *SessionService) User(token string) (*User, error) {

	tokenHash := ss.hash(token)
	u := User{}

	dbUser, err := ss.DB.UserBySession(context.Background(), tokenHash)
	if err != nil {
		return nil, err
	}
	u.Email = dbUser.Email
	u.ID = dbUser.UserID
	u.PasswordHash = dbUser.PasswordHash

	return &u, nil
}

func (ss SessionService) hash(token string) string {
	sum := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(sum[:])
}
