from flask import Flask, request, jsonify
import os
import random
import time
import sys

app = Flask(__name__)

# Configuration variables for testing Istio Circuit Breaking
# e.g., Set FAILURE_RATE=0.5 to make 50% of requests return a 503 error
FAILURE_RATE = float(os.environ.get('FAILURE_RATE', '0.0')) 
SIMULATE_DELAY = float(os.environ.get('SIMULATE_DELAY', '0.0'))

@app.route('/reserve', methods=['POST'])
def reserve_inventory():
    print("[Inventory-Service] Received reservation request.", file=sys.stderr)
    
    # 1. Simulate Latency (Tests Istio timeout configurations)
    if SIMULATE_DELAY > 0:
        time.sleep(SIMULATE_DELAY)
        
    # 2. Simulate Random Failures (Tests Istio 5xx Outlier Detection)
    if random.random() < FAILURE_RATE:
        print(f"[Inventory-Service] Simulating failure! (Rate: {FAILURE_RATE})", file=sys.stderr)
        return jsonify({"error": "Simulated Database Connection Failure"}), 503
        
    # 3. Normal Processing
    data = request.json or {}
    item = data.get('item', 'Unknown Item')
    amount = data.get('amount', 1)
    
    print(f"[Inventory-Service] Successfully reserved {amount}x {item}", file=sys.stderr)
    
    return jsonify({
        "status": "Inventory Reserved",
        "item": item,
        "amount": amount,
        "backend": "Python/Flask (Tier 3)"
    }), 200

@app.route('/health', methods=['GET'])
def health_check():
    # Kubernetes readiness/liveness probe endpoint
    return "OK", 200

if __name__ == '__main__':
    port = int(os.environ.get('PORT', 8080))
    print(f"[Inventory-Service] Starting on port {port}...", file=sys.stderr)
    app.run(host='0.0.0.0', port=port)