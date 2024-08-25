package auth

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

const (
	UploadPhotoTokenType = "upload_photo"
)

// PermissionSession данные прав
type PermissionSession struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// AccessSession данные авторизации
type AccessSession struct {
	UserID      uuid.UUID           `json:"user_id"`
	Permissions []PermissionSession `json:"permissions"`
}

type session struct {
	UserID       string
	AccessToken  *string
	RefreshToken string
}

// accessSessionClaims данные авторизованного токена
type accessSessionClaims struct {
	jwt.StandardClaims
	AccessSession
}

// Valid валидность claims
func (c *accessSessionClaims) Valid() error {
	err := c.StandardClaims.Valid()
	if err != nil {
		return err
	}
	return nil
}

// TokenSession данные токена
type TokenSession struct {
	UserID    uuid.UUID
	TokenType string
}
