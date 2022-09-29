package certificates

import (
	"errors"
)

var (
	UnknownCertKind = errors.New("Unkown certificate kind")
	WrongCertKind   = errors.New("Wrong certificate kind")
)
