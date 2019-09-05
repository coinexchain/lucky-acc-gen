package accgen

import (
	"fmt"
	"sync"
	"strings"
	"sync/atomic"
	"github.com/cosmos/cosmos-sdk/crypto/keys/hd"
	sdk "github.com/cosmos/cosmos-sdk/types"

	bip39 "github.com/cosmos/go-bip39"

	"github.com/tendermint/tendermint/crypto/secp256k1"
	"golang.org/x/crypto/blake2b"
)

type Result struct {
	found    bool
	addr     string
	mnemonic string
}

func TryAddressParallel(prefix string, numCpu int) (string, string) {
	var totalTry float64
	totalTry = 1.0
	n := len(prefix)-len("coinex1")
	for i:=0; i<n; i++ {
		totalTry *= 32.0
	}
	resPtr := &Result{}
	var globalCounter uint64
	var resAtomic atomic.Value
	resAtomic.Store(resPtr)
	var wg sync.WaitGroup
	wg.Add(numCpu)
	for i:=0; i<numCpu; i++ {
		go tryAddress(prefix, resAtomic, &wg, &globalCounter, totalTry)
	}
	wg.Wait()
	return resPtr.addr, resPtr.mnemonic
}

const BatchCount = 1000
const BigBatchCount = 10*BatchCount

func tryAddress(prefix string, resAtomic atomic.Value, wg *sync.WaitGroup, globalCounter *uint64, totalTry float64) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		panic(err.Error())
	}
	counter := 0
	for {
		if counter%BatchCount==0 {
			resPtr := resAtomic.Load().(*Result)
			if resPtr.found {
				break
			}
			count := atomic.AddUint64(globalCounter, BatchCount)
			if count%BigBatchCount==0 {
				percent := 100.0*float64(count)/totalTry
				fmt.Printf("%d times have been tried, estimated progress: %.2f%%\n", count, percent)
			}
		}
		addr, mnemonic, err := getAddress(entropy)
		if err != nil {
			panic(err.Error())
		}
		if strings.HasPrefix(addr, prefix) {
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

// NewFundraiserParams creates a BIP 44 parameter object from the params:
// m / 44' / coinType' / account' / 0 / address_index
// The fixed parameters (purpose', coin_type', and change) are determined by what was used in the fundraiser.
//func NewFundraiserParams(account, coinType, addressIdx uint32) *BIP44Params {
//	return NewParams(44, coinType, account, false, addressIdx)
//}

func getAddress(entropy []byte) (string,string,error) {
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", mnemonic, err
	}

	DefaultBIP39Passphrase := ""
	seed := bip39.NewSeed(mnemonic, DefaultBIP39Passphrase)
	fullHdPath := "44'/688'/0'/0/0" //coinType=688 account=0 addressIdx=0
	masterPriv, ch := hd.ComputeMastersFromSeed(seed)
	derivedPriv, err := hd.DerivePrivateKeyForPath(masterPriv, ch, fullHdPath)
	pubk := secp256k1.PrivKeySecp256k1(derivedPriv).PubKey()
	addr := pubk.Address()
	acc := sdk.AccAddress(addr)
	return acc.String(), mnemonic, nil
}

func getRandAddress() (string,string,error) {
	entropy, err := bip39.NewEntropy(32*4)
	if err != nil {
		return "", "", err
	}
	return getAddress(entropy)
}

