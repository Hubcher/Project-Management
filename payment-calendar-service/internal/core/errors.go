package core

import "errors"

var (
	ErrPaymentNotFound = errors.New("payment not found")
	ErrInvalidPayment  = errors.New("invalid payment")
	ErrProjectNotFound = errors.New("project not found")
	ErrStageNotFound   = errors.New("project stage not found")
	ErrAlreadyExists   = errors.New("resource already exists")
)
