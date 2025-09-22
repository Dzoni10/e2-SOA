FROM golang:1.25-alpine AS builder

WORKDIR /build

COPY . .

WORKDIR /build/Blog

RUN go mod download

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /build/Blog/main .

EXPOSE 8082

CMD ["./main"]