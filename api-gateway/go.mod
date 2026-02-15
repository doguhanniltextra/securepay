module securepay/api-gateway

go 1.24.6

require (
	github.com/golang-jwt/jwt/v5 v5.3.1
	github.com/spiffe/go-spiffe/v2 v2.6.0
	google.golang.org/grpc v1.79.1
	securepay/proto v0.0.0-00010101000000-000000000000
)

require (
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/go-jose/go-jose/v4 v4.1.3 // indirect
	golang.org/x/net v0.48.0 // indirect
	golang.org/x/sys v0.39.0 // indirect
	golang.org/x/text v0.32.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251202230838-ff82c1b0f217 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace securepay/proto => ../proto
