//go:build js && wasm

package sc

import "unsafe"

// WebAssembly runtime implementations
// These functions are imported from the Massa WebAssembly runtime

//go:wasmimport massa assembly_script_print
//go:noescape
func wasmPrint(ptr uintptr)

//go:wasmimport massa assembly_script_call
//go:noescape
func wasmCall(addressPtr, funcPtr, paramPtr uintptr, coins uint64) uintptr

//go:wasmimport massa assembly_script_local_call
//go:noescape
func wasmLocalCall(addressPtr, funcPtr, paramPtr uintptr) uintptr

//go:wasmimport massa assembly_script_local_execution
//go:noescape
func wasmLocalExecution(bytecodePtr, funcPtr, paramPtr uintptr) uintptr

//go:wasmimport massa assembly_script_get_bytecode
//go:noescape
func wasmGetBytecode() uintptr

//go:wasmimport massa assembly_script_get_bytecode_for
//go:noescape
func wasmGetBytecodeOf(addressPtr uintptr) uintptr

//go:wasmimport massa assembly_script_caller_has_write_access
//go:noescape
func wasmCallerHasWriteAccess() uint32

//go:wasmimport massa assembly_script_function_exists
//go:noescape
func wasmFunctionExists(addressPtr, funcPtr uintptr) uint32

//go:wasmimport massa assembly_script_get_remaining_gas
//go:noescape
func wasmRemainingGas() uint64

//go:wasmimport massa assembly_script_create_sc
//go:noescape
func wasmCreateSC(bytecodePtr uintptr) uintptr

//go:wasmimport massa assembly_script_get_keys
//go:noescape
func wasmGetKeys(prefixPtr uintptr) uintptr

//go:wasmimport massa assembly_script_get_keys_for
//go:noescape
func wasmGetKeysOf(addressPtr, prefixPtr uintptr) uintptr

//go:wasmimport massa assembly_script_set_data
//go:noescape
func wasmSetData(keyPtr, valuePtr uintptr)

//go:wasmimport massa assembly_script_set_data_for
//go:noescape
func wasmSetDataFor(addressPtr, keyPtr, valuePtr uintptr)

//go:wasmimport massa assembly_script_get_data
//go:noescape
func wasmGetData(keyPtr uintptr) uintptr

//go:wasmimport massa assembly_script_get_data_for
//go:noescape
func wasmGetDataFor(addressPtr, keyPtr uintptr) uintptr

//go:wasmimport massa assembly_script_delete_data
//go:noescape
func wasmDeleteData(keyPtr uintptr)

//go:wasmimport massa assembly_script_delete_data_for
//go:noescape
func wasmDeleteDataFor(addressPtr, keyPtr uintptr)

//go:wasmimport massa assembly_script_append_data
//go:noescape
func wasmAppendData(keyPtr, valuePtr uintptr)

//go:wasmimport massa assembly_script_append_data_for
//go:noescape
func wasmAppendDataFor(addressPtr, keyPtr, valuePtr uintptr)

//go:wasmimport massa assembly_script_has_data
//go:noescape
func wasmHasData(keyPtr uintptr) uint32

//go:wasmimport massa assembly_script_has_data_for
//go:noescape
func wasmHasDataFor(addressPtr, keyPtr uintptr) uint32

//go:wasmimport massa assembly_script_get_owned_addresses
//go:noescape
func wasmGetOwnedAddresses() uintptr

//go:wasmimport massa assembly_script_get_call_stack
//go:noescape
func wasmGetCallStack() uintptr

//go:wasmimport massa assembly_script_generate_event
//go:noescape
func wasmGenerateEvent(eventPtr uintptr)

//go:wasmimport massa assembly_script_transfer_coins
//go:noescape
func wasmTransferCoins(toPtr uintptr, amount uint64)

//go:wasmimport massa assembly_script_transfer_coins_for
//go:noescape
func wasmTransferCoinsFor(fromPtr, toPtr uintptr, amount uint64)

//go:wasmimport massa assembly_script_get_balance
//go:noescape
func wasmGetBalance() uint64

//go:wasmimport massa assembly_script_get_balance_for
//go:noescape
func wasmGetBalanceFor(addressPtr uintptr) uint64

//go:wasmimport massa assembly_script_get_call_coins
//go:noescape
func wasmGetCallCoins() uint64

//go:wasmimport massa assembly_script_hash
//go:noescape
func wasmBlake3(bytecodePtr uintptr) uintptr

//go:wasmimport massa assembly_script_signature_verify
//go:noescape
func wasmSignatureVerify(dataPtr, signaturePtr, publicKeyPtr uintptr) uint32

//go:wasmimport massa assembly_script_evm_signature_verify
//go:noescape
func wasmEvmSignatureVerify(dataPtr, signaturePtr, publicKeyPtr uintptr) uint32

//go:wasmimport massa assembly_script_evm_get_address_from_pubkey
//go:noescape
func wasmEvmGetAddressFromPubkey(publicKeyPtr uintptr) uintptr

