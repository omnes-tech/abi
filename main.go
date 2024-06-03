package main

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/omnes-tech/abi/abi"
)

func main() {
	sender := common.HexToAddress("0xb79B3Ea1C9e85FdfFcBEa16868edA7a7E7623125")
	nonce := big.NewInt(10)
	initCode := []byte{}
	callData := common.Hex2Bytes("51945447000000000000000000000000e36bd65609c08cd17b53520293523cf4560533d00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000044a9059cbb00000000000000000000000049422c33d1fc8565410e93aca8574429a4a7eceb0000000000000000000000000000000000000000000000008ac7230489e8000000000000000000000000000000000000000000000000000000000000")
	callGasLimit := big.NewInt(105399)
	verificationGasLimit := big.NewInt(116961)
	preVerificationGas := big.NewInt(66096)
	maxFeePerGas := big.NewInt(35650000000)
	maxPriorityFeePerGas := big.NewInt(35650000000)
	paymasterAndData := common.Hex2Bytes("e3dc822d77f8ca7ac74c30b0dffea9fcdcaaa32100000000000000000000000000000000000000000000000000000000665a463a0000000000000000000000000000000000000000000000000000000000000000e349eeccd9325b24f3b048fb7725a4099a9281e20f9ef96182e9fda34f6149b01c2886e7d977ec84ba012d3dc105233976156b49f941deeeaa487d6a20c153e91b")
	signature := common.Hex2Bytes("00000000c1b18522482071eda4e09198f2b7db72a81b30496925dc9ac249c1144865de3c4891c60f00f39a73c3219b2407e5c586305040940fd743a6627bfcc10223dbfc1b")
	beneficiary := common.HexToAddress("0x4337bB3F3c9645a2A76b6652C7A9F6dF25BC39A8")

	types := []string{"(address,uint256,bytes,bytes,uint256,uint256,uint256,uint256,uint256,bytes,bytes)[]", "address"}

	encoded, err := abi.Encode(
		types,
		[]any{
			[]any{
				sender,
				nonce,
				initCode,
				callData,
				callGasLimit,
				verificationGasLimit,
				preVerificationGas,
				maxFeePerGas,
				maxPriorityFeePerGas,
				paymasterAndData,
				signature,
			},
			[]any{
				sender,
				nonce,
				initCode,
				callData,
				callGasLimit,
				verificationGasLimit,
				preVerificationGas,
				maxFeePerGas,
				maxPriorityFeePerGas,
				paymasterAndData,
				signature,
			},
		},
		beneficiary,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(common.Bytes2Hex(encoded))

	dencoded, err := abi.Decode(
		types,
		encoded,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(dencoded)

}
