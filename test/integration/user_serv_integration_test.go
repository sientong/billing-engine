package integration_test

import (
	"billing-engine/protobuff/pb"
	"context"
	"database/sql"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"billing-engine/repository"
	"billing-engine/service"
	"billing-engine/test"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
)

var client pb.UserServiceClient
var loanClient pb.LoanServiceClient
var paymentClient pb.PaymentServiceClient
var db *sql.DB
var grpcServer *grpc.Server

func TestMain(m *testing.M) {

	// set up database
	db = test.SetupTestDB()
	defer db.Close()

	// set up gRPC server
	lis, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer = grpc.NewServer()
	userRepo := repository.NewUserRepository()
	billingRepo := repository.NewBillingScheduleRepo()
	paymentRepo := repository.NewPaymentRepository()
	loanRepo := repository.NewLoanRepository()

	userService := service.NewUserService(userRepo, billingRepo, db)
	loanService := service.NewLoanService(loanRepo, billingRepo, db)
	paymentService := service.NewPaymentService(paymentRepo, loanRepo, billingRepo, db)

	pb.RegisterUserServiceServer(grpcServer, userService)
	pb.RegisterLoanServiceServer(grpcServer, loanService)
	pb.RegisterPaymentServiceServer(grpcServer, paymentService)

	go grpcServer.Serve(lis)

	// gRPC client
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(3*time.Second))
	if err != nil {
		log.Fatalf("Failed to dial gRPC: %v", err)
	}
	defer conn.Close()

	client = pb.NewUserServiceClient(conn)
	loanClient = pb.NewLoanServiceClient(conn)
	paymentClient = pb.NewPaymentServiceClient(conn)

	// run tests
	code := m.Run()

	grpcServer.Stop()
	db.Close()
	os.Exit(code)
}

func CreateUser() *pb.CreateUserRequest {
	return &pb.CreateUserRequest{
		Name:           "John Doe",
		IdentityNumber: "INTEG-123",
	}
}

func TestUserIntegration_Create(t *testing.T) {
	ctx := context.Background()

	req := CreateUser()

	res, err := client.CreateUser(ctx, req)
	require.NoError(t, err)
	require.Equal(t, "John Doe", res.Name)
	require.Equal(t, "INTEG-123", res.IdentityNumber)
	require.True(t, res.IsActive)
	require.False(t, res.IsDelinquent)
	require.NotEmpty(t, res.Id)

	test.TruncateAllTables(db)
}

func TestUserIntegration_CheckDeliquencyStatus(t *testing.T) {
	ctx := context.Background()

	req := CreateUser()

	_, err := client.CreateUser(ctx, req)
	require.NoError(t, err)

	secondReq := &pb.GetDeliquencyStatusRequest{
		IdentityNumber: "INTEG-123",
	}

	res, err := client.IsDelinquent(ctx, secondReq)
	require.NoError(t, err)
	require.False(t, res.IsDelinquent)

	test.TruncateAllTables(db)
}

func TestUserIntegration_UpgradeDeliquencyStatus(t *testing.T) {
	ctx := context.Background()

	req := CreateUser()

	_, err := client.CreateUser(ctx, req)
	require.NoError(t, err)

	secondReq := &pb.UpdateDeliquentStatusRequest{
		IdentityNumber: "INTEG-123",
	}

	res, err := client.UpdateDeliquentStatus(ctx, secondReq)
	require.NoError(t, err)
	require.False(t, res.IsDelinquent)

	test.TruncateAllTables(db)
}
