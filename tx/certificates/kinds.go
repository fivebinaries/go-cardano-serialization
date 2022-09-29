package certificates

import (
	"github.com/fivebinaries/go-cardano-serialization/address"
)

type CertificateKind byte

const (
	StakeRegistration   CertificateKind = iota // // = (0, stake_credential)
	StakeDeregistration                        // = (1, stake_credential)
	StakeDelegation                            // = (2, stake_credential, pool_keyhash)
	PoolRegistration                           // = (3, pool_params)
	PoolRetirement                             // = (4, pool_keyhash, epoch)

	//genesis_key_delegation = (5, genesishash, genesis_delegate_hash, vrf_keyhash)
	//move_instantaneous_rewards_cert = (6, move_instantaneous_reward)

)

// for encoding structs, doesn't implement any iface
type (
	stakeCert struct {
		_    struct{} `cbor:",toarray"`
		Kind CertificateKind
		Cred *address.StakeCredential
	}

	delegationCert struct {
		_    struct{} `cbor:",toarray"`
		Kind CertificateKind
		Cred *address.StakeCredential
		Pool [28]byte
	}

	poolRegistrationCert struct {
		_      struct{} `cbor:",toarray"`
		Kind   CertificateKind
		Params interface{}
	}

	poolRetirementCert struct {
		_     struct{} `cbor:",toarray"`
		Kind  CertificateKind
		Pool  [28]byte
		Epoch uint
	}
)

// implicit ifaces implementors
type (
	certWithCred   struct{ *Certificate }
	certWithPool   struct{ *Certificate }
	certWithEpoch  struct{ *Certificate }
	certWithParams struct{ *Certificate }
)

// expliicit ifaces implementor wrappers
type (
	asStakeCert struct {
		iWithKind
		iWithCred
	}

	asDelegationCert struct {
		iWithKind
		iWithCred
		iWithPool
	}

	asPoolRegistrationCert struct {
		iWithKind
		iWithParams
	}

	asPoolRetirementCert struct {
		iWithKind
		iWithPool
		iWithEpoch
	}
)

var (
	_ StakeCert            = &asStakeCert{}
	_ DelegationCert       = &asDelegationCert{}
	_ PoolRegistrationCert = &asPoolRegistrationCert{}
	_ PoolRetirementCert   = &asPoolRetirementCert{}
)

func (c *certWithCred) Cred() address.StakeCredential { return *c.Certificate.cred }
func (c *certWithPool) Pool() [28]byte {
	pool := [28]byte{0}
	copy(pool[:], c.Certificate.poolKeyHash[:])
	return pool
}
func (c *certWithParams) Params() interface{} { return c.Certificate.poolParams }
func (c *certWithEpoch) Epoch() uint          { return c.Certificate.epoch }
