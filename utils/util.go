package utils

import (
	"crypto/sha256"
	"fmt"

	"github.com/ethereum/go-ethereum/common"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/satori/go.uuid"
)

//EmptyHash all zero,invalid
var EmptyHash = common.Hash{}

//EmptyAddress all zero,invalid
var EmptyAddress = common.Address{}

func NewRandomHash() (hash common.Hash) {
	u2, err := uuid.NewV4()
	if err != nil {
		panic(fmt.Sprintf("Something went wrong: %s", err))
	}
	copy(hash[:], crypto.Keccak256(u2.Bytes()))
	return
}

//ShaSecret is short for sha256
func ShaSecret(data []byte) common.Hash {
	//	return crypto.Keccak256Hash(data...)
	return sha256.Sum256(data)
}
