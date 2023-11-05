package ryde

import "fmt"

var (
	ErrInvalidDepositType = fmt.Errorf("invalid deposit type, only FULL or DIFF are allowed")
)
