# # Build stage
# FROM golang:1.24-alpine AS builder
# WORKDIR /app

# ENV CGO_ENABLED=0 GOOS=linux
# RUN apk add --no-cache git
# COPY go.mod go.sum ./
# RUN go mod download

# COPY . .
# RUN go build -o gateway main.go

# # Final stage
# FROM alpine:latest
# WORKDIR /app
# COPY --from=builder /app/gateway .
# EXPOSE 3000
# ENTRYPOINT ["./gateway"]

FROM golang:1.24 AS builder
WORKDIR /app

ENV CGO_ENABLED=0 GOOS=linux

# ติดตั้ง git สำหรับ go mod download
RUN apk add --no-cache git

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o gateway main.go

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/gateway .
EXPOSE 3000
ENTRYPOINT ["./gateway"]
