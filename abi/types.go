package abi

import (
	"fmt"
	"math"
	"math/big"
)

// ParamType specifies the byte length and
// minimum and maximum values for a type.
type ParamType struct {
	ByteLength int // 0 means there is no length restriction
	Min        *big.Int
	Max        *big.Int
}

// minusTwo big.Int for -2
var minusTwo = big.NewInt(-2)

// zero big.Int for 0
var zero = big.NewInt(0)

// floatZero big.Float for 0
var floatZero = big.NewFloat(0)

// one big.Int for 1
var one = big.NewInt(1)

// two big.Int for 2
var two = big.NewInt(2)

// floatTen big.Float for 10
var floatTen = big.NewFloat(10)

// validCoreTypes maps type to its byte length and
// minimum and maximum value restrictions.
var validCoreTypes = map[string]ParamType{
	"uint8":   {1, zero, big.NewInt(255)},
	"uint16":  {2, zero, big.NewInt(65535)},
	"uint24":  {3, zero, big.NewInt(16777215)},
	"uint32":  {4, zero, big.NewInt(4294967295)},
	"uint40":  {5, zero, big.NewInt(1099511628000)},
	"uint48":  {6, zero, big.NewInt(281474976710655)},
	"uint56":  {7, zero, big.NewInt(72057594037927935)},
	"uint64":  {8, zero, convertStringToBigInt("18446744073709551615")},
	"uint72":  {9, zero, convertStringToBigInt("4722366482869645213695")},
	"uint80":  {10, zero, convertStringToBigInt("1208925819614629174706175")},
	"uint88":  {11, zero, convertStringToBigInt("309485009821345068724781055")},
	"uint96":  {12, zero, convertStringToBigInt("79228162514264337593543950335")},
	"uint104": {13, zero, convertStringToBigInt("20282409603651670423947251286015")},
	"uint112": {14, zero, convertStringToBigInt("5192296858534827628530496329220095")},
	"uint120": {15, zero, convertStringToBigInt("1329227995784915872903807060280344575")},
	"uint128": {16, zero, convertStringToBigInt("340282366920938463463374607431768211455")},
	"uint136": {17, zero, convertStringToBigInt("87112285931760246646623899502532662132735")},
	"uint144": {18, zero, convertStringToBigInt("22300745198530623141535718272648361505980415")},
	"uint152": {19, zero, convertStringToBigInt("5708990770823839524233143877797980545530986495")},
	"uint160": {20, zero, convertStringToBigInt("1461501637330902918203684832716283019655932542975")},
	"uint168": {21, zero, convertStringToBigInt("374144419156711147060143317175368453031918731001855")},
	"uint176": {22, zero, convertStringToBigInt("95780971304118053647396689196894323976171195136475135")},
	"uint184": {23, zero, convertStringToBigInt("24519928653854221733733552434404946937899825954937634815")},
	"uint192": {24, zero, convertStringToBigInt("6277101735386680763835789423207666416102355444464034512895")},
	"uint200": {25, zero, convertStringToBigInt("1606938044258990275541962092341162602522202993782792835301375")},
	"uint208": {26, zero, convertStringToBigInt("411376139330301510538742295639337626245683966408394965837152255")},
	"uint216": {27, zero, convertStringToBigInt("105312291668557186697918027683670432318895095400549111254310977535")},
	"uint224": {28, zero, convertStringToBigInt("26959946667150639794667015087019630673637144422540572481103610249215")},
	"uint232": {29, zero, convertStringToBigInt("6901746346790563787434755862277025452451108972170386555162524223799295")},
	"uint240": {30, zero, convertStringToBigInt("1766847064778384329583297500742918515827483896875618958121606201292619775")},
	"uint248": {31, zero, convertStringToBigInt("452312848583266388373324160190187140051835877600158453279131187530910662655")},
	"uint256": {32, zero, convertStringToBigInt("115792089237316195423570985008687907853269984665640564039457584007913129639935")},
	"int8":    {1, big.NewInt(-128), big.NewInt(127)},
	"int16":   {2, big.NewInt(-32768), big.NewInt(32767)},
	"int24":   {3, big.NewInt(-8388608), big.NewInt(8388607)},
	"int32":   {4, big.NewInt(-2147483648), big.NewInt(2147483647)},
	"int40":   {5, big.NewInt(-549755813888), big.NewInt(549755813887)},
	"int48":   {6, big.NewInt(-140737488355328), big.NewInt(140737488355327)},
	"int56":   {7, big.NewInt(-36028797018963968), big.NewInt(36028797018963967)},
	"int64":   {8, big.NewInt(-9223372036854775808), big.NewInt(9223372036854775807)},
	"int72":   {9, convertStringToBigInt("-2361183241434822606848"), convertStringToBigInt("2361183241434822606847")},
	"int80":   {10, convertStringToBigInt("-604462909807314587353088"), convertStringToBigInt("604462909807314587353087")},
	"int88":   {11, convertStringToBigInt("-154742504910672534362390528"), convertStringToBigInt("154742504910672534362390527")},
	"int96":   {12, convertStringToBigInt("-39614081257132168796771975168"), convertStringToBigInt("39614081257132168796771975167")},
	"int104":  {13, convertStringToBigInt("-10141204801825835211973625643008"), convertStringToBigInt("10141204801825835211973625643007")},
	"int112":  {14, convertStringToBigInt("-2596148429267413814265248164610048"), convertStringToBigInt("2596148429267413814265248164610047")},
	"int120":  {15, convertStringToBigInt("-664613997892457936451903530140172288"), convertStringToBigInt("664613997892457936451903530140172287")},
	"int128":  {16, convertStringToBigInt("-170141183460469231731687303715884105728"), convertStringToBigInt("170141183460469231731687303715884105727")},
	"int136":  {17, convertStringToBigInt("-43556142965880123323311949751266331066368"), convertStringToBigInt("43556142965880123323311949751266331066367")},
	"int144":  {18, convertStringToBigInt("-11150372599265311570767859136324180752990208"), convertStringToBigInt("11150372599265311570767859136324180752990207")},
	"int152":  {19, convertStringToBigInt("-2854495385411919762116571938898990272765493248"), convertStringToBigInt("2854495385411919762116571938898990272765493247")},
	"int160":  {20, convertStringToBigInt("-730750818665451459101842416358141509827966271488"), convertStringToBigInt("730750818665451459101842416358141509827966271487")},
	"int168":  {21, convertStringToBigInt("-187072209578355573530071658587684226515959365500928"), convertStringToBigInt("187072209578355573530071658587684226515959365500927")},
	"int176":  {22, convertStringToBigInt("-47890485652059026823698344598447161988085597568237568"), convertStringToBigInt("47890485652059026823698344598447161988085597568237567")},
	"int184":  {23, convertStringToBigInt("-12259964326927110866866776217202473468949912977468817408"), convertStringToBigInt("12259964326927110866866776217202473468949912977468817407")},
	"int192":  {24, convertStringToBigInt("-3138550867693340381917894711603833208051177722232017256448"), convertStringToBigInt("3138550867693340381917894711603833208051177722232017256447")},
	"int200":  {25, convertStringToBigInt("-803469022129495137770981046170581301261101496891396417650688"), convertStringToBigInt("803469022129495137770981046170581301261101496891396417650687")},
	"int208":  {26, convertStringToBigInt("-205688069665150755269371147819668813122841983204197482918576128"), convertStringToBigInt("205688069665150755269371147819668813122841983204197482918576127")},
	"int216":  {27, convertStringToBigInt("-52656145834278593348959013841835216159447547700274555627155488768"), convertStringToBigInt("52656145834278593348959013841835216159447547700274555627155488767")},
	"int224":  {28, convertStringToBigInt("-13479973333575319897333507543509815336818572211270286240551805124608"), convertStringToBigInt("13479973333575319897333507543509815336818572211270286240551805124607")},
	"int232":  {29, convertStringToBigInt("-3450873173395281893717377931138512726225554486085193277581262111899648"), convertStringToBigInt("3450873173395281893717377931138512726225554486085193277581262111899647")},
	"int240":  {30, convertStringToBigInt("-883423532389192164791648750371459257913741948437809479060803100646309888"), convertStringToBigInt("883423532389192164791648750371459257913741948437809479060803100646309887")},
	"int248":  {31, convertStringToBigInt("-226156424291633194186662080095093570025917938800079226639565593765455331328"), convertStringToBigInt("226156424291633194186662080095093570025917938800079226639565593765455331327")},
	"int256":  {32, convertStringToBigInt("-57896044618658097711785492504343953926634992332820282019728792003956564819968"), convertStringToBigInt("57896044618658097711785492504343953926634992332820282019728792003956564819967")},
	"bytes1":  {1, zero, zero},
	"bytes2":  {2, zero, zero},
	"bytes3":  {3, zero, zero},
	"bytes4":  {4, zero, zero},
	"bytes5":  {5, zero, zero},
	"bytes6":  {6, zero, zero},
	"bytes7":  {7, zero, zero},
	"bytes8":  {8, zero, zero},
	"bytes9":  {9, zero, zero},
	"bytes10": {10, zero, zero},
	"bytes11": {11, zero, zero},
	"bytes12": {12, zero, zero},
	"bytes13": {13, zero, zero},
	"bytes14": {14, zero, zero},
	"bytes15": {15, zero, zero},
	"bytes16": {16, zero, zero},
	"bytes17": {17, zero, zero},
	"bytes18": {18, zero, zero},
	"bytes19": {19, zero, zero},
	"bytes20": {20, zero, zero},
	"bytes21": {21, zero, zero},
	"bytes22": {22, zero, zero},
	"bytes23": {23, zero, zero},
	"bytes24": {24, zero, zero},
	"bytes25": {25, zero, zero},
	"bytes26": {26, zero, zero},
	"bytes27": {27, zero, zero},
	"bytes28": {28, zero, zero},
	"bytes29": {29, zero, zero},
	"bytes30": {30, zero, zero},
	"bytes31": {31, zero, zero},
	"bytes32": {32, zero, zero},
	"bytes":   {0, zero, zero},
	"string":  {0, zero, zero},
	"bool":    {1, zero, one},
	"address": {20, zero, zero},
}

