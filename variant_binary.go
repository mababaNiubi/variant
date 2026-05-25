package variant

import (
	"bytes"
	"fmt"
	"unsafe"

	"github.com/vmihailenco/msgpack/v5"
)

// msgpackFormatMarker is the 1-byte prefix for the WAL binary format.
// It discriminates from legacy JSON (which starts with '{', '"', digits, etc.).
const msgpackFormatMarker = 0x01

// ─── msgpack.Marshaler / Unmarshaler (for JsonEncoder, column compression) ───

// MarshalMsgpack implements msgpack.Marshaler. Encodes the variant without
// a format marker — suitable for use inside larger msgpack structures.
func (v Variant) MarshalMsgpack() ([]byte, error) {
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	if err := encodeVariant(enc, v); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// UnmarshalMsgpack implements msgpack.Unmarshaler.
func (v *Variant) UnmarshalMsgpack(data []byte) error {
	dec := msgpack.NewDecoder(bytes.NewReader(data))
	decoded, err := decodeVariant(dec)
	if err != nil {
		return err
	}
	*v = decoded
	return nil
}

// ─── WAL binary format (with format marker) ───

// MarshalBinary encodes the variant with a 1-byte format marker for WAL storage.
func (v Variant) MarshalBinary() ([]byte, error) {
	raw, err := v.MarshalMsgpack()
	if err != nil {
		return nil, err
	}
	out := make([]byte, 1+len(raw))
	out[0] = msgpackFormatMarker
	copy(out[1:], raw)
	return out, nil
}

// UnmarshalBinary decodes a msgpack-encoded variant with format marker.
func UnmarshalBinary(data []byte) (Variant, int, error) {
	if len(data) < 2 || data[0] != msgpackFormatMarker {
		return Variant{}, 0, fmt.Errorf("not a msgpack binary variant")
	}
	var v Variant
	if err := v.UnmarshalMsgpack(data[1:]); err != nil {
		return Variant{}, 0, err
	}
	return v, len(data), nil
}

// IsBinaryFormat reports whether data starts with the msgpack format marker.
func IsBinaryFormat(data []byte) bool {
	return len(data) > 0 && data[0] == msgpackFormatMarker
}

// ─── Recursive encoder (zero-allocation write path) ───

// encodeVariant encodes a single Variant directly to the msgpack encoder.
func encodeVariant(enc *msgpack.Encoder, v Variant) error {
	switch v.variantType {
	case TypeEmpty:
		return enc.EncodeNil()
	case TypeBool:
		return enc.EncodeBool(v.numberValue != 0)
	case TypeInt64:
		return enc.EncodeInt64(v.numberValue)
	case TypeUInt64:
		return enc.EncodeUint64(uint64(v.numberValue))
	case TypeFloat64:
		return enc.EncodeFloat64(*(*float64)(unsafe.Pointer(&v.numberValue)))
	case TypeString:
		return enc.EncodeString(v.complexValue.(string))
	case TypeList:
		list := v.complexValue.([]Variant)
		if err := enc.EncodeArrayLen(len(list)); err != nil {
			return err
		}
		for i := range list {
			if err := encodeVariant(enc, list[i]); err != nil {
				return err
			}
		}
		return nil
	case TypeMap:
		mp := v.complexValue.(map[string]Variant)
		if err := enc.EncodeMapLen(len(mp)); err != nil {
			return err
		}
		for k, val := range mp {
			if err := enc.EncodeString(k); err != nil {
				return err
			}
			if err := encodeVariant(enc, val); err != nil {
				return err
			}
		}
		return nil
	default:
		return fmt.Errorf("unknown variant type: %v", v.variantType)
	}
}

// ─── Recursive decoder ───

// decodeVariant decodes a single Variant from the msgpack decoder.
func decodeVariant(dec *msgpack.Decoder) (Variant, error) {
	raw, err := dec.DecodeInterface()
	if err != nil {
		return Variant{}, err
	}
	return decodeInterface(raw)
}

// decodeInterface converts a Go value from DecodeInterface to Variant.
func decodeInterface(raw interface{}) (Variant, error) {
	switch v := raw.(type) {
	case nil:
		return NewEmpty(), nil
	case bool:
		return NewBool(v), nil
	case int64:
		return NewInt64(v), nil
	case int8:
		return NewInt64(int64(v)), nil
	case uint64:
		return NewUInt64(v), nil
	case float64:
		return NewFloat64(v), nil
	case string:
		return NewString(v), nil
	case []interface{}:
		list := make([]Variant, len(v))
		for i, item := range v {
			var err error
			list[i], err = decodeInterface(item)
			if err != nil {
				return Variant{}, err
			}
		}
		return NewValueList(list), nil
	case map[string]interface{}:
		mp := make(map[string]Variant, len(v))
		for k, item := range v {
			var err error
			mp[k], err = decodeInterface(item)
			if err != nil {
				return Variant{}, err
			}
		}
		return NewValueMap(mp), nil
	default:
		return Variant{}, fmt.Errorf("unsupported type: %T", raw)
	}
}
