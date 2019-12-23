package eth

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/stretchr/testify/assert"
)

func TestDeriver(t *testing.T) {
	deriver := NewBaseDeriver()
	assert.Equal(t, accounts.DefaultBaseDerivationPath, deriver.base)
	assert.Equal(t, uint32(1), deriver.next())
	assert.Equal(t, uint32(2), deriver.next())

	expected := accounts.DefaultBaseDerivationPath
	expected[DerivationPathAddressIndex] = 3
	assert.Equal(t, expected, deriver.DeriveNext())

	expected[DerivationPathAddressIndex] = 4
	assert.Equal(t, expected, deriver.DeriveNext())

	expected[DerivationPathAddressIndex] = 5
	assert.Equal(t, expected, deriver.DeriveNext())

	assert.Equal(t, accounts.DefaultBaseDerivationPath, deriver.Base())
}