// convertStringToBigInt converts string to big.Int value.
func convertStringToBigInt(s string) *big.Int {
	bigInt, ok := new(big.Int).SetString(s, 10)
	if !ok {
		fmt.Printf("error converting string to big.Int: %v", s)
		return &big.Int{}
	}

	return bigInt
}

// computeSignedIntegerBounds computes the bound range
// for signed integer based on its bits number.
func computeSignedIntegerBounds(numBits int64) (*big.Int, *big.Int) {
	resultLowest := new(big.Int)
	resultLowest.Exp(minusTwo, big.NewInt(numBits-1), nil)

	resultHighest := new(big.Int)
	resultHighest.Exp(two, big.NewInt(numBits-1), nil)

	return resultLowest, resultHighest.Sub(resultHighest, one)
}

// computeUnsignedIntegerBounds computes the bound range
// for unsigned integer based on its bits number.
func computeUnsignedIntegerBounds(numBits int64) (*big.Int, *big.Int) {
	result := new(big.Int)
	result.Exp(two, big.NewInt(numBits-1), nil)

	return zero, result.Sub(result, one)
}

// computeSignedFixedBounds computes the bound range
// for signed fixed point type based on its bits number
// and number of decimals.
func computeSignedFixedBounds(numbBits int64, fracPlaces float64) (*big.Float, *big.Float) {
	lowest, highest := computeSignedIntegerBounds(numbBits)
	floatHighest, _ := highest.Float64()
	floatLowest, _ := lowest.Float64()

	exp := math.Pow(10, -fracPlaces)

	return big.NewFloat(floatLowest * exp), big.NewFloat(floatHighest * exp)
}

// computeUnsignedFixedBounds computes the bound range
// for unsigned fixed point type based on its bits number
// and number of decimals.
func computeUnsignedFixedBounds(numbBits int64, fracPlaces float64) (*big.Float, *big.Float) {
	_, highest := computeUnsignedIntegerBounds(numbBits)
	floatHighest, _ := highest.Float64()

	return floatZero, big.NewFloat(floatHighest * math.Pow(10, -fracPlaces))
}
