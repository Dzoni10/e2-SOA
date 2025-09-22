module follower-service

go 1.25.0

require (
	github.com/gorilla/mux v1.8.1
	github.com/neo4j/neo4j-go-driver/v5 v5.14.0
	github.com/rs/cors v1.11.1
	local/common v0.0.0
)

require (
	github.com/gorilla/websocket v1.5.3
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/nats-io/nats.go v1.45.0 // indirect
	github.com/nats-io/nkeys v0.4.11 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	golang.org/x/crypto v0.37.0 // indirect
	golang.org/x/sys v0.32.0 // indirect
)

replace local/common => ../common
