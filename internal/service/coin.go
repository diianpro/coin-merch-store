package service

import (
	"context"

	"diianpro/coin-merch-store/internal/repo"
)

type SrcCoin struct {
	coinRepo repo.Coin
}

func NewCoinService(coinRepo repo.Coin) *SrcCoin {
	return &SrcCoin{
		coinRepo: coinRepo,
	}
}

func (cr *SrcCoin) TransferCoins(ctx context.Context, from, to, amount int) error {
	err := cr.coinRepo.Do(ctx, func(c context.Context) error {
		err := cr.coinRepo.DecreaseBalance(c, from, amount)
		if err != nil {
			return err
		}
		err = cr.coinRepo.IncreaseBalance(c, to, amount)
		if err != nil {
			return err
		}
		err = cr.coinRepo.AddOperationTransaction(c, from, to, amount)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		return err
	}
	return nil
}
