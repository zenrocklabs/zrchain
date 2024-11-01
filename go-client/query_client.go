package client

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	insecurecreds "google.golang.org/grpc/credentials/insecure"
)

// QueryClient holds a query client for the auth and treasury modules.
type QueryClient struct {
	*AuthQueryClient
	*TreasuryQueryClient
	*PolicyQueryClient
	*ValidationQueryClient
	*ZenBTCQueryClient
	conn *grpc.ClientConn
}

// NewQueryClient returns a QueryClient. The supplied url must be a GRPC compatible endpoint for zenrockd.
func NewQueryClient(url string, insecure bool) (*QueryClient, error) {
	grpcConn, err := NewClientConn(url, insecure)
	if err != nil {
		return nil, err
	}
	return NewQueryClientWithConn(grpcConn), nil
}

func NewClientConn(url string, insecure bool) (*grpc.ClientConn, error) {
	opts := []grpc.DialOption{}
	if insecure {
		opts = append(opts, grpc.WithTransportCredentials(insecurecreds.NewCredentials()))
	} else {
		certPool := x509.NewCertPool()
		serverCert, err := ioutil.ReadFile("server.crt")
		if err != nil {
			panic("'server.crt' file is not in root dir")
		}
		certPool.AppendCertsFromPEM(serverCert)

		// Create the credentials and connect to server
		config := &tls.Config{
			RootCAs: certPool,
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(config)))
	}
	grpcConn, err := grpc.NewClient(url, opts...)
	if err != nil {
		return nil, err
	}
	return grpcConn, nil
}

// NewQueryClientWithConn returns a QueryClient with the supplied GRPC client connection.
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
