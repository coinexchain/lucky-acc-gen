package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/lucky-acc-gen/accgen"
)

func main() {
	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")

	coreCount := runtime.NumCPU()
	for {
		fmt.Printf("Please enter the number of cpu cores you want to use (you have %d cores, press enter to use all the cores): ", coreCount)
		var input string
		fmt.Scanln(&input)
		input = strings.TrimSpace(input)
		if len(input) == 0 {
			break
		}
		if n, err := strconv.ParseUint(input, 10, 32); err == nil {
			if int(n) < coreCount {
				coreCount = int(n)
			}
			break
		} else {
			fmt.Printf("Invalid input. Please enter a digit.\n")
		}
	}

	prefix := "coinex1"
	validChars := make(map[rune]bool)
	for _, c := range "023456789acdefghjklmnpqrstuvwxyz" {
		validChars[c] = true
	}
	for {
		fmt.Print("Please the enter several characters after \"coinex1\": ")
		var input string
		fmt.Scanln(&input)
		input = strings.TrimSpace(input)
		if len(input) == 0 {
			fmt.Printf("Nothting is entered!\n")
			continue
		}
		isValid := true
		for _, c := range input {
			if _, ok := validChars[c]; !ok {
				fmt.Printf("Invalid character: %s\n", strconv.QuoteRune(c))
				isValid = false
				break
			}
		}
		if isValid {
			if len(input) > 7 {
				fmt.Printf("\nWARNING! you specified %d characters. It would take very long time to compute!\n")
			}
			prefix = prefix + input
			break
		}
	}

	addr, mnemonic := accgen.TryAddressParallel(prefix, coreCount)
	fmt.Printf("Mnemonic: %s\n", mnemonic)
	fmt.Printf("Addr: %s\n", addr)
	fmt.Print("Press Enter to Exit")
	var input string
	fmt.Scanln(&input)
}
