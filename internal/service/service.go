package service

import (
	"context"
	"time"

	"diianpro/coin-merch-store/internal/repo"
	"diianpro/coin-merch-store/pkg/hasher"
)

type AuthCreateUserInput struct {
	Username string
	Password string
}

type AuthGenerateTokenInput struct {
	Username string
	Password string
}

type Auth interface {
	CreateUser(ctx context.Context, input AuthCreateUserInput) (int32, error)
	GenerateToken(ctx context.Context, input AuthGenerateTokenInput) (string, error)
	ParseToken(token string) (int, error)
}

type Coin interface {
	TransferCoins(ctx context.Context, from, to, amount int) error
}

type Merch interface {
	OrderMerch(ctx context.Context, userID int, merch string) error
}

type Services struct {
	Auth  Auth
	Coin  Coin
	Merch Merch
}

type ServicesDependencies struct {
	Repos  *repo.Repositories
	Hasher hasher.PasswordHasher

	SignKey  string
	TokenTTL time.Duration
}

func NewServices(deps ServicesDependencies) *Services {
	return &Services{
		Auth:  NewAuthService(deps.Repos.User, deps.Hasher, deps.SignKey, deps.TokenTTL),
		Coin:  NewCoinService(deps.Repos.Coin),
		Merch: NewMerchService(deps.Repos.Merch, deps.Repos.Coin),
	}
}
