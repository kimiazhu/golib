// Description: utils/rand.go random toolkits
// Author: ZHU HAIHUA
// Since: 2016-02-26 19:08
package util

import (
	"crypto/rand"
)

var charTable = []rune("abcdefghijkmnpqrstuvwxyz23456789")

// RandStrN return a random lower case string which length is n.
// this string will NOT contain characters [0/1/o/l]
func RandStrN(n int) string {
	random := make([]byte, n)
	result := make([]rune, n)
	rand.Read(random[:])
	for i := 0; i < len(random); i++ {
		result[i] = charTable[uint(random[i]>>3)]
	}
	return string(result)
}
