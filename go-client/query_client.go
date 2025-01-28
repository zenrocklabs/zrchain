package client

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	insecurecreds "google.golang.org/grpc/credentials/insecure"
)

// Package client provides a unified gRPC client implementation for interacting with various
// Zenrock blockchain modules. It offers high-level query interfaces for auth, treasury,
// policy, validation, and ZenBTC functionalities.

// QueryClient is the main client interface that aggregates all module-specific query clients.
// It provides a unified access point for querying different aspects of the Zenrock blockchain.
type QueryClient struct {
	*AuthQueryClient                        // For querying account-related information
	*TreasuryQueryClient                    // For querying treasury module data
	*PolicyQueryClient                      // For querying policy-related information
	*ValidationQueryClient                  // For querying validator information
	*ZenBTCQueryClient                      // For querying ZenBTC-related data
	conn                   *grpc.ClientConn // The underlying gRPC connection
}

// NewQueryClient creates a new QueryClient instance by establishing a gRPC connection
// to the specified Zenrock node endpoint.
//
// Parameters:
//   - url: The gRPC endpoint URL of the Zenrock node
//   - insecure: If true, creates an insecure connection. For development use only
//
// Returns:
//   - *QueryClient: A new query client instance with all module clients initialized
//   - error: An error if the connection fails
//
// Example:
//
//	client, err := NewQueryClient("localhost:9090", false)
//	if err != nil {
//	    // Handle error
//	}
//	defer client.conn.Close()
func NewQueryClient(url string, insecure bool) (*QueryClient, error) {
	grpcConn, err := NewClientConn(url, insecure)
	if err != nil {
		return nil, err
	}
	return NewQueryClientWithConn(grpcConn), nil
}

// NewClientConn establishes a new gRPC connection to the specified endpoint.
// It handles both secure and insecure connection types based on the insecure parameter.
//
// Parameters:
//   - url: The gRPC endpoint URL
//   - insecure: If true, creates an insecure connection
//
// Returns:
//   - *grpc.ClientConn: The established gRPC connection
//   - error: An error if the connection fails
//
// Note: When using secure connections, expects a 'server.crt' file in the root directory
func NewClientConn(url string, insecure bool) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{}
	if insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecurecreds.NewCredentials()))
	} else {
		systemCertPool, err := x509.SystemCertPool()
		if err != nil {
			return nil, fmt.Errorf("failed to load system cert pool: %v", err)
		}
		// Create the credentials and connect to server
		config := &tls.Config{
			RootCAs: systemCertPool,
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	}
	grpcConn, err := grpc.NewClient(url, opts...)
	if err != nil {
		return nil, err
	}
	return grpcConn, nil
}

// NewQueryClientWithConn creates a new QueryClient using an existing gRPC connection.
// This is useful when you want to manage the gRPC connection lifecycle separately
// or share a connection among multiple clients.
//
// Parameters:
//   - conn: An established gRPC client connection
//
// Returns:
//   - *QueryClient: A new query client instance using the provided connection
//
// Example:
//
//	conn, _ := grpc.Dial("localhost:9090", grpc.WithInsecure())
//	client := NewQueryClientWithConn(conn)
//	defer conn.Close()
func NewQueryClientWithConn(conn *grpc.ClientConn) *QueryClient {
	return &QueryClient{
		AuthQueryClient:       NewAuthQueryClient(conn),
		TreasuryQueryClient:   NewTreasuryQueryClient(conn),
		PolicyQueryClient:     NewPolicyQueryClient(conn),
		ValidationQueryClient: NewValidationQueryClient(conn),
		ZenBTCQueryClient:     NewZenBTCQueryClient(conn),
		conn:                  conn,
	}
}
