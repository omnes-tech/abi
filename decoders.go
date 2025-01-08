package abi

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
)

// DecodeWithSelector decodes bytecode restricted to given selector.
func DecodeWithSelector(selector []byte, typeStrs []string, data []byte) ([]any, error) {
	if !isSelectorIsEqual(selector, data[:4]) {
		return []any{}, fmt.Errorf("invalid selector")
	}

	return Decode(typeStrs, data[4:])
}

// DecodeWithSignature decodes bytecode based on given signature.
func DecodeWithSignature(funcSignature string, data []byte) ([]any, error) {
	typeStrs, err := GetSigTypes(funcSignature)
	if err != nil {
		return []any{}, err
	}

	selector := EncodeSignature(funcSignature)
	if !isSelectorIsEqual(selector, data[:4]) {
		return []any{}, fmt.Errorf("invalid selector")
	}

	return Decode(typeStrs, data[4:])
}

// DecodePacked decodes bytecode following packed format.
// It supports only one dynamic type (either string or bytes)
// as last item in typeStrs array.
func DecodePacked(typeStrs []string, data []byte) ([]any, error) {
	var result []any
	var byteCursor uint64
	for i, typeStr := range typeStrs {
		isTypeDynamic := IsDynamic(typeStr, false)
		isTypeArray, _, err := IsArray(typeStr)
		if err != nil {
			return []any{}, err
		}

		if isTypeDynamic && !isTypeArray && i != len(typeStrs)-1 {
			return []any{}, fmt.Errorf("supports only one dynamic type as last type")
		}

		if typeStr[:3] == "int" {
			if len(typeStr) == 3 {
				typeStr += "256"
			}
		} else {
			if len(typeStr) == 4 {
				typeStr += "256"
			}
		}
		var byteLength int
		if !isTypeDynamic {
			byteLength = validCoreTypes[typeStr].ByteLength
		} else {
			byteLength = len(data[byteCursor:])
		}

		val, err := decodePacked(typeStr, data[byteCursor:byteCursor+uint64(byteLength)])
		if err != nil {
			return []any{}, err
		}

		_, ok := val.([]byte)
		if ok {
			val = common.Bytes2Hex(val.([]byte))
		}

		result = append(result, val)
		byteCursor += uint64(byteLength)
	}

	return result, nil
}

// Decode decodes bytecode to given type strings
func Decode(typeStrs []string, data []byte) ([]any, error) {

	var result []any
	var byteCursor uint64
	for _, typeStr := range typeStrs {

		isTypeTuple, splitedTypes, err := IsTuple(typeStr)
		if err != nil {
			return []any{}, err
		}

		var offset int
		var init int
		var end int
		isTypeDynamic := IsDynamic(typeStr, isTypeTuple)
		if isTypeDynamic {
			offset = int(bigInt.SetBytes(data[byteCursor : byteCursor+32]).Uint64())

			var dynamicSize int
			if !isTypeTuple {
				dynamicSize = int(bigInt.SetBytes(data[offset : offset+32]).Uint64())
			}
			init = offset
			end = offset + dynamicSize

		} else {
			init = int(byteCursor)
			end = int(byteCursor + 32)
		}

		isTypeArray, givenArraySize, err := IsArray(typeStr)
		if err != nil {
			return []any{}, err
		}

		if isTypeArray {
			if givenArraySize != 0 {
				offset = int(bigInt.SetBytes(data[byteCursor : byteCursor+32]).Uint64())
				typeStr = typeStr[:strings.LastIndex(typeStr, "[")]

				var arrayTypeStrs []string
				for j := 0; j < givenArraySize; j++ {
					arrayTypeStrs = append(arrayTypeStrs, typeStr)
				}

				innerResult, err := Decode(arrayTypeStrs, data[offset:offset+32*givenArraySize])
				if err != nil {
					return []any{}, err
				}

				result = append(result, innerResult)
			} else {
				arraySize := int(bigInt.SetBytes(data[offset : offset+32]).Uint64())
				typeStr = typeStr[:strings.LastIndex(typeStr, "[")]

				var arrayTypeStrs []string
				for j := 0; j < arraySize; j++ {
					arrayTypeStrs = append(arrayTypeStrs, typeStr)
				}

				innerResult, err := Decode(arrayTypeStrs, data[offset+32:offset+32*(arraySize+1)])
				if err != nil {
					return []any{}, err
				}

				result = append(result, innerResult)
			}

		} else if isTypeTuple {
			innerResult, err := Decode(splitedTypes, data[init:end])
			if err != nil {
				return []any{}, err
			}

			result = append(result, innerResult)
		} else {
			val, err := decode(typeStr, data[init:end])
			if err != nil {
				return []any{}, err
			}

			result = append(result, val)
		}

		byteCursor += 32
	}

	return result, nil
}

