FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o main -ldflags="-s -w" .

FROM alpine:3.20

WORKDIR /app

# Kopiramo binarni fajl iz build stage-a
COPY --from=builder /app/main .

EXPOSE 8081

CMD ["./main"]