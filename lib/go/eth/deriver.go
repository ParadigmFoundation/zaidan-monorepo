package eth

import (
	"github.com/ethereum/go-ethereum/accounts"
)

const (
	DerivationPathPurposeIndex = iota
	DerivationPathCoinTypeIndex
	DerivationPathAccountIndex
	DerivationPathChangeIndex
	DerivationPathAddressIndex
)

const DerivationPathLength = 5

// Deriver is a type that simplifies generating derivation paths
type Deriver struct {
	base accounts.DerivationPath
	last uint32
}

// NewDeriver creates a derivation path deriver given a base path
func NewDeriver(basePath accounts.DerivationPath) *Deriver {
	return &Deriver{
		base: basePath,
		last: basePath[DerivationPathAddressIndex],
	}
}

// NewBaseDeriver creates a deriver from the standard base path
func NewBaseDeriver() *Deriver {
	return NewDeriver(accounts.DefaultBaseDerivationPath)
}

// DeriveNext derives and returns the next path, by incrementing the address index
func (d *Deriver) DeriveNext() accounts.DerivationPath {
	next := d.copyBase()
	next[DerivationPathAddressIndex] = d.next()
	return next
}

// DeriveAt returns a copy of the base derivation path with the address index set to index
func (d *Deriver) DeriveAt(index uint32) accounts.DerivationPath {
	next := d.copyBase()
	next[DerivationPathAddressIndex] = index
	return next
}

// Base returns a copy of the deriver's base derivation path
func (d *Deriver) Base() accounts.DerivationPath { return d.copyBase() }

// next increments last and returns the next value for the address index
func (d *Deriver) next() uint32 {
	d.last++
	return d.last
}

// copy base returns a copy of the base path
func (d *Deriver) copyBase() accounts.DerivationPath {
	dst := make(accounts.DerivationPath, DerivationPathLength)
	copy(dst, d.base)
	return dst
}
