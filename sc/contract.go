package sc

// Contract provides functions for interacting with other smart contract functions
// and manipulating their bytecode.

// Args represents serialized arguments for contract calls
type Args struct {
	data []byte
}

// NewArgs creates a new Args from byte data
func NewArgs(data []byte) *Args {
	return &Args{data: data}
}

// NewArgsFromString creates a new Args from a string
func NewArgsFromString(s string) *Args {
	return &Args{data: []byte(s)}
}

// Serialize returns the serialized arguments as bytes
func (a *Args) Serialize() []byte {
	return a.data
}

// String returns the arguments as a string
func (a *Args) String() string {
	return string(a.data)
}

// Call invokes a function of a smart contract deployed at a given address
func Call(at *Address, functionName string, args *Args, coins uint64) []byte {
	if args == nil {
		args = NewArgs(nil)
	}

	if isWasmRuntime() {
		addressPtr := newWasmString(at.String())
		funcPtr := newWasmString(functionName)
		paramPtr := newWasmBytes(args.Serialize())
		resultPtr := wasmCall(addressPtr, funcPtr, paramPtr, coins)
		return wasmBytesFromPtr(resultPtr)
	}

	// In non-WASM mode, return empty for testing
	return nil
}

// CallString is a convenience function for calling with string arguments
func CallString(at *Address, functionName string, args string, coins uint64) string {
	result := Call(at, functionName, NewArgsFromString(args), coins)
	return string(result)
}

// LocalCall calls a function of a smart contract in the current context
// without creating a new execution context
func LocalCall(at *Address, functionName string, args *Args) []byte {
	if args == nil {
		args = NewArgs(nil)
	}

	if isWasmRuntime() {
		addressPtr := newWasmString(at.String())
		funcPtr := newWasmString(functionName)
		paramPtr := newWasmBytes(args.Serialize())
		resultPtr := wasmLocalCall(addressPtr, funcPtr, paramPtr)
		return wasmBytesFromPtr(resultPtr)
	}

	// In non-WASM mode, return empty for testing
	return nil
}

// LocalCallString is a convenience function for local calling with string arguments
func LocalCallString(at *Address, functionName string, args string) string {
	result := LocalCall(at, functionName, NewArgsFromString(args))
	return string(result)
}

// LocalExecution calls a function from a contract's bytecode in the current context
func LocalExecution(bytecode []byte, functionName string, args *Args) []byte {
	if args == nil {
		args = NewArgs(nil)
	}

	if isWasmRuntime() {
		bytecodePtr := newWasmBytes(bytecode)
		funcPtr := newWasmString(functionName)
		paramPtr := newWasmBytes(args.Serialize())
		resultPtr := wasmLocalExecution(bytecodePtr, funcPtr, paramPtr)
		return wasmBytesFromPtr(resultPtr)
	}

	// In non-WASM mode, return empty for testing
	return nil
}

// LocalExecutionString is a convenience function for local execution with string arguments
func LocalExecutionString(bytecode []byte, functionName string, args string) string {
	result := LocalExecution(bytecode, functionName, NewArgsFromString(args))
	return string(result)
}

// FunctionExists checks if a function exists in a smart contract
func FunctionExists(address *Address, functionName string) bool {
	if isWasmRuntime() {
		addressPtr := newWasmString(address.String())
		funcPtr := newWasmString(functionName)
		result := wasmFunctionExists(addressPtr, funcPtr)
		return result != 0
	}

	// In non-WASM mode, return true for testing
	return true
}

// CallerHasWriteAccess determines if the caller has write access to the data
// stored in the called smart contract
func CallerHasWriteAccess() bool {
	if isWasmRuntime() {
		result := wasmCallerHasWriteAccess()
		return result != 0
	}

	// In non-WASM mode, return true for testing
	return true
}

// CreateSC creates a new smart contract and returns its address
func CreateSC(bytecode []byte) (*Address, error) {
	if isWasmRuntime() {
		bytecodePtr := newWasmBytes(bytecode)
		resultPtr := wasmCreateSC(bytecodePtr)
		addressStr := wasmStringFromPtr(resultPtr)
		return NewAddress(addressStr)
	}

	// In non-WASM mode, return a mock address for testing
	return NewAddress("AU12UBnqTHDQALpocVBnkPNy7y5CndUJQTLutaVDDFgMJcq5kQiKq")
}

// GetBytecode retrieves the bytecode of the current smart contract
func GetBytecode() []byte {
	if isWasmRuntime() {
		resultPtr := wasmGetBytecode()
		return wasmBytesFromPtr(resultPtr)
	}

	// In non-WASM mode, return empty for testing
	return nil
}

