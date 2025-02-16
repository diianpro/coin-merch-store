package pgdb

import (
	"context"

	"diianpro/coin-merch-store/internal/repo/models"
	"diianpro/coin-merch-store/pkg/postgres"
)

type CoinRepo struct {
	*postgres.Repository
}

func NewCoinRepo(pg *postgres.Repository) *CoinRepo {
	return &CoinRepo{pg}
}

func (c *CoinRepo) CreateWallet(ctx context.Context, userID, amount int) error {
	query := `
        INSERT INTO coins (user_id, amount)
        VALUES ($1, $2)
    `
	if _, err := c.DB.Exec(ctx, query, userID, amount); err != nil {
		return err
	}
	return nil
}

func (c *CoinRepo) GetBalance(ctx context.Context, userID int) (int64, error) {
	var amount int64
	query := `SELECT amount FROM coins WHERE user_id = $1`
	err := c.DB.QueryRow(ctx, query, userID).Scan(&amount)
	if err != nil {
		return 0, err
	}
	return amount, nil
}

func (c *CoinRepo) DecreaseBalance(ctx context.Context, userID, amount int) error {
	query := `
		UPDATE coins
		SET amount = amount - $2
		WHERE user_id = $1`
	if _, err := c.DB.Exec(ctx, query, userID, amount); err != nil {
		return err
	}
	return nil
}

func (c *CoinRepo) IncreaseBalance(ctx context.Context, userID, amount int) error {
	query := `
		UPDATE coins
		SET amount = amount + $2
		WHERE user_id = $1`
	if _, err := c.DB.Exec(ctx, query, userID, amount); err != nil {
		return err
	}
	return nil
}

func (c *CoinRepo) AddOperationTransaction(ctx context.Context, fromUserID, toUserID int, amount int) error {
	query := `
		INSERT INTO operations (from_user_id, to_user_id, amount, transaction_date)
		VALUES ($1, $2, $3, NOW())`
	if _, err := c.DB.Exec(ctx, query, fromUserID, toUserID, amount); err != nil {
		return err
	}
	return nil
}

func (c *CoinRepo) GetCoinFromTransactionHistory(ctx context.Context, from int) ([]models.Info, error) {
	resultsInfo := make([]models.Info, 0)
	query := `SELECT from_user_id, to_user_id, amount, transaction_date FROM operations WHERE from_user_id = $1`
	rows, err := c.DB.Query(ctx, query, from)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var result models.Info
		err = rows.Scan(&result.FromUserID, &result.ToUserID, &result.Amount, &result.Date)
		if err != nil {
			return nil, err
		}
		resultsInfo = append(resultsInfo, result)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()
	return resultsInfo, nil
}

func (c *CoinRepo) GetCoinToTransactionHistory(ctx context.Context, to int) ([]models.Info, error) {
	resultsInfo := make([]models.Info, 0)
	query := `SELECT from_user_id, to_user_id, amount, transaction_date FROM operations WHERE to_user_id = $1`
	rows, err := c.DB.Query(ctx, query, to)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var result models.Info
		err = rows.Scan(&result.FromUserID, &result.ToUserID, &result.Amount, &result.Date)
		if err != nil {
			return nil, err
		}
		resultsInfo = append(resultsInfo, result)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	rows.Close()
	return resultsInfo, nil
}
