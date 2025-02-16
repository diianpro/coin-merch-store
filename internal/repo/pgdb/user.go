package pgdb

import (
	"context"
	"fmt"
	"time"

	"diianpro/coin-merch-store/internal/repo/models"
	"diianpro/coin-merch-store/internal/repo/utils"
	"diianpro/coin-merch-store/pkg/postgres"

	"github.com/jackc/pgx/v5"
	"github.com/pkg/errors"
)

type UserRepo struct {
	*postgres.Repository
}

func NewUserRepo(pg *postgres.Repository) *UserRepo {
	return &UserRepo{pg}
}

func (u *UserRepo) CreateUser(ctx context.Context, user *models.User) (int32, error) {
	query := `INSERT INTO users (username, password, created_at) VALUES ($1, $2, $3) RETURNING user_id`
	var userId int
	err := u.DB.QueryRow(ctx, query, user.Username, user.Password, time.Now()).Scan(&userId)
	if err != nil {
		return 0, nil
	}
	return int32(userId), nil
}
func (u *UserRepo) GetUserByUsernameAndPassword(ctx context.Context, username, password string) (*models.User, error) {
	query := `SELECT user_id, username, password, created_at FROM users WHERE username = $1 AND password = $2`
	var user models.User
	err := u.DB.QueryRow(ctx, query, username, password).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("UserRepo.GetUserByUsernameAndPassword - r.Pool.QueryRow: %v", err)
	}
	return &user, nil
}
func (u *UserRepo) GetUserById(ctx context.Context, id int64) (*models.User, error) {
	query := `SELECT user_id, username, password, created_at FROM users WHERE user_id = $1`
	var user models.User
	err := u.DB.QueryRow(ctx, query, id).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("UserRepo.GetUserById - r.Pool.QueryRow: %v", err)
	}
	return &user, nil
}
func (u *UserRepo) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `SELECT user_id, username, password, created_at FROM users WHERE username = $1`
	var user models.User
	err := u.DB.QueryRow(ctx, query, username).Scan(
		&user.Id,
		&user.Username,
		&user.Password,
		&user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, utils.ErrNotFound
		}
		return nil, fmt.Errorf("UserRepo.GetUserByUsername - r.Pool.QueryRow: %v", err)
	}
	return &user, nil
}
