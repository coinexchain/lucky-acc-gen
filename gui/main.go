package main

import (
	"fmt"
	"runtime"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/lucky-acc-gen/accgen"
)

var mainwin *ui.Window

var prefixEntry *ui.Entry
var suffixEntry *ui.Entry
var multiEntry *ui.MultilineEntry
var runBtn *ui.Button

func run() {
	if !runBtn.Enabled() {
		return
	}
	prefix := prefixEntry.Text()
	s, ok := accgen.CheckValid(prefix)
	if !ok {
		ui.MsgBoxError(mainwin, "Invalid Prefix!", s)
		return
	}
	prefix = accgen.AddrPrefix + prefix

	suffix := suffixEntry.Text()
	s, ok = accgen.CheckValid(suffix)
	if !ok {
		ui.MsgBoxError(mainwin, "Invalid Suffix!", s)
		return
	}

	if n := len(prefix + suffix); n > len(accgen.AddrPrefix)+7 {
		s := fmt.Sprintf("You specified %d characters totally. It would take very long time to compute!", n)
		ui.MsgBox(mainwin, "Warning!", s)
	}
	coreCount := runtime.NumCPU()
	runBtn.Disable()
	go func() {
		addr, mnemonic := accgen.TryAddressParallel(prefix, suffix, func(s string) {
			ui.QueueMain(func() {
				multiEntry.Append(s)
			})
		}, coreCount)
		ui.QueueMain(func() {
			multiEntry.Append(fmt.Sprintf("Mnemonic: %s\n", mnemonic))
			multiEntry.Append(fmt.Sprintf("Addr: %s\n", addr))
			runBtn.Enable()
		})
	}()
}

func setupUI() {
	mainwin = ui.NewWindow("Generate Lucky Account Address", 640, 480, true)
	mainwin.OnClosing(func(*ui.Window) bool {
		ui.Quit()
		return true
	})
	ui.OnShouldQuit(func() bool {
		mainwin.Destroy()
		return true
	})

	vbox := ui.NewVerticalBox()
	vbox.SetPadded(true)
	vbox.Append(ui.NewLabel("Please enter the prefix and suffix of your desired address below."), false)
	vbox.Append(ui.NewLabel("Prefix is the characters coming immediately after \"coinex1\"."), false)
	vbox.Append(ui.NewLabel("Suffix is the last few characters at the ending of an address."), false)
	vbox.Append(ui.NewLabel("Valid characters in a bech32 address are '023456789acdefghjklmnpqrstuvwxyz'."), false)
	vbox.Append(ui.NewLabel("Please note 'b', 'i', 'o' and '1' are not valid after \"coinex1\"."), false)

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	prefixEntry = ui.NewEntry()
	entryForm.Append(`Prefix:`, prefixEntry, false)
	suffixEntry = ui.NewEntry()
	entryForm.Append(`Suffix:`, suffixEntry, false)
	multiEntry = ui.NewMultilineEntry()
	entryForm.Append("Progress:", multiEntry, true)
	multiEntry.SetReadOnly(true)
	vbox.Append(entryForm, false)
	runBtn = ui.NewButton("Generate!")
	vbox.Append(runBtn, false)
	runBtn.OnClicked(func(_ *ui.Button) {
		run()
	})

	mainwin.SetChild(vbox)
	mainwin.SetMargined(true)

	mainwin.Show()
}

func main() {
	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	ui.Main(setupUI)
}

