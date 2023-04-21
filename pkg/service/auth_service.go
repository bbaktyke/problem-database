package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"git.01.alem.school/bbaktyke/test.project.git/pkg/models"
	"git.01.alem.school/bbaktyke/test.project.git/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	salt       = "qwewqydgasgay213fdds"
	tokenTTL   = 12 * time.Hour
	signingKey = "adsersd213dssf"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserID int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (auth *AuthService) CreateUserService(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return auth.repo.CreateUserRepo(user)
}

func (auth *AuthService) GenerateTokenService(username, password string) (string, error) {
	user, err := auth.repo.GetUserRepo(username, generatePasswordHash(password))
	if err != nil {
		return "", nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserID: user.ID,
	})

	return token.SignedString([]byte(signingKey))
}

func (auth *AuthService) ParseTokenService(tokenString string) (int, error) {
	parsedToken, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return 0, err
	}
	claims, ok := parsedToken.Claims.(*tokenClaims)
	if !ok {
		return 0, errors.New("token claims are not of type *tokenClaims")
	}
	return claims.UserID, nil
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
