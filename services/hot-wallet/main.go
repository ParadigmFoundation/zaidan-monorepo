package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ParadigmFoundation/zaidan-monorepo/services/hw/zeroex"
)

func signOrder(w http.ResponseWriter, r *http.Request) {
	var order zeroex.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	hash, err := order.ComputeOrderHash()
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to compute hash, likely an invalid order: %v", err), http.StatusInternalServerError)
	}

	fmt.Fprintf(w, `{"orderHash":"%s"}`, hash.Hex())
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/order/sign", signOrder)

	err := http.ListenAndServe(":8000", mux)
	log.Fatal(err)
}
