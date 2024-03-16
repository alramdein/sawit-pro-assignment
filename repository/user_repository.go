package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
)

// InsertUser insert new user
func (r *Repository) InsertUser(ctx context.Context, input InsertUserInput) (err error) {
	query := "INSERT INTO users (id, name, phone_number, password) VALUES ($1, $2, $3, $4)"
	fields := []interface{}{
		input.Id, input.Name, input.Phone, input.Password,
	}
	_, err = r.Db.ExecContext(ctx, query, fields...)
	if err != nil {
		return err
	}

	return nil
}

// GetUserByID retrieves user profile details from the database by user ID
func (r *Repository) GetUserByID(ctx context.Context, userID string) (*model.User, error) {
	var user model.User
	err := r.Db.QueryRowContext(ctx, "SELECT name, phone_number FROM users WHERE id = $1", userID).Scan(&user.Name, &user.Phone)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByPhone retrieves user profile details from the database by user phone
func (r *Repository) GetUserByPhone(ctx context.Context, phone string) (*model.User, error) {
	var user model.User
	query := `
		SELECT id, name, phone_number, password
		FROM users
		WHERE phone_number = $1
	`
	err := r.Db.QueryRowContext(ctx, query, phone).Scan(&user.ID, &user.Name, &user.Phone, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // User not found
		}
		return nil, err
	}
	return &user, nil
}

// UpdateUserProfile updates the user profile in the database
func (r *Repository) UpdateUser(ctx context.Context, userID string, input generated.UpdateProfileJSONBody) error {
	query, data := buildUpdateProfileQuery(userID, input)
	_, err := r.Db.ExecContext(ctx, query, data...)
	return err
}

func buildUpdateProfileQuery(userId string, in generated.UpdateProfileJSONBody) (string, []interface{}) {
	var fields, args []string
	var data []interface{}
	if in.Name != nil {
		fields = append(fields, "name")
		data = append(data, in.Name)
	}
	if in.Phone != nil {
		fields = append(fields, "phone_number")
		data = append(data, in.Phone)
	}

	fields = append(fields, "updated_at")
	data = append(data, time.Now())

	for i, v := range fields {
		args = append(args, fmt.Sprintf("%s=$%d", v, i+1))
	}

	data = append(data, userId)
	query := "UPDATE users SET " + strings.Join(args, ",") + " WHERE id = $" + fmt.Sprintf("%d", len(fields)+1)

	return query, data
}
