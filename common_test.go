package abi_test

import (
	"fmt"

	"github.com/omnes-tech/abi"
)

func ExampleGetSigTypes() {
	signature := "funcName((address, uint256, bytes)[], uint256, string, bool)"
	sigTypes, err := abi.GetSigTypes(signature)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(sigTypes)

	// Output: [(address, uint256, bytes)[]  uint256  string  bool]
}

func ExampleSplitParams() {
	typesStr := "(address, uint256, bytes)[], uint256, string, bool"
	splittedTypes := abi.SplitParams(typesStr)

	fmt.Println(splittedTypes)

	// Output: [(address, uint256, bytes)[]  uint256  string  bool]
}
