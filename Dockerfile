FROM golang:1.21-alpine

WORKDIR /app

# ✅ 1. Copy go.mod ก่อน
COPY go.mod go.sum ./

# ✅ 2. run go mod tidy ได้เลย
RUN go mod tidy

# ✅ 3. ค่อย copy โค้ดที่เหลือ
COPY . .

# ✅ 4. build
RUN go build -o main .

CMD ["./main"]
