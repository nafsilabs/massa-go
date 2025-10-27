package wallet

import (
	"github.com/zeebo/blake3"
)

// blake3Hash computes the BLAKE3 hash of the input data
// This is used for address generation in the Massa blockchain
func blake3Hash(data []byte) []byte {
	hasher := blake3.New()
	hasher.Write(data)
	return hasher.Sum(nil)
}
