// Package random is copied from
// https://stackoverflow.com/questions/22892120/how-to-generate-a-random-string-of-a-fixed-length-in-go
package random

import (
	"math/rand"
	"strings"
)

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
const (
	// 6 bits to represent a letter index
	letterIdxBits = 6
	// All 1-bits, as many as letterIdxBits
	letterIdxMask = 1<<letterIdxBits - 1
	// # of letter indices fitting in 63 bits
	letterIdxMax = 63 / letterIdxBits
)

// Random returns a random string with fixed length
func Random(n int) string {
	sb := strings.Builder{}
	sb.Grow(n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, rand.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = rand.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			if err := sb.WriteByte(letterBytes[idx]); err != nil {
				panic(err)
			}
			i--
		}
		cache >>= letterIdxBits
		remain--
	}

	return sb.String()
}
