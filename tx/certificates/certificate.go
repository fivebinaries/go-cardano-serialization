package certificates

import (
	"encoding/hex"
	"github.com/fivebinaries/go-cardano-serialization/address"
	"github.com/fivebinaries/go-cardano-serialization/internal/bech32"
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

func parseStakeCred(in interface{}) (*address.StakeCredential, error) {
	var cred address.StakeCredential
	switch v := in.(type) {
	case string:
		{
			var (
				addr  address.Address
				raddr *address.RewardAddress
				err   error
				ok    bool
			)
			addr, err = address.NewAddress(v)
			if err != nil {
				addr, err = address.NewAddressFromHex(v)
				if err != nil {
					return nil, BadStakeCredential
				}
			}
			raddr, ok = addr.(*address.RewardAddress)
			if !ok {
				return nil, BadStakeCredential
			}
			cred = raddr.Stake
		}
	case *address.StakeCredential:
		cred = *v
	case address.StakeCredential:
		cred = v
	}
	return &cred, nil
}

func parsePool(in interface{}) ([28]byte, error) {
	pool := [28]byte{0}
	switch v := in.(type) {
	case [28]byte:
		// assuming it is already raw bytes
		copy(pool[:], v[:])
	case string:
		if len(v) < 5 {
			return [28]byte{0}, BadPoolId
		}
		if v[:5] == "pool1" {
			// is bech32 encoded
			hrp, data, err := bech32.Decode(v)
			if err != nil {
				return [28]byte{0}, err
			}
			if hrp != "pool" {
				return [28]byte{0}, BadPoolId
			}
			copy(pool[:], data)
		} else {
			// assuming this string is the HEX of the data
			if hex.DecodedLen(len(v)) != 28 {
				return [28]byte{0}, BadPoolId
			}
			data, err := hex.DecodeString(v)
			if err != nil {
				return [28]byte{0}, err
			}
			if len(data) != 28 {
				return [28]byte{0}, BadPoolId
			}
			copy(pool[:], data)
		}
	}
	return pool, nil
}

func NewCertificate(kind CertificateKind, args ...interface{}) (*Certificate, error) {
	switch kind {
	case StakeRegistration, StakeDeregistration:
		if len(args) != 1 {
			return nil, BadCertificateParams
		}
		cred, err := parseStakeCred(args[0])
		if err != nil {
			return nil, err
		}
		return &Certificate{kind: kind, cred: cred}, nil
	case StakeDelegation:
		if len(args) != 2 {
			return nil, BadCertificateParams
		}
		cred, err := parseStakeCred(args[0])
		if err != nil {
			return nil, err
		}
		pool, err := parsePool(args[1])
		if err != nil {
			return nil, err
		}
		return &Certificate{kind: StakeDelegation, cred: cred, poolKeyHash: pool}, nil

	case PoolRegistration:
		if len(args) != 1 {
			return nil, BadCertificateParams
		}
		// nothing to parse yet
		return &Certificate{kind: PoolRegistration, poolParams: args[0]}, nil

	case PoolRetirement:
		if len(args) != 2 {
			return nil, BadCertificateParams
		}
		pool, err := parsePool(args[0])
		if err != nil {
			return nil, err
		}
		epoch, ok := args[1].(uint)
		if !ok {
			return nil, BadCertificateParams
		}
		return &Certificate{kind: PoolRegistration, poolKeyHash: pool, epoch: epoch}, nil
	default:
		return nil, UnknownCertKind
	}
}

func NewStakeRegistrationCertificate(in interface{}) (*Certificate, error) {
	return NewCertificate(StakeRegistration, in)
}

func NewStakeDeregistrationCertificate(in interface{}) (*Certificate, error) {
	return NewCertificate(StakeDeregistration, in)
}

func NewStakeDelegationCertificate(credIn, poolIn interface{}) (*Certificate, error) {
	return NewCertificate(StakeDelegation, credIn, poolIn)
}

func NewPoolRegistrationCertificate(params interface{}) (*Certificate, error) {
	return NewCertificate(PoolRegistration, params)
}

func NewPoolRetirmentCertificate(poolIn, epochIn interface{}) (*Certificate, error) {
	return NewCertificate(PoolRegistration, poolIn, epochIn)
}
