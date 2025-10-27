package sc

import (
	"encoding/json"
	"fmt"
)

// Events and Logging functions for Massa smart contracts

// GenerateEvent generates an event on the blockchain with the given message
func GenerateEvent(message string) {
	if isWasmRuntime() {
		messagePtr := newWasmString(message)
		wasmGenerateEvent(messagePtr)
	}
	// In non-WASM mode, this would log to mock event system for testing
}

// Print outputs a message for debugging purposes
// This is visible in the smart contract execution logs
func Print(message string) {
	if isWasmRuntime() {
		messagePtr := newWasmString(message)
		wasmPrint(messagePtr)
	}
	// In non-WASM mode, this would log to console for testing
}

// Printf formats and prints a message for debugging purposes
func Printf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	Print(message)
}

// Event represents a structured event that can be emitted
type Event struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

// CreateEvent creates a structured event and emits it
func CreateEvent(eventType string, data interface{}) error {
	event := Event{
		Type: eventType,
		Data: data,
	}

	jsonData, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	GenerateEvent(string(jsonData))
	return nil
}

// Logging levels
const (
	LogLevelDebug = "DEBUG"
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARN"
	LogLevelError = "ERROR"
)

// LogEntry represents a structured log entry
type LogEntry struct {
	Level     string      `json:"level"`
	Message   string      `json:"message"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp uint64      `json:"timestamp"`
}

// Log emits a structured log entry
func Log(level string, message string, data interface{}) {
	entry := LogEntry{
		Level:     level,
		Message:   message,
		Data:      data,
		Timestamp: Timestamp(),
	}

	jsonData, err := json.Marshal(entry)
	if err != nil {
		// Fallback to simple message if JSON marshaling fails
		Print(fmt.Sprintf("[%s] %s", level, message))
		return
	}

	GenerateEvent(string(jsonData))
}

// Convenience logging functions

// LogDebug logs a debug message
func LogDebug(message string, data ...interface{}) {
	var logData interface{}
	if len(data) > 0 {
		if len(data) == 1 {
			logData = data[0]
		} else {
			logData = data
		}
	}
	Log(LogLevelDebug, message, logData)
}

// LogInfo logs an info message
func LogInfo(message string, data ...interface{}) {
	var logData interface{}
	if len(data) > 0 {
		if len(data) == 1 {
			logData = data[0]
		} else {
			logData = data
		}
	}
	Log(LogLevelInfo, message, logData)
}

// LogWarn logs a warning message
func LogWarn(message string, data ...interface{}) {
	var logData interface{}
	if len(data) > 0 {
		if len(data) == 1 {
			logData = data[0]
		} else {
			logData = data
		}
	}
	Log(LogLevelWarn, message, logData)
}

// LogError logs an error message
func LogError(message string, data ...interface{}) {
	var logData interface{}
	if len(data) > 0 {
		if len(data) == 1 {
			logData = data[0]
		} else {
			logData = data
		}
	}
	Log(LogLevelError, message, logData)
}

// Event types for common smart contract events

// ContractDeployedEvent represents a contract deployment event
type ContractDeployedEvent struct {
	Address   *Address `json:"address"`
	Deployer  *Address `json:"deployer"`
	Timestamp uint64   `json:"timestamp"`
}

// EmitContractDeployed emits a contract deployed event
func EmitContractDeployed(address, deployer *Address) error {
	return CreateEvent("contract_deployed", ContractDeployedEvent{
		Address:   address,
		Deployer:  deployer,
		Timestamp: Timestamp(),
	})
}

// FunctionCallEvent represents a function call event
type FunctionCallEvent struct {
	Contract    *Address `json:"contract"`
	Function    string   `json:"function"`
	Caller      *Address `json:"caller"`
	Args        []byte   `json:"args,omitempty"`
	CoinsAmount uint64   `json:"coins_amount,omitempty"`
	Timestamp   uint64   `json:"timestamp"`
}

// EmitFunctionCall emits a function call event
func EmitFunctionCall(contract *Address, function string, args []byte, coinsAmount uint64) error {
	return CreateEvent("function_call", FunctionCallEvent{
		Contract:    contract,
		Function:    function,
		Caller:      Caller(),
		Args:        args,
		CoinsAmount: coinsAmount,
		Timestamp:   Timestamp(),
	})
}

// ErrorEvent represents an error event
type ErrorEvent struct {
	Message   string      `json:"message"`
	Code      string      `json:"code,omitempty"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp uint64      `json:"timestamp"`
	Context   *Address    `json:"context"`
}

// EmitError emits an error event
func EmitError(message string, code string, data interface{}) error {
	return CreateEvent("error", ErrorEvent{
		Message:   message,
		Code:      code,
		Data:      data,
		Timestamp: Timestamp(),
		Context:   Callee(),
	})
}

// OwnershipTransferEvent represents an ownership transfer event
type OwnershipTransferEvent struct {
	From      *Address `json:"from"`
	To        *Address `json:"to"`
	Asset     string   `json:"asset,omitempty"`
	Timestamp uint64   `json:"timestamp"`
}

// EmitOwnershipTransfer emits an ownership transfer event
func EmitOwnershipTransfer(from, to *Address, asset string) error {
	return CreateEvent("ownership_transfer", OwnershipTransferEvent{
		From:      from,
		To:        to,
		Asset:     asset,
		Timestamp: Timestamp(),
	})
}

// PermissionChangeEvent represents a permission change event
type PermissionChangeEvent struct {
	User       *Address `json:"user"`
	Permission string   `json:"permission"`
	Granted    bool     `json:"granted"`
	Grantor    *Address `json:"grantor"`
	Timestamp  uint64   `json:"timestamp"`
}

// EmitPermissionChange emits a permission change event
func EmitPermissionChange(user *Address, permission string, granted bool) error {
	return CreateEvent("permission_change", PermissionChangeEvent{
		User:       user,
		Permission: permission,
		Granted:    granted,
		Grantor:    Caller(),
		Timestamp:  Timestamp(),
	})
}

// Custom event emitters

// EmitCustomEvent emits a custom event with arbitrary data
func EmitCustomEvent(eventType string, data map[string]interface{}) error {
	// Add timestamp if not provided
	if _, exists := data["timestamp"]; !exists {
		data["timestamp"] = Timestamp()
	}

	return CreateEvent(eventType, data)
}

// EmitJSON emits an event with pre-formatted JSON data
func EmitJSON(jsonData string) {
	GenerateEvent(jsonData)
}

// Debug helpers

// AssertTrue asserts that a condition is true, logging an error if not
func AssertTrue(condition bool, message string) {
	if !condition {
		LogError("Assertion failed: " + message)
		EmitError("Assertion failed", "ASSERTION_ERROR", map[string]interface{}{
			"message": message,
			"caller":  Caller().String(),
		})
	}
}

// AssertFalse asserts that a condition is false, logging an error if not
func AssertFalse(condition bool, message string) {
	AssertTrue(!condition, message)
}

// AssertEqual asserts that two values are equal
func AssertEqual(a, b interface{}, message string) {
	if a != b {
		LogError(fmt.Sprintf("Assertion failed: %s (expected %v, got %v)", message, a, b))
		EmitError("Assertion failed", "ASSERTION_ERROR", map[string]interface{}{
			"message":  message,
			"expected": a,
			"actual":   b,
			"caller":   Caller().String(),
		})
	}
}

// TraceExecution logs function entry and exit for debugging
func TraceExecution(functionName string) func() {
	LogDebug(fmt.Sprintf("Entering function: %s", functionName))
	return func() {
		LogDebug(fmt.Sprintf("Exiting function: %s", functionName))
	}
}
