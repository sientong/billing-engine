package service_test

import (
	"billing-engine/model/domain"
	"context"
	"database/sql"
)

type MockUserRepository struct {
	CalledWith domain.User
	ReturnUser domain.User
}

func (m *MockUserRepository) Create(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
	m.CalledWith = user
	return m.ReturnUser
}

func (m *MockUserRepository) Update(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
	m.CalledWith = user
	return m.ReturnUser
}

func (m *MockUserRepository) Delete(ctx context.Context, tx *sql.Tx, user domain.User) domain.User {
	m.CalledWith = user
	return m.ReturnUser
}

func (m *MockUserRepository) FindByIdentityNumber(ctx context.Context, tx *sql.Tx, identityNumber string) (domain.User, error) {
	return m.ReturnUser, nil
}
