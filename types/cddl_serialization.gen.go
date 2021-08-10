package types

import (
	"github.com/fivebinaries/go-cardano-serialization/hash_map"
)

type Hash28 []byte

type Hash32 []byte

type Vkey []byte

type VrfVkey []byte

type VrfCert struct {
	_  interface{} `cbor:",toarray"`
	V1 []byte
	V2 []byte
}

type Natural []byte

type KesVkey []byte

type KesSignature []byte

type SignkeyKES []byte

type Signature []byte

type Nonce struct {
	V1 *int32

	V2 *NonceComposition0
}

type UnitInterval struct {
	_  interface{} `cbor:",toarray"`
	V1 uint
	V2 uint
}

type Rational struct {
	_  interface{} `cbor:",toarray"`
	V1 uint
	V2 uint
}

type Block struct {
	_                      interface{} `cbor:",toarray"`
	V1                     Header
	TransactionBodies      []TransactionBody
	TransactionWitnessSets []TransactionWitnessSet
	AuxiliaryDataSet       hash_map.HashMap //map[TransactionIndex]AuxiliaryData
}

type Transaction struct {
	_  interface{} `cbor:",toarray"`
	V1 TransactionBody
	V2 TransactionWitnessSet
	V3 *AuxiliaryData
}

type TransactionIndex uint

type Header struct {
	_             interface{} `cbor:",toarray"`
	V1            HeaderBody
	BodySignature KesSignature
}

type HeaderBody struct {
	_             interface{} `cbor:",toarray"`
	BlockNumber   uint
	Slot          uint
	PrevHash      *Hash32
	IssuerVkey    Vkey
	VrfVkey       VrfVkey
	NonceVrf      VrfCert
	LeaderVrf     VrfCert
	BlockBodySize uint
	BlockBodyHash Hash32
	V10           OperationalCert
	V11           ProtocolVersion
}

type OperationalCert struct {
	V1HotVkey        KesVkey
	V2SequenceNumber uint
	V3KesPeriod      uint
	V4Sigma          Signature
}

type ProtocolVersion struct {
	V1Uint uint
	V2Uint uint
}

type TransactionBody struct {
	V1SetTransactionInput   SetTransactionInput `cbor:"0,keyasint"`
	V2TransactionOutputList []TransactionOutput `cbor:"1,keyasint"`
	V3Coin                  Coin                `cbor:"2,keyasint"`
	V4Uint                  *uint               `cbor:"3,keyasint,omitempty"`
	V5CertificateList       *[]Certificate      `cbor:"4,keyasint,omitempty"`
	V6Withdrawals           *Withdrawals        `cbor:"5,keyasint,omitempty"`
	V7Update                *Update             `cbor:"6,keyasint,omitempty"`
	V8MetadataHash          *MetadataHash       `cbor:"7,keyasint,omitempty"`
	V9Uint                  *uint               `cbor:"8,keyasint,omitempty"`
	V10Mint                 *Mint               `cbor:"9,keyasint,omitempty"`
}

type TransactionInput struct {
	_             interface{} `cbor:",toarray"`
	TransactionId Hash32
	Index         uint
}

type TransactionOutput struct {
	_      interface{} `cbor:",toarray"`
	V1     Address
	Amount Value
}

type RewardAccount []byte

type Certificate struct {
	V1 *StakeRegistration

	V2 *StakeDeregistration

	V3 *StakeDelegation

	V4 *PoolRegistration

	V5 *PoolRetirement

	V6 *GenesisKeyDelegation

	V7 *MoveInstantaneousRewardsCert
}

type MoveInstantaneousReward struct {
	_  interface{} `cbor:",toarray"`
	V1 int
	V2 hash_map.HashMap //map[StakeCredential]Coin
}

type PoolParams struct {
	V1Operator      PoolKeyhash
	V2VrfKeyhash    VrfKeyhash
	V3Pledge        Coin
	V4Cost          Coin
	V5Margin        UnitInterval
	V6RewardAccount RewardAccount
	V7PoolOwners    SetAddrKeyhash
	V8Relays        []Relay
	V9PoolMetadata  *PoolMetadata
}

