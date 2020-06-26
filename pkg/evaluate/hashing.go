package evaluate

import (
	"crypto/md5"
	"encoding/hex"
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
