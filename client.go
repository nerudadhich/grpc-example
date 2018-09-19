package main

import (
	"log"
	"os"
	"strconv"
	"time"

	pb "protoc"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBankClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	operation := os.Args[1]

	switch operation {
	case "balance":
		account := os.Args[2]
		r, err := c.GetBalance(ctx, &pb.CheckBalanceRequest{Account: account})
		if err != nil {
			log.Fatalf("Error : %v", err)
		}
		log.Printf("Balance: %v", r.Balance)
	case "deposite":
		account := os.Args[2]
		amount, err := strconv.ParseFloat(os.Args[3], 32)
		r, err := c.Deposite(ctx, &pb.DepositeRequest{Account: account, Amount: float32(amount)})
		if err != nil {
			log.Fatalf("could not get")
		}
		log.Printf("Deposited : %v", r.Status)
	case "withdraw":
		account := os.Args[2]
		amount, err := strconv.ParseFloat(os.Args[3], 32)
		r, err := c.Withdraw(ctx, &pb.WithdrawRequest{Account: account, Amount: float32(amount)})
		if err != nil {
			log.Fatalf("could not get")
		}
		log.Printf("Withdrawn : %v and Remaining Balance: %v", r.Status, r.Balance)
	}
}
