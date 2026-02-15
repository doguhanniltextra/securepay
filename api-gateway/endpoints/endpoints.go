package endpoints

// APIPrefix is the base path for protected API endpoints.
const APIPrefix = "/api/v1/"

// HealthCheckPath is the path for the health check endpoint.
const HealthCheckPath = "/health"

// Protected Endpoints Patterns (including HTTP Method for Go 1.22+ mux)

// InitiatePaymentPathPattern is the route pattern for initiating a payment.
const InitiatePaymentPathPattern = "POST " + APIPrefix + "payments"

// GetPaymentPathPattern is the route pattern for retrieving payment details.
const GetPaymentPathPattern = "GET " + APIPrefix + "payments/{id}"

// CheckBalancePathPattern is the route pattern for checking account balance.
const CheckBalancePathPattern = "GET " + APIPrefix + "accounts/{id}/balance"
