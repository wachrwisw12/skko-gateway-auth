# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build for Linux amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o gateway main.go

# Final stage
FROM alpine:3.18

WORKDIR /app
COPY --from=builder /app/gateway .

# ทำให้ binary เป็น executable
RUN chmod +x ./gateway

ENTRYPOINT ["./gateway"]
