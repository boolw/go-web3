package abi

import (
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"
	"strconv"

	"github.com/boolw/go-web3"
	"github.com/mitchellh/mapstructure"
)

// Decode decodes the input with a given type
func Decode(t *Type, input []byte) (val interface{}, err error) {
	//defer func() {
	//	if e := recover(); e != nil {
	//		err = fmt.Errorf("decode 0x%s error : %v", hex.EncodeToString(input),e)
	//	}
	//}()
	val, _, err = decode(t, input)
	return val, err
}

// DecodeStruct decodes the input with a type to a struct
func DecodeStruct(t *Type, input []byte, out interface{}) error {
	val, err := Decode(t, input)
	if err != nil {
		return err
	}
	if err := mapstructure.Decode(val, out); err != nil {
		return err
	}
	return nil
}

func decode(t *Type, input []byte) (interface{}, []byte, error) {
	var length int
	var err error

	if t.isVariableInput() {
		length, err = readLength(input)
		if err != nil {
			return nil, nil, err
		}
	}

	switch t.kind {
	case KindTuple:
		return decodeTuple(t, input)

	case KindSlice:
		if data,e := readSlice(input,32,0);e != nil {
			return nil, nil, e
		}else{
			return decodeArraySlice(t, data, length)
		}
	case KindArray:
		return decodeArraySlice(t, input, t.size)
	}

	var val interface{}
	switch t.kind {
	case KindBool:
		if data,e := readSlice(input,0,32);e != nil {
			return nil, nil, e
		}else{
			val, err = decodeBool(data)
		}
	case KindInt, KindUInt:
		if data,e := readSlice(input,0,32);e != nil {
			return nil, nil, e
		}else{
			val = readInteger(t, data)
		}
	case KindString:
		if data,e := readSlice(input,32,length);e != nil {
			return nil, nil, e
		}else{
			val = string(data)
		}
	case KindBytes:
		if data,e := readSlice(input,32,length);e != nil {
			return nil, nil, e
		}else{
			val = data
		}
	case KindAddress:
		if data,e := readSlice(input,0,32);e != nil {
			return nil, nil, e
		}else{
			val, err = readAddr(data)
		}
	case KindFixedBytes:
		if data,e := readSlice(input,0,32);e != nil {
			return nil, nil, e
		}else{
			val, err =readFixedBytes(t, data)
		}
	case KindFunction:
		if data,e := readSlice(input,0,32);e != nil {
			return nil, nil, e
		}else{
			val, err =readFunctionType(t, data)
		}
	default:
		return nil, nil, fmt.Errorf("decoding not available for type '%s'", t.kind)
	}

	if data,e := readSlice(input,32,0);e != nil {
		return nil, nil, e
	}else{
		return val, data, err
	}
}

var (
	maxUint256 = big.NewInt(0).Add(
		big.NewInt(0).Exp(big.NewInt(2), big.NewInt(256), nil),
		big.NewInt(-1))
	maxInt256 = big.NewInt(0).Add(
		big.NewInt(0).Exp(big.NewInt(2), big.NewInt(255), nil),
		big.NewInt(-1))
)

func readAddr(b []byte) (web3.Address, error) {
	res := web3.Address{}
	if len(b) != 32 {
		return res, fmt.Errorf("len is not correct")
	}
	copy(res[:], b[12:])
	return res, nil
}

func readSlice(input []byte, start int, length int) ([]byte, error) {
	end := len(input)
	if start < 0{
		start = end + start
	}
	if length < 0 {
		end = end + length
	} else if length > 0 {
		end = start + length
	}
	if len(input) < end{
		return nil,fmt.Errorf("input %d read [%d:%d] error",len(input),start,length)
	}
	return input[start:end], nil
}

func readInteger(t *Type, b []byte) interface{} {
	switch t.t.Kind() {
	case reflect.Uint8:
		return b[len(b)-1]

	case reflect.Uint16:
		return binary.BigEndian.Uint16(b[len(b)-2:])

	case reflect.Uint32:
		return binary.BigEndian.Uint32(b[len(b)-4:])

	case reflect.Uint64:
		return binary.BigEndian.Uint64(b[len(b)-8:])

	case reflect.Int8:
		return int8(b[len(b)-1])

	case reflect.Int16:
		return int16(binary.BigEndian.Uint16(b[len(b)-2:]))

	case reflect.Int32:
		return int32(binary.BigEndian.Uint32(b[len(b)-4:]))

	case reflect.Int64:
		return int64(binary.BigEndian.Uint64(b[len(b)-8:]))

	default:
		ret := new(big.Int).SetBytes(b)
		if t.kind == KindUInt {
			return ret
		}

		if ret.Cmp(maxInt256) > 0 {
			ret.Add(maxUint256, big.NewInt(0).Neg(ret))
			ret.Add(ret, big.NewInt(1))
			ret.Neg(ret)
		}
		return ret
	}
}

