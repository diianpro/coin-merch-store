package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"diianpro/coin-merch-store/internal/repo"
	"diianpro/coin-merch-store/internal/repo/models"
	repoErr "diianpro/coin-merch-store/internal/repo/utils"
	"diianpro/coin-merch-store/internal/service/utils"
	"diianpro/coin-merch-store/pkg/hasher"

	"github.com/golang-jwt/jwt"
	"github.com/pkg/errors"
)

type TokenClaims struct {
	jwt.StandardClaims
	UserId int
}

type AuthService struct {
	userRepo       repo.User
	passwordHasher hasher.PasswordHasher
	signKey        string
	tokenTTL       time.Duration
}

func NewAuthService(userRepo repo.User, passwordHasher hasher.PasswordHasher, signKey string, tokenTTL time.Duration) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		passwordHasher: passwordHasher,
		signKey:        signKey,
		tokenTTL:       tokenTTL,
	}
}

func (s *AuthService) CreateUser(ctx context.Context, input AuthCreateUserInput) (int32, error) {
	hashPassword, err := s.passwordHasher.Hash(input.Password)
	if err != nil {
		return 0, err
	}
	user := models.User{
		Username: input.Username,
		Password: hashPassword,
	}

	userId, err := s.userRepo.CreateUser(ctx, &user)
	if err != nil {
		if errors.Is(err, repoErr.ErrAlreadyExists) {
			return 0, utils.ErrAccountAlreadyExists
		}
		slog.Error("AuthService.CreateUser - c.userRepo.CreateUser: %v", err)
		return 0, utils.ErrCannotCreateUser
	}
	return userId, nil
}

func (s *AuthService) GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error) {
	hashPassword, err := s.passwordHasher.Hash(input.Password)
	if err != nil {
		return "", err
	}
	// get user from DB
	user, err := s.userRepo.GetUserByUsernameAndPassword(ctx, input.Username, hashPassword)
	if err != nil {
		if errors.Is(err, repoErr.ErrNotFound) {
			return "", utils.ErrAccountNotFound
		}
		slog.Error("AuthService.GenerateToken: cannot get user: %v", err)
		return "", utils.ErrCannotGetUser
	}

	// generate token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		UserId: user.Id,
	})

	// sign token
	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		slog.Error("AuthService.GenerateToken: cannot sign token: %v", err)
		return "", utils.ErrCannotSignToken
	}

	return tokenString, nil
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})

	if err != nil {
		return 0, utils.ErrCannotParseToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, utils.ErrCannotParseToken
	}

	return claims.UserId, nil
}