type Port uint

type Ipv4 []byte

type Ipv6 []byte

type DnsName string

type Relay struct {
	V1 *SingleHostAddr

	V2 *SingleHostName

	V3 *MultiHostName
}

type PoolMetadata struct {
	_  interface{} `cbor:",toarray"`
	V1 Url
	V2 MetadataHash
}

type Url string

type Withdrawals hash_map.HashMap //map[RewardAccount]Coin

type Update struct {
	_  interface{} `cbor:",toarray"`
	V1 ProposedProtocolParameterUpdates
	V2 Epoch
}

type ProposedProtocolParameterUpdates hash_map.HashMap //map[Genesishash]ProtocolParamUpdate

type ProtocolParamUpdate struct {
	V1Uint                 *uint              `cbor:"0,keyasint,omitempty"`
	V2Uint                 *uint              `cbor:"1,keyasint,omitempty"`
	V3Uint                 *uint              `cbor:"2,keyasint,omitempty"`
	V4Uint                 *uint              `cbor:"3,keyasint,omitempty"`
	V5Uint                 *uint              `cbor:"4,keyasint,omitempty"`
	V6Coin                 *Coin              `cbor:"5,keyasint,omitempty"`
	V7Coin                 *Coin              `cbor:"6,keyasint,omitempty"`
	V8Epoch                *Epoch             `cbor:"7,keyasint,omitempty"`
	V9Uint                 *uint              `cbor:"8,keyasint,omitempty"`
	V10Rational            *Rational          `cbor:"9,keyasint,omitempty"`
	V11UnitInterval        *UnitInterval      `cbor:"10,keyasint,omitempty"`
	V12UnitInterval        *UnitInterval      `cbor:"11,keyasint,omitempty"`
	V13UnitInterval        *UnitInterval      `cbor:"12,keyasint,omitempty"`
	V14Nonce               *Nonce             `cbor:"13,keyasint,omitempty"`
	V15ProtocolVersionList *[]ProtocolVersion `cbor:"14,keyasint,omitempty"`
	V16Coin                *Coin              `cbor:"15,keyasint,omitempty"`
}

type TransactionWitnessSet struct {
	V1VkeywitnessList      *[]Vkeywitness      `cbor:"0,keyasint,omitempty"`
	V2NativeScriptList     *[]NativeScript     `cbor:"1,keyasint,omitempty"`
	V3BootstrapWitnessList *[]BootstrapWitness `cbor:"2,keyasint,omitempty"`
}

type TransactionMetadatum struct {
	V1TransactionMetadatumMap  *TransactionMetadatumAdditionalType0
	V2TransactionMetadatumList *TransactionMetadatumAdditionalType1
	V3Int                      *int
	V4Bytes                    *Bytes
	V5Text                     *Text
}

type TransactionMetadatumLabel uint

type AuxiliaryData struct {
	V1TransactionMetadatumLabelMap *AuxiliaryDataAdditionalType0
	V2SomeArray                    *AuxiliaryDataAdditionalType1
}

type Vkeywitness struct {
	_  interface{} `cbor:",toarray"`
	V1 Vkey
	V2 Signature
}

type BootstrapWitness struct {
	_          interface{} `cbor:",toarray"`
	PublicKey  Vkey
	Signature  Signature
	ChainCode  []byte
	Attributes []byte
}

type NativeScript struct {
	V1 *ScriptPubkey

	V2 *ScriptAll

	V3 *ScriptAny

	V4 *ScriptNOfK

	V5 *InvalidBefore

	V6 *InvalidHereafter
}

type Coin uint

type PolicyId Scripthash

type AssetName []byte

type Value struct {
	V1Coin      *Coin
	V2SomeArray *ValueAdditionalType0
}

type Mint MultiassetInt64

type Epoch uint

type AddrKeyhash Hash28

type Scripthash Hash28

type GenesisDelegateHash Hash28

