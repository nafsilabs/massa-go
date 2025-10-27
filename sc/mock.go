//go:build !js || !wasm

package sc

import "fmt"

// Mock implementations for non-WebAssembly environments
// These are used for testing and development outside of the Massa runtime

// Mock WebAssembly functions - these provide basic functionality for testing

func wasmPrint(ptr uintptr) {
	fmt.Println("Mock print called")
}

func wasmCall(addressPtr, funcPtr, paramPtr uintptr, coins uint64) uintptr {
	return 0
}

func wasmLocalCall(addressPtr, funcPtr, paramPtr uintptr) uintptr {
	return 0
}

func wasmLocalExecution(bytecodePtr, funcPtr, paramPtr uintptr) uintptr {
	return 0
}

func wasmGetBytecode() uintptr {
	return 0
}

func wasmGetBytecodeOf(addressPtr uintptr) uintptr {
	return 0
}

func wasmCallerHasWriteAccess() uint32 {
	return 1
}

func wasmFunctionExists(addressPtr, funcPtr uintptr) uint32 {
	return 1
}

func wasmRemainingGas() uint64 {
	return 1000000
}

func wasmCreateSC(bytecodePtr uintptr) uintptr {
	return 0
}

func wasmGetKeys(prefixPtr uintptr) uintptr {
	return 0
}

func wasmGetKeysOf(addressPtr, prefixPtr uintptr) uintptr {
	return 0
}

func wasmSetData(keyPtr, valuePtr uintptr) {}

func wasmSetDataFor(addressPtr, keyPtr, valuePtr uintptr) {}

func wasmGetData(keyPtr uintptr) uintptr {
	return 0
}

func wasmGetDataFor(addressPtr, keyPtr uintptr) uintptr {
	return 0
}

func wasmDeleteData(keyPtr uintptr) {}

func wasmDeleteDataFor(addressPtr, keyPtr uintptr) {}

func wasmAppendData(keyPtr, valuePtr uintptr) {}

func wasmAppendDataFor(addressPtr, keyPtr, valuePtr uintptr) {}

func wasmHasData(keyPtr uintptr) uint32 {
	return 0
}

func wasmHasDataFor(addressPtr, keyPtr uintptr) uint32 {
	return 0
}

func wasmGetOwnedAddresses() uintptr {
	return 0
}

func wasmGetCallStack() uintptr {
	return 0
}

func wasmGenerateEvent(eventPtr uintptr) {}

func wasmTransferCoins(toPtr uintptr, amount uint64) {}

func wasmTransferCoinsFor(fromPtr, toPtr uintptr, amount uint64) {}

func wasmGetBalance() uint64 {
	return 1000000000
}

func wasmGetBalanceFor(addressPtr uintptr) uint64 {
	return 1000000000
}

func wasmGetCallCoins() uint64 {
	return 0
}

func wasmBlake3(bytecodePtr uintptr) uintptr {
	return 0
}

func wasmSignatureVerify(dataPtr, signaturePtr, publicKeyPtr uintptr) uint32 {
	return 1
}

func wasmEvmSignatureVerify(dataPtr, signaturePtr, publicKeyPtr uintptr) uint32 {
	return 1
}

func wasmEvmGetAddressFromPubkey(publicKeyPtr uintptr) uintptr {
	return 0
}

func wasmEvmGetPubkeyFromSignature(hashPtr, signaturePtr uintptr) uintptr {
	return 0
}

func wasmIsAddressEoa(addressPtr uintptr) uint32 {
	return 1
}

func wasmPublicKeyToAddress(publicKeyPtr uintptr) uintptr {
	return 0
}

func wasmGetTime() uint64 {
	return 1640995200000 // Mock timestamp
}

func wasmUnsafeRandom() int64 {
	return 42 // Mock random value
}

func wasmSendMessage(addressPtr, handlerPtr uintptr, validityStartPeriod uint64, validityStartThread uint32, validityEndPeriod uint64, validityEndThread uint32, maxGas uint64, rawFee uint64, coins uint64, dataPtr, filterAddressPtr, filterKeyPtr uintptr) {
}

func wasmGetOriginOperationId() uintptr {
	return 0
}

func wasmGetCurrentPeriod() uint64 {
	return 12345
}

func wasmGetCurrentThread() uint32 {
	return 0
}

func wasmSetBytecode(bytecodePtr uintptr) {}

func wasmSetBytecodeFor(addressPtr, bytecodePtr uintptr) {}

func wasmGetOpKeys() uintptr {
	return 0
}

func wasmGetOpKeysPrefix(prefixPtr uintptr) uintptr {
	return 0
}

func wasmHasOpKey(keyPtr uintptr) uintptr {
	return 0
}

func wasmGetOpData(keyPtr uintptr) uintptr {
	return 0
}

func wasmSha256(bytecodePtr uintptr) uintptr {
	return 0
}

func wasmMimc(dataPtr uintptr) uintptr {
	return 0
}

func wasmKeccak256(dataPtr uintptr) uintptr {
	return 0
}

func wasmValidateAddress(addressPtr uintptr) uint32 {
	return 1
}

func wasmChainId() uint64 {
	return 77658366
}

func wasmGetDeferredCallQuote(period uint64, thread uint32, maxGas uint64, paramsSize uint64) uint64 {
	return 1000
}

func wasmDeferredCallRegister(targetAddressPtr, targetFunctionPtr uintptr, targetPeriod uint64, targetThread uint32, maxGas uint64, paramsPtr uintptr, rawCoins uint64) uintptr {
	return 0
}

func wasmDeferredCallExists(idPtr uintptr) uint32 {
	return 0
}

func wasmDeferredCallCancel(idPtr uintptr) {}

// Memory management helpers for non-WebAssembly environments

func newWasmString(s string) uintptr {
	return uintptr(len(s)) // Mock implementation
}

func newWasmBytes(data []byte) uintptr {
	return uintptr(len(data)) // Mock implementation
}

func wasmStringFromPtr(ptr uintptr) string {
	return "mock_string"
}

func wasmBytesFromPtr(ptr uintptr) []byte {
	return []byte("mock_bytes")
}

// isWasmRuntime returns false when not running in WebAssembly
func isWasmRuntime() bool {
	return false
}
