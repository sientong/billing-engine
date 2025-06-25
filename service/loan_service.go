package service

import (
	pb "billing-engine/protobuff/pb"
	"context"
)

type LoanService interface {
	CreateNewLoan(ctx context.Context, req *pb.CreateNewLoanRequest) (*pb.LoanResponse, error)
	GetOutstanding(ctx context.Context, req *pb.GetOutstandingRequest) (*pb.OutstandingResponse, error)
}
