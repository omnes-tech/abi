package abi

import (
	"fmt"
	"strconv"
	"strings"
)

// IsDynamic checks whether given type string is a dynamic,
// i.e. if it is either a string, bytes, or an array.
func IsDynamic(typeStr string, isTuple bool) bool {

	if strings.Contains(typeStr, "[") ||
		strings.Contains(typeStr, "string") ||
		typeStr == "bytes" ||
		(isTuple && (strings.Contains(typeStr, "bytes,") || strings.Contains(typeStr, "bytes)"))) {
		return true
	}

	return false
}

// IsArray checks whether given type string is an array.
// Returns whether or not it is an array and the array length
// for bounded arrays (i.e. `uint256[3]`). Array length can
// be 0 wheter an error occured, it is not an array,
// or it is an unbounded array (i.e. `uint256[]`).
func IsArray(typeStr string) (bool, int, error) {
	if strings.Count(typeStr, "[") != strings.Count(typeStr, "]") {
		return false, 0, fmt.Errorf("invalid array definition")
	}

	if strings.Count(typeStr, "[") > 0 {
		openBracketIndex := strings.LastIndex(typeStr, "[")
		closeBracketIndex := strings.LastIndex(typeStr, "]")
		closeParenthesisIndex := strings.LastIndex(typeStr, ")")

		if openBracketIndex < closeParenthesisIndex || closeBracketIndex < closeParenthesisIndex {
			return false, 0, nil
		}

		var arraySize int
		var err error
		if closeBracketIndex > openBracketIndex+1 {
			arraySize, err = strconv.Atoi(typeStr[openBracketIndex+1 : closeBracketIndex])
			if err != nil {
				return false, 0, fmt.Errorf("invalid array definition")
			}
		}

		return true, arraySize, nil
	}

	return false, 0, nil
}

// IsTuple checks whether given type string is a tuple (i.e. `(uint256,bytes,address)`).
// Also returns the array of type strings in the tuple (i.e. [uint256,bytes,address]).
func IsTuple(typeStr string) (bool, []string, error) {
	if strings.Count(typeStr, "(") != strings.Count(typeStr, ")") {
		return false, nil, fmt.Errorf("invalid tuple definition")
	}

	if strings.Count(typeStr, "(") > 0 {
		openParenthesisIndex := strings.Index(typeStr, "(")
		closeParenthesisIndex := strings.LastIndex(typeStr, ")")
		innerCloseParenthesisIndex := strings.LastIndex(typeStr[:closeParenthesisIndex], ")")
		innerCloseBracketsIndex := strings.LastIndex(typeStr[:closeParenthesisIndex], "]")

		var splitTypes []string
		if innerCloseParenthesisIndex != -1 &&
			innerCloseBracketsIndex != -1 {
			innerOpenParenthesisIndex := strings.Index(typeStr[openParenthesisIndex+1:closeParenthesisIndex], "(") + len(typeStr[:openParenthesisIndex]) + 1
			if innerCloseParenthesisIndex > innerCloseBracketsIndex {
				splitTypes = strings.Split(typeStr[openParenthesisIndex+1:innerOpenParenthesisIndex], ",")
				splitTypes = append(splitTypes, typeStr[innerOpenParenthesisIndex:innerCloseParenthesisIndex+1])
				splitTypes = append(splitTypes, strings.Split(typeStr[innerCloseParenthesisIndex+2:closeParenthesisIndex], ",")...)
			} else {
				splitTypes = strings.Split(typeStr[openParenthesisIndex+1:innerOpenParenthesisIndex], ",")
				splitTypes = append(splitTypes, typeStr[innerOpenParenthesisIndex:innerCloseBracketsIndex+1])
				splitTypes = append(splitTypes, strings.Split(typeStr[innerCloseBracketsIndex+2:closeParenthesisIndex], ",")...)
			}
		} else {
			splitTypes = strings.Split(typeStr[openParenthesisIndex+1:closeParenthesisIndex], ",")
		}

		return true, splitTypes, nil
	}

	return false, nil, nil
}

// GetSigTypes gets the input parameters type from given function
// signature.
// Calls splitParams function.
func GetSigTypes(funcSig string) ([]string, error) {
	openParIndex := strings.Index(funcSig, "(")
	if openParIndex == -1 {
		return []string{}, fmt.Errorf("no opening parenthesis found in function signature")
	}

	closeParIndex := strings.LastIndex(funcSig, ")")
	if closeParIndex == -1 {
		return []string{}, fmt.Errorf("no closing parenthesis found in function signature")
	}

	return SplitParams(funcSig[openParIndex+1 : closeParIndex]), nil
}

// SplitParams splits parameters type from given type string
func SplitParams(typesStr string) []string {

	if len(typesStr) == 0 {
		return nil
	}

	var result []string
	var buffer []string
	insideParentheses := 0
	for _, char := range typesStr {
		if string(char) == "(" {
			insideParentheses += 1
		} else if string(char) == ")" {
			insideParentheses -= 1
		}

		if string(char) == "," && insideParentheses <= 0 {
			result = append(result, strings.Join(buffer, ""))
			buffer = []string{}
		} else {
			buffer = append(buffer, string(char))
		}
	}

	result = append(result, strings.Join(buffer, ""))

	return result
}
