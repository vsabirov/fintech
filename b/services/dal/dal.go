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

type RefundParams struct {
	TransferID string
	Amount     float64

	Sender   string
	Receiver string
}

type DataSource interface {
	prepareStatements(ctx context.Context, tx pgx.Tx) error

	Transfer(ctx context.Context, params TransferParams) error
	Refund(ctx context.Context, params RefundParams) error
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

	_, err = tx.Prepare(ctx, "select-account-id", "SELECT (id) FROM accounts WHERE id = $1")
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

	_, err = tx.Prepare(ctx, "insert-transfer", "INSERT INTO transfers (id, amount, sender, receiver) VALUES ($1, $2, $3, $4)")
	if err != nil {
		return err
	}

	_, err = tx.Prepare(ctx, "remove-transfer", "DELETE FROM transfers WHERE id = $1")

	return err
}

func (dataSource PostgresDataSource) Refund(ctx context.Context, params RefundParams) error {
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

	var account entity.Account

	row := tx.QueryRow(ctx, "select-account-id", params.Sender)
	err = row.Scan(&account.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("No sender account found in the database.")
		}

		return err
	}

	row = tx.QueryRow(ctx, "select-account-total", params.Receiver)
	err = row.Scan(&account.Total)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("No receiver account found in the database.")
		}

		return err
	}

	if account.Total < params.Amount {
		return errors.New("Not enough balance.")
	}

	_, err = tx.Exec(ctx, "subtract-total", params.Receiver, params.Amount)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "add-total", params.Sender, params.Amount)
	if err != nil {
		return err
	}

	_, err = tx.Exec(ctx, "remove-transfer", params.TransferID)

	return err
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

	row = tx.QueryRow(ctx, "select-account-id", params.Receiver)
	err = row.Scan(&account.ID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return errors.New("No receiver account found in the database.")
		}

		return err
	}

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

	_, err = tx.Exec(ctx, "insert-transfer", params.ID, params.Amount, params.Sender, params.Receiver)

	return err
}
