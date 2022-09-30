package tx

import (
	"github.com/fivebinaries/go-cardano-serialization/tx/certificates"
)

type Certificate = certificates.Certificate

type CertificateKind = certificates.CertificateKind

const (
	StakeRegistration   = certificates.StakeRegistration
	StakeDeregistration = certificates.StakeDeregistration
	StakeDelegation     = certificates.StakeDelegation
	PoolRegistration    = certificates.PoolRegistration
	PoolRetirement      = certificates.PoolRetirement
)

type StakeCertificate = certificates.StakeCert
type DelegationCertificate = certificates.DelegationCert
type PoolRegistrationCertificate = certificates.PoolRegistrationCert
type PoolRetirementCertificate = certificates.PoolRetirementCert

var (
	NewCertificate = certificates.NewCertificate

	NewStakeRegistrationCertificate   = certificates.NewStakeRegistrationCertificate
	NewStakeDeregistrationCertificate = certificates.NewStakeDeregistrationCertificate
	NewStakeDelegationCertificate     = certificates.NewStakeDelegationCertificate
	NewPoolRegistrationCertificate    = certificates.NewPoolRegistrationCertificate
	NewPoolRetirmentCertificate       = certificates.NewPoolRetirmentCertificate
)
