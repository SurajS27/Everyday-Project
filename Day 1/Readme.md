# Zero-Trust Multi-Service Mesh Demo

This project demonstrates a fully functional, 3-tier microservice architecture deployed on Kubernetes using Istio as a Service Mesh. It implements enterprise-grade traffic management, zero-trust security, and a complete observability stack.

## Table of Contents

- [🌟 Key Features Implemented](#-key-features-implemented)
- [🏗️ Architecture Overview](#️-architecture-overview)
- [📁 Project Structure](#-project-structure)
- [📋 Prerequisites](#-prerequisites)
- [🚀 Step-by-Step Implementation Guide](#-step-by-step-implementation-guide)
- [🎮 How to Use and Test the Mesh](#-how-to-use-and-test-the-mesh)
- [🧹 Cleanup](#-cleanup)

Zero-Trust Security (mTLS): All pod-to-pod communication is strictly encrypted and mutually authenticated. Unencrypted traffic is rejected at the proxy level.

Canary Deployments: Traffic from the Tier 1 Gateway is dynamically split: 90% is routed to the v1 Order Service, and 10% is routed to the v2 Canary release.

Circuit Breaking: The Tier 3 Inventory Service is protected by connection pool limits and Outlier Detection. If it returns consecutive 5xx errors, Istio automatically ejects the failing pod from the load balancing pool to prevent cascading failures.

Complete Observability (PLG Stack): Integrated Prometheus (Metrics), Loki (Logs), Promtail (Log Forwarding), and Grafana (Visualization) to monitor mesh traffic and application health.

🏗️ Architecture Overview

Tier 1: Web Frontend (Node.js) - Acts as the API Gateway and serves a simple UI.

Tier 2: Order Service (Go) - Processes orders and demonstrates Istio traffic splitting (v1 and v2).

Tier 3: Inventory Service (Python) - Backend service equipped with artificial failure injection to test Istio Circuit Breaking.

Service Mesh (Istio) - Envoy sidecars are injected into every pod to intercept and manage all network traffic.

## 📁 Project Structure

```
README.md
apps/
	inventory-service/
		Dockerfile
		requirements.txt
		src/
			app.py
	order-service/
		Dockerfile
		go.mod
		main.go
	web-frontend/
		Dockerfile
		package.json
		src/
			server.js
deploy/
	base/
		inventory-service.yaml
		namespaces.yaml
		order-service-v1.yaml
		order-service-v2.yaml
		web-frontend.yaml
	istio/
		canary-routing.yaml
		circuit-breaker.yaml
		gateway.yaml
		peer-auth.yaml
	observability/
		Grafana.yaml
		loki.yaml
		Prometheus.yaml
		Promtail.yaml
```

## 📋 Prerequisites

To run this project locally, ensure you have the following installed on your machine:

Docker Desktop (Make sure the Docker daemon is running)

Minikube (Local Kubernetes cluster)

kubectl (Kubernetes command-line tool)

istioctl (Istio command-line tool)

PowerShell (The commands below are tailored for a Windows environment)

🚀 Step-by-Step Implementation Guide

Follow these steps from the root of your project directory to deploy the complete stack.

1. Start the Local Kubernetes Cluster

Istio and the PLG stack require adequate resources. Start Minikube with at least 8GB of RAM and 4 CPUs:

```powershell
minikube start --memory=8192 --cpus=4
```

2. Install the Istio Service Mesh

Install Istio using the demo profile, which includes the necessary Ingress Gateway:

```powershell
istioctl install --set profile=demo -y
```

3. Build Docker Images Locally

Point your local Docker CLI to Minikube's internal Docker registry, then build the application images. This ensures Kubernetes can find the images without needing a remote registry like DockerHub.

```powershell
# Point PowerShell to Minikube's Docker daemon
minikube docker-env | Invoke-Expression

# Build the microservices
docker build -t your-registry/web-frontend:latest ./apps/web-frontend
docker build -t your-registry/order-service:latest ./apps/order-service
docker build -t your-registry/inventory-service:latest ./apps/inventory-service
```

4. Deploy Base Applications

First, create the namespace to enable Istio's automatic sidecar injection, then deploy the apps:

```powershell
# Create namespace and enable sidecar injection
kubectl apply -f deploy/base/namespaces.yaml

# Deploy all applications
kubectl apply -f deploy/base/
```

Tip: Wait until all pods show 2/2 under the READY column by running kubectl get pods -n zero-trust-mesh -w.

5. Apply Istio Mesh Policies

Apply the Custom Resource Definitions (CRDs) that enforce mTLS, Canary routing, and Circuit Breaking:

```powershell
kubectl apply -f deploy/istio/
```

6. Deploy the Observability Stack

Spin up Prometheus, Loki, Promtail, and Grafana:

```powershell
kubectl apply -f deploy/observability/
```

🎮 How to Use and Test the Mesh

Accessing the Web Frontend & Testing Canary Routing

To access the application through the Istio Ingress Gateway, open a new PowerShell window and run:

```powershell
minikube tunnel
```


Leave this terminal running, open your web browser, and navigate to: http://localhost

Test the Canary Split: Click the "Place an Order" button multiple times. In the JSON response, you will see the processed_by_version field change between v1 (roughly 90% of the time) and v2 (roughly 10% of the time).

Accessing Grafana (Logs & Metrics)

To view your mesh telemetry, open another PowerShell window and port-forward the Grafana service:

```powershell
kubectl port-forward svc/grafana 3000:3000 -n observability
```

Navigate your browser to: http://localhost:3000

Username: admin

Password: admin

Go to the Explore tab on the left menu. Use the dropdown at the top to switch between the Prometheus data source (for checking metrics like istio_requests_total) and the Loki data source (for viewing real-time pod logs aggregated by Promtail).

🧹 Cleanup

When you are finished testing, you can tear down the local cluster to free up resources:

```powershell
minikube delete
```
