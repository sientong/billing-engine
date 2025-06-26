package main

import (
	app "billing-engine/app"
	"billing-engine/protobuff/pb"
	"billing-engine/repository"
	"billing-engine/service"
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Step 1: Setup DB
	db := app.SetupDB()

	// Step 2: Setup cancellable context
	_, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Step 3: Start TCP listener
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Step 4: Init gRPC server
	grpcServer := grpc.NewServer()

	// Step 5: Init repositories and services
	userRepo := repository.NewUserRepository()
	loanRepo := repository.NewLoanRepository()
	paymentRepo := repository.NewPaymentRepository()
	billingRepo := repository.NewBillingScheduleRepo()

	userService := service.NewUserService(userRepo, billingRepo, db)
	loanService := service.NewLoanService(loanRepo, billingRepo, db)
	paymentService := service.NewPaymentService(paymentRepo, loanRepo, billingRepo, db)

	// Step 6: Register services
	pb.RegisterUserServiceServer(grpcServer, userService)
	pb.RegisterLoanServiceServer(grpcServer, loanService)
	pb.RegisterPaymentServiceServer(grpcServer, paymentService)

	reflection.Register(grpcServer)

	// Step 7: Start gRPC server in goroutine
	go func() {
		log.Println("gRPC server started on :50051")
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	// Step 8: Listen for shutdown signal
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop // block until signal received
	log.Println("Shutdown signal received")

	// Step 9: Gracefully stop gRPC server
	grpcServer.GracefulStop()
	log.Println("gRPC server gracefully stopped")

	// Step 10: Clean up resources
	cancel()
	_ = db.Close()
	log.Println("Resources cleaned up")
}