type PoolKeyhash Hash28

type Genesishash Hash28

type VrfKeyhash Hash32

type MetadataHash Hash32

type TransactionBodyAllegra struct {
	V1SetTransactionInput          SetTransactionInput        `cbor:"0,keyasint"`
	V2TransactionOutputAllegraList []TransactionOutputAllegra `cbor:"1,keyasint"`
	V3Coin                         Coin                       `cbor:"2,keyasint"`
	V4Uint                         *uint                      `cbor:"3,keyasint,omitempty"`
	V5CertificateList              *[]Certificate             `cbor:"4,keyasint,omitempty"`
	V6Withdrawals                  *Withdrawals               `cbor:"5,keyasint,omitempty"`
	V7Update                       *Update                    `cbor:"6,keyasint,omitempty"`
	V8MetadataHash                 *MetadataHash              `cbor:"7,keyasint,omitempty"`
	V9Uint                         *uint                      `cbor:"8,keyasint,omitempty"`
}

type TransactionOutputAllegra struct {
	_      interface{} `cbor:",toarray"`
	V1     Address
	Amount Coin
}

type AuxiliaryDataAdditionalType0 hash_map.HashMap //map[TransactionMetadatumLabel]TransactionMetadatum

type AuxiliaryDataAdditionalType1 struct {
	_                   interface{}      `cbor:",toarray"`
	TransactionMetadata hash_map.HashMap //map[TransactionMetadatumLabel]TransactionMetadatum
	AuxiliaryScripts    []NativeScript
}

type ValueAdditionalType0 struct {
	_  interface{} `cbor:",toarray"`
	V1 Coin
	V2 MultiassetUint
}

type TransactionMetadatumAdditionalType0 hash_map.HashMap //map[TransactionMetadatum]TransactionMetadatum

type TransactionMetadatumAdditionalType1 []TransactionMetadatum

type Bytes []byte

type Text string

type MultiHostName struct {
	V0 int32
	V1 DnsName
}
type ScriptAll struct {
	V0 int32
	V1 []NativeScript
}
type NonceComposition0 struct {
	V0 int32
	V1 []byte
}
type StakeRegistration struct {
	V0 int32
	V1 StakeCredential
}
type MoveInstantaneousRewardsCert struct {
	V0 int32
	V1 MoveInstantaneousReward
}
type InvalidBefore struct {
	V0 int32
	V1 uint
}
type StakeCredentialComposition0 struct {
	V0 int32
	V1 AddrKeyhash
}
type ScriptAny struct {
	V0 int32
	V1 []NativeScript
}
type StakeCredentialComposition1 struct {
	V0 int32
	V1 Scripthash
}
type SingleHostName struct {
	V0 int32
	V1 *Port
	V2 DnsName
}
type ScriptPubkey struct {
	V0 int32
	V1 AddrKeyhash
}
type StakeDeregistration struct {
	V0 int32
	V1 StakeCredential
}
type StakeDelegation struct {
	V0 int32
	V1 StakeCredential
	V2 PoolKeyhash
}
type PoolRegistration struct {
	V0 int32
	V1 PoolParams
}
type PoolRetirement struct {
	V0 int32
	V1 PoolKeyhash
	V2 Epoch
}
type GenesisKeyDelegation struct {
	V0 int32
	V1 Genesishash
	V2 GenesisDelegateHash
	V3 VrfKeyhash
}
type SingleHostAddr struct {
	V0 int32
	V1 *Port
	V2 *Ipv4
	V3 *Ipv6
}
type ScriptNOfK struct {
	V0 int32
	V1 uint
	V2 []NativeScript
}
type InvalidHereafter struct {
	V0 int32
	V1 uint
}
type SetTransactionInput []TransactionInput
type SetAddrKeyhash []AddrKeyhash
type MultiassetInt64 hash_map.HashMap //map[PolicyId]hash_map.HashMap //map[AssetName]int64
type MultiassetUint hash_map.HashMap  //map[PolicyId]hash_map.HashMap //map[AssetName]uint
