package tc

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/mdlayher/netlink"
)

const (
	tcaHfscUnspec = iota
	tcaHfscRsc
	tcaHfscFsc
	tcaHfscUsc
)

// Hfsc contains attributes of the hfsc discipline
type Hfsc struct {
	Rsc *ServiceCurve
	Fsc *ServiceCurve
	Usc *ServiceCurve
}

// unmarshalHfsc parses the Hfsc-encoded data and stores the result in the value pointed to by info.
func unmarshalHfsc(data []byte, info *Hfsc) error {
	ad, err := netlink.NewAttributeDecoder(data)
	if err != nil {
		return err
	}
	ad.ByteOrder = nativeEndian
	for ad.Next() {
		switch ad.Type() {
		case tcaHfscRsc:
			curve := &ServiceCurve{}
			if err := extractServiceCurve(ad.Bytes(), curve); err != nil {
				return err
			}
			info.Rsc = curve
		case tcaHfscFsc:
			curve := &ServiceCurve{}
			if err := extractServiceCurve(ad.Bytes(), curve); err != nil {
				return err
			}
			info.Fsc = curve
		case tcaHfscUsc:
			curve := &ServiceCurve{}
			if err := extractServiceCurve(ad.Bytes(), curve); err != nil {
				return err
			}
			info.Usc = curve
		default:
			return fmt.Errorf("unmarshalHfsc()\t%d\n\t%v", ad.Type(), ad.Bytes())
		}
	}
	return nil
}

// marshalHfsc returns the binary encoding of Hfsc
func marshalHfsc(info *Hfsc) ([]byte, error) {
	options := []tcOption{}

	if info == nil {
		return []byte{}, fmt.Errorf("Hfsc options are missing")
	}

	// TODO: improve logic and check combinations

	if info.Rsc != nil {
		data, err := validateServiceCurve(info.Rsc)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaHfscRsc, Data: data})
	}
	if info.Fsc != nil {
		data, err := validateServiceCurve(info.Fsc)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaHfscFsc, Data: data})
	}
	if info.Usc != nil {
		data, err := validateServiceCurve(info.Usc)
		if err != nil {
			return []byte{}, err
		}
		options = append(options, tcOption{Interpretation: vtBytes, Type: tcaHfscUsc, Data: data})
	}

	return marshalAttributes(options)
}

// ServiceCurve from include/uapi/linux/pkt_sched.h
type ServiceCurve struct {
	M1 uint32
	D  uint32
	M2 uint32
}

func extractServiceCurve(data []byte, info *ServiceCurve) error {
	b := bytes.NewReader(data)
	return binary.Read(b, nativeEndian, info)
}

func validateServiceCurve(info *ServiceCurve) ([]byte, error) {
	var buf bytes.Buffer
	err := binary.Write(&buf, nativeEndian, *info)
	return buf.Bytes(), err
}
