package main

import (
	"fmt"
	"math/big"

	"github.com/nafsilabs/massa-go/utils"
)

func main() {
	fmt.Println("=== Massa Units Conversion Examples ===")

	// Example 1: Convert Massa to nanoMassa
	fmt.Println("1. Converting Massa to nanoMassa:")
	massaAmount := 1.5
	nanoMassa := utils.FromMAS(massaAmount)
	fmt.Printf("   %f Massa = %s nanoMassa\n\n", massaAmount, nanoMassa.String())

	// Example 2: Convert nanoMassa to Massa
	fmt.Println("2. Converting nanoMassa to Massa:")
	nanoMassaAmount := big.NewInt(2_500_000_000)
	massa := utils.ToMAS(nanoMassaAmount)
	fmt.Printf("   %s nanoMassa = %f Massa\n\n", nanoMassaAmount.String(), massa)

	// Example 3: Round trip conversion
	fmt.Println("3. Round trip conversion:")
	original := 10.123456789
	converted := utils.FromMAS(original)
	backToMassa := utils.ToMAS(converted)
	fmt.Printf("   Original: %f Massa\n", original)
	fmt.Printf("   As nanoMassa: %s\n", converted.String())
	fmt.Printf("   Back to Massa: %f\n\n", backToMassa)

	// Example 4: Using FromNanoMAS
	fmt.Println("4. FromNanoMAS (uint64 to big.Int):")
	var nanoMassaUint uint64 = 1_000_000_000
	bigInt := utils.FromNanoMAS(nanoMassaUint)
	fmt.Printf("   %d nanoMassa (uint64) = %s (big.Int)\n\n", nanoMassaUint, bigInt.String())

	// Example 5: Using ToNanoMAS
	fmt.Println("5. ToNanoMAS (big.Int to uint64):")
	bigIntAmount := big.NewInt(500_000_000)
	uint64Amount := utils.ToNanoMAS(bigIntAmount)
	fmt.Printf("   %s (big.Int) = %d nanoMassa (uint64)\n\n", bigIntAmount.String(), uint64Amount)

	// Example 6: ParseMassa
	fmt.Println("6. ParseMassa (float64 to uint64):")
	massaFloat := 3.75
	parsed := utils.ParseMassa(massaFloat)
	fmt.Printf("   %f Massa = %d nanoMassa (uint64)\n\n", massaFloat, parsed)

	// Example 7: Common amounts
	fmt.Println("7. Common Massa amounts:")
	amounts := []float64{0.001, 1.0, 10.0, 100.0, 1000.0}
	for _, amount := range amounts {
		nano := utils.FromMAS(amount)
		fmt.Printf("   %8.3f Massa = %15s nanoMassa\n", amount, nano.String())
	}

	// Example 8: Precision with small amounts
	fmt.Println("\n8. Small amount precision:")
	smallAmount := 0.000000001 // 1 nanoMassa
	nanoSmall := utils.FromMAS(smallAmount)
	fmt.Printf("   %f Massa = %s nanoMassa\n", smallAmount, nanoSmall.String())

	// Example 9: Transaction example
	fmt.Println("\n9. Transaction scenario:")
	senderBalance := utils.FromMAS(100.0) // 100 Massa
	transferAmount := utils.FromMAS(25.5) // 25.5 Massa
	fee := utils.FromMAS(0.01)            // 0.01 Massa

	// Calculate remaining balance
	remaining := new(big.Int).Sub(senderBalance, transferAmount)
	remaining.Sub(remaining, fee)

	fmt.Printf("   Sender balance: %.2f Massa\n", utils.ToMAS(senderBalance))
	fmt.Printf("   Transfer amount: %.2f Massa\n", utils.ToMAS(transferAmount))
	fmt.Printf("   Fee: %.2f Massa\n", utils.ToMAS(fee))
	fmt.Printf("   Remaining: %.2f Massa (%s nanoMassa)\n",
		utils.ToMAS(remaining), remaining.String())
}
