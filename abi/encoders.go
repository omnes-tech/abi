package abi

import (
	"fmt"
	"math/big"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func EncodeWithSelector(selector [4]byte, typeStrs []string, params ...any) ([]byte, error) {

	if len(typeStrs) != len(params) {
		return []byte{}, fmt.Errorf("number of parameter types and given paramenters mismatch")
	}

	encodedParams, err := Encode(typeStrs, params...)
	if err != nil {
		return []byte{}, err
	}

	encoded, err := EncodePacked([]string{"bytes4", "bytes"}, selector, encodedParams)
	if err != nil {
		return []byte{}, err
	}

	return encoded, nil

}

func EncodeWithSignature(funcSignature string, params ...any) ([]byte, error) {
	if funcSignature == "" {
		return []byte{}, nil
	}

	selector := EncodeSelector(funcSignature)
	paramTypes, err := getSigTypes(funcSignature)
	if err != nil {
		return []byte{}, err
	}

	if len(paramTypes) != len(params) {
		return []byte{}, fmt.Errorf("number of parameter types and given paramenters mismatch")
	}

	var encodedParams []byte
	if len(paramTypes) > 0 {
		encodedParams, err = Encode(paramTypes, params...)
		if err != nil {
			return []byte{}, err
		}
	}

	encoded, err := EncodePacked([]string{"bytes4", "bytes"}, selector, encodedParams)
	if err != nil {
		return []byte{}, err
	}

	return encoded, nil

}

func EncodeSelector(funcSignature string) []byte {
	return crypto.Keccak256([]byte(funcSignature))[:4]
}

func Encode(typeStrs []string, values ...any) ([]byte, error) {
	if len(typeStrs) != len(values) {
		return []byte{}, fmt.Errorf("typeStrs and values must have the same length. typeStrs: %v, values: %v", typeStrs, values)
	}

	var rawHeadChunks [][]byte
	var tailChunks [][]byte
	for i, typeStr := range typeStrs {
		var encoded []byte

		isTypeTuple, splitedTypes, err := isTuple(typeStr)
		if err != nil {
			return []byte{}, err
		}

		isTypeArray, arraySize, err := isArray(typeStr)
		if err != nil {
			return []byte{}, err
		}

		if isTypeArray {
			if arraySize != 0 && len(values[i].([]any)) != arraySize {
				return nil, fmt.Errorf("array size mismatch")
			}
			openBracketIndex := strings.Index(typeStr, "[")

			var arrayTypes []string
			arrayValues := values[i].([]any)
			for j := 0; j < len(arrayValues); j++ {
				arrayTypes = append(arrayTypes, typeStr[:openBracketIndex])
			}

			encoded, err = Encode(arrayTypes, arrayValues...)
			if err != nil {
				return []byte{}, err
			}

			arraySize := big.NewInt(int64(len(arrayValues)))
			encoded = append(common.LeftPadBytes(arraySize.Bytes(), 32), encoded...)
		} else if isTypeTuple {
			encoded, err = Encode(splitedTypes, values[i].([]any)...)
			if err != nil {
				return []byte{}, err
			}
		} else {
			encoded, err = encode(typeStr, values[i])
			if err != nil {
				return []byte{}, err
			}
		}

		if !isDynamic(typeStr) {
			rawHeadChunks = append(rawHeadChunks, encoded)
			tailChunks = append(tailChunks, nil)
		} else {
			rawHeadChunks = append(rawHeadChunks, nil)
			tailChunks = append(tailChunks, encoded)
		}
	}

	headLength := calculateHeadLength(rawHeadChunks)
	tailOffsets := calculateTailOffsets(tailChunks)
	headChunks, err := buildHeadChunks(rawHeadChunks, tailOffsets, headLength)
	if err != nil {
		return []byte{}, err
	}

	final := joinChunks(headChunks, tailChunks)

	return final, nil
}

func EncodePacked(typeStrs []string, values ...any) ([]byte, error) {
	if len(typeStrs) != len(values) {
		return []byte{}, fmt.Errorf("typeStrs and values must have the same length. typeStrs: %v, values: %v", typeStrs, values)
	}

	var result []byte
	for i, typeStr := range typeStrs {
		var encoded []byte

		isTypeTuple, splitedTypes, err := isTuple(typeStr)
		if err != nil {
			return []byte{}, err
		}

		isTypeArray, arraySize, err := isArray(typeStr)
		if err != nil {
			return []byte{}, err
		}

		if isTypeArray {
			if arraySize != 0 && len(values[i].([]any)) != arraySize {
				return nil, fmt.Errorf("array size mismatch")
			}
			openBracketIndex := strings.Index(typeStr, "[")

			var arrayTypes []string
			arrayValues := values[i].([]any)
			for j := 0; j < len(arrayValues); j++ {
				arrayTypes = append(arrayTypes, typeStr[:openBracketIndex])
			}

			encoded, err = EncodePacked(arrayTypes, arrayValues...)
			if err != nil {
				return []byte{}, err
			}
		} else if isTypeTuple {
			encoded, err = EncodePacked(splitedTypes, values[i].([]any)...)
			if err != nil {
				return []byte{}, err
			}
		} else {
			encoded, err = encodePacked(typeStr, values[i])
			if err != nil {
				return []byte{}, err
			}
		}

		result = append(result, encoded...)
	}

	return result, nil
}

func encode(typeStr string, value any) ([]byte, error) {
	encoded, err := encodePacked(typeStr, value)
	if err != nil {
		return []byte{}, err
	}

	if typeStr == "string" || typeStr == "bytes" {
		bytesLengthBigInt := big.NewInt(int64(len(encoded)))
		bytesLength := common.LeftPadBytes(bytesLengthBigInt.Bytes(), 32)

		for len(encoded)%32 != 0 {
			encoded = append(encoded, 0x0)
		}

		encoded = append(bytesLength, encoded...)

	} else {
		encoded = common.LeftPadBytes(encoded[:], 32)
	}

	return encoded, nil
}

func encodePacked(typeStr string, value any) ([]byte, error) {

	bytes := make([]byte, 0)
	switch typeStr {
	case "address":
		val, ok := value.(common.Address)
		if !ok {
			return []byte{}, fmt.Errorf("invalid parameter type: %v, %T", typeStr, value)
		}
		bytes = append(bytes, val[:]...)

	case "bool":
		val, ok := value.(bool)
		if !ok {
			return []byte{}, fmt.Errorf("invalid parameter type: %v, %T", typeStr, value)
		}
		if val {
			bytes = append(bytes, []byte{0x0, 0x1}...)
		} else {
			bytes = append(bytes, []byte{0x0, 0x0}...)
		}
	case "string":
		val, ok := value.(string)
		if !ok {
			return []byte{}, fmt.Errorf("invalid parameter type: %v, %T", typeStr, value)
		}
		bytes = append(bytes, []byte(val)...)
	default:
		if typeStr[:3] == "int" || typeStr[:4] == "uint" {
			val, ok := value.(*big.Int)
			if !ok {
				return []byte{}, fmt.Errorf("invalid parameter type: %v, %T", typeStr, value)
			}

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

			if val.Cmp(VALID_CORE_TYPES[typeStr].Max) == 1 || val.Cmp(VALID_CORE_TYPES[typeStr].Min) == -1 {
				return []byte{}, fmt.Errorf("value out of allowed range: %v, %v", typeStr, val)
			}

			bytes = append(bytes, common.LeftPadBytes(val.Bytes(), bits/8)...)

		} else if typeStr[:5] == "bytes" {
			val, ok := value.([]byte)
			if !ok {
				return []byte{}, fmt.Errorf("invalid parameter type: %v, %T", typeStr, value)
			}

			if len(typeStr) > 5 {
				bytesSize, err := strconv.Atoi(typeStr[5:])
				if err != nil {
					return []byte{}, fmt.Errorf("error getting bytes size: %v", typeStr)
				}

				if bytesSize < 1 || bytesSize > 32 {
					return []byte{}, fmt.Errorf("invalid byte size: %v", typeStr)
				}

				if len(val) > bytesSize {
					return []byte{}, fmt.Errorf("value and type bytes size mismatch: type %v; value bytes size %v", typeStr, len(val))
				}

				for len(val)%bytesSize != 0 {
					val = append(val, 0x0)
				}
			}

			bytes = append(bytes, val...)
		} else if typeStr[:5] == "fixed" || typeStr[:6] == "ufixed" { // @note differences in result with fixed/ufixed types
			val, ok := value.(*big.Float)
			if !ok {
				return []byte{}, fmt.Errorf("invalid parameter type: %v, %T", typeStr, value)
			}

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

			fracPlaces, err := strconv.Atoi(sizes[1])
			if err != nil {
				return []byte{}, fmt.Errorf("error getting frac. places from %s: %v", typeStr, err)
			}

			if bits%8 != 0 {
				return []byte{}, fmt.Errorf("invalid bits value: %v, bits = %v", typeStr, bits)
			}

			var min *big.Float
			var max *big.Float
			if typeStr[:5] == "fixed" {
				min, max = computeSignedFixedBounds(int64(bits), float64(fracPlaces))
			} else {
				min, max = computeUnsignedFixedBounds(int64(bits), float64(fracPlaces))
			}

			if val.Cmp(min) == -1 || val.Cmp(max) == 1 {
				return []byte{}, fmt.Errorf("value out of allowed range: %v, %v", typeStr, val)
			}

			scaledValue := new(big.Float)
			scaledValue.Mul(val, pow(FLOAT_TEN, uint64(fracPlaces)))
			bigIntValue := new(big.Int)
			scaledValue.Int(bigIntValue)

			bytes = append(bytes, common.LeftPadBytes(bigIntValue.Bytes(), bits/8)...)

		} else {
			return []byte{}, fmt.Errorf("invalid parameter type: %v, %T", typeStr, value)
		}
	}

	return bytes, nil
}

func calculateHeadLength(rawHeadChunks [][]byte) uint64 {
	headLength := uint64(0)
	for _, chunk := range rawHeadChunks {
		if chunk != nil {
			headLength += uint64(len(chunk))
		} else {
			headLength += 32
		}
	}

	return headLength
}

func calculateTailOffsets(tailChunks [][]byte) []uint64 {

	tailOffsets := []uint64{0}
	accSum := uint64(0)
	for i := 0; i < len(tailChunks)-1; i++ {
		accSum += uint64(len(tailChunks[i]))
		tailOffsets = append(tailOffsets, accSum)
	}

	return tailOffsets
}

func buildHeadChunks(rawHeadChunks [][]byte, tailOffsets []uint64, headLength uint64) ([][]byte, error) {
	if len(rawHeadChunks) != len(tailOffsets) {
		return nil, fmt.Errorf("tailOffsets and rawHeadChunks must be the same length")
	}

	var headChunks [][]byte
	for i, chunk := range rawHeadChunks {
		if chunk != nil {
			headChunks = append(headChunks, chunk)
		} else {
			headChunks = append(headChunks, encodeUint256(big.NewInt(int64(headLength+tailOffsets[i]))))
		}
	}

	return headChunks, nil
}

func encodeUint256(value *big.Int) []byte {
	return common.LeftPadBytes(value.Bytes(), 32)
}

func joinChunks(headChunks [][]byte, tailChunks [][]byte) []byte {

	joined := make([]byte, 0)

	for _, chunk := range headChunks {
		joined = append(joined, chunk...)
	}

	for _, chunk := range tailChunks {
		joined = append(joined, chunk...)
	}

	return joined
}

func pow(a *big.Float, e uint64) *big.Float {
	result := zero().Copy(a)
	for i := uint64(0); i < e-1; i++ {
		result = mul(result, a)
	}
	return result
}

func zero() *big.Float {
	r := big.NewFloat(0.0)
	r.SetPrec(256)
	return r
}

func mul(a, b *big.Float) *big.Float {
	return zero().Mul(a, b)
}
