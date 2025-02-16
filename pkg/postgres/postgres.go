package postgres

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"path"
	"runtime"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	DB *pgxpool.Pool
}

func New(ctx context.Context, cfg *Config) (*Repository, error) {
	connConfig, err := pgxpool.ParseConfig(cfg.URL)
	if err != nil {
		return nil, err
	}
	connConfig.MinConns = cfg.MinConns
	connConfig.MaxConns = cfg.MaxConns

	conn, err := pgxpool.New(ctx, connConfig.ConnString())
	if err != nil {
		return nil, err
	}

	return &Repository{DB: conn}, nil
}

func (r *Repository) Close() {
	r.DB.Close()
}

func (r *Repository) database(ctx context.Context) *pgx.Tx {
	return DefaultTrOrDB(ctx, r.DB)
}

func (r *Repository) txFactory(ctx context.Context, options pgx.TxOptions) (pgx.Tx, error) {
	return r.DB.BeginTx(ctx, options)
}

func (r *Repository) Do(ctx context.Context, fn func(c context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel: pgx.RepeatableRead,
	}
	tx, err := r.txFactory(ctx, opts)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			err = errors.Join(tx.Rollback(ctx))
		} else {
			err = errors.Join(tx.Commit(ctx))
		}
	}()

	c := context.WithValue(ctx, defaultCtxKey, tx)
	err = fn(c)
	if err != nil {
		return err
	}
	return nil
}

type ctxKey struct{}

var defaultCtxKey = ctxKey{}

func DefaultTrOrDB(ctx context.Context, db *pgxpool.Pool) *pgx.Tx {
	if tr, ok := ctx.Value(defaultCtxKey).(*pgx.Tx); ok {
		return tr
	}
	return nil
}

func ApplyMigrate(databaseUrl, migrationsDir string) error {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("could not find migration path")
	}
	dir := path.Join(path.Dir(filename), migrationsDir)

	mig, err := migrate.New(
		fmt.Sprintf("file://%s", dir),
		databaseUrl)
	if err != nil {
		slog.Error("failed to create migrations instance: %v", err)
		return err
	}

	if err := mig.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		slog.Error("could not exec migration: %v", err)
		return err
	}
	return nil
}
