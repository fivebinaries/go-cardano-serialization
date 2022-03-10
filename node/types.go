package node

// NetworkTip contains parameters from the tip of the network
type NetworkTip struct {
	Slot  uint `json:"slot,omitempty"`
	Epoch uint `json:"epoch,omitempty"`
	Block uint `json:"block,omitempty"`
}
