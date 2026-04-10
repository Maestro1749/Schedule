package models

import "errors"

var (
	ErrInvalidWeekType  = errors.New("invalid week type: must be 1 or 2")
	ErrInvalidWeekday   = errors.New("invalid weekday: must be between 1 and 7")
	ErrInternalServer   = errors.New("internal server error")
	ErrInvalidDataInput = errors.New("Invalid data input")
	ErrAlreadyExists    = errors.New("entity already exists")
)
