package repository

import (
	"context"
)

func (r *Repository) IncrementLoginCount(ctx context.Context, userID string) error {
	query := `
		UPDATE users
		SET login_count = login_count + 1,
		updated_at = now()
		WHERE id = $1
	`
	_, err := r.Db.ExecContext(ctx, query, userID)
	if err != nil {
		return err
	}
	return nil
}
