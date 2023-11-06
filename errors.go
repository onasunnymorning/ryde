package ryde

import "fmt"

var (
	ErrInvalidDepositType     = fmt.Errorf("invalid deposit type, only FULL or DIFF are allowed")
	ErrInvalidDepositFileName = fmt.Errorf("invalid deposit file name, must end with .xml")
	ErrNoXMLReader            = fmt.Errorf("XMLFile.osFile is nil, try calling OpenXMLFile() first")
	ErrNoXMLDecoder           = fmt.Errorf("XMLFile.Decoder is nil, try calling CreateXMLDecoder() first")
	ErrNoDepositTagInFile     = fmt.Errorf("reached EOF before finding a <rde:deposit> start element")
)
