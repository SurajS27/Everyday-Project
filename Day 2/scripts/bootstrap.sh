#!/bin/bash

# Bootstrap script for GitOps platform

# Exit on error
set -e

echo "🚀 Bootstrapping GitOps platform..."

# Step 1: Create ArgoCD namespace
echo "📦 Creating ArgoCD namespace..."
kubectl apply -f - <<EOF
apiVersion: v1
kind: Namespace
metadata:
  name: argocd
EOF

# Step 2: Install ArgoCD
echo "🐙 Installing ArgoCD..."
kubectl apply -n argocd -f https://raw.githubusercontent.com/argoproj/argo-cd/stable/manifests/install.yaml

# Wait for ArgoCD to be ready
echo "⏳ Waiting for ArgoCD to be ready..."
kubectl wait --for=condition=available --timeout=300s deployment/argocd-server -n argocd

# Step 3: Deploy root application
echo "📋 Deploying root application..."
kubectl apply -f infrastructure/argocd/root-app.yaml

echo "✅ Bootstrap complete!"
echo
echo "🔑 ArgoCD UI can be accessed by port-forwarding:"
echo "   kubectl port-forward svc/argocd-server -n argocd 8080:443"
echo "   Then open https://localhost:8080"
echo
ARGO_PASS=$(kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d)
echo "👤 ArgoCD admin username: admin"
echo "🔐 ArgoCD admin password: $ARGO_PASS"