package accgen

import (
	"strconv"
	"strings"
	"sync"
	"sync/atomic"

	"golang.org/x/crypto/blake2b"

	bip39 "github.com/cosmos/go-bip39"

	"github.com/coinexchain/polarbear/keybase"
)

func CheckValid(str string) (string, bool) {
	for _, c := range str {
		if _, ok := bech32Chars[c]; !ok {
			return strconv.QuoteRune(c), false
		}
	}
	return "", true
}

const AddrPrefix = "coinex1"

var bech32Chars map[rune]bool

func init() {
	bech32Chars = make(map[rune]bool)
	for _, c := range "023456789acdefghjklmnpqrstuvwxyz" {
		bech32Chars[c] = true
	}
}

type Result struct {
	found    bool
	addr     string
	mnemonic string
}

func TryAddressParallel(prefix, suffix string, repFn func(uint64, float64), numCpu int) (string, string) {
	var totalTry float64
	totalTry = 1.0
	n := len(prefix+suffix) - len("coinex1")
	for i := 0; i < n; i++ {
		totalTry *= 32.0
	}
	resPtr := &Result{}
	var globalCounter uint64
	var resAtomic atomic.Value
	resAtomic.Store(resPtr)
	var wg sync.WaitGroup
	wg.Add(numCpu)
	for i := 0; i < numCpu; i++ {
		go tryAddress(prefix, suffix, repFn, resAtomic, &wg, &globalCounter, totalTry)
	}
	wg.Wait()
	return resPtr.addr, resPtr.mnemonic
}

const BatchCount = 1000
const BigBatchCount = 10 * BatchCount

func tryAddress(prefix, suffix string, repFn func(uint64, float64),
	resAtomic atomic.Value, wg *sync.WaitGroup, globalCounter *uint64, totalTry float64) {

	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		panic(err.Error())
	}
	counter := 0
	for {
		if counter%BatchCount == 0 {
			resPtr := resAtomic.Load().(*Result)
			if resPtr.found {
				break
			}
			count := atomic.AddUint64(globalCounter, BatchCount)
			if count%BigBatchCount == 0 {
				percent := 100.0 * float64(count) / totalTry
				repFn(count, percent)
			}
		}
		addr, mnemonic, err := keybase.GetAddressFromEntropy(entropy)
		if err != nil {
			panic(err.Error())
		}
		if strings.HasPrefix(addr, prefix) && strings.HasSuffix(addr, suffix) {
			resPtr := resAtomic.Load().(*Result)
			resPtr.found = true
			resPtr.addr = addr
			resPtr.mnemonic = mnemonic
			resAtomic.Store(resPtr)
			break
		}
		sum := blake2b.Sum256(entropy)
		entropy = sum[:]
		counter++
	}
	wg.Done()
}
