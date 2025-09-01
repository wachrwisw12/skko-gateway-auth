FROM golang:1.24 AS builder
WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux

# ติดตั้ง git บน Debian/Ubuntu
RUN apt-get update && apt-get install -y git && rm -rf /var/lib/apt/lists/*

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o gateway main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/gateway .
EXPOSE 3000
ENTRYPOINT ["./gateway"]
