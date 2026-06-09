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

	// Verify account before paying out
	result, err := client.BankInquiry.Verify(context.Background(), &surepay.BankInquiryRequest{
		BankCode:      "VCB",
		AccountNumber: "1234567890",
	})
	if err != nil {
		log.Fatalf("bank inquiry: %v", err)
	}
	fmt.Printf("Account name: %s\n", result.AccountName)

	ctx := surepay.WithIdempotencyKey(context.Background(), "PAY-20260101-001")

	pay, err := client.Payouts.Create(ctx, &surepay.CreatePayoutRequest{
		Amount:      200_000,
		BankCode:    "VCB",
		BankAccount: "1234567890",
		FullName:    result.AccountName,
		Description: "Salary January 2026",
	})
	if err != nil {
		log.Fatalf("create payout: %v", err)
	}

	fmt.Printf("Payout ID: %s\n", pay.ID)
	fmt.Printf("Status:    %s\n", pay.Status)
}
