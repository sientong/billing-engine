package repository

import (
	"billing-engine/helper"
	"billing-engine/model/domain"
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type UserRepositoryImpl struct{}

func NewUserRepository() UserRepository {
	return &UserRepositoryImpl{}
}

func (repository *UserRepositoryImpl) Create(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error) {

	existingUser, err := repository.FindByIdentityNumber(ctx, tx, user.IdentityNumber)
	if err != nil && err != sql.ErrNoRows {
		return domain.User{}, err
	}

	if existingUser.ID != "" {
		err := fmt.Errorf("user with identity number %s already exists", user.IdentityNumber)
		return domain.User{}, err
	}

	// If user does not exist, create a new user
	SQL := "INSERT INTO users(name, identity_number, is_delinquent, is_active) VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at"

	err = tx.QueryRowContext(ctx, SQL, user.Name, user.IdentityNumber, user.IsDelinquent, user.IsActive).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
	if err != nil {
		helper.CheckErrorOrReturn(fmt.Errorf("error inserting user: %w", err))
		return domain.User{}, err
	}

	return user, nil
}

func (repository *UserRepositoryImpl) Update(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error) {
	existingUser, err := repository.FindByIdentityNumber(ctx, tx, user.IdentityNumber)
	if err != nil && err != sql.ErrNoRows {
		helper.CheckErrorOrReturn(fmt.Errorf("error finding user by identity number: %w", err))
		return domain.User{}, err
	}

	if existingUser.ID == "" {
		helper.CheckErrorOrReturn(fmt.Errorf("user with identity number %s does not exist", user.IdentityNumber))
		return domain.User{}, nil
	}

	SQL := "UPDATE users SET name = $1, is_delinquent = $2, is_active = $3, updated_at = NOW() WHERE id = $4"

	_, err = tx.ExecContext(ctx, SQL, user.Name, user.IsDelinquent, user.IsActive, user.ID)
	if err != nil {
		helper.CheckErrorOrReturn(fmt.Errorf("error updating user: %w", err))
		return domain.User{}, err
	}

	return user, nil
}

func (repository *UserRepositoryImpl) Delete(ctx context.Context, tx *sql.Tx, user domain.User) (domain.User, error) {

	existingUser, err := repository.FindByIdentityNumber(ctx, tx, user.IdentityNumber)
	if err != nil && err != sql.ErrNoRows {
		helper.CheckErrorOrReturn(fmt.Errorf("error finding user by identity number: %w", err))
	}
	if existingUser.ID == "" {
		helper.CheckErrorOrReturn(fmt.Errorf("user with identity number %s does not exist", user.IdentityNumber))
	}

	// Soft delete the user by setting is_active to FALSE
	SQL := "UPDATE users SET is_active = FALSE, updated_at = NOW() WHERE id = $1"
	_, err = tx.ExecContext(ctx, SQL, user.ID)
	if err != nil {
		helper.CheckErrorOrReturn(fmt.Errorf("error deleting user: %w", err))
		return domain.User{}, err
	}

	return user, err
}

func (repository *UserRepositoryImpl) FindByIdentityNumber(ctx context.Context, tx *sql.Tx, identityNumber string) (domain.User, error) {
	SQL := "SELECT id, name, identity_number, is_delinquent, is_active, created_at, updated_at FROM users WHERE identity_number = $1"
	rows, err := tx.QueryContext(ctx, SQL, identityNumber)
	if err != nil {
		helper.CheckErrorOrReturn(fmt.Errorf("error finding user by identity number: %w", err))
		return domain.User{}, err
	}
	defer rows.Close()

	if rows.Next() {
		var user domain.User
		err := rows.Scan(&user.ID, &user.Name, &user.IdentityNumber, &user.IsDelinquent, &user.IsActive, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return domain.User{}, err
		}
		return user, nil
	}

	return domain.User{}, sql.ErrNoRows
}
