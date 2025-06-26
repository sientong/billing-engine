package repository

import (
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type UserRepository interface {
	Create(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error)
	Update(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error)
	Delete(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error)
	FindByIdentityNumber(ctx context.Context, tx *sql.Tx, identityNumber string) (domain.User, error)
}
