package product

import "errors"

var (
	ErrEventNotFound   = errors.New("event not found")
	ErrProductNotFound = errors.New("product not found")
	ErrUnauthorized    = errors.New("unauthorized access to this event")
	ErrNoUpdateFields  = errors.New("no fields provided for update")
	ErrInvalidSeller   = errors.New("invalid seller ID")
	ErrInvalidEvent    = errors.New("invalid event ID")
)
