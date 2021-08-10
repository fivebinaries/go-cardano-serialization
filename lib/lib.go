package lib

type Slot uint32

type TimelockStart struct {
	Slot Slot
}
type TimelockExpiry struct {
	Slot Slot
}

// CertificateIndex index of a cert within a tx
type CertificateIndex uint32
