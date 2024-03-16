// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"

	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/model"
)

type RepositoryInterface interface {
	InsertUser(ctx context.Context, in InsertUserInput) error
	IncrementLoginCount(ctx context.Context, userID string) error
	GetUserByID(ctx context.Context, userID string) (*model.User, error)
	GetUserByPhone(ctx context.Context, phoneNumber string) (*model.User, error)
	UpdateUser(ctx context.Context, userID string, input generated.UpdateProfileJSONBody) error
}
