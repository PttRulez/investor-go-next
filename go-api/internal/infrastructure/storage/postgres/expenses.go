package postgres

import (
	"context"
	"fmt"

	"github.com/pttrulez/investor-go-next/go-api/internal/domain"
	"github.com/pttrulez/investor-go-next/go-api/internal/infrastructure/storage"
)

func (pg *Repository) DeleteExpense(ctx context.Context, id int, userID int) error {
	const op = "Repository.DeleteExpense"

	queryString := `DELETE FROM expenses WHERE id = $1 AND user_id = $2;`
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

func (pg *Repository) GetExpenseList(ctx context.Context, portfolioID int) (
	[]domain.Expense, error) {
	const op = "Repository.GetExpenseList"

	queryString := `SELECT amount, date, description, id, portfolio_id FROM expenses
		WHERE portfolio_id = $1 ORDER BY date DESC, id DESC;`

	rows, err := pg.db.QueryContext(ctx, queryString, portfolioID)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}
	defer rows.Close()

	var expenses []domain.Expense
	for rows.Next() {
		var e domain.Expense
		err := rows.Scan(
			&e.Amount,
			&e.Date,
			&e.Description,
			&e.ID,
			&e.PortfolioID,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", op, err)
		}
		expenses = append(expenses, e)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", op, rows.Err())
	}

	return expenses, nil
}

func (pg *Repository) InsertExpense(ctx context.Context, c domain.Expense,
	userID int) error {
	const op = "Repository.InsertExpense"

	queryString := `INSERT INTO expenses (amount, date, description, portfolio_id, user_id)
    VALUES ($1, $2, $3, $4, $5);`

	_, err := pg.db.ExecContext(ctx, queryString,
		c.Amount,
		c.Date,
		c.Description,
		c.PortfolioID,
		userID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