func readFunctionType(t *Type, word []byte) ([24]byte, error) {
	res := [24]byte{}
	if !allZeros(word[24:32]) {
		return res, fmt.Errorf("function type expects the last 8 bytes to be empty but found: %b", word[24:32])
	}
	copy(res[:], word[0:24])
	return res, nil
}

func readFixedBytes(t *Type, word []byte) (interface{}, error) {
	array := reflect.New(t.t).Elem()
	reflect.Copy(array, reflect.ValueOf(word[0:t.size]))
	return array.Interface(), nil
}

func decodeTuple(t *Type, data []byte) (interface{}, []byte, error) {
	res := make(map[string]interface{})

	orig := data
	origLen := len(orig)
	for indx, arg := range t.tuple {
		entry := data
		if arg.Elem.isDynamicType() {
			offset, err := readOffset(data, origLen)
			if err != nil {
				return nil, nil, err
			}
			entry = orig[offset:]
		}

		val, tail, err := decode(arg.Elem, entry)
		if err != nil {
			return nil, nil, err
		}

		if !arg.Elem.isDynamicType() {
			data = tail
		} else {
			data = data[32:]
		}

		name := arg.Name
		if name == "" {
			name = strconv.Itoa(indx)
		}
		if _, ok := res[name]; !ok {
			res[name] = val
		} else {
			return nil, nil, fmt.Errorf("tuple with repeated values")
		}
	}
	return res, data, nil
}

func decodeArraySlice(t *Type, data []byte, size int) (interface{}, []byte, error) {
	if size < 0 {
		return nil, nil, fmt.Errorf("size is lower than zero")
	}
	if 32*size > len(data) {
		return nil, nil, fmt.Errorf("size is too big")
	}

	var res reflect.Value
	if t.kind == KindSlice {
		res = reflect.MakeSlice(t.t, size, size)
	} else if t.kind == KindArray {
		res = reflect.New(t.t).Elem()
	}

	orig := data
	origLen := len(orig)
	for indx := 0; indx < size; indx++ {
		isDynamic := t.elem.isDynamicType()

		entry := data
		if isDynamic {
			offset, err := readOffset(data, origLen)
			if err != nil {
				return nil, nil, err
			}
			entry = orig[offset:]
		}

		val, tail, err := decode(t.elem, entry)
		if err != nil {
			return nil, nil, err
		}

		if !isDynamic {
			data = tail
		} else {
			data = data[32:]
		}
		res.Index(indx).Set(reflect.ValueOf(val))
	}
	return res.Interface(), data, nil
}

func decodeBool(data []byte) (interface{}, error) {
	switch data[31] {
	case 0:
		return false, nil
	case 1:
		return true, nil
	default:
		return false, fmt.Errorf("bad boolean")
	}
}

func readOffset(data []byte, len int) (int, error) {
	offsetBig := big.NewInt(0).SetBytes(data[0:32])
	if offsetBig.BitLen() > 63 {
		return 0, fmt.Errorf("offset larger than int64: %v", offsetBig.Int64())
	}
	offset := int(offsetBig.Int64())
	if offset > len {
		return 0, fmt.Errorf("offset insufficient %v require %v", len, offset)
	}
	return offset, nil
}

func readLength(data []byte) (int, error) {
	input ,err := readSlice(data,0,32)
	if err != nil {
		return 0, err
	}
	lengthBig := big.NewInt(0).SetBytes(input)
	if lengthBig.BitLen() > 63 {
		return 0, fmt.Errorf("length larger than int64: %v", lengthBig.Int64())
	}
	length := int(lengthBig.Uint64())
	if length > len(data) {
		return 0, fmt.Errorf("length insufficient %v require %v", len(data), length)
	}
	return length, nil
}

func allZeros(b []byte) bool {
	for _, i := range b {
		if i != 0 {
			return false
		}
	}
	return true
}
