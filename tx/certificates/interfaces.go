package certificates

import (
	"github.com/fivebinaries/go-cardano-serialization/address"
)

// implicit for code-sharing
type (
	iWithKind interface {
		Kind() CertificateKind
	}

	iWithCred interface {
		Cred() address.StakeCredential
	}

	iWithPool interface {
		Pool() [28]byte
	}

	iWithEpoch interface {
		Epoch() uint
	}

	iWithParams interface {
		Params() interface{}
	}
)

// explicit ifaces
type (
	StakeCert interface {
		Kind() CertificateKind
		Cred() address.StakeCredential
	}

	DelegationCert interface {
		Kind() CertificateKind
		Cred() address.StakeCredential
		Pool() [28]byte
	}

	PoolRegistrationCert interface {
		Kind() CertificateKind
		Params() interface{}
	}

	PoolRetirementCert interface {
		Kind() CertificateKind
		Pool() [28]byte
		Epoch() uint
	}
)
