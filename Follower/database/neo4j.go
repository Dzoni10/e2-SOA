package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

var Driver neo4j.DriverWithContext

func InitNeo4j() error {
	uri := os.Getenv("NEO4J_URI")
	if uri == "" {
		uri = "bolt://localhost:7687"
	}

	username := os.Getenv("NEO4J_USERNAME")
	if username == "" {
		username = "neo4j"
	}

	password := os.Getenv("NEO4J_PASSWORD")
	if password == "" {
		password = "password"
	}

	var err error
	Driver, err = neo4j.NewDriverWithContext(
		uri,
		neo4j.BasicAuth(username, password, ""),
		func(config *neo4j.Config) {
			config.MaxConnectionLifetime = 30 * time.Minute
			config.MaxConnectionPoolSize = 50
			config.ConnectionAcquisitionTimeout = 2 * time.Minute
		},
	)

	if err != nil {
		return fmt.Errorf("failed to create Neo4j driver: %v", err)
	}

	// Test connection
	ctx := context.Background()
	err = Driver.VerifyConnectivity(ctx)
	if err != nil {
		return fmt.Errorf("failed to verify Neo4j connectivity: %v", err)
	}

	log.Println("Successfully connected to Neo4j database")

	// Create constraints and indexes
	err = createConstraintsAndIndexes(ctx)
	if err != nil {
		log.Printf("Warning: Failed to create constraints and indexes: %v", err)
	}

	return nil
}

func createConstraintsAndIndexes(ctx context.Context) error {
	session := Driver.NewSession(ctx, neo4j.SessionConfig{
		AccessMode: neo4j.AccessModeWrite,
	})
	defer session.Close(ctx)

	// Create unique constraint for User username (primary key)
	constraintQuery := `
		CREATE CONSTRAINT user_username_unique IF NOT EXISTS 
		FOR (u:User) REQUIRE u.username IS UNIQUE
	`

	_, err := session.Run(ctx, constraintQuery, nil)
	if err != nil {
		return fmt.Errorf("failed to create username constraint: %v", err)
	}

	// Create index for userId (for legacy support)
	userIdIndexQuery := `
		CREATE INDEX user_id_index IF NOT EXISTS 
		FOR (u:User) ON (u.userId)
	`

	_, err = session.Run(ctx, userIdIndexQuery, nil)
	if err != nil {
		return fmt.Errorf("failed to create userId index: %v", err)
	}

	// Create index for username (for faster lookups)
	usernameIndexQuery := `
		CREATE INDEX user_username_index IF NOT EXISTS 
		FOR (u:User) ON (u.username)
	`

	_, err = session.Run(ctx, usernameIndexQuery, nil)
	if err != nil {
		return fmt.Errorf("failed to create username index: %v", err)
	}

	// Create index for name (for search functionality)
	nameIndexQuery := `
		CREATE INDEX user_name_index IF NOT EXISTS 
		FOR (u:User) ON (u.name)
	`

	_, err = session.Run(ctx, nameIndexQuery, nil)
	if err != nil {
		return fmt.Errorf("failed to create name index: %v", err)
	}

	log.Println("Neo4j constraints and indexes created successfully")
	return nil
}

func CloseNeo4j(ctx context.Context) {
	if Driver != nil {
		Driver.Close(ctx)
		log.Println("Neo4j connection closed")
	}
}
