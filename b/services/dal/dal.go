package dal

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/vsabirov/fintech/b/services/entity"
)

type TransferParams struct {
	ID     string
	Amount float64

	Sender   string
	Receiver string
}

type DataSource interface {
	prepareStatements(ctx context.Context, tx pgx.Tx) error

	Transfer(ctx context.Context, params TransferParams) error
}

type PostgresDataSource struct {
	Pool *pgxpool.Pool
}

func (dataSource PostgresDataSource) prepareStatements(ctx context.Context, tx pgx.Tx) error {
	_, err := tx.Prepare(ctx, "select-transfer-id", "SELECT (id) FROM transfers WHERE id = $1")
	if err != nil {
		return err
	}

	_, err = tx.Prepare(ctx, "select-account-total", "SELECT (total) FROM accounts WHERE id = $1")
	if err != nil {
		return err
	}

	_, err = tx.Prepare(ctx, "subtract-total", "UPDATE accounts SET total = total - $2 WHERE id = $1")
	if err != nil {
		return err
	}

	_, err = tx.Prepare(ctx, "add-total", "UPDATE accounts SET total = total + $2 WHERE id = $1")
	if err != nil {
		return err
	}

	return nil
}

func (dataSource PostgresDataSource) Transfer(ctx context.Context, params TransferParams) error {
	tx, err := dataSource.Pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		} else {
			tx.Commit(ctx)
		}
	}()

	err = dataSource.prepareStatements(ctx, tx)
	if err != nil {
		return err
	}

	var existingTransferId string

	row := tx.QueryRow(ctx, "select-transfer-id", params.ID)
	err = row.Scan(&existingTransferId)
	if err == nil {
		return errors.New("Transfer with this ID already exists.")
	}

	var account entity.Account

	row = tx.QueryRow(ctx, "select-account-total", params.Sender)
	err = row.Scan(&account.Total)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("No sender account found in the database.")
		}

		return err
	}

	if account.Total < params.Amount {
		return errors.New("Not enough balance.")
	}

	_, err = tx.Exec(ctx, "subtract-total", params.Sender, params.Amount)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "add-total", params.Receiver, params.Amount)
	if err != nil {
		return err
	}

	return err
}
