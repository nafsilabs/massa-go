package sc

import (
	"encoding/json"
)

// Storage provides functions for interacting with the key-value datastore
// used for persistent storage of data on the blockchain.

// Set stores a value in the datastore under the given key
func Set(key []byte, value []byte) {
	if isWasmRuntime() {
		keyPtr := newWasmBytes(key)
		valuePtr := newWasmBytes(value)
		wasmSetData(keyPtr, valuePtr)
	}
	// In non-WASM mode, this would store in a mock storage for testing
}

// SetString stores a string value in the datastore under the given string key
func SetString(key string, value string) {
	Set([]byte(key), []byte(value))
}

// Get retrieves a value from the datastore for the given key
func Get(key []byte) []byte {
	if isWasmRuntime() {
		keyPtr := newWasmBytes(key)
		resultPtr := wasmGetData(keyPtr)
		return wasmBytesFromPtr(resultPtr)
	}
	// In non-WASM mode, return empty for testing
	return nil
}

// GetString retrieves a string value from the datastore for the given string key
func GetString(key string) string {
	data := Get([]byte(key))
	return string(data)
}

// Has checks if a key exists in the datastore
func Has(key []byte) bool {
	if isWasmRuntime() {
		keyPtr := newWasmBytes(key)
		result := wasmHasData(keyPtr)
		return result != 0
	}
	// In non-WASM mode, return false for testing
	return false
}

// HasString checks if a string key exists in the datastore
func HasString(key string) bool {
	return Has([]byte(key))
}

// Delete removes a key-value pair from the datastore
func Delete(key []byte) {
	if isWasmRuntime() {
		keyPtr := newWasmBytes(key)
		wasmDeleteData(keyPtr)
	}
	// In non-WASM mode, this would delete from mock storage for testing
}

// DeleteString removes a string key from the datastore
func DeleteString(key string) {
	Delete([]byte(key))
}

// Append appends data to an existing value in the datastore
func Append(key []byte, value []byte) {
	if isWasmRuntime() {
		keyPtr := newWasmBytes(key)
		valuePtr := newWasmBytes(value)
		wasmAppendData(keyPtr, valuePtr)
	}
	// In non-WASM mode, this would append to mock storage for testing
}

// AppendString appends a string to an existing value in the datastore
func AppendString(key string, value string) {
	Append([]byte(key), []byte(value))
}

// GetKeys retrieves all keys with the given prefix
func GetKeys(prefix []byte) [][]byte {
	if isWasmRuntime() {
		prefixPtr := newWasmBytes(prefix)
		resultPtr := wasmGetKeys(prefixPtr)
		data := wasmBytesFromPtr(resultPtr)
		// The result is serialized as JSON array of base64 strings
		var keys []string
		if err := json.Unmarshal(data, &keys); err != nil {
			return nil
		}
		result := make([][]byte, len(keys))
		for i, key := range keys {
			result[i] = []byte(key)
		}
		return result
	}
	// In non-WASM mode, return empty for testing
	return nil
}

// GetKeysString retrieves all string keys with the given string prefix
func GetKeysString(prefix string) []string {
	keys := GetKeys([]byte(prefix))
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = string(key)
	}
	return result
}

// Storage operations for other addresses
// These functions can only be called at smart contract generation time
// by the parent smart contract to write to or delete data from the child's datastore

// SetOf stores a value in another address's datastore
func SetOf(address *Address, key []byte, value []byte) {
	if isWasmRuntime() {
		addressPtr := newWasmString(address.String())
		keyPtr := newWasmBytes(key)
		valuePtr := newWasmBytes(value)
		wasmSetDataFor(addressPtr, keyPtr, valuePtr)
	}
	// In non-WASM mode, this would store in mock storage for testing
}

// SetOfString stores a string value in another address's datastore
func SetOfString(address *Address, key string, value string) {
	SetOf(address, []byte(key), []byte(value))
}

