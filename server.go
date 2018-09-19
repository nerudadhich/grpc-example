package main

import (
	"log"
	"net"

	pb "./protoc"
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	port = ":50051"
)

var bankAccounts map[string]float32

type server struct{}

func (s *server) GetBalance(ctx context.Context, in *pb.CheckBalanceRequest) (*pb.CheckBalanceResponse, error) {
	balance, exists := bankAccounts[in.Account]
	if exists {
		return &pb.CheckBalanceResponse{Balance: balance}, nil
	}
	return nil, errors.New("Account not available")
}

func (s *server) Deposite(ctx context.Context, in *pb.DepositeRequest) (*pb.DipositeResponse, error) {
	balance, exists := bankAccounts[in.Account]
	if exists {
		bankAccounts[in.Account] = balance + in.Amount
		return &pb.DipositeResponse{Status: true}, nil
	}
	return nil, errors.New("Account not available")
}

func (s *server) Withdraw(ctx context.Context, in *pb.WithdrawRequest) (*pb.WithdrawResponse, error) {
	balance, exists := bankAccounts[in.Account]
	if exists {
		bankAccounts[in.Account] = balance - in.Amount
		return &pb.WithdrawResponse{Status: true, Balance: bankAccounts[in.Account]}, nil
	}
	return nil, errors.New("Account not available")
}

func main() {
	bankAccounts = make(map[string]float32)
	bankAccounts["101"] = 2000.00
	bankAccounts["102"] = 2000.00

	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBankServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
