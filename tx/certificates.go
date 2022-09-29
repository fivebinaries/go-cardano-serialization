package tx

import (
	"github.com/fivebinaries/go-cardano-serialization/tx/certificates"
)

type Certificate = certificates.Certificate

// type Certificates []*Certificate

type CertificateKind = certificates.CertificateKind

type StakeCertificate = certificates.StakeCert
type DelegationCertificate = certificates.DelegationCert
type PoolRegistrationCertificate = certificates.PoolRegistrationCert
type PoolRetirementCertificate = certificates.PoolRetirementCert
