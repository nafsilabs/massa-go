package utils

import (
	"encoding/binary"
)

type NetworkType int64

const (
	MAINNET   NetworkType = 77658377
	BUILDNET  NetworkType = 77658366
	SECURENET NetworkType = 77658383
	LABNET    NetworkType = 77658376
	SANDBOX   NetworkType = 77
)

// Serialize serializes the NetworkType into an 8-byte big-endian byte slice.
func (nt NetworkType) Serialize() []byte {
	data := make([]byte, 8)
	binary.BigEndian.PutUint64(data, uint64(nt))
	return data
}
