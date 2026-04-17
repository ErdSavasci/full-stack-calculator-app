# --- Stage 1: The Build Environment ---
FROM golang:1.26-alpine AS builder

WORKDIR /app

# 1. Handle dependencies
COPY go.mod go.sum ./
RUN go mod download

# 2. Copy the source code AND the config file into the builder
COPY . .

# 3. Compile the binary
RUN GOMEMLIMIT=2GiB CGO_ENABLED=0 GOOS=linux go build -p 2 -o main ./services/multiplication

# --- Stage 2: The Production Image ---
FROM alpine:latest

WORKDIR /root/

# 4. Copy the binary from the builder stage
COPY --from=builder /app/main .

# 5. Copy the config file from the builder stage
# (This ensures Viper finds it at ./config/config.yaml)
RUN mkdir -p config
COPY --from=builder /app/config/config.yaml ./config/

# 6. Run the application
CMD ["./main"]