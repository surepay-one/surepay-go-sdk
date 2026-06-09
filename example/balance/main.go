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

	bal, err := client.Balance.Get(context.Background())
	if err != nil {
		log.Fatalf("get balance: %v", err)
	}

	fmt.Printf("Balance:   %d %s\n", bal.Balance, bal.Currency)
	fmt.Printf("Hold:      %d %s\n", bal.Hold, bal.Currency)
	fmt.Printf("Available: %d %s\n", bal.Available, bal.Currency)
}
