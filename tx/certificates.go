package tx

import (
	"errors"

	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fxamacker/cbor/v2"
)

var UnknownCertKind = errors.New("Unkown certificate kind")

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
	Kind CertificateKind

	cred        *address.StakeCredential
	poolKeyHash [28]byte
	poolParams  interface{}
	epoch       uint
}

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

func (c *Certificate) UnmarshalCBOR(in []byte) (err error) {
	certKind := CertificateKind(in[1])
	switch certKind {
	case StakeRegistration, StakeDeregistration:
		x := &stakeCert{}
		err = cbor.Unmarshal(in, x)
		if err != nil {
			return
		}
		c.Kind = certKind
		c.cred = x.Cred
	case StakeDelegation:
		x := &delegationCert{}
		err = cbor.Unmarshal(in, x)
		if err != nil {
			return
		}
		c.Kind = certKind
		c.cred = x.Cred
		copy(c.poolKeyHash[:], x.Pool[:])
	case PoolRegistration:
		x := &poolRegistrationCert{}
		err = cbor.Unmarshal(in, x)
		if err != nil {
			return
		}
		c.Kind = certKind
		c.poolParams = x.Params
	case PoolRetirement:
		x := &poolRetirementCert{}
		err = cbor.Unmarshal(in, x)
		if err != nil {
			return
		}
		c.Kind = certKind
		copy(c.poolKeyHash[:], x.Pool[:])
		c.epoch = x.Epoch
	}
	return
}

func (c *Certificate) MarshalCBOR() ([]byte, error) {
	switch c.Kind {
	case StakeRegistration, StakeDeregistration:
		x := &stakeCert{
			Kind: c.Kind,
			Cred: c.cred,
		}
		return cbor.Marshal(x)
	case StakeDelegation:
		x := &delegationCert{
			Kind: c.Kind,
			Cred: c.cred,
		}
		copy(x.Pool[:], c.poolKeyHash[:])
		return cbor.Marshal(x)
	case PoolRegistration:
		x := &poolRegistrationCert{
			Kind:   c.Kind,
			Params: c.poolParams,
		}
		return cbor.Marshal(x)
	case PoolRetirement:
		x := &poolRetirementCert{
			Kind:  c.Kind,
			Epoch: c.epoch,
		}
		copy(x.Pool[:], c.poolKeyHash[:])
		return cbor.Marshal(x)
	}
	return nil, UnknownCertKind
}
