package client

import (
	"bytes"
	"testing"
)

func TestIsContractDeployDatastore(t *testing.T) {
	tests := []struct {
		name       string
		datastore  []ContractDatastoreEntry
		isDeploySC bool
	}{
		{
			name: "Valid contract deploy datastore",
			datastore: []ContractDatastoreEntry{
				{Key: []byte{0}, Value: []byte{1}},
				{Key: getContractByteCodeKey(), Value: []byte{1, 2, 3}},
				{Key: getArgsKey(), Value: []byte{4, 5, 6}},
				{Key: getCoinsKey(), Value: []byte{7, 8, 9}},
			},
			isDeploySC: true,
		},
		{
			name: "Invalid contract deploy datastore - wrong length",
			datastore: []ContractDatastoreEntry{
				{Key: []byte{0}, Value: []byte{1}},
				{Key: getContractByteCodeKey(), Value: []byte{1, 2, 3}},
			},
			isDeploySC: false,
		},
		{
			name: "Invalid contract deploy datastore - wrong smart contract number key",
			datastore: []ContractDatastoreEntry{
				{Key: []byte{1}, Value: []byte{1}},
				{Key: getContractByteCodeKey(), Value: []byte{1, 2, 3}},
				{Key: getArgsKey(), Value: []byte{4, 5, 6}},
				{Key: getCoinsKey(), Value: []byte{7, 8, 9}},
			},
			isDeploySC: false,
		},
		{
			name: "Invalid contract deploy datastore - wrong bytecode key",
			datastore: []ContractDatastoreEntry{
				{Key: []byte{0}, Value: []byte{1}},
				{Key: []byte("wrong key"), Value: []byte{1, 2, 3}},
				{Key: getArgsKey(), Value: []byte{4, 5, 6}},
				{Key: getCoinsKey(), Value: []byte{7, 8, 10}},
			},
			isDeploySC: false,
		},
		{
			name: "Invalid contract deploy datastore - wrong Args key",
			datastore: []ContractDatastoreEntry{
				{Key: []byte{0}, Value: []byte{1}},
				{Key: getContractByteCodeKey(), Value: []byte{1, 2, 3}},
				{Key: []byte("wrong key"), Value: []byte{4, 5, 6}},
				{Key: getCoinsKey(), Value: []byte{7, 8, 10}},
			},
			isDeploySC: false,
		},
		{
			name: "Invalid contract deploy datastore - wrong coins key",
			datastore: []ContractDatastoreEntry{
				{Key: []byte{0}, Value: []byte{1}},
				{Key: getContractByteCodeKey(), Value: []byte{1, 2, 3}},
				{Key: getArgsKey(), Value: []byte{4, 5, 6}},
				{Key: []byte("wrong key"), Value: []byte{7, 8, 10}},
			},
			isDeploySC: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isContractDeployDatastore(tt.datastore)
			if result != tt.isDeploySC {
				t.Errorf("isContractDeployDatastore() = %v, isDeploySC %v", result, tt.isDeploySC)
			}
		})
	}
}

func TestDeSerializeDatastore(t *testing.T) {
	tests := []struct {
		name     string
		input    []byte
		expected []ContractDatastoreEntry
	}{
		{
			name: "Valid datastore",
			input: func() []byte {
				datastore := []ContractDatastoreEntry{
					{Key: []byte{0}, Value: []byte{1}},
					{Key: []byte{1}, Value: []byte{2, 3}},
				}
				serialized, _ := SerializeDatastore(datastore)

				return serialized
			}(),
			expected: []ContractDatastoreEntry{
				{Key: []byte{0}, Value: []byte{1}},
				{Key: []byte{1}, Value: []byte{2, 3}},
			},
		},
		{
			name:     "Empty datastore",
			input:    []byte{},
			expected: nil,
		},
		{
			name: "Large datastore",
			input: func() []byte {
				datastore := []ContractDatastoreEntry{
					{
						Key:   []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
						Value: []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
					},
				}
				serialized, _ := SerializeDatastore(datastore)

				return serialized
			}(),
			expected: []ContractDatastoreEntry{
				{
					Key:   []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
					Value: []byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := DeserializeDatastore(test.input)
			if err != nil {
				t.Errorf("Got unexpected error = %v", err)

				return
			}

			if !equalDatastoreEntries(result, test.expected) {
				t.Errorf("DeSerializeDatastore() = %+v, expected: %+v", result, test.expected)
			}
		})
	}
}

func equalDatastoreEntries(entries1, entries2 []ContractDatastoreEntry) bool {
	if len(entries1) != len(entries2) {
		return false
	}

	for i := range entries1 {
		if !bytes.Equal(entries1[i].Key, entries2[i].Key) || !bytes.Equal(entries1[i].Value, entries2[i].Value) {
			return false
		}
	}

	return true
}
