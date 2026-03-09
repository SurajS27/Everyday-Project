#!/bin/bash

# Test script for self-healing mechanism

# NOTE: This assumes the 'my-self-healing-app' is deployed in the 'dev' namespace.
APP_NAMESPACE="dev"

echo "🛡️ Testing self-healing mechanism..."

# Scale down the deployment to 0 replicas
echo "📉 Scaling deployment to 0 replicas..."
kubectl scale deployment my-self-healing-app -n ${APP_NAMESPACE} --replicas=0

# Watch pods for 2 minutes to see ArgoCD restore the deployment
echo "👀 Watching pods for self-healing (2 minutes)..."
timeout 120 kubectl get pods -n ${APP_NAMESPACE} -w

echo "✅ Self-healing test complete!"
echo "Expected result: Deployment should be restored to 1 replica (dev environment)"