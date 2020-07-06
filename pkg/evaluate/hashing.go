package evaluate

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/big"
)

//GetMD5Hash hashes a string with MD5 algo
func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

//CohortFraction returns the fraction calculated by the hash function
func CohortFraction(salt string) float64 {
	bi := big.NewInt(0)
	hexStr := GetMD5Hash(salt)
	bi.SetString(hexStr[:6], 16)
	largest6DigHex := big.NewFloat(0xFFFFFF)
	testSVal := new(big.Float).SetInt(bi)
	fraction := new(big.Float).Quo(testSVal, largest6DigHex)
	comparison, _ := fraction.Float64()
	return comparison
}

func abTest(s, salt string, control float64) bool {
	bi := big.NewInt(0)

	hexStr := GetMD5Hash(fmt.Sprintf("%s-%s", s, salt))
	bi.SetString(hexStr[:6], 16)
	largest6DigHex := big.NewFloat(0xFFFFFF)
	testSVal := new(big.Float).SetInt(bi)
	fraction := new(big.Float).Quo(testSVal, largest6DigHex)
	comparison, _ := fraction.Float64()
	if comparison > control {
		return true
	}
	return false
}

func run() {
	fmt.Println("marko", abTest("marko", "test", 0.7))
	fmt.Println("hana", abTest("hana", "test", 0.7))
	fmt.Println("jimmy", abTest("jimmy", "test", 0.7))
	fmt.Println("browna", abTest("browna", "test", 0.7))
	fmt.Println("lyn", abTest("lyn", "test", 0.7))
	fmt.Println("th00", abTest("th00", "test", 0.7))
	fmt.Println("sokcina", abTest("sokcina", "test", 0.7))
	fmt.Println("gutman", abTest("gutman", "test", 0.7))

}