// decode decodes give bytecode slice to specified type.
func decode(typeStr string, data []byte) (any, error) {
	var decoded any
	var err error
	if typeStr == "string" || typeStr == "bytes" {
		byteLengthBigInt := new(big.Int)
		byteLength := data[:32]
		byteLengthBigInt.SetBytes(byteLength)

		decoded, err = decodePacked(typeStr, data[32:32+byteLengthBigInt.Uint64()])
		if err != nil {
			return nil, err
		}
	} else {
		decoded, err = decodePacked(typeStr, data)
		if err != nil {
			return nil, err
		}
	}

	return decoded, nil
}

// decodePacked decodes bytecode slice to given type considering
// packed format.
func decodePacked(typeStr string, data []byte) (any, error) {

	switch typeStr {
	case "address":
		if len(data) < validCoreTypes[typeStr].ByteLength {
			return nil, fmt.Errorf("data byte size is too short for %v. Length: %d", typeStr, len(data))
		}

		return common.BytesToAddress(data).Hex(), nil

	case "bool":
		if len(data) < validCoreTypes[typeStr].ByteLength {
			return nil, fmt.Errorf("data byte size is too short for %v. Length: %d", typeStr, len(data))
		}

		return data[len(data)-1] == 1, nil
	case "string": // @follow-up check this later
		return string(data), nil
	default:
		if typeStr[:3] == "int" || typeStr[:4] == "uint" {
			var index int
			if typeStr[:3] == "int" {
				index = 3
				if len(typeStr) == 3 {
					typeStr += "256"
				}
			} else {
				index = 4
				if len(typeStr) == 4 {
					typeStr += "256"
				}
			}

			bits, err := strconv.Atoi(typeStr[index:])
			if err != nil {
				return []byte{}, fmt.Errorf("error getting bits from %s: %v", typeStr, err)
			}
			if bits%8 != 0 {
				return []byte{}, fmt.Errorf("invalid bits value: %v, bits = %v", typeStr, bits)
			}

			if len(data) < validCoreTypes[typeStr].ByteLength {
				return nil, fmt.Errorf("data byte size is too short for %v. Length: %d", typeStr, len(data))
			}

			decoded := new(big.Int)
			return decoded.SetBytes(data), nil
		} else if typeStr[:5] == "bytes" {
			if len(typeStr) > 5 {
				bytesSize, err := strconv.Atoi(typeStr[5:])
				if err != nil {
					return []byte{}, fmt.Errorf("error getting bytes size: %v", typeStr)
				}

				if bytesSize < 1 || bytesSize > 32 {
					return []byte{}, fmt.Errorf("invalid byte size: %v", typeStr)
				}

				if len(data) < validCoreTypes[typeStr].ByteLength {
					return nil, fmt.Errorf("data byte size is too short for %v. Length: %d", typeStr, len(data))
				}
			}

			return data, nil
		} else if typeStr[:5] == "fixed" || typeStr[:6] == "ufixed" { // @note differences in result with fixed/ufixed types
			var index int
			if typeStr[:5] == "fixed" {
				index = 5
				if len(typeStr) == 5 {
					typeStr += "128x18"
				}
			} else {
				index = 6
				if len(typeStr) == 6 {
					typeStr += "128x18"
				}
			}

			sizes := strings.Split(typeStr[index:], "x")
			bits, err := strconv.Atoi(sizes[0])
			if err != nil {
				return []byte{}, fmt.Errorf("error getting bits from %s: %v", typeStr, err)
			}

			if bits%8 != 0 {
				return []byte{}, fmt.Errorf("invalid bits value: %v, bits = %v", typeStr, bits)
			}

			if len(data) < bits/8 {
				return nil, fmt.Errorf("data byte size is too short for %v. Length: %d", typeStr, len(data))
			}

			decoded := new(big.Float)
			converted, ok := decoded.SetString(string(data))
			if !ok {
				return nil, fmt.Errorf("error converting to big.Float: %v", data)
			}

			return converted, nil

		} else {
			return nil, fmt.Errorf("invalid parameter type: %v", typeStr)
		}
	}
}

// isSelectorIsEqual checks whether given selector is equal to given
// bytecode slice.
func isSelectorIsEqual(selector []byte, data []byte) bool {
	for i := 0; i < len(selector); i++ {
		if selector[i] != data[i] {
			return false
		}
	}
	return true
}
