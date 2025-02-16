package pgdb

import (
	"context"
	"fmt"
	"math/rand"
	"testing"
	"time"

	"diianpro/coin-merch-store/internal/repo/models"
	"diianpro/coin-merch-store/pkg/postgres"

	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	db        *postgres.Repository
	repoUser  *UserRepo
	repoCoin  *CoinRepo
	repoMerch *MerchRepo
	container *postgres.Container
	suite.Suite
}

func (i *IntegrationTestSuite) SetupSuite() {
	var err error

	ctx := context.Background()

	cfg := &postgres.Config{
		MinConns: 1,
		MaxConns: 2,
	}
	i.container, err = postgres.NewContainer(cfg, func() error {
		i.db, err = postgres.New(ctx, cfg)
		if err != nil {
			return err
		}
		err = postgres.ApplyMigrate(cfg.URL, "../../migrations")
		if err != nil {
			return err
		}
		return nil
	})
	i.Require().NoError(err)

	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	err = i.db.DB.Ping(ctxPing)
	if err != nil {
		return
	}
}

func (i *IntegrationTestSuite) TearDownSuite() {
	err := i.container.Purge()
	i.Assert().NoError(err)
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}

func (i *IntegrationTestSuite) TestCreateUserAndGetUserById() {
	i.repoUser = NewUserRepo(i.db)

	ctx := context.Background()

	userID, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	id, err := i.repoUser.GetUserById(ctx, int64(userID))
	i.Require().NoError(err)
	i.Require().Equal(int(userID), id.Id)

	info, err := i.repoUser.GetUserByUsername(ctx, id.Username)
	i.Require().NoError(err)
	i.Require().Equal(id.Username, info.Username)
}

func (i *IntegrationTestSuite) TestGetUserByUsername() {
	i.repoUser = NewUserRepo(i.db)

	ctx := context.Background()

	userID, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	id, err := i.repoUser.GetUserById(ctx, int64(userID))
	i.Require().NoError(err)
	i.Require().Equal(int(userID), id.Id)

	info, err := i.repoUser.GetUserByUsername(ctx, id.Username)
	i.Require().NoError(err)
	i.Require().Equal(id.Username, info.Username)
}

func (i *IntegrationTestSuite) TestGetUserByUsernameAndPassword() {
	i.repoUser = NewUserRepo(i.db)

	ctx := context.Background()

	userID, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	id, err := i.repoUser.GetUserById(ctx, int64(userID))
	i.Require().NoError(err)
	i.Require().Equal(int(userID), id.Id)

	info, err := i.repoUser.GetUserByUsernameAndPassword(ctx, id.Username, id.Password)
	i.Require().NoError(err)
	i.Require().Equal(id.Id, info.Id)
}

func (i *IntegrationTestSuite) TestCreateWallet() {
	i.repoUser = NewUserRepo(i.db)
	i.repoCoin = NewCoinRepo(i.db)

	ctx := context.Background()

	userID, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	err = i.repoCoin.CreateWallet(ctx, int(userID), 0)
	i.Require().NoError(err)
}

func (i *IntegrationTestSuite) TestGetBalance() {
	i.repoUser = NewUserRepo(i.db)
	i.repoCoin = NewCoinRepo(i.db)

	ctx := context.Background()

	userID, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	balanceAmount := 10
	err = i.repoCoin.CreateWallet(ctx, int(userID), balanceAmount)
	i.Require().NoError(err)

	amount, err := i.repoCoin.GetBalance(ctx, int(userID))
	i.Require().NoError(err)
	i.Require().Equal(balanceAmount, int(amount))
}

func (i *IntegrationTestSuite) TestDecreaseBalance() {
	i.repoUser = NewUserRepo(i.db)
	i.repoCoin = NewCoinRepo(i.db)

	ctx := context.Background()

	userID, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	balanceAmount := 10
	err = i.repoCoin.CreateWallet(ctx, int(userID), balanceAmount)
	i.Require().NoError(err)

	decreaseBalanceOn := 3
	err = i.repoCoin.DecreaseBalance(ctx, int(userID), decreaseBalanceOn)
	i.Require().NoError(err)

	amount, err := i.repoCoin.GetBalance(ctx, int(userID))
	i.Require().NoError(err)
	i.Require().Equal(balanceAmount-decreaseBalanceOn, int(amount))
}

func (i *IntegrationTestSuite) TestIncreaseBalance() {
	i.repoUser = NewUserRepo(i.db)
	i.repoCoin = NewCoinRepo(i.db)

	ctx := context.Background()

	userID, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	balanceAmount := 10
	err = i.repoCoin.CreateWallet(ctx, int(userID), balanceAmount)
	i.Require().NoError(err)

	increaseBalanceOn := 3
	err = i.repoCoin.IncreaseBalance(ctx, int(userID), increaseBalanceOn)
	i.Require().NoError(err)

	amount, err := i.repoCoin.GetBalance(ctx, int(userID))
	i.Require().NoError(err)
	i.Require().Equal(balanceAmount+increaseBalanceOn, int(amount))
}

