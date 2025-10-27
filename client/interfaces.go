package client

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

// ReadOnlyCallParams is the struct used to send a read only callSC to the node.
type ReadOnlyCallParams struct {
	MaxGas         int           `json:"max_gas"`
	Coins          string        `json:"coins"`
	Fee            string        `json:"fee"`
	TargetAddress  string        `json:"target_address"`
	TargetFunction string        `json:"target_function"`
	Parameter      JSONableSlice `json:"parameter"`
	CallerAddress  string        `json:"caller_address"`
}

// ReadOnlyExecuteParams is the struct used to send a read only executeSC to the node.
type ReadOnlyExecuteParams struct {
	MaxGas             int           `json:"max_gas"`
	Coins              string        `json:"coins"`
	Fee                string        `json:"fee"`
	Address            string        `json:"address"`
	Bytecode           JSONableSlice `json:"bytecode"`
	OperationDatastore JSONableSlice `json:"operation_datastore"`
}

type ReadOnlyResult struct {
	ExecutedAt   Timestamp   `json:"executed_at"`
	Result       Result      `json:"result"`
	OutputEvents []Event     `json:"output_events"`
	GasCost      int         `json:"gas_cost"`
	StateChanges StateChange `json:"state_changes"`
}

// Result is a struct that can hold both 'Ok' and 'Error' fields.
//
//nolint:tagliatelle
type Result struct {
	Ok    []interface{} `json:"Ok,omitempty"`
	Error string        `json:"Error,omitempty"`
}

// Custom UnmarshalJSON to handle both cases for the Result field.
func (r *Result) UnmarshalJSON(data []byte) error {
	var res map[string]json.RawMessage
	if err := json.Unmarshal(data, &res); err != nil {
		return fmt.Errorf("unmarshaling result: %w", err)
	}

	if message, ok := res["Ok"]; ok {
		var okSlice []interface{}
		if err := json.Unmarshal(message, &okSlice); err != nil {
			return fmt.Errorf("unmarshaling ok slice: %w", err)
		}

		r.Ok = okSlice

		return nil
	}

	if message, ok := res["Error"]; ok {
		var errorMessage string
		if err := json.Unmarshal(message, &errorMessage); err != nil {
			return fmt.Errorf("unmarshaling error message: %w", err)
		}

		r.Error = errorMessage

		return nil
	}

	return nil
}

type Timestamp struct {
	Period int `json:"period"`
	Thread int `json:"thread"`
}

type Event struct {
	Context EventContext `json:"context"`
	Data    string       `json:"data"`
}

type EventContext struct {
	Slot              Timestamp   `json:"slot"`
	Block             interface{} `json:"block"` // Assuming `null` can be represented as an `interface{}`
	ReadOnly          bool        `json:"read_only"`
	IndexInSlot       int         `json:"index_in_slot"`
	CallStack         []string    `json:"call_stack"`
	OriginOperationID interface{} `json:"origin_operation_id"` // Assuming `null` can be represented as an `interface{}`
	IsFinal           bool        `json:"is_final"`
	IsError           bool        `json:"is_error"`
}

type EventSearchCriteria struct {
	Start                 *Slot   `json:"start"`
	End                   *Slot   `json:"end"`
	EmitterAddress        *string `json:"emitter_address"`
	OriginalCallerAddress *string `json:"original_caller_address"`
	OriginalOperationID   *string `json:"original_operation_id"`
}

//nolint:lll
type StateChange struct {
	LedgerChanges                map[string]LedgerEntryChange `json:"ledger_changes"`
	AsyncPoolChanges             []interface{}                `json:"async_pool_changes"` // Empty in provided JSON, assumed to be a map
	PosChanges                   PosChanges                   `json:"pos_changes"`
	ExecutedOpsChanges           map[string]interface{}       `json:"executed_ops_changes"`
	ExecutedDenunciationsChanges []interface{}                `json:"executed_denunciations_changes"` // Empty array in provided JSON
	ExecutionTrailHashChange     interface{}                  `json:"execution_trail_hash_change"`
}

// ReadOnlyResultToBytes extracts the returned payload bytes from a ReadOnlyResult.
// It supports several JSON shapes returned by the node's `Ok` field:
// - top-level numeric array: Ok: [14,0,0,0]
// - nested numeric array: Ok: [[14,0,0,0]]
// - base64-encoded string: Ok: ["<base64>"]
// - raw string: Ok: ["raw"]
// - []byte when already decoded
// Returns an error if no payload is present or the shape is unsupported.
func ReadOnlyResultToBytes(r *ReadOnlyResult) ([]byte, error) {
	if r == nil {
		return nil, fmt.Errorf("nil ReadOnlyResult")
	}
	if r.Result.Ok == nil || len(r.Result.Ok) == 0 {
		return nil, fmt.Errorf("no Ok value in result")
	}

	// If the top-level Ok slice contains numeric values directly (e.g., [14,0,0,0]),
	// treat the entire slice as bytes.
	first := r.Result.Ok[0]
	switch first.(type) {
	case float64, int:
		out := make([]byte, len(r.Result.Ok))
		for i, iv := range r.Result.Ok {
			switch n := iv.(type) {
			case float64:
				out[i] = byte(uint8(n))
			case int:
				out[i] = byte(n)
			default:
				out[i] = 0
			}
		}
		return out, nil
	}

	// Otherwise inspect the first element for nested shapes.
	switch v := first.(type) {
	case []interface{}:
		out := make([]byte, len(v))
		for i, iv := range v {
			switch n := iv.(type) {
			case float64:
				out[i] = byte(uint8(n))
			case int:
				out[i] = byte(n)
			default:
				out[i] = 0
			}
		}
		return out, nil
	case string:
		// try base64 first, otherwise return raw bytes
		if b, err := base64.StdEncoding.DecodeString(v); err == nil {
			return b, nil
		}
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return nil, fmt.Errorf("unsupported Ok[0] type: %T", first)
	}
}

//nolint:tagliatelle
type LedgerEntryChange struct {
	Update LedgerUpdate `json:"Update"`
}

type LedgerUpdate struct {
	Balance   ChangeSet     `json:"balance"`
	Bytecode  interface{}   `json:"bytecode"`
	Datastore []interface{} `json:"datastore"`
}

//nolint:tagliatelle
type ChangeSet struct {
	Set interface{} `json:"Set"` // Set can be of different types, using `interface{}`
	// Depending on usage, more fields representing other change types may be necessary.
}

type PosChanges struct {
	SeedBits        SeedBitsInfo           `json:"seed_bits"`
	RollChanges     map[string]interface{} `json:"roll_changes"`     // Empty in provided JSON, assumed to be a map
	ProductionStats map[string]interface{} `json:"production_stats"` // Empty in provided JSON, assumed to be a map
	DeferredCredits DeferredCreditsInfo    `json:"deferred_credits"`
}

type SeedBitsInfo struct {
	Order string         `json:"order"`
	Head  BitVecHeadInfo `json:"head"`
	Bits  int            `json:"bits"`
	Data  []interface{}  `json:"data"` // Empty array in provided JSON
}

type BitVecHeadInfo struct {
	Width int `json:"width"`
	Index int `json:"index"`
}

type DeferredCreditsInfo struct {
	Credits map[string]interface{} `json:"credits"` // Empty in provided JSON, assumed to be a map
}
