const express = require('express');
const axios = require('axios');

const app = express();
const PORT = process.env.PORT || 3000;

// The URL to the Order Service (Tier 2). 
// In Kubernetes, this resolves via internal DNS to the K8s Service name.
const ORDER_SERVICE_URL = process.env.ORDER_SERVICE_URL || 'http://order-service:8080';

app.use(express.json());

// Main entry point for external traffic (Serves a simple UI)
app.get('/', (req, res) => {
    res.send(`
        <html>
        <head>
            <title>Zero-Trust Mesh Store</title>
            <style>
                body { font-family: Arial, sans-serif; padding: 2rem; }
                button { padding: 10px 20px; font-size: 16px; cursor: pointer; }
                pre { background: #eee; padding: 15px; border-radius: 5px; }
            </style>
        </head>
        <body>
            <h1>Zero-Trust Mesh Store</h1>
            <p>Welcome to the Web Frontend (Tier 1). Traffic from this UI will route to Tier 2.</p>
            <button onclick="placeOrder()">Place an Order</button>
            <pre id="result">Waiting for action...</pre>

            <script>
                async function placeOrder() {
                    document.getElementById('result').innerText = 'Forwarding request to Order Service...';
                    try {
                        const response = await fetch('/api/checkout', { method: 'POST' });
                        const data = await response.json();
                        document.getElementById('result').innerText = JSON.stringify(data, null, 2);
                    } catch (error) {
                        document.getElementById('result').innerText = 'Error: ' + error.message;
                    }
                }
            </script>
        </body>
        </html>
    `);
});

// API Gateway route that forwards requests downstream to the Order Service
app.post('/api/checkout', async (req, res) => {
    console.log('[Web-Frontend] Received checkout request. Forwarding to Order Service...');
    
    try {
        // Forward the request to the Order Service
        const response = await axios.post(`${ORDER_SERVICE_URL}/order`, {
            item: "Mesh-T-Shirt",
            quantity: 1
        });
        
        console.log(`[Web-Frontend] Order Service responded with status: ${response.status}`);
        res.json({
            status: "Success",
            gateway: "Web Frontend (Tier 1)",
            downstreamResponse: response.data
        });
        
    } catch (error) {
        console.error('[Web-Frontend] Error calling Order Service:', error.message);
        
        // Return details if it fails. This is helpful for seeing Circuit Breaker or 5xx errors.
        const status = error.response ? error.response.status : 503;
        res.status(status).json({
            error: "Failed to process checkout downstream",
            details: error.message
        });
    }
});

// Health check endpoint for Kubernetes Liveness/Readiness probes
app.get('/health', (req, res) => {
    res.status(200).send('OK');
});

app.listen(PORT, () => {
    console.log(`[Web-Frontend] Listening on port ${PORT}`);
});