func (i *IntegrationTestSuite) TestAddOperationTransaction() {
	i.repoUser = NewUserRepo(i.db)
	i.repoCoin = NewCoinRepo(i.db)

	ctx := context.Background()

	userID, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	userIDAlt, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	err = i.repoCoin.CreateWallet(ctx, int(userID), 10)
	i.Require().NoError(err)

	err = i.repoCoin.CreateWallet(ctx, int(userIDAlt), 0)
	i.Require().NoError(err)

	err = i.repoCoin.AddOperationTransaction(ctx, int(userID), int(userIDAlt), 1)
	i.Require().NoError(err)
}

func (i *IntegrationTestSuite) TestGetCoinFromTransactionHistory() {
	i.repoUser = NewUserRepo(i.db)
	i.repoCoin = NewCoinRepo(i.db)

	ctx := context.Background()

	userID, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	userIDAlt, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	err = i.repoCoin.CreateWallet(ctx, int(userID), 10)
	i.Require().NoError(err)

	err = i.repoCoin.CreateWallet(ctx, int(userIDAlt), 0)
	i.Require().NoError(err)

	err = i.repoCoin.AddOperationTransaction(ctx, int(userID), int(userIDAlt), 1)
	i.Require().NoError(err)

	results, err := i.repoCoin.GetCoinFromTransactionHistory(ctx, int(userID))
	i.Require().NoError(err)
	i.Require().NotNil(results)
}

func (i *IntegrationTestSuite) TestGetCoinToTransactionHistory() {
	i.repoUser = NewUserRepo(i.db)
	i.repoCoin = NewCoinRepo(i.db)

	ctx := context.Background()

	userID, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	userIDAlt, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	err = i.repoCoin.CreateWallet(ctx, int(userID), 10)
	i.Require().NoError(err)

	err = i.repoCoin.CreateWallet(ctx, int(userIDAlt), 0)
	i.Require().NoError(err)

	err = i.repoCoin.AddOperationTransaction(ctx, int(userID), int(userIDAlt), 1)
	i.Require().NoError(err)

	results, err := i.repoCoin.GetCoinToTransactionHistory(ctx, int(userIDAlt))
	i.Require().NoError(err)
	i.Require().NotNil(results)
}

func (i *IntegrationTestSuite) TestCoinDo() {
	ctx := context.Background()

	i.repoUser = NewUserRepo(i.db)
	i.repoCoin = NewCoinRepo(i.db)

	userID, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	userIDAlt, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	err = i.repoCoin.CreateWallet(ctx, int(userID), 10)
	i.Require().NoError(err)

	err = i.repoCoin.CreateWallet(ctx, int(userIDAlt), 0)
	i.Require().NoError(err)

	errDo := i.repoCoin.Do(ctx, func(c context.Context) error {
		err = i.repoCoin.IncreaseBalance(c, int(userIDAlt), 500)
		i.Require().NoError(err)

		err = i.repoCoin.DecreaseBalance(c, int(userID), 500)
		if err != nil {
			i.Require().Error(err)
			return err
		}

		return nil
	})
	if errDo != nil {
		i.Require().Error(errDo)
		return
	}

	balance, err := i.repoCoin.GetBalance(ctx, int(userID))
	i.Require().NoError(err)
	i.Require().Equal(int(balance), 10)

	balance, err = i.repoCoin.GetBalance(ctx, int(userIDAlt))
	i.Require().NoError(err)
	i.Require().Equal(int(balance), 0)
}

func (i *IntegrationTestSuite) TestGetMerchByID() {
	ctx := context.Background()

	i.repoMerch = NewMerchRepo(i.db)

	name, price, err := i.repoMerch.GetMerchIDByName(ctx, "t-shirt")
	i.Require().NoError(err)
	i.Require().NotNil(price)
	i.Require().Equal(1, name)
}

func (i *IntegrationTestSuite) TestGetOrderHistory() {
	ctx := context.Background()

	i.repoUser = NewUserRepo(i.db)
	i.repoMerch = NewMerchRepo(i.db)

	userID, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	merchID, price, err := i.repoMerch.GetMerchIDByName(ctx, "powerbank")
	i.Require().NoError(err)
	i.Require().NotNil(merchID)
	i.Require().NotNil(price)

	err = i.repoMerch.OrderMerch(ctx, int(userID), merchID)
	i.Require().NoError(err)

	history, err := i.repoMerch.GetOrderHistory(ctx, int(userID))
	i.Require().NoError(err)
	i.Require().NotNil(history)
}

func (i *IntegrationTestSuite) TestGetOrderHistoryIfEmpty() {
	ctx := context.Background()

	i.repoUser = NewUserRepo(i.db)
	i.repoMerch = NewMerchRepo(i.db)

	userID, err := i.repoUser.CreateUser(ctx, generateRandomUser())
	i.Require().NoError(err)

	history, err := i.repoMerch.GetOrderHistory(ctx, int(userID))
	i.Require().NoError(err)
	i.Require().Empty(history)
}

var (
	randomSeed = time.Now().UnixNano()
	randGen    = rand.New(rand.NewSource(randomSeed))
)

func generateRandomUser() *models.User {
	username := fmt.Sprintf("test-%d", randGen.Intn(10000))
	password := fmt.Sprintf("password-%d", randGen.Intn(10000))
	return &models.User{
		Username:  username,
		Password:  password,
		CreatedAt: time.Now(),
	}
}
