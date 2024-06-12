package abi_test

import (
	"fmt"

	"github.com/omnes-tech/abi"
)

func ExampleGetSigTypes() {
	signature := "funcName((address,uint256,bytes)[],uint256,string,bool)"
	sigTypes, err := abi.GetSigTypes(signature)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(sigTypes)

	// Output: [(address,uint256,bytes)[] uint256 string bool]
}

func ExampleSplitParams() {
	typesStr := "(address,uint256,bytes)[],uint256,string,bool,address"
	splittedTypes := abi.SplitParams(typesStr)

	fmt.Println(splittedTypes)

	// Output: [(address,uint256,bytes)[] uint256 string bool address]
}

func ExampleIsDynamic() {
	typeStr := "(address,uint256,bytes)[]"
	isDynamic := abi.IsDynamic(typeStr, true)

	fmt.Println(isDynamic)

	// Output: true
}

func ExampleIsTuple() {
	typeStr := "(address,uint256,bytes)[]"
	isTuple, types, err := abi.IsTuple(typeStr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(isTuple, types)

	// Output: true [address uint256 bytes]
}

func ExampleIsArray() {
	typeStr := "(address,uint256,bytes)[]"
	isArray, size, err := abi.IsArray(typeStr)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(isArray, size)

	// Output: true 0
}
