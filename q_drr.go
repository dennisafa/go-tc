package tc

import (
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaDrrUnspec = iota
	tcaDrrQuantum
)

// Drr contains attributes of the drr discipline
type Drr struct {
	Quantum uint32
}

// unmarshalDrr parses the Drr-encoded data and stores the result in the value pointed to by info.
func unmarshalDrr(data []byte, info *Drr) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaDrrQuantum:
			info.Quantum = ad.Uint32()
		default:
			return fmt.Errorf("UnmarshalDrr()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalDrr returns the binary encoding of Qfq
func marshalDrr(info *Drr) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Drr options are missing")
	}

	// TODO: improve logic and check combinations
	if info.Quantum != 0 {
		options = append(options, tcOption{Interpretation: vtUint32, Type: tcaDrrQuantum, Data: info.Quantum})
	}
	return marshalAttributes(options)
}
