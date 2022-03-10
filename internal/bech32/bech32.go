//
// From https://github.com/Bitrue-exchange/libada-go/blob/4a319a1e2371ffcdb6877c4f5a76aec0493b151a/internal/bech32/bech32.go
//

package bech32

import (
	"fmt"
	"strings"
)

var charset = "qpzry9x8gf2tvdw0s3jn54khce6mua7l"

var generator = []uint32{0x3b6a57b2, 0x26508e6d, 0x1ea119fa, 0x3d4233dd, 0x2a1462b3}

func polymod(values []byte) uint32 {
	chk := uint32(1)
	for _, v := range values {
		top := chk >> 25
		chk = (chk & 0x1ffffff) << 5
		chk = chk ^ uint32(v)
		for i := 0; i < 5; i++ {
			bit := top >> i & 1
			if bit == 1 {
				chk ^= generator[i]
			}
		}
	}
	return chk
}

func hrpExpand(hrp string) []byte {
	h := []byte(strings.ToLower(hrp))
	var ret []byte
	for _, c := range h {
		ret = append(ret, c>>5)
	}
	ret = append(ret, 0)
	for _, c := range h {
		ret = append(ret, c&31)
	}
	return ret
}

func verifyChecksum(hrp string, data []byte) bool {
	return polymod(append(hrpExpand(hrp), data...)) == 1
}

func createChecksum(hrp string, data []byte) []byte {
	values := append(hrpExpand(hrp), data...)
	values = append(values, []byte{0, 0, 0, 0, 0, 0}...)
	mod := polymod(values) ^ 1
	ret := make([]byte, 6)
	for p := range ret {
		shift := 5 * (5 - p)
		ret[p] = byte(mod>>shift) & 31
	}
	return ret
}

func convertBits(data []byte, frombits, tobits byte, pad bool) ([]byte, error) {
	var ret []byte
	acc := uint32(0)
	bits := byte(0)
	maxv := byte(1<<tobits - 1)
	for idx, value := range data {
		if value>>frombits != 0 {
			return nil, fmt.Errorf("invalid data range: data[%d]=%d (frombits=%d)", idx, value, frombits)
		}
		acc = acc<<frombits | uint32(value)
		bits += frombits
		for bits >= tobits {
			bits -= tobits
			ret = append(ret, byte(acc>>bits)&maxv)
		}
	}
	if pad {
		if bits > 0 {
			ret = append(ret, byte(acc<<(tobits-bits))&maxv)
		}
	} else if bits >= frombits {
		return nil, fmt.Errorf("illegal zero padding")
	} else if byte(acc<<(tobits-bits))&maxv != 0 {
		return nil, fmt.Errorf("non-zero padding")
	}
	return ret, nil
}

func Encode(hrp string, data []byte) (string, error) {
	values, err := convertBits(data, 8, 5, true)
	if err != nil {
		return "", err
	}
	if len(hrp) < 1 {
		return "", fmt.Errorf("invalid HRP: %q", hrp)
	}
	for p, c := range hrp {
		if c < 33 || c > 126 {
			return "", fmt.Errorf("invalid HRP character: hrp[%d]=%d", p, c)
		}
	}
	if strings.ToUpper(hrp) != hrp && strings.ToLower(hrp) != hrp {
		return "", fmt.Errorf("mixed case HRP: %q", hrp)
	}
	lower := strings.ToLower(hrp) == hrp
	hrp = strings.ToLower(hrp)
	var ret strings.Builder
	ret.WriteString(hrp)
	ret.WriteString("1")
	for _, p := range values {
		ret.WriteByte(charset[p])
	}
	for _, p := range createChecksum(hrp, values) {
		ret.WriteByte(charset[p])
	}
	if lower {
		return ret.String(), nil
	}
	return strings.ToUpper(ret.String()), nil
}

func Decode(s string) (hrp string, data []byte, err error) {
	if strings.ToLower(s) != s && strings.ToUpper(s) != s {
		return "", nil, fmt.Errorf("mixed case")
	}
	pos := strings.LastIndex(s, "1")
	if pos < 1 || pos+7 > len(s) {
		return "", nil, fmt.Errorf("separator '1' at invalid position: pos=%d, len=%d", pos, len(s))
	}
	hrp = s[:pos]
	for p, c := range hrp {
		if c < 33 || c > 126 {
			return "", nil, fmt.Errorf("invalid character human-readable part: s[%d]=%d", p, c)
		}
	}
	s = strings.ToLower(s)
	for p, c := range s[pos+1:] {
		d := strings.IndexRune(charset, c)
		if d == -1 {
			return "", nil, fmt.Errorf("invalid character data part: s[%d]=%v", p, c)
		}
		data = append(data, byte(d))
	}
	if !verifyChecksum(hrp, data) {
		return "", nil, fmt.Errorf("invalid checksum")
	}
	data, err = convertBits(data[:len(data)-6], 5, 8, false)
	if err != nil {
		return "", nil, err
	}
	return hrp, data, nil
}