// GetBytecodeOf retrieves the bytecode of another smart contract
func GetBytecodeOf(address *Address) []byte {
	if isWasmRuntime() {
		addressPtr := newWasmString(address.String())
		resultPtr := wasmGetBytecodeOf(addressPtr)
		return wasmBytesFromPtr(resultPtr)
	}

	// In non-WASM mode, return empty for testing
	return nil
}

// SetBytecode sets the bytecode of the current smart contract
func SetBytecode(bytecode []byte) {
	if isWasmRuntime() {
		bytecodePtr := newWasmBytes(bytecode)
		wasmSetBytecode(bytecodePtr)
	}
	// In non-WASM mode, this would set bytecode in mock storage for testing
}

// SetBytecodeOf sets the bytecode of another smart contract
func SetBytecodeOf(address *Address, bytecode []byte) {
	if isWasmRuntime() {
		addressPtr := newWasmString(address.String())
		bytecodePtr := newWasmBytes(bytecode)
		wasmSetBytecodeFor(addressPtr, bytecodePtr)
	}
	// In non-WASM mode, this would set bytecode in mock storage for testing
}

// MessageParams represents parameters for sending a message
type MessageParams struct {
	Address             *Address
	Handler             string
	ValidityStartPeriod uint64
	ValidityStartThread uint8
	ValidityEndPeriod   uint64
	ValidityEndThread   uint8
	MaxGas              uint64
	RawFee              uint64
	Coins               uint64
	Data                []byte
	FilterAddress       *Address
	FilterKey           []byte
}

// SendMessage sends a message to another address
func SendMessage(params MessageParams) {
	if isWasmRuntime() {
		addressPtr := newWasmString(params.Address.String())
		handlerPtr := newWasmString(params.Handler)
		dataPtr := newWasmBytes(params.Data)

		var filterAddressPtr uintptr
		if params.FilterAddress != nil {
			filterAddressPtr = newWasmString(params.FilterAddress.String())
		}

		filterKeyPtr := newWasmBytes(params.FilterKey)

		wasmSendMessage(
			addressPtr,
			handlerPtr,
			params.ValidityStartPeriod,
			uint32(params.ValidityStartThread),
			params.ValidityEndPeriod,
			uint32(params.ValidityEndThread),
			params.MaxGas,
			params.RawFee,
			params.Coins,
			dataPtr,
			filterAddressPtr,
			filterKeyPtr,
		)
	}
	// In non-WASM mode, this would send a message via mock system for testing
}

// DeferredCallParams represents parameters for a deferred call
type DeferredCallParams struct {
	TargetAddress  *Address
	TargetFunction string
	TargetPeriod   uint64
	TargetThread   uint8
	MaxGas         uint64
	Params         []byte
	RawCoins       uint64
}

// GetDeferredCallQuote gets a quote for a deferred call
func GetDeferredCallQuote(period uint64, thread uint8, maxGas uint64, paramsSize uint64) uint64 {
	if isWasmRuntime() {
		return wasmGetDeferredCallQuote(period, uint32(thread), maxGas, paramsSize)
	}

	// In non-WASM mode, return a mock quote for testing
	return 1000
}

// RegisterDeferredCall registers a deferred call and returns its ID
func RegisterDeferredCall(params DeferredCallParams) string {
	if isWasmRuntime() {
		targetAddressPtr := newWasmString(params.TargetAddress.String())
		targetFunctionPtr := newWasmString(params.TargetFunction)
		paramsPtr := newWasmBytes(params.Params)

		resultPtr := wasmDeferredCallRegister(
			targetAddressPtr,
			targetFunctionPtr,
			params.TargetPeriod,
			uint32(params.TargetThread),
			params.MaxGas,
			paramsPtr,
			params.RawCoins,
		)

		return wasmStringFromPtr(resultPtr)
	}

	// In non-WASM mode, return a mock call ID for testing
	return "deferred_call_123"
}

// DeferredCallExists checks if a deferred call exists
func DeferredCallExists(callId string) bool {
	if isWasmRuntime() {
		idPtr := newWasmString(callId)
		result := wasmDeferredCallExists(idPtr)
		return result != 0
	}

	// In non-WASM mode, return true for testing
	return true
}

// CancelDeferredCall cancels a deferred call
func CancelDeferredCall(callId string) {
	if isWasmRuntime() {
		idPtr := newWasmString(callId)
		wasmDeferredCallCancel(idPtr)
	}
	// In non-WASM mode, this would cancel in mock system for testing
}