//go:wasmimport massa assembly_script_evm_get_pubkey_from_signature
//go:noescape
func wasmEvmGetPubkeyFromSignature(hashPtr, signaturePtr uintptr) uintptr

//go:wasmimport massa assembly_script_is_address_eoa
//go:noescape
func wasmIsAddressEoa(addressPtr uintptr) uint32

//go:wasmimport massa assembly_script_address_from_public_key
//go:noescape
func wasmPublicKeyToAddress(publicKeyPtr uintptr) uintptr

//go:wasmimport massa assembly_script_get_time
//go:noescape
func wasmGetTime() uint64

//go:wasmimport massa assembly_script_unsafe_random
//go:noescape
func wasmUnsafeRandom() int64

//go:wasmimport massa assembly_script_send_message
//go:noescape
func wasmSendMessage(addressPtr, handlerPtr uintptr, validityStartPeriod uint64, validityStartThread uint32, validityEndPeriod uint64, validityEndThread uint32, maxGas uint64, rawFee uint64, coins uint64, dataPtr, filterAddressPtr, filterKeyPtr uintptr)

//go:wasmimport massa assembly_script_get_origin_operation_id
//go:noescape
func wasmGetOriginOperationId() uintptr

//go:wasmimport massa assembly_script_get_current_period
//go:noescape
func wasmGetCurrentPeriod() uint64

//go:wasmimport massa assembly_script_get_current_thread
//go:noescape
func wasmGetCurrentThread() uint32

//go:wasmimport massa assembly_script_set_bytecode
//go:noescape
func wasmSetBytecode(bytecodePtr uintptr)

//go:wasmimport massa assembly_script_set_bytecode_for
//go:noescape
func wasmSetBytecodeFor(addressPtr, bytecodePtr uintptr)

//go:wasmimport massa assembly_script_get_op_keys
//go:noescape
func wasmGetOpKeys() uintptr

//go:wasmimport massa assembly_script_get_op_keys_prefix
//go:noescape
func wasmGetOpKeysPrefix(prefixPtr uintptr) uintptr

//go:wasmimport massa assembly_script_has_op_key
//go:noescape
func wasmHasOpKey(keyPtr uintptr) uintptr

//go:wasmimport massa assembly_script_get_op_data
//go:noescape
func wasmGetOpData(keyPtr uintptr) uintptr

//go:wasmimport massa assembly_script_hash_sha256
//go:noescape
func wasmSha256(bytecodePtr uintptr) uintptr

//go:wasmimport massa assembly_script_hash_mimc
//go:noescape
func wasmMimc(dataPtr uintptr) uintptr

//go:wasmimport massa assembly_script_keccak256_hash
//go:noescape
func wasmKeccak256(dataPtr uintptr) uintptr

//go:wasmimport massa assembly_script_validate_address
//go:noescape
func wasmValidateAddress(addressPtr uintptr) uint32

//go:wasmimport massa assembly_script_chain_id
//go:noescape
func wasmChainId() uint64

//go:wasmimport massa assembly_script_get_deferred_call_quote
//go:noescape
func wasmGetDeferredCallQuote(period uint64, thread uint32, maxGas uint64, paramsSize uint64) uint64

//go:wasmimport massa assembly_script_deferred_call_register
//go:noescape
func wasmDeferredCallRegister(targetAddressPtr, targetFunctionPtr uintptr, targetPeriod uint64, targetThread uint32, maxGas uint64, paramsPtr uintptr, rawCoins uint64) uintptr

//go:wasmimport massa assembly_script_deferred_call_exists
//go:noescape
func wasmDeferredCallExists(idPtr uintptr) uint32

//go:wasmimport massa assembly_script_deferred_call_cancel
//go:noescape
func wasmDeferredCallCancel(idPtr uintptr)

// Memory management helpers for WebAssembly

// newWasmString allocates a string in WebAssembly memory and returns a pointer to it
func newWasmString(s string) uintptr {
	if len(s) == 0 {
		return 0
	}
	return uintptr(unsafe.Pointer(unsafe.StringData(s)))
}

// newWasmBytes allocates a byte slice in WebAssembly memory and returns a pointer to it
func newWasmBytes(data []byte) uintptr {
	if len(data) == 0 {
		return 0
	}
	return uintptr(unsafe.Pointer(&data[0]))
}

// wasmStringFromPtr converts a WebAssembly memory pointer to a Go string
func wasmStringFromPtr(ptr uintptr) string {
	if ptr == 0 {
		return ""
	}
	// In actual WebAssembly runtime, this would properly read from memory
	// This is a placeholder implementation
	return ""
}

// wasmBytesFromPtr converts a WebAssembly memory pointer to a Go byte slice
func wasmBytesFromPtr(ptr uintptr) []byte {
	if ptr == 0 {
		return nil
	}
	// In actual WebAssembly runtime, this would properly read from memory
	// This is a placeholder implementation
	return nil
}

// isWasmRuntime returns true when running in WebAssembly
func isWasmRuntime() bool {
	return true
}
