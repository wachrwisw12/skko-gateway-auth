FROM golang:1.21-alpine

WORKDIR /app

# ✅ ต้อง copy ก่อน tidy
COPY go.mod go.sum ./
RUN go mod tidy

# แล้วค่อย copy โค้ด
COPY . .

RUN go build -o main .

CMD ["./main"]