// GetOf retrieves a value from another address's datastore
func GetOf(address *Address, key []byte) []byte {
	if isWasmRuntime() {
		addressPtr := newWasmString(address.String())
		keyPtr := newWasmBytes(key)
		resultPtr := wasmGetDataFor(addressPtr, keyPtr)
		return wasmBytesFromPtr(resultPtr)
	}
	// In non-WASM mode, return empty for testing
	return nil
}

// GetOfString retrieves a string value from another address's datastore
func GetOfString(address *Address, key string) string {
	data := GetOf(address, []byte(key))
	return string(data)
}

// HasOf checks if a key exists in another address's datastore
func HasOf(address *Address, key []byte) bool {
	if isWasmRuntime() {
		addressPtr := newWasmString(address.String())
		keyPtr := newWasmBytes(key)
		result := wasmHasDataFor(addressPtr, keyPtr)
		return result != 0
	}
	// In non-WASM mode, return false for testing
	return false
}

// HasOfString checks if a string key exists in another address's datastore
func HasOfString(address *Address, key string) bool {
	return HasOf(address, []byte(key))
}

// DeleteOf removes a key-value pair from another address's datastore
func DeleteOf(address *Address, key []byte) {
	if isWasmRuntime() {
		addressPtr := newWasmString(address.String())
		keyPtr := newWasmBytes(key)
		wasmDeleteDataFor(addressPtr, keyPtr)
	}
	// In non-WASM mode, this would delete from mock storage for testing
}

// DeleteOfString removes a string key from another address's datastore
func DeleteOfString(address *Address, key string) {
	DeleteOf(address, []byte(key))
}

// AppendOf appends data to an existing value in another address's datastore
func AppendOf(address *Address, key []byte, value []byte) {
	if isWasmRuntime() {
		addressPtr := newWasmString(address.String())
		keyPtr := newWasmBytes(key)
		valuePtr := newWasmBytes(value)
		wasmAppendDataFor(addressPtr, keyPtr, valuePtr)
	}
	// In non-WASM mode, this would append to mock storage for testing
}

// AppendOfString appends a string to an existing value in another address's datastore
func AppendOfString(address *Address, key string, value string) {
	AppendOf(address, []byte(key), []byte(value))
}

// GetKeysOf retrieves all keys with the given prefix from another address's datastore
func GetKeysOf(address *Address, prefix []byte) [][]byte {
	if isWasmRuntime() {
		addressPtr := newWasmString(address.String())
		prefixPtr := newWasmBytes(prefix)
		resultPtr := wasmGetKeysOf(addressPtr, prefixPtr)
		data := wasmBytesFromPtr(resultPtr)
		// The result is serialized as JSON array of base64 strings
		var keys []string
		if err := json.Unmarshal(data, &keys); err != nil {
			return nil
		}
		result := make([][]byte, len(keys))
		for i, key := range keys {
			result[i] = []byte(key)
		}
		return result
	}
	// In non-WASM mode, return empty for testing
	return nil
}

// GetKeysOfString retrieves all string keys with the given string prefix from another address's datastore
func GetKeysOfString(address *Address, prefix string) []string {
	keys := GetKeysOf(address, []byte(prefix))
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = string(key)
	}
	return result
}

// Storage helper functions for typed data

// SetJSON stores a JSON-serializable value in the datastore
func SetJSON(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	SetString(key, string(data))
	return nil
}

// GetJSON retrieves and deserializes a JSON value from the datastore
func GetJSON(key string, value interface{}) error {
	data := GetString(key)
	if data == "" {
		return nil // Return nil for missing keys
	}
	return json.Unmarshal([]byte(data), value)
}

// SetJSONOf stores a JSON-serializable value in another address's datastore
func SetJSONOf(address *Address, key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	SetOfString(address, key, string(data))
	return nil
}

// GetJSONOf retrieves and deserializes a JSON value from another address's datastore
func GetJSONOf(address *Address, key string, value interface{}) error {
	data := GetOfString(address, key)
	if data == "" {
		return nil // Return nil for missing keys
	}
	return json.Unmarshal([]byte(data), value)
}
