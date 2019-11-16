package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/lucky-acc-gen/accgen"
)

const addrPrefix = "coinex1"

var bech32Chars map[rune]bool

func init() {
	bech32Chars = make(map[rune]bool)
	for _, c := range "023456789acdefghjklmnpqrstuvwxyz" {
		bech32Chars[c] = true
	}
}

func main() {
	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")

	coreCount := askCoreCount()
	prefix := askPrefix()
	suffix := askSuffix()

	if n := len(prefix + suffix); n > len(addrPrefix)+7 {
		fmt.Printf("\nWARNING! you specified %d characters totally. It would take very long time to compute!\n", n)
	}

	addr, mnemonic := accgen.TryAddressParallel(prefix, suffix, func(s string) {
		fmt.Print(s)
	}, coreCount)
	fmt.Printf("Mnemonic: %s\n", mnemonic)
	fmt.Printf("Addr: %s\n", addr)
	fmt.Print("Press Enter to Exit")
	var input string
	fmt.Scanln(&input)
}

func askCoreCount() int {
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
	return coreCount
}

func askPrefix() string {
	prefix := addrPrefix

	for {
		fmt.Print("Please enter several characters after \"coinex1\": ")
		var input string
		fmt.Scanln(&input)
		input = strings.TrimSpace(input)
		if len(input) == 0 {
			fmt.Printf("Nothting is entered!\n")
			//continue
		}
		if isValid(input) {
			prefix = prefix + input
			break
		}
	}
	return prefix
}

func askSuffix() string {
	suffix := ""

	for {
		fmt.Print("Please enter address postfix: ")
		var input string
		fmt.Scanln(&input)
		input = strings.TrimSpace(input)
		if len(input) == 0 {
			fmt.Printf("Nothting is entered!\n")
			//continue
		}
		if isValid(input) {
			suffix = input
			break
		}
	}
	return suffix
}

func isValid(str string) bool {
	for _, c := range str {
		if _, ok := bech32Chars[c]; !ok {
			fmt.Printf("Invalid character: %s\n", strconv.QuoteRune(c))
			return false
		}
	}
	return true
}
