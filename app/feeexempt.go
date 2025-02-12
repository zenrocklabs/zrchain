package app

import "context"

// FeeExempt is an interface that defines the methods for a fee exempt.
type FeeExempt interface {
	IsAllowed(ctx context.Context, typeURL string) (bool, error)
}
