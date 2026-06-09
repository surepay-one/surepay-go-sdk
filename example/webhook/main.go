package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	surepay "github.com/surepay-one/surepay-go-sdk"
)

func main() {
	secret := os.Getenv("SUREPAY_WEBHOOK_SECRET")

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "read error", http.StatusBadRequest)
			return
		}

		sig := r.Header.Get("X-Surepay-Signature")
		if !surepay.VerifyWebhookSignature(secret, body, sig) {
			http.Error(w, "invalid signature", http.StatusUnauthorized)
			return
		}

		var evt struct {
			Event string `json:"event"`
		}
		if err := json.Unmarshal(body, &evt); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}

		switch evt.Event {
		case "deposit.success", "deposit.failed":
			var dep surepay.DepositEvent
			json.Unmarshal(body, &dep)
			fmt.Printf("deposit event: id=%s status=%s\n", dep.ID, dep.Status)
		case "payout.success", "payout.failed":
			var pay surepay.PayoutEvent
			json.Unmarshal(body, &pay)
			fmt.Printf("payout event: id=%s status=%s\n", pay.PayoutID, pay.Status)
		}

		w.WriteHeader(http.StatusOK)
	})

	log.Println("listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
