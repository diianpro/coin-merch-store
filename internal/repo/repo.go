package repo

import (
	"context"

	"diianpro/coin-merch-store/internal/repo/models"
	"diianpro/coin-merch-store/internal/repo/pgdb"
	"diianpro/coin-merch-store/pkg/postgres"
)

type User interface {
	CreateUser(ctx context.Context, user *models.User) (int32, error)
	GetUserByUsernameAndPassword(ctx context.Context, username, password string) (*models.User, error)
	GetUserById(ctx context.Context, id int64) (*models.User, error)
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

type Coin interface {
	CreateWallet(ctx context.Context, userID, amount int) error
	GetBalance(ctx context.Context, userID int) (int64, error)
	DecreaseBalance(ctx context.Context, userID, amount int) error
	IncreaseBalance(ctx context.Context, userID, amount int) error

	AddOperationTransaction(ctx context.Context, fromUserID, toUserID int, amount int) error
	GetCoinFromTransactionHistory(ctx context.Context, from int) ([]models.Info, error)
	GetCoinToTransactionHistory(ctx context.Context, to int) ([]models.Info, error)

	Do(ctx context.Context, fn func(c context.Context) error) error
}

type Merch interface {
	OrderMerch(ctx context.Context, userID, merchID int) error
	GetMerchIDByName(ctx context.Context, name string) (int, int, error)
	GetOrderHistory(ctx context.Context, userID int) ([]models.Merch, error)
}

type Repositories struct {
	User
	Coin
	Merch
}

func NewRepositories(pg *postgres.Repository) *Repositories {
	return &Repositories{
		User:  pgdb.NewUserRepo(pg),
		Coin:  pgdb.NewCoinRepo(pg),
		Merch: pgdb.NewMerchRepo(pg),
	}
}
