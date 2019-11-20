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

	coreCount := askCoreCount()
	prefix := askPrefix()
	suffix := askSuffix()

	if n := len(prefix + suffix); n > len(accgen.AddrPrefix)+7 {
		fmt.Printf("\nWARNING! you specified %d characters totally. It would take very long time to compute!\n", n)
	}

	addr, mnemonic := accgen.TryAddressParallel(prefix, suffix, func(count uint64, percent float64) {
		fmt.Printf("%d times have been tried, estimated progress: %.2f%%\n", count, percent)
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
	prefix := accgen.AddrPrefix

	for {
		fmt.Print("Please enter several characters after \"coinex1\": ")
		var input string
		fmt.Scanln(&input)
		input = strings.TrimSpace(input)
		if len(input) == 0 {
			fmt.Printf("Nothting is entered!\n")
			//continue
		}
		s, ok := accgen.CheckValid(input)
		if ok {
			prefix = prefix + input
			break
		} else {
			fmt.Sprintf("Invalid character: %s\n", s)
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
		s, ok := accgen.CheckValid(input)
		if ok {
			suffix = input
			break
		} else {
			fmt.Println(s)
		}
	}
	return suffix
}

