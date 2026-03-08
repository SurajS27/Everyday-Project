# Build Stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copy module files and download dependencies
COPY go.mod ./
# (If we had a go.sum we'd copy it here: COPY go.mod go.sum ./)
RUN go mod download

# Copy the source code
COPY main.go ./

# Build the statically linked Go binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o order-service main.go

# Production Stage (Minimal Alpine image)
FROM alpine:latest

WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/order-service .

# Expose the application port
EXPOSE 8080

# Run the binary
CMD ["./order-service"]