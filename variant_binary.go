package variant

import (
	"encoding/binary"
	"fmt"
	"math"
)

// binaryFormatMarker discriminates binary format from legacy JSON (which starts with '{', '"', digits, etc.).
const binaryFormatMarker = 0x01

// ─── Encoder ──────────────────────────────────────────────────────────────────

// AppendBinary appends the binary encoding of v (with format marker) to dst.
// For batch encoding, reuse dst across calls to avoid allocations.
func (v Variant) AppendBinary(dst []byte) []byte {
	dst = append(dst, binaryFormatMarker)
	return v.appendValue(dst)
}

// appendValue appends [type][payload] without format marker. Used recursively
// for List/Map elements so nested values don't carry redundant markers.
func (v Variant) appendValue(dst []byte) []byte {
	dst = append(dst, byte(v.variantType))
	switch v.variantType {
	case TypeEmpty:
		// no payload
	case TypeBool:
		if v.numberValue != 0 {
			dst = append(dst, 1)
		} else {
			dst = append(dst, 0)
		}
	case TypeInt64, TypeUInt64, TypeFloat64:
		var b [8]byte
		binary.BigEndian.PutUint64(b[:], uint64(v.numberValue))
		dst = append(dst, b[:]...)
	case TypeString:
		s := v.complexValue.(string)
		dst = appendString(dst, s)
	case TypeList:
		list := v.complexValue.([]Variant)
		var lb [4]byte
		binary.BigEndian.PutUint32(lb[:], uint32(len(list)))
		dst = append(dst, lb[:]...)
		for i := range list {
			dst = list[i].appendValue(dst)
		}
	case TypeMap:
		// Raw (map[string]any from New) and typed (map[string]Variant) both
		// serialize the same way: key, then value bytes.
		switch m := v.complexValue.(type) {
		case map[string]any:
			var lb [4]byte
			binary.BigEndian.PutUint32(lb[:], uint32(len(m)))
			dst = append(dst, lb[:]...)
			for k, val := range m {
				dst = appendString(dst, k)
				dst = appendRawValue(dst, val)
			}
		case map[string]Variant:
			var lb [4]byte
			binary.BigEndian.PutUint32(lb[:], uint32(len(m)))
			dst = append(dst, lb[:]...)
			for k, val := range m {
				dst = appendString(dst, k)
				dst = val.appendValue(dst)
			}
		}
	}
	return dst
}

func appendString(dst []byte, s string) []byte {
	var lb [4]byte
	binary.BigEndian.PutUint32(lb[:], uint32(len(s)))
	dst = append(dst, lb[:]...)
	dst = append(dst, s...)
	return dst
}

// appendRawValue serializes a raw map[string]any field value. Strings are
// written directly (no string→interface boxing); other values fall through to
// NewRawValue.appendValue.
func appendRawValue(dst []byte, val any) []byte {
	switch t := val.(type) {
	case string:
		return appendString(dst, t)
	case bool:
		return NewBool(t).appendValue(dst)
	case float64:
		return NewFloat64(t).appendValue(dst)
	case float32:
		return NewFloat64(float64(t)).appendValue(dst)
	case int:
		return NewInt(t).appendValue(dst)
	case int8:
		return NewInt(int(t)).appendValue(dst)
	case int16:
		return NewInt(int(t)).appendValue(dst)
	case int32:
		return NewInt(int(t)).appendValue(dst)
	case int64:
		return NewInt64(t).appendValue(dst)
	case uint:
		return NewUInt64(uint64(t)).appendValue(dst)
	case uint8:
		return NewUInt64(uint64(t)).appendValue(dst)
	case uint16:
		return NewUInt64(uint64(t)).appendValue(dst)
	case uint32:
		return NewUInt64(uint64(t)).appendValue(dst)
	case uint64:
		return NewUInt64(t).appendValue(dst)
	case []byte:
		return appendString(dst, string(t))
	default:
		return NewRawValue(val).appendValue(dst)
	}
}

// MarshalBinary encodes the variant with format marker. Convenience wrapper
// around AppendBinary.
func (v Variant) MarshalBinary() ([]byte, error) {
	return v.AppendBinary(nil), nil
}

