package main

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	// version distinguishes between v1 and v2 for our Canary testing
	version             = os.Getenv("APP_VERSION") 
	inventoryServiceURL = os.Getenv("INVENTORY_SERVICE_URL")
)

type OrderRequest struct {
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

func init() {
	if version == "" {
		version = "v1" // Default to v1
	}
	if inventoryServiceURL == "" {
		inventoryServiceURL = "http://inventory-service:8080"
	}
}

func main() {
	http.HandleFunc("/order", handleOrder)
	
	// Health check endpoint for Kubernetes probes
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("[Order-Service %s] Listening on port %s\n", version, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleOrder(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var orderReq OrderRequest
	err := json.NewDecoder(r.Body).Decode(&orderReq)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("[Order-Service %s] Received order for %d x %s\n", version, orderReq.Quantity, orderReq.Item)

	// Forward the request downstream to the Inventory Service (Tier 3)
	inventoryReqBody, _ := json.Marshal(map[string]interface{}{
		"item":   orderReq.Item,
		"action": "reserve",
		"amount": orderReq.Quantity,
	})

	resp, err := http.Post(inventoryServiceURL+"/reserve", "application/json", bytes.NewBuffer(inventoryReqBody))
	if err != nil {
		log.Printf("[Order-Service %s] Error calling Inventory Service: %v\n", version, err)
		// Returning 503 triggers Istio's OutlierDetection (Circuit Breaker) if it happens frequently
		http.Error(w, "Failed to reach Inventory Service", http.StatusServiceUnavailable)
		return
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)

	// If Inventory fails, propagate the error upward
	if resp.StatusCode != http.StatusOK {
		log.Printf("[Order-Service %s] Inventory Service returned status %d\n", version, resp.StatusCode)
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
		return
	}

	// Success response containing the version, proving which Canary received the traffic
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "Order Processed",
		"processed_by_version": version,
		"inventory_response": json.RawMessage(respBody),
	})
}