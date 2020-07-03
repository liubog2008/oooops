package random

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomStringLengthAndCharset(t *testing.T) {
	strLength := 10
	s := Random(strLength)
	assert.Equal(t, len(s), strLength)
	for _, r := range s {
		c := r
		assert.Condition(t, func() bool {
			return (c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')
		}, "invalid char %c", c)
	}
}

func TestRandomStringIsNotEqual(t *testing.T) {
	strLength := 10
	a := Random(strLength)
	b := Random(strLength)

	assert.NotEqual(t, a, b)
}
