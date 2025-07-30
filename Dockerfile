# Step 1: ติดตั้ง Go 1.24.4 แบบ custom
FROM alpine:latest AS go-installer

RUN apk add --no-cache curl tar gcc musl-dev git \
    && curl -LO https://go.dev/dl/go1.24.4.src.tar.gz \
    && tar -C /usr/local -xzf go1.24.4.src.tar.gz

ENV GOROOT=/usr/local/go
ENV PATH=$GOROOT/bin:$PATH

WORKDIR /go-src
COPY . .

RUN go mod tidy && go build -o main .

# Step 2: Run แบบ minimal
FROM alpine:latest

WORKDIR /app
COPY --from=go-installer /go-src/main .

EXPOSE 3000
CMD ["./main"]
