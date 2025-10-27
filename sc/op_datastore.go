package sc

import (
	"encoding/json"
)

// Operation Datastore provides functions for accessing operation-level data storage
// This is different from regular storage as it's associated with the operation itself

// GetOpKeys retrieves all operation keys
func GetOpKeys() [][]byte {
	if isWasmRuntime() {
		resultPtr := wasmGetOpKeys()
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

// GetOpKeysString retrieves all operation keys as strings
func GetOpKeysString() []string {
	keys := GetOpKeys()
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = string(key)
	}
	return result
}

// GetOpKeysWithPrefix retrieves operation keys with the given prefix
func GetOpKeysWithPrefix(prefix []byte) [][]byte {
	if isWasmRuntime() {
		prefixPtr := newWasmBytes(prefix)
		resultPtr := wasmGetOpKeysPrefix(prefixPtr)
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

// GetOpKeysWithPrefixString retrieves operation keys with the given string prefix
func GetOpKeysWithPrefixString(prefix string) []string {
	keys := GetOpKeysWithPrefix([]byte(prefix))
	result := make([]string, len(keys))
	for i, key := range keys {
		result[i] = string(key)
	}
	return result
}

// HasOpKey checks if an operation key exists
func HasOpKey(key []byte) bool {
	if isWasmRuntime() {
		keyPtr := newWasmBytes(key)
		resultPtr := wasmHasOpKey(keyPtr)
		// The result is a boolean encoded as bytes
		data := wasmBytesFromPtr(resultPtr)
		return len(data) > 0 && data[0] != 0
	}
	// In non-WASM mode, return false for testing
	return false
}

// HasOpKeyString checks if an operation key exists (string version)
func HasOpKeyString(key string) bool {
	return HasOpKey([]byte(key))
}

// GetOpData retrieves data for an operation key
func GetOpData(key []byte) []byte {
	if isWasmRuntime() {
		keyPtr := newWasmBytes(key)
		resultPtr := wasmGetOpData(keyPtr)
		return wasmBytesFromPtr(resultPtr)
	}
	// In non-WASM mode, return empty for testing
	return nil
}

// GetOpDataString retrieves data for an operation key as a string
func GetOpDataString(key string) string {
	data := GetOpData([]byte(key))
	return string(data)
}

// Operation datastore utilities

// GetOpDataJSON retrieves and deserializes JSON data from the operation datastore
func GetOpDataJSON(key string, value interface{}) error {
	data := GetOpDataString(key)
	if data == "" {
		return nil // Return nil for missing keys
	}
	return json.Unmarshal([]byte(data), value)
}

// IterateOpKeys iterates over all operation keys and calls the provided function for each
func IterateOpKeys(fn func(key []byte) bool) {
	keys := GetOpKeys()
	for _, key := range keys {
		if !fn(key) {
			break
		}
	}
}

// IterateOpKeysString iterates over all operation keys as strings
func IterateOpKeysString(fn func(key string) bool) {
	IterateOpKeys(func(key []byte) bool {
		return fn(string(key))
	})
}

// IterateOpKeysWithPrefix iterates over operation keys with a given prefix
func IterateOpKeysWithPrefix(prefix []byte, fn func(key []byte) bool) {
	keys := GetOpKeysWithPrefix(prefix)
	for _, key := range keys {
		if !fn(key) {
			break
		}
	}
}

// IterateOpKeysWithPrefixString iterates over operation keys with a given string prefix
func IterateOpKeysWithPrefixString(prefix string, fn func(key string) bool) {
	IterateOpKeysWithPrefix([]byte(prefix), func(key []byte) bool {
		return fn(string(key))
	})
}

// GetOpKeyValuePairs retrieves all key-value pairs from the operation datastore
func GetOpKeyValuePairs() map[string][]byte {
	result := make(map[string][]byte)
	keys := GetOpKeys()
	for _, key := range keys {
		data := GetOpData(key)
		result[string(key)] = data
	}
	return result
}

// GetOpKeyValuePairsString retrieves all key-value pairs as strings
func GetOpKeyValuePairsString() map[string]string {
	result := make(map[string]string)
	pairs := GetOpKeyValuePairs()
	for key, value := range pairs {
		result[key] = string(value)
	}
	return result
}

// GetOpKeyValuePairsWithPrefix retrieves key-value pairs with a given prefix
func GetOpKeyValuePairsWithPrefix(prefix []byte) map[string][]byte {
	result := make(map[string][]byte)
	keys := GetOpKeysWithPrefix(prefix)
	for _, key := range keys {
		data := GetOpData(key)
		result[string(key)] = data
	}
	return result
}

// GetOpKeyValuePairsWithPrefixString retrieves key-value pairs with a given string prefix
func GetOpKeyValuePairsWithPrefixString(prefix string) map[string]string {
	result := make(map[string]string)
	pairs := GetOpKeyValuePairsWithPrefix([]byte(prefix))
	for key, value := range pairs {
		result[key] = string(value)
	}
	return result
}

// Utility functions for operation datastore analysis

// CountOpKeys returns the total number of operation keys
func CountOpKeys() int {
	return len(GetOpKeys())
}

// CountOpKeysWithPrefix returns the number of operation keys with a given prefix
func CountOpKeysWithPrefix(prefix []byte) int {
	return len(GetOpKeysWithPrefix(prefix))
}

// CountOpKeysWithPrefixString returns the number of operation keys with a given string prefix
func CountOpKeysWithPrefixString(prefix string) int {
	return CountOpKeysWithPrefix([]byte(prefix))
}

// FilterOpKeys filters operation keys by a predicate function
func FilterOpKeys(predicate func(key []byte) bool) [][]byte {
	var result [][]byte
	keys := GetOpKeys()
	for _, key := range keys {
		if predicate(key) {
			result = append(result, key)
		}
	}
	return result
}

// FilterOpKeysString filters operation keys by a predicate function (string version)
func FilterOpKeysString(predicate func(key string) bool) []string {
	var result []string
	keys := GetOpKeysString()
	for _, key := range keys {
		if predicate(key) {
			result = append(result, key)
		}
	}
	return result
}

// FindOpKey finds the first operation key that matches a predicate
func FindOpKey(predicate func(key []byte) bool) []byte {
	keys := GetOpKeys()
	for _, key := range keys {
		if predicate(key) {
			return key
		}
	}
	return nil
}

// FindOpKeyString finds the first operation key that matches a predicate (string version)
func FindOpKeyString(predicate func(key string) bool) string {
	result := FindOpKey(func(key []byte) bool {
		return predicate(string(key))
	})
	return string(result)
}

// Operation datastore constants and helpers

// Common operation datastore key prefixes
const (
	OpKeyPrefixConfig     = "config:"
	OpKeyPrefixMetadata   = "metadata:"
	OpKeyPrefixTempData   = "temp:"
	OpKeyPrefixUserData   = "user:"
	OpKeyPrefixSystemData = "system:"
)

// GetConfigOpData retrieves configuration data from the operation datastore
func GetConfigOpData(key string) string {
	return GetOpDataString(OpKeyPrefixConfig + key)
}

// GetMetadataOpData retrieves metadata from the operation datastore
func GetMetadataOpData(key string) string {
	return GetOpDataString(OpKeyPrefixMetadata + key)
}

// GetTempOpData retrieves temporary data from the operation datastore
func GetTempOpData(key string) string {
	return GetOpDataString(OpKeyPrefixTempData + key)
}

// GetUserOpData retrieves user data from the operation datastore
func GetUserOpData(key string) string {
	return GetOpDataString(OpKeyPrefixUserData + key)
}

// GetSystemOpData retrieves system data from the operation datastore
func GetSystemOpData(key string) string {
	return GetOpDataString(OpKeyPrefixSystemData + key)
}

// GetAllConfigOpData retrieves all configuration data
func GetAllConfigOpData() map[string]string {
	return GetOpKeyValuePairsWithPrefixString(OpKeyPrefixConfig)
}

// GetAllMetadataOpData retrieves all metadata
func GetAllMetadataOpData() map[string]string {
	return GetOpKeyValuePairsWithPrefixString(OpKeyPrefixMetadata)
}

// GetAllTempOpData retrieves all temporary data
func GetAllTempOpData() map[string]string {
	return GetOpKeyValuePairsWithPrefixString(OpKeyPrefixTempData)
}

// GetAllUserOpData retrieves all user data
func GetAllUserOpData() map[string]string {
	return GetOpKeyValuePairsWithPrefixString(OpKeyPrefixUserData)
}

// GetAllSystemOpData retrieves all system data
func GetAllSystemOpData() map[string]string {
	return GetOpKeyValuePairsWithPrefixString(OpKeyPrefixSystemData)
}

// Operation datastore debugging and inspection

// DumpOpDatastore dumps all operation datastore contents for debugging
func DumpOpDatastore() map[string]interface{} {
	dump := map[string]interface{}{
		"total_keys": CountOpKeys(),
		"config":     GetAllConfigOpData(),
		"metadata":   GetAllMetadataOpData(),
		"temp":       GetAllTempOpData(),
		"user":       GetAllUserOpData(),
		"system":     GetAllSystemOpData(),
		"other":      make(map[string]string),
	}

	// Get all other keys that don't match known prefixes
	allData := GetOpKeyValuePairsString()
	otherData := dump["other"].(map[string]string)

	for key, value := range allData {
		if !hasKnownPrefix(key) {
			otherData[key] = value
		}
	}

	return dump
}

// hasKnownPrefix checks if a key has a known prefix
func hasKnownPrefix(key string) bool {
	prefixes := []string{
		OpKeyPrefixConfig,
		OpKeyPrefixMetadata,
		OpKeyPrefixTempData,
		OpKeyPrefixUserData,
		OpKeyPrefixSystemData,
	}

	for _, prefix := range prefixes {
		if len(key) >= len(prefix) && key[:len(prefix)] == prefix {
			return true
		}
	}

	return false
}

// LogOpDatastoreContents logs the operation datastore contents for debugging
func LogOpDatastoreContents() {
	dump := DumpOpDatastore()
	LogDebug("Operation datastore contents", dump)
}
