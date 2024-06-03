package abi

import (
	"fmt"
	"strconv"
	"strings"
)

func isDynamic(typeStr string) bool {

	if strings.Contains(typeStr, "[") ||
		strings.Contains(typeStr, "string") ||
		strings.Contains(typeStr, "bytes") {
		return true
	}

	return false
}

func isArray(typeStr string) (bool, int, error) {
	if strings.Count(typeStr, "[") != strings.Count(typeStr, "]") {
		return false, 0, fmt.Errorf("invalid array definition")
	}

	if strings.Count(typeStr, "[") > 0 {
		openBracketIndex := strings.Index(typeStr, "[")
		closeBracketIndex := strings.Index(typeStr, "[")

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

func isTuple(typeStr string) (bool, []string, error) {
	if strings.Count(typeStr, "(") != strings.Count(typeStr, ")") {
		return false, nil, fmt.Errorf("invalid tuple definition")
	}

	if strings.Count(typeStr, "(") > 0 {
		openParenthesisIndex := strings.Index(typeStr, "(")
		closeParenthesisIndex := strings.LastIndex(typeStr, ")")

		splitTypes := strings.Split(typeStr[openParenthesisIndex+1:closeParenthesisIndex], ",")

		return true, splitTypes, nil
	}

	return false, nil, nil
}
