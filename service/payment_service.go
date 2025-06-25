package service

import (
	"billing-engine/protobuff/pb"
	"context"
)

type PaymentService interface {
	MakePayment(ctx context.Context, req *pb.MakePaymentRequest) (*pb.MakePaymentResponse, error)
}
