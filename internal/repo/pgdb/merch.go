package pgdb

import (
	"context"
	"database/sql"
	"fmt"

	"diianpro/coin-merch-store/internal/repo/models"
	"diianpro/coin-merch-store/pkg/postgres"

	"github.com/pkg/errors"
)

type MerchRepo struct {
	*postgres.Repository
}

func NewMerchRepo(pg *postgres.Repository) *MerchRepo {
	return &MerchRepo{pg}
}

func (m *MerchRepo) OrderMerch(ctx context.Context, userID int, merchID int) error {
	query := `INSERT INTO purchases (user_id, merch_id, purchase_date) VALUES ($1, $2, NOW())`
	_, err := m.DB.Exec(ctx, query, userID, merchID)
	if err != nil {
		return err
	}
	return nil
}

func (m *MerchRepo) GetOrderHistory(ctx context.Context, userID int) ([]models.Merch, error) {
	resultsInfo := make([]models.Merch, 0)
	query := `
		SELECT m.name, m.price, p.purchase_date
		FROM merch m
		JOIN purchases p ON m.merch_id = p.merch_id
		WHERE p.user_id = $1;`
	rows, err := m.DB.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var result models.Merch
		err = rows.Scan(&result.Name, &result.Amount, &result.Date)
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

func (m *MerchRepo) GetMerchIDByName(ctx context.Context, name string) (int, int, error) {
	var merchID, amount int
	err := m.DB.QueryRow(ctx, `
        SELECT merch_id, price FROM merch WHERE name = $1
    `, name).Scan(&merchID, &amount)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, 0, fmt.Errorf("merch with name %s not found", name)
		}
		return 0, 0, fmt.Errorf("failed to get merch ID by name: %w", err)
	}
	return merchID, amount, nil
}
