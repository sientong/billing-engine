package repository

import (
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user domain.User) domain.User
	Update(ctx context.Context, tx *sql.Tx, user domain.User) domain.User
	Delete(ctx context.Context, tx *sql.Tx, user domain.User) domain.User
	FindByIdentityNumber(ctx context.Context, tx *sql.Tx, identityNumber string) (domain.User, error)
}
