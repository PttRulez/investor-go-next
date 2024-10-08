package postgres

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go/internal/domain"
	"github.com/pttrulez/investor-go/internal/infrastructure/storage"
)

func (pg *Repository) DeleteTransaction(ctx context.Context, id int, userID int) error {
	const op = "Repository.DeleteTransaction"

	queryString := `DELETE FROM transactions WHERE id = $1 AND user_id = $2;`
	result, err := pg.db.ExecContext(ctx, queryString, id, userID)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	if rowsAffected == 0 {
		return storage.ErrNotFound
	}

	return nil
}

func (pg *Repository) GetTransaction(ctx context.Context, id int, userID int) (domain.Transaction, error) {
	const op = "Repository.GetTransaction"

	queryString := `SELECT * FROM transactions WHERE id = $1 AND user_id = $2;`
	row := pg.db.QueryRowContext(ctx, queryString, id, userID)

	var t domain.Transaction
	err := row.Scan(&t.ID, &t.Amount, &t.Date, &t.PortfolioID, &t.Type)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("%s: %w", op, err)
	}

	return t, nil
}

func (pg *Repository) GetTransactionList(ctx context.Context, id int, userID int) (
	[]domain.Transaction, error) {
	const op = "Repository.GetPortfolioTransactionList"

	queryString := `SELECT * FROM transactions WHERE portfolio_id = $1 AND user_id = $2;`
	rows, err := pg.db.QueryContext(ctx, queryString, id, userID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var ts []domain.Transaction
	for rows.Next() {
		var t domain.Transaction
		err = rows.Scan(&t.ID, &t.Amount, &t.Date, &t.PortfolioID, &t.Type, &t.UserID)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		ts = append(ts, t)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return ts, nil
}

func (pg *Repository) InsertTransaction(ctx context.Context, t domain.Transaction) (domain.Transaction, error) {
	const op = "Repository.InsertTransaction"

	queryString := `INSERT INTO transactions (amount, date, portfolio_id, type, user_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, amount, date, portfolio_id, type;`

	var res domain.Transaction
	err := pg.db.QueryRowContext(ctx, queryString, t.Amount, t.Date, t.PortfolioID, t.Type,
		t.UserID).
		Scan(&res.ID, &res.Amount, &res.Date, &res.PortfolioID, &res.Type)
	if err != nil {
		return domain.Transaction{}, fmt.Errorf("%s: %w", op, err)
	}
	return res, nil
}
