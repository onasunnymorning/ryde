package ryde

import "fmt"

var (
	ErrInvalidDepositType     = fmt.Errorf("invalid deposit type, only FULL or DIFF are allowed")
	ErrInvalidDepositFileName = fmt.Errorf("invalid deposit file name, must end with .xml")
)
