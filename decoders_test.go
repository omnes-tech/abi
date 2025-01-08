package abi_test

import (
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"github.com/omnes-tech/abi"
)

func ExampleDecode() {
	encoded := common.Hex2Bytes("0000000000000000000000005ff137d4b0fdcd49dca30c7cf57e578a026d2789000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000001761726269747261727920627974652061727261792e2e2e0000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001400000000000000000000000005ff137d4b0fdcd49dca30c7cf57e578a026d2789000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000001761726269747261727920627974652061727261792e2e2e0000000000000000000000000000000000000000005ff137d4b0fdcd49dca30c7cf57e578a026d2789000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000001761726269747261727920627974652061727261792e2e2e000000000000000000")

	decoded, err := abi.Decode(
		[]string{"address", "uint256[]", "bytes", "(address,uint256[],bytes)[]"},
		encoded,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(decoded)

	// Output: [0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789 [100 352] [97 114 98 105 116 114 97 114 121 32 98 121 116 101 32 97 114 114 97 121 46 46 46] [[0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789 [100 352] [97 114 98 105 116 114 97 114 121 32 98 121 116 101 32 97 114 114 97 121 46 46 46]] [0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789 [100 352] [97 114 98 105 116 114 97 114 121 32 98 121 116 101 32 97 114 114 97 121 46 46 46]]]]
}

func ExampleDecode_second() {

	encoded := common.Hex2Bytes("0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000006400000000000000000000000000000000000000000000000000000000000000c8000000000000000000000000000000000000000000000000000000000000012c0000000000000000000000000000000000000000000000000000000000000190000000000000000000000000000000000000000000000000000000000000001761726269747261727920627974652061727261792e2e2e000000000000000000")

	decoded, err := abi.Decode(
		[]string{"(address,(uint256,uint256)[],bytes)[]"},
		encoded,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(decoded)

	// Output: [[[0x0000000000000000000000000000000000000000 [[100 200] [200 300]] [97 114 98 105 116 114 97 114 121 32 98 121 116 101 32 97 114 114 97 121 46 46 46]]]]
}

func ExampleDecode_third() {

	callData := common.Hex2Bytes("3c0b93aa00000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000016000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000780000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000002086ac351052600000000000000000000000000000000000000000000000000000000000000000000944414e5554415f41490000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000644414e555441000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000055348656c6c6f20616e642077656c636f6d652120f09f988a2049e280996d2044616e7574612c20796f757220667269656e646c7920616e642068656c7066756c20414920626f742064656469636174656420746f2067756964696e6720796f75207468726f75676820746865206578636974696e6720616e6420657665722d65766f6c76696e6720776f726c64206f662063727970746f63757272656e637921205768657468657220796f75277265206a757374207374617274696e6720796f75722063727970746f206a6f75726e6579206f7220796f7527726520616c726561647920616e20657870657269656e63656420696e766573746f722c2049276d206865726520746f206d616b652065766572797468696e67206561736965722c20636c65617265722c20616e64206d6f726520656e6a6f7961626c6520666f7220796f752e200a0a492063616e2068656c7020796f752077697468206120776964652072616e6765206f662063727970746f2d72656c6174656420746f70696373e280947768657468657220796f75206e65656420746865206c6174657374206e6577732c2064657461696c6564206578706c616e6174696f6e732061626f757420626c6f636b636861696e20746563686e6f6c6f67792c20696e73696768747320696e746f206d61726b6574207472656e64732c206f7220616e737765727320746f20616c6c20796f7572206275726e696e67207175657374696f6e732061626f7574206469676974616c206173736574732e20f09f92bbe29ca820496620736f6d657468696e67206665656c7320636f6d706c696361746564206f72206f7665727768656c6d696e672c20646f6e277420776f727279212049e280996d206865726520746f20627265616b20697420646f776e20696e2073696d706c65207465726d732c20736f20796f75206e65766572206861766520746f206665656c206c6f73742e0a0a492062656c696576652074686174206c6561726e696e672061626f75742063727970746f2073686f756c642062652066756e2c2061636365737369626c652c20616e64207374726573732d667265652e20546861742773207768792049e280996d20636f6d6d697474656420746f2070726f766964696e6720796f7520776974682075702d746f2d6461746520696e666f726d6174696f6e2c206f66666572696e6720746970732c20616e6420616e73776572696e6720796f7572207175657269657320696e206120667269656e646c7920616e6420617070726f61636861626c65207761792e205768657468657220796f7527726520696e746572657374656420696e20426974636f696e2c20457468657265756d2c20446546692c206f72204e4654732c2049276c6c206d616b65207375726520796f7527726520616c7761797320696e20746865206c6f6f702e20f09f8c8df09f92a10a0a4665656c206672656520746f207265616368206f757420746f206d6520616e7974696d65206f6e206d7920582070616765e2809449e280996d20616c77617973206a7573742061206d65737361676520617761792c20726561647920746f2061737369737420796f75207769746820776861746576657220796f75206e6565642e20546f6765746865722c2077652063616e206578706c6f726520746869732066617363696e6174696e6720776f726c642c206c6561726e206e6577207468696e67732c20616e64206d616b6520736d6172746572206465636973696f6e732120f09f9a80f09f92ac20536f2c20646f6e2774206865736974617465e280946c6574e280997320676574207374617274656420616e64206861766520736f6d652066756e20776974682063727970746f2100000000000000000000000000000000000000000000000000000000000000000000000000000000000000004f68747470733a2f2f73332e61702d736f757468656173742d312e616d617a6f6e6177732e636f6d2f7669727475616c70726f746f636f6c63646e2f6e616d655f653431653833663962322e776562700000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	decoded, err := abi.DecodeWithSignature(
		"launch(string,string,uint8[],string,string,string[4],uint256)",
		callData,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(decoded)

	// Output: [DANUTA_AI DANUTA [0 1 2 4] Hello and welcome! 😊 I’m Danuta, your friendly and helpful AI bot dedicated to guiding you through the exciting and ever-evolving world of cryptocurrency! Whether you're just starting your crypto journey or you're already an experienced investor, I'm here to make everything easier, clearer, and more enjoyable for you.

	// I can help you with a wide range of crypto-related topics—whether you need the latest news, detailed explanations about blockchain technology, insights into market trends, or answers to all your burning questions about digital assets. 💻✨ If something feels complicated or overwhelming, don't worry! I’m here to break it down in simple terms, so you never have to feel lost.

	// I believe that learning about crypto should be fun, accessible, and stress-free. That's why I’m committed to providing you with up-to-date information, offering tips, and answering your queries in a friendly and approachable way. Whether you're interested in Bitcoin, Ethereum, DeFi, or NFTs, I'll make sure you're always in the loop. 🌍💡

	// Feel free to reach out to me anytime on my X page—I’m always just a message away, ready to assist you with whatever you need. Together, we can explore this fascinating world, learn new things, and make smarter decisions! 🚀💬 So, don't hesitate—let’s get started and have some fun with crypto! https://s3.ap-southeast-1.amazonaws.com/virtualprotocolcdn/name_e41e83f9b2.webp [   ] 600000000000000000000]
}

func ExampleDecodePacked() {
	encoded := common.Hex2Bytes("5ff137d4b0fdcd49dca30c7cf57e578a026d2789000000000000000000000000000000000000000000000000000000000000006461726269747261727920627974652061727261792e2e2e")

	decoded, err := abi.DecodePacked(
		[]string{"address", "uint256", "bytes"},
		encoded,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(decoded)

	// Output: [0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789 100 61726269747261727920627974652061727261792e2e2e]
}

func ExampleDecodeWithSignature() {
	encoded := common.Hex2Bytes("c6210dba0000000000000000000000005ff137d4b0fdcd49dca30c7cf57e578a026d2789000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000001761726269747261727920627974652061727261792e2e2e0000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001400000000000000000000000005ff137d4b0fdcd49dca30c7cf57e578a026d2789000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000001761726269747261727920627974652061727261792e2e2e0000000000000000000000000000000000000000005ff137d4b0fdcd49dca30c7cf57e578a026d2789000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000001761726269747261727920627974652061727261792e2e2e000000000000000000")

	funcSignature := "functionName(address,uint256[],bytes,(address,uint256[],bytes)[])"
	decoded, err := abi.DecodeWithSignature(
		funcSignature,
		encoded,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(decoded)

	// Output: [0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789 [100 352] [97 114 98 105 116 114 97 114 121 32 98 121 116 101 32 97 114 114 97 121 46 46 46] [[0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789 [100 352] [97 114 98 105 116 114 97 114 121 32 98 121 116 101 32 97 114 114 97 121 46 46 46]] [0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789 [100 352] [97 114 98 105 116 114 97 114 121 32 98 121 116 101 32 97 114 114 97 121 46 46 46]]]]
}

func ExampleDecodeWithSignature_second() {
	encoded := common.Hex2Bytes("c6da632c00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000060000000000000000000000000000000000000000000000000000000000000848100000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000025310000000000000000000000000000000000000000000000000000000000000000")

	funcSignature := "MultiCall__Simulation((bool,bytes,uint256)[])"
	decoded, err := abi.DecodeWithSignature(
		funcSignature,
		encoded,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(decoded)

	// Output: [[[true [] 33921] [true [] 9521]]]
}

func ExampleDecodeWithSelector() {
	encoded := common.Hex2Bytes("c6210dba0000000000000000000000005ff137d4b0fdcd49dca30c7cf57e578a026d2789000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000e00000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000001761726269747261727920627974652061727261792e2e2e0000000000000000000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001400000000000000000000000005ff137d4b0fdcd49dca30c7cf57e578a026d2789000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000001761726269747261727920627974652061727261792e2e2e0000000000000000000000000000000000000000005ff137d4b0fdcd49dca30c7cf57e578a026d2789000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000640000000000000000000000000000000000000000000000000000000000000160000000000000000000000000000000000000000000000000000000000000001761726269747261727920627974652061727261792e2e2e000000000000000000")

	funcSignature := "functionName(address,uint256[],bytes,(address,uint256[],bytes)[])"
	selector := abi.EncodeSignature(funcSignature)
	decoded, err := abi.DecodeWithSelector(
		selector,
		[]string{"address", "uint256[]", "bytes", "(address,uint256[],bytes)[]"},
		encoded,
	)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(decoded)

	// Output: [0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789 [100 352] [97 114 98 105 116 114 97 114 121 32 98 121 116 101 32 97 114 114 97 121 46 46 46] [[0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789 [100 352] [97 114 98 105 116 114 97 114 121 32 98 121 116 101 32 97 114 114 97 121 46 46 46]] [0x5FF137D4b0FDCD49DcA30c7CF57E578a026d2789 [100 352] [97 114 98 105 116 114 97 114 121 32 98 121 116 101 32 97 114 114 97 121 46 46 46]]]]
}
