package certificates

import (
	"errors"
)

var (
	UnknownCertKind      = errors.New("Unkown certificate kind")
	WrongCertKind        = errors.New("Wrong certificate kind")
	BadCertificateParams = errors.New("Unable to parse certificate parameters")
	BadStakeCredential   = errors.New("Unable to parse stake credential")
	BadPoolId            = errors.New("Unable to parse stake pool identifier")
)
