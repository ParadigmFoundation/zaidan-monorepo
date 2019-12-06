package zeroex

import (
	"github.com/ethereum/go-ethereum/common"
)

// ECSignature contains the parameters of an elliptic curve signature
type ECSignature struct {
	V byte
	R common.Hash
	S common.Hash
}

type Signer interface {
	EthSign(message []byte, signerAddress common.Address) (*ECSignature, error)
}
