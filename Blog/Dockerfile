FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main -ldflags="-s -w" .

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 8082

CMD ["./main"]