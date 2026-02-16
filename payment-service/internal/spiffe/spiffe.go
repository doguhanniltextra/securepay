package spiffe

import (
	"context"
	"fmt"

	"github.com/spiffe/go-spiffe/v2/spiffeid"
	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// InitSPIFFESource initializes and returns a new X.509 source connected to the SPIRE Agent.
// It is the caller's responsibility to close the source when done.
func InitSPIFFESource(ctx context.Context, socketPath string) (*workloadapi.X509Source, error) {
	// Create a new X.509 source with the configured socket path.
	// The source will automatically fetch and renew SVIDs.
	clientOptions := workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath))
	source, err := workloadapi.NewX509Source(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to create X509Source: %w", err)
	}

	return source, nil
}

// AccountServiceCredentials returns gRPC dial options with mTLS credentials
// specifically for connecting to the Account Service.
// It enforces that the server presents a valid SVID with the Account Service's SPIFFE ID.
func AccountServiceCredentials(source *workloadapi.X509Source) (grpc.DialOption, error) {
	// Define the expected SPIFFE ID for the Account Service
	accountServiceID := spiffeid.RequireFromString("spiffe://securepay.dev/account-service")

	// Create mTLS client configuration
	// - source: provides our client certificate (SVID)
	// - source: provides the trust bundle to verify the server's certificate
	// - AuthorizeID: ensures the server has the specific Account Service ID
	tlsConfig := tlsconfig.MTLSClientConfig(source, source, tlsconfig.AuthorizeID(accountServiceID))

	return grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)), nil
}

// PaymentServiceServerCredentials returns gRPC server options with mTLS credentials.
// It authorizes any client within the trust domain (e.g., API Gateway).
func PaymentServiceServerCredentials(source *workloadapi.X509Source) grpc.ServerOption {
	// Authorize any client with a valid SVID from the trust domain.
	// This allows API Gateway and others to connect securely.
	tlsConfig := tlsconfig.MTLSServerConfig(source, source, tlsconfig.AuthorizeAny())
	return grpc.Creds(credentials.NewTLS(tlsConfig))
}
