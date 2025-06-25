package service_test

import (
	"billing-engine/model/domain"
	"billing-engine/protobuff/pb"
	"billing-engine/service"
	"context"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

var returnUser = domain.User{
	ID:             "abc123",
	Name:           "John Doe",
	IdentityNumber: "123456789",
	IsActive:       true,
	IsDelinquent:   false,
	CreatedAt:      "2024-01-01T00:00:00Z",
	UpdatedAt:      "2024-01-01T00:00:00Z",
}

var firstBillingSchedule = domain.BillingSchedule{
	ID:             "billing1",
	LoanID:         "loan1",
	PaymentDueDate: "2024-01-01T00:00:00Z",
	AmountDue:      100.0,
	Status:         "pending",
}

var secondBillingSchedule = domain.BillingSchedule{
	ID:             "billing2",
	LoanID:         "loan1",
	PaymentDueDate: "2024-01-01T00:00:00Z",
	AmountDue:      100.0,
	Status:         "pending",
}

var returnBillingSchedules = []domain.BillingSchedule{firstBillingSchedule, secondBillingSchedule}

func TestUserService_CreateNewUser(t *testing.T) {
	mockRepo := &MockUserRepository{
		ReturnUser: returnUser,
	}

	mockBillingRepo := &MockBillingScheduleRepo{
		ReturnBillingSchedules: returnBillingSchedules,
	}

	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()
	service := service.NewUserService(mockRepo, mockBillingRepo, mockDB)

	req := &pb.CreateUserRequest{
		Name:           "John Doe",
		IdentityNumber: "123456789",
	}

	resp, err := service.CreateUser(context.Background(), req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.Name != "John Doe" {
		t.Errorf("Expected name to be 'John Doe', got %s", resp.Name)
	}
	if resp.Id != "abc123" {
		t.Errorf("Expected ID to be 'abc123', got %s", resp.Id)
	}
}

func TestUserService_UpdateDelinquentStatus(t *testing.T) {
	returnUser.IsDelinquent = true

	mockBillingRepo := &MockBillingScheduleRepo{
		ReturnBillingSchedules: returnBillingSchedules,
	}

	mockRepo := &MockUserRepository{
		ReturnUser: returnUser,
	}

	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()
	service := service.NewUserService(mockRepo, mockBillingRepo, mockDB)

	req := &pb.UpdateDeliquentStatusRequest{
		IdentityNumber: "123456789",
	}

	resp, err := service.UpdateDeliquentStatus(context.Background(), req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.IdentityNumber != "123456789" {
		t.Errorf("Expected identity number to be '123456789', got %s", resp.IdentityNumber)
	}
	if resp.IsDelinquent != true {
		t.Errorf("Expected deliquent status to be true, got %t", resp.IsDelinquent)
	}
}

func TestUserService_IsDelinquent(t *testing.T) {

	mockBillingRepo := &MockBillingScheduleRepo{
		ReturnBillingSchedules: returnBillingSchedules,
	}

	mockRepo := &MockUserRepository{
		ReturnUser: returnUser,
	}

	mockDB, mock, _ := sqlmock.New()
	defer mockDB.Close()

	mock.ExpectBegin()
	mock.ExpectCommit()
	service := service.NewUserService(mockRepo, mockBillingRepo, mockDB)

	req := &pb.GetDeliquencyStatusRequest{
		IdentityNumber: "123456789",
	}

	resp, err := service.IsDelinquent(context.Background(), req)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if resp.IsDelinquent != false {
		t.Errorf("Expected deliquent status to be false, got %t", resp.IsDelinquent)
	}
}
