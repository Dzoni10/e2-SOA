FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY . .

WORKDIR /build/Follower

RUN go mod download

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /build/Follower/main .

EXPOSE 8083

CMD ["./main"]