package validator

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	pb "securepay/proto/gen/go/payment/v1"
)

// Validator handles request validation logic
type Validator struct{}

// New creates a new Validator
func New() *Validator {
	return &Validator{}
}

// UUID regex pattern
var uuidRegex = regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)

// ValidateInitiatePayment validates the InitiatePaymentRequest
func (v *Validator) ValidateInitiatePayment(req *pb.InitiatePaymentRequest) error {
	if req == nil {
		return errors.New("request cannot be nil")
	}

	// 1. Amount > 0
	if req.Amount <= 0 {
		return errors.New("amount must be greater than 0")
	}

	// 2. Currency valid (TRY, USD, EUR)
	validCurrencies := map[string]bool{
		"TRY": true,
		"USD": true,
		"EUR": true,
	}
	if !validCurrencies[strings.ToUpper(req.Currency)] {
		return fmt.Errorf("invalid currency: %s (supported: TRY, USD, EUR)", req.Currency)
	}

	// 3. from_account and to_account in UUID format
	if !uuidRegex.MatchString(req.FromAccount) {
		return fmt.Errorf("invalid from_account format: %s", req.FromAccount)
	}
	if !uuidRegex.MatchString(req.ToAccount) {
		return fmt.Errorf("invalid to_account format: %s", req.ToAccount)
	}

	// 4. from_account == to_account (forbidden)
	if req.FromAccount == req.ToAccount {
		return errors.New("from_account and to_account cannot be the same")
	}

	return nil
}
