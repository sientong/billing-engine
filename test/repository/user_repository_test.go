package test

import (
	"billing-engine/repository"
	"billing-engine/test"
	"context"
	"database/sql"
	"testing"
)

func TestUserRepo_CreateUserWithCorrectData(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Unexpected panic occurred: %v", r)
		}
	}()

	db := test.SetupTestDB()
	defer db.Close()

	result, err := test.CreateUser(db)
	if err != nil {
		t.Errorf("Unexpected error occurred: %v", err)
	}

	if result.ID == "" {
		t.Error("Expected user ID to be set, got empty string")
	}
	if result.CreatedAt == "" {
		t.Error("Expected user CreatedAt to be set, got empty string")
	}
	if result.UpdatedAt == "" {
		t.Error("Expected user UpdatedAt to be set, got empty string")
	}
	if !result.IsActive {
		t.Error("Expected user IsActive to be true, got false")
	}
	if result.IsDelinquent {
		t.Error("Expected user IsDelinquent to be false, got true")
	}
	if result.Name != "John Doe" {
		t.Errorf("Expected user Name to be 'John Doe', got '%s'", result.Name)
	}
	if result.IdentityNumber != "123456789" {
		t.Errorf("Expected user IdentityNumber to be '123456789', got '%s'", result.IdentityNumber)
	}

	test.TruncateUsersTable(db)

}

func TestUserRepo_CreateUserWithExistingUser(t *testing.T) {
	db := test.SetupTestDB()
	defer db.Close()

	_, err := test.CreateUser(db)
	if err != nil {
		t.Errorf("Unexpected error occurred: %v", err)
	}

	// Attempt to create the same user again
	test.ExpectPanicMessage(t, func() {
		_, err = test.CreateUser(db)
	}, "user with identity number 123456789 already exists")

	test.TruncateUsersTable(db)
}

func TestUserRepo_FindUserByIdentityNumber(t *testing.T) {
	db := test.SetupTestDB()
	defer db.Close()

	user, err := test.CreateUser(db)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()
	repo := repository.NewUserRepository()
	foundUser, err := repo.FindByIdentityNumber(ctx, tx, user.IdentityNumber)
	if err != nil {
		t.Fatalf("Failed to find user by identity number: %v", err)
	}
	if foundUser.ID != user.ID {
		t.Errorf("Expected found user ID to be %s, got %s", user.ID, foundUser.ID)
	}
}

func TestUserRepo_UpdateUser(t *testing.T) {
	db := test.SetupTestDB()
	defer db.Close()
	test.TruncateUsersTable(db)

	// Create a user to update
	user, err := test.CreateUser(db)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	// Update the user's name
	user.Name = "Jane Doe"
	user.IsDelinquent = true

	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	defer tx.Rollback()

	repo := repository.NewUserRepository()
	updatedUser := repo.Update(ctx, tx, user)

	if updatedUser.Name != "Jane Doe" {
		t.Errorf("Expected updated user Name to be 'Jane Doe', got '%s'", updatedUser.Name)
	}
	if !updatedUser.IsDelinquent {
		t.Error("Expected updated user IsDelinquent to be true, got false")
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	test.TruncateUsersTable(db)
}

func TestUserRepo_DeleteUser(t *testing.T) {
	db := test.SetupTestDB()
	defer db.Close()

	test.TruncateUsersTable(db)
	// Create a user to delete
	user, err := test.CreateUser(db)
	if err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}
	ctx := context.Background()
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("Failed to begin transaction: %v", err)
	}
	repo := repository.NewUserRepository()
	// Delete the user
	repo.Delete(ctx, tx, user)

	// Verify the user is deleted
	user, err = repo.FindByIdentityNumber(ctx, tx, user.IdentityNumber)
	if err == nil && user.IsActive {
		t.Errorf("Expected user with identity number %s to be deleted, but it is still active", user.IdentityNumber)
	} else if err == sql.ErrNoRows {
		t.Errorf("Unexpected error occurred while finding deleted user: %v", err)
	}

	if err := tx.Commit(); err != nil {
		t.Fatalf("Failed to commit transaction: %v", err)
	}

	test.TruncateUsersTable(db)
}
