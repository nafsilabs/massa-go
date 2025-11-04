# Massa Units Conversion

This package provides utilities for converting between Massa and nanoMassa units in the Massa blockchain.

## Overview

In the Massa blockchain:
- 1 Massa = 1,000,000,000 nanoMassa (10^9)
- nanoMassa is the smallest unit used in blockchain operations
- Massa is the human-readable denomination

This package uses Go's `math/big` package for arbitrary precision arithmetic to ensure no loss of precision when working with large amounts.

## Installation

```go
import "github.com/nafsilabs/massa-go/utils"
```

## Core Functions

### FromMAS
Convert Massa (float64) to nanoMassa (big.Int):
```go
amount := utils.FromMAS(1.5)
// Returns: 1500000000 nanoMassa
```

### ToMAS
Convert nanoMassa (big.Int) to Massa (float64):
```go
nanoMassa := big.NewInt(2_500_000_000)
massa := utils.ToMAS(nanoMassa)
// Returns: 2.5 Massa
```

### FromNanoMAS
Convert nanoMassa (uint64) to big.Int:
```go
amount := utils.FromNanoMAS(1_000_000_000)
// Returns: big.Int representing 1000000000
```

### ToNanoMAS
Convert big.Int to nanoMassa (uint64):
```go
bigInt := big.NewInt(500_000_000)
nano := utils.ToNanoMAS(bigInt)
// Returns: 500000000 (uint64)
```

### ParseMassa
Parse Massa amount (float64) directly to uint64 nanoMassa:
```go
nano := utils.ParseMassa(3.75)
// Returns: 3750000000 (uint64)
```

## Usage Examples

### Basic Conversion
```go
package main

import (
    "fmt"
    "github.com/nafsilabs/massa-go/utils"
)

func main() {
    // Convert 1.5 Massa to nanoMassa
    nanoMassa := utils.FromMAS(1.5)
    fmt.Printf("1.5 Massa = %s nanoMassa\n", nanoMassa.String())
    
    // Convert back to Massa
    massa := utils.ToMAS(nanoMassa)
    fmt.Printf("%s nanoMassa = %f Massa\n", nanoMassa.String(), massa)
}
```

### Transaction Calculation
```go
package main

import (
    "fmt"
    "math/big"
    "github.com/nafsilabs/massa-go/utils"
)

func main() {
    // Sender has 100 Massa
    balance := utils.FromMAS(100.0)
    
    // Transfer 25.5 Massa
    transfer := utils.FromMAS(25.5)
    
    // Fee is 0.01 Massa
    fee := utils.FromMAS(0.01)
    
    // Calculate remaining balance
    remaining := new(big.Int).Sub(balance, transfer)
    remaining.Sub(remaining, fee)
    
    fmt.Printf("Remaining balance: %.2f Massa\n", utils.ToMAS(remaining))
    // Output: Remaining balance: 74.49 Massa
}
```

### Working with Protocol Operations
```go
package main

import (
    "github.com/nafsilabs/massa-go/client"
    "github.com/nafsilabs/massa-go/utils"
)

func main() {
    // Create transaction with 10 Massa
    amount := utils.FromMAS(10.0)
    
    // Convert to uint64 for protocol operations
    amountNano := utils.ToNanoMAS(amount)
    
    // Use in transaction
    tx := client.NewTransactionOp(
        "recipient_address",
        amountNano,
        0, // expiry period
    )
}
```

## Constants

```go
const (
    DecimalScale      = 9                      // 9 decimal places
    NanoMassaPerMassa = 1_000_000_000         // 10^9 nanoMassa per Massa
)
```

## Precision Considerations

- **FromMAS/ToMAS**: Use for user-facing amounts, handles floating point with appropriate precision
- **FromNanoMAS/ToNanoMAS**: Use for protocol operations, provides uint64 ↔ big.Int conversion
- **ParseMassa**: Convenience function for quick conversions to uint64

### Floating Point Precision
Due to floating point limitations, very small amounts may experience rounding:
```go
// 1 nanoMassa
small := utils.FromMAS(0.000000001)
fmt.Println(small.String()) // "1"

// Round trip maintains precision for reasonable amounts
original := 10.123456789
converted := utils.FromMAS(original)
back := utils.ToMAS(converted)
// back == 10.123456789 (within floating point precision)
```

## Testing

Run the test suite:
```bash
cd utils
go test -v
```

See the example:
```bash
cd examples/massa_units
go run main.go
```

## Comparison with Other Languages

### Dart Implementation
This Go implementation mirrors the Dart massa_units.dart package:
- `fromMAS` → `FromMAS`
- `toMAS` → `ToMAS`
- Uses `big.Int` instead of Dart's `BigInt`
- Same conversion factor: 10^9

### JavaScript/TypeScript
For JavaScript/TypeScript implementations, consider using:
- `BigInt` for integer arithmetic
- Libraries like `bignumber.js` or `decimal.js` for decimal precision

## Best Practices

1. **Always use big.Int for calculations**: Prevents overflow and maintains precision
2. **Convert to float64 only for display**: Use `ToMAS()` when showing amounts to users
3. **Use uint64 for protocol operations**: Most blockchain operations expect uint64
4. **Validate user input**: Always check that user-provided amounts are valid before conversion

## Common Patterns

### Validating Sufficient Balance
```go
balance := utils.FromMAS(100.0)
required := utils.FromMAS(150.0)

if balance.Cmp(required) < 0 {
    return errors.New("insufficient balance")
}
```

### Calculating Total with Fee
```go
amount := utils.FromMAS(10.0)
fee := utils.FromMAS(0.01)

total := new(big.Int).Add(amount, fee)
fmt.Printf("Total: %.2f Massa\n", utils.ToMAS(total))
```

### Splitting Amounts
```go
total := utils.FromMAS(100.0)
parts := 3

perPart := new(big.Int).Div(total, big.NewInt(int64(parts)))
fmt.Printf("Each part: %.2f Massa\n", utils.ToMAS(perPart))
```

## See Also

- [Massa Documentation](https://docs.massa.net)
- [massa-go Client](../client/README.md)
- [massa-go Wallet](../wallet/README.md)
