package tx

import (
	"errors"

	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fxamacker/cbor/v2"
)

var (
	UnknownCertKind = errors.New("Unkown certificate kind")
	WrongCertKind   = errors.New("Wrong certificate kind")
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

type Certificate struct {
	kind CertificateKind

	cred        *address.StakeCredential
	poolKeyHash [28]byte
	poolParams  interface{}
	epoch       uint
}

type (
	i_cert_WithKind interface {
		Kind() CertificateKind
	}

	i_cert_WithCred interface {
		Cred() address.StakeCredential
	}
	cert_WithCred struct{ *Certificate }

	i_cert_WithPool interface {
		Pool() [28]byte
	}
	cert_WithPool struct{ *Certificate }

	i_cert_WithEpoch interface {
		Epoch() uint
	}
	cert_WithEpoch struct{ *Certificate }

	i_cert_WithParams interface {
		Params() interface{}
	}
	cert_WithParams struct{ *Certificate }
)

type (
	StakeCert interface {
		Kind() CertificateKind
		Cred() address.StakeCredential
	}
	stakeCert struct {
		_    struct{} `cbor:",toarray"`
		Kind CertificateKind
		Cred *address.StakeCredential
	}
	asStakeCert struct {
		i_cert_WithKind
		i_cert_WithCred
	}

	DelegationCert interface {
		Kind() CertificateKind
		Cred() address.StakeCredential
		Pool() [28]byte
	}
	delegationCert struct {
		_    struct{} `cbor:",toarray"`
		Kind CertificateKind
		Cred *address.StakeCredential
		Pool [28]byte
	}
	asDelegationCert struct {
		i_cert_WithKind
		i_cert_WithCred
		i_cert_WithPool
	}
)

type (
	PoolRegistrationCert interface {
		Kind() CertificateKind
		Params() interface{}
	}
	poolRegistrationCert struct {
		_      struct{} `cbor:",toarray"`
		Kind   CertificateKind
		Params interface{}
	}
	asPoolRegistrationCert struct {
		i_cert_WithKind
		i_cert_WithParams
	}

	PoolRetirementCert interface {
		Kind() CertificateKind
		Pool() [28]byte
		Epoch() uint
	}
	poolRetirementCert struct {
		_     struct{} `cbor:",toarray"`
		Kind  CertificateKind
		Pool  [28]byte
		Epoch uint
	}
	asPoolRetirementCert struct {
		i_cert_WithKind
		i_cert_WithPool
		i_cert_WithEpoch
	}
)

var (
	_ i_cert_WithKind      = &Certificate{}
	_ StakeCert            = &asStakeCert{}
	_ DelegationCert       = &asDelegationCert{}
	_ PoolRegistrationCert = &asPoolRegistrationCert{}
	_ PoolRetirementCert   = &asPoolRetirementCert{}
)

func (c *Certificate) Kind() CertificateKind { return c.kind }

func (c *Certificate) AsCert() interface{} {
	switch c.kind {
	case StakeRegistration, StakeDeregistration:
		return &asStakeCert{i_cert_WithKind: c, i_cert_WithCred: &cert_WithCred{Certificate: c}}
	case StakeDelegation:
		return &asDelegationCert{i_cert_WithKind: c, i_cert_WithCred: &cert_WithCred{Certificate: c}, i_cert_WithPool: &cert_WithPool{Certificate: c}}
	case PoolRegistration:
		return &asPoolRegistrationCert{i_cert_WithKind: c, i_cert_WithParams: &cert_WithParams{c}}
	case PoolRetirement:
		return &asPoolRetirementCert{i_cert_WithKind: c, i_cert_WithPool: &cert_WithPool{Certificate: c}, i_cert_WithEpoch: &cert_WithEpoch{c}}
	}
	return UnknownCertKind
}

func (c *Certificate) AsStakeCert() (StakeCert, error) {
	if c.kind != StakeRegistration && c.kind != StakeDeregistration {
		return nil, WrongCertKind
	}
	return c.AsCert().(StakeCert), nil
}
func (c *Certificate) AsDelegationCert() (DelegationCert, error) {
	if c.kind != StakeDelegation {
		return nil, WrongCertKind
	}
	return c.AsCert().(DelegationCert), nil
}
func (c *Certificate) AsPoolRegistrationCert() (PoolRegistrationCert, error) {
	if c.kind != PoolRegistration {
		return nil, WrongCertKind
	}
	return c.AsCert().(PoolRegistrationCert), nil
}
func (c *Certificate) AsPoolRetirementCert() (PoolRetirementCert, error) {
	if c.kind != PoolRetirement {
		return nil, WrongCertKind
	}
	return c.AsCert().(PoolRetirementCert), nil
}

func (c *Certificate) UnmarshalCBOR(in []byte) (err error) {
	certKind := CertificateKind(in[1])
	switch certKind {
	case StakeRegistration, StakeDeregistration:
		x := &stakeCert{}
		err = cbor.Unmarshal(in, x)
		if err != nil {
			return
		}
		c.kind = certKind
		c.cred = x.Cred
	case StakeDelegation:
		x := &delegationCert{}
		err = cbor.Unmarshal(in, x)
		if err != nil {
			return
		}
		c.kind = certKind
		c.cred = x.Cred
		copy(c.poolKeyHash[:], x.Pool[:])
	case PoolRegistration:
		x := &poolRegistrationCert{}
		err = cbor.Unmarshal(in, x)
		if err != nil {
			return
		}
		c.kind = certKind
		c.poolParams = x.Params
	case PoolRetirement:
		x := &poolRetirementCert{}
		err = cbor.Unmarshal(in, x)
		if err != nil {
			return
		}
		c.kind = certKind
		copy(c.poolKeyHash[:], x.Pool[:])
		c.epoch = x.Epoch
	}
	return
}

func (c *Certificate) MarshalCBOR() ([]byte, error) {
	switch c.kind {
	case StakeRegistration, StakeDeregistration:
		x := &stakeCert{
			Kind: c.kind,
			Cred: c.cred,
		}
		return cbor.Marshal(x)
	case StakeDelegation:
		x := &delegationCert{
			Kind: c.kind,
			Cred: c.cred,
		}
		copy(x.Pool[:], c.poolKeyHash[:])
		return cbor.Marshal(x)
	case PoolRegistration:
		x := &poolRegistrationCert{
			Kind:   c.kind,
			Params: c.poolParams,
		}
		return cbor.Marshal(x)
	case PoolRetirement:
		x := &poolRetirementCert{
			Kind:  c.kind,
			Epoch: c.epoch,
		}
		copy(x.Pool[:], c.poolKeyHash[:])
		return cbor.Marshal(x)
	}
	return nil, UnknownCertKind
}

func (c *cert_WithCred) Cred() address.StakeCredential { return *c.Certificate.cred }
func (c *cert_WithPool) Pool() [28]byte {
	pool := [28]byte{0}
	copy(pool[:], c.Certificate.poolKeyHash[:])
	return pool
}
func (c *cert_WithParams) Params() interface{} { return c.Certificate.poolParams }
func (c *cert_WithEpoch) Epoch() uint          { return c.Certificate.epoch }
