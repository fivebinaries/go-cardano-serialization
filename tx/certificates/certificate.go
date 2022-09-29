package certificates

import (
	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fxamacker/cbor/v2"
)

// opaque polymorphic certificate struct
type Certificate struct {
	kind CertificateKind

	cred        *address.StakeCredential
	poolKeyHash [28]byte
	poolParams  interface{}
	epoch       uint
}

var _ iWithKind = &Certificate{}

func (c *Certificate) Kind() CertificateKind { return c.kind }

func (c *Certificate) AsCert() interface{} {
	switch c.kind {
	case StakeRegistration, StakeDeregistration:
		return &asStakeCert{iWithKind: c, iWithCred: &certWithCred{Certificate: c}}
	case StakeDelegation:
		return &asDelegationCert{iWithKind: c, iWithCred: &certWithCred{Certificate: c}, iWithPool: &certWithPool{Certificate: c}}
	case PoolRegistration:
		return &asPoolRegistrationCert{iWithKind: c, iWithParams: &certWithParams{c}}
	case PoolRetirement:
		return &asPoolRetirementCert{iWithKind: c, iWithPool: &certWithPool{Certificate: c}, iWithEpoch: &certWithEpoch{c}}
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
