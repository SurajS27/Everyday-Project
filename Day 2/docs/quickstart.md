# Quick Start Guide

## Prerequisites

- Docker
- kubectl
- kind (Kubernetes in Docker)

## Step 1: Create Kind Cluster

```bash
kind create cluster --config cluster/kind-config.yaml
```

## Step 2: Bootstrap the Platform

```bash
./scripts/bootstrap.sh
```

This will:
1. Install ArgoCD
2. Deploy the root application
3. Set up all components automatically

## Step 3: Access the UI

### ArgoCD UI
```bash
kubectl port-forward svc/argocd-server -n argocd 8080:443
```
Open: https://localhost:8080
Password: `kubectl -n argocd get secret argocd-initial-admin-secret -o jsonpath="{.data.password}" | base64 -d`

### Grafana
```bash
kubectl port-forward svc/grafana -n monitoring 3000:3000
```
Open: http://localhost:3000
Username: admin
Password: admin

## Step 4: Test Self-Healing

```bash
./scripts/test-self-healing.sh
```

## Architecture Overview

- **ArgoCD**: GitOps controller
- **Prometheus**: Metrics collection
- **Grafana**: Visualization
- **Loki**: Log aggregation
- **Promtail**: Log shipping
- **My Self-Healing App**: Sample application with dev/prod overlays