// ─── Decoder ──────────────────────────────────────────────────────────────────

// IsBinaryFormat reports whether data starts with the binary format marker.
func IsBinaryFormat(data []byte) bool {
	return len(data) > 0 && data[0] == binaryFormatMarker
}

// UnmarshalBinary decodes a binary-encoded variant with format marker.
// Returns the variant and the total number of bytes consumed.
func UnmarshalBinary(data []byte) (Variant, int, error) {
	if len(data) < 2 || data[0] != binaryFormatMarker {
		return Variant{}, 0, fmt.Errorf("not a binary variant")
	}
	v, n, err := readValue(data[1:])
	if err != nil {
		return Variant{}, 0, err
	}
	return v, n + 1, nil // +1 for the marker byte
}

// readValue reads a [type][payload] variant from data. Does not expect a format marker.
func readValue(data []byte) (Variant, int, error) {
	if len(data) < 1 {
		return Variant{}, 0, fmt.Errorf("empty variant data")
	}
	t := Type(data[0])
	switch t {
	case TypeEmpty:
		return NewEmpty(), 1, nil
	case TypeBool:
		if len(data) < 2 {
			return Variant{}, 0, fmt.Errorf("short bool")
		}
		return NewBool(data[1] != 0), 2, nil
	case TypeInt64:
		if len(data) < 9 {
			return Variant{}, 0, fmt.Errorf("short int64")
		}
		return NewInt64(int64(binary.BigEndian.Uint64(data[1:9]))), 9, nil
	case TypeUInt64:
		if len(data) < 9 {
			return Variant{}, 0, fmt.Errorf("short uint64")
		}
		return NewUInt64(binary.BigEndian.Uint64(data[1:9])), 9, nil
	case TypeFloat64:
		if len(data) < 9 {
			return Variant{}, 0, fmt.Errorf("short float64")
		}
		bits := binary.BigEndian.Uint64(data[1:9])
		return NewFloat64(math.Float64frombits(bits)), 9, nil
	case TypeString:
		return readStringValue(data)
	case TypeList:
		return readListValue(data)
	case TypeMap:
		return readMapValue(data)
	default:
		return Variant{}, 0, fmt.Errorf("unknown variant type: %d", t)
	}
}

func readStringValue(data []byte) (Variant, int, error) {
	if len(data) < 5 {
		return Variant{}, 0, fmt.Errorf("short string header")
	}
	slen := binary.BigEndian.Uint32(data[1:5])
	if uint32(len(data)) < 5+slen {
		return Variant{}, 0, fmt.Errorf("short string data")
	}
	// unsafe string conversion: the backing array of data lives as long as the caller needs it,
	// but Variant may outlive the input buffer, so we must copy.
	s := string(data[5 : 5+slen])
	return NewString(s), 5 + int(slen), nil
}

func readListValue(data []byte) (Variant, int, error) {
	if len(data) < 5 {
		return Variant{}, 0, fmt.Errorf("short list header")
	}
	count := binary.BigEndian.Uint32(data[1:5])
	offset := 5
	list := make([]Variant, count)
	for i := uint32(0); i < count; i++ {
		elem, n, err := readValue(data[offset:])
		if err != nil {
			return Variant{}, 0, err
		}
		list[i] = elem
		offset += n
	}
	return NewValueList(list), offset, nil
}

func readMapValue(data []byte) (Variant, int, error) {
	if len(data) < 5 {
		return Variant{}, 0, fmt.Errorf("short map header")
	}
	count := binary.BigEndian.Uint32(data[1:5])
	offset := 5
	mp := make(map[string]Variant, count)
	for i := uint32(0); i < count; i++ {
		if len(data)-offset < 4 {
			return Variant{}, 0, fmt.Errorf("short map key header")
		}
		klen := binary.BigEndian.Uint32(data[offset : offset+4])
		offset += 4
		if len(data)-offset < int(klen) {
			return Variant{}, 0, fmt.Errorf("short map key data")
		}
		key := string(data[offset : offset+int(klen)])
		offset += int(klen)
		val, n, err := readValue(data[offset:])
		if err != nil {
			return Variant{}, 0, err
		}
		mp[key] = val
		offset += n
	}
	return NewValueMap(mp), offset, nil
}
