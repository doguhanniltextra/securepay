package spiffe

import (
	"context"
	"fmt"

	"github.com/spiffe/go-spiffe/v2/spiffetls/tlsconfig"
	"github.com/spiffe/go-spiffe/v2/workloadapi"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// InitSPIFFESource initializes and returns a new X.509 source connected to the SPIRE Agent.
// It is the caller's responsibility to close the source when done.
func InitSPIFFESource(ctx context.Context, socketPath string) (*workloadapi.X509Source, error) {
	clientOptions := workloadapi.WithClientOptions(workloadapi.WithAddr(socketPath))
	source, err := workloadapi.NewX509Source(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("unable to create X509Source: %w", err)
	}

	return source, nil
}

// ServerCredentials returns gRPC server options with mTLS credentials.
// It authorizes any client within the trust domain (e.g., API Gateway, Payment Service).
func ServerCredentials(source *workloadapi.X509Source) grpc.ServerOption {
	tlsConfig := tlsconfig.MTLSServerConfig(source, source, tlsconfig.AuthorizeAny())
	return grpc.Creds(credentials.NewTLS(tlsConfig))
}
