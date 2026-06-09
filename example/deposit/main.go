package main

import (
	"context"
	"fmt"
	"log"
	"os"

	surepay "github.com/surepay-one/surepay-go-sdk"
)

func main() {
	client := surepay.New(
		os.Getenv("SUREPAY_API_KEY"),
		os.Getenv("SUREPAY_API_SECRET"),
	)

	ctx := surepay.WithIdempotencyKey(context.Background(), "ORD-20260101-001")

	dep, err := client.Deposits.Create(ctx, &surepay.CreateDepositRequest{
		Amount:    100_000,
		RequestID: "ORD-20260101-001",
	})
	if err != nil {
		log.Fatalf("create deposit: %v", err)
	}

	fmt.Printf("Deposit ID:   %s\n", dep.ID)
	fmt.Printf("Status:       %s\n", dep.Status)
	fmt.Printf("Checkout URL: %s\n", dep.CheckoutURL)
	fmt.Printf("Expires at:   %s\n", dep.ExpiresAt)

	// Poll or use webhook — here we just show how to Get by ID
	fetched, err := client.Deposits.Get(context.Background(), dep.ID)
	if err != nil {
		log.Fatalf("get deposit: %v", err)
	}
	fmt.Printf("Re-fetched status: %s\n", fetched.Status)
}
