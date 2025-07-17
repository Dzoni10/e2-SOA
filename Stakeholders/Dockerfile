# Start from the official Go image
FROM golang:1.21-alpine

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Copy the entire source
COPY . .

# Build the Go app
RUN go build -o main .

# Expose port (npr. 8080 ako koristiš njega u main.go)
EXPOSE 8080

# Run the binary
CMD ["./main"]