# 1️⃣ Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# ตั้งค่า environment เพื่อ build Linux binary
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
ENV GOPROXY=https://proxy.golang.org,direct

# copy go.mod / go.sum ก่อนเพื่อ cache layer
COPY go.mod go.sum ./
RUN go mod download

# copy source code
COPY . .

# build binary ชื่อ gateway
RUN go build -o gateway main.go

# 2️⃣ Final stage
FROM alpine:latest
WORKDIR /app

# copy binary จาก builder stage
COPY --from=builder /app/gateway .

# expose port ตามที่แอปคุณใช้งาน
EXPOSE 3000

# รัน binary
ENTRYPOINT ["./gateway"]
