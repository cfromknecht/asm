package main

import (
	"fmt"
	"github.com/cfromknecht/asm"
)

func main() {
	numValues := (1 << 16) - 1
	acc := asm.NewAsyncAcc()
	witnesses := []asm.WitnessPath{}

	fmt.Print("[ADDING VALUES] |")
	for i := 0; i < numValues; i++ {
		wit := acc.Add(fmt.Sprintf("%d", i))
		witnesses = append(witnesses, wit)

		// Print progress bar
		divisor := numValues / 100
		if i%(10*divisor) == 0 && i != 0 {
			if i/divisor != 100 {
				fmt.Print(i / divisor)
			}
		} else if i%divisor == 0 {
			fmt.Print("=")
		}
	}

	for i, wit := range witnesses {
		fmt.Println("witness", i, wit)
	}

	fmt.Println("|\naccumulator for", numValues, "values:", acc)
}
