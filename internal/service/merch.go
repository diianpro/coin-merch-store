package service

import (
	"context"

	"diianpro/coin-merch-store/internal/repo"
)

type SrcMerch struct {
	merchRepo repo.Merch
	coinRepo  repo.Coin
}

func NewMerchService(merchRepo repo.Merch, coinRepo repo.Coin) *SrcMerch {
	return &SrcMerch{
		merchRepo: merchRepo,
		coinRepo:  coinRepo,
	}
}

func (m *SrcMerch) OrderMerch(ctx context.Context, userID int, merch string) error {
	merchID, price, err := m.merchRepo.GetMerchIDByName(ctx, merch)
	if err != nil {
		return err
	}
	err = m.coinRepo.DecreaseBalance(ctx, userID, price)
	if err != nil {
		return err
	}
	err = m.merchRepo.OrderMerch(ctx, userID, merchID)
	if err != nil {
		return err
	}
	return nil
}
