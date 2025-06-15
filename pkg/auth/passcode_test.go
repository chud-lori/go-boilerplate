package auth_test

import (
	"testing"

	"github.com/chud-lori/go-boilerplate/pkg/auth"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePasscode(t *testing.T) {
	passcode := auth.GeneratePasscode()
	assert.Len(t, passcode, 8)

	alpha := passcode[:4]
	num := passcode[4:]

	for _, ch := range alpha {
		assert.True(t, ch >= 'A' && ch <= 'Z', "expected capital letter, got: %c", ch)
	}
	for _, ch := range num {
		assert.True(t, ch >= '0' && ch <= '9', "expected digit, got: %c", ch)
	}
}
