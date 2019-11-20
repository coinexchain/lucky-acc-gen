package main

import (
	"fmt"
	"runtime"

	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cloudfoundry/jibber_jabber"
	"github.com/qor/i18n"

	"github.com/coinexchain/lucky-acc-gen/accgen"
)

var mainwin *ui.Window

var prefixEntry *ui.Entry
var suffixEntry *ui.Entry
var multiEntry *ui.MultilineEntry
var runBtn *ui.Button

var LC string

var I18n = i18n.New()

func T(s string) string {
	return string(I18n.T(LC, s))
}

func add(key, en, cn string) {
	I18n.AddTranslation(&i18n.Translation{
		Key: key,
		Locale: "en-US",
		Value: en,
	})
	I18n.AddTranslation(&i18n.Translation{
		Key: key,
		Locale: "zh-CN",
		Value: cn,
	})
}

func init() {
	add("invalid_prefix", "Invalid Prefix!", "非法前缀！")
	add("invalid_char", "Invalid Character: %s\n", "非法字符：%s\n")
	add("invalid_suffix", "Invalid Suffix!", "非法后缀！")
	add("warn", "Warning", "警告")
	add("long_run_time", "You specified %d characters totally. It would take very long time to compute!", "您总共指定了%d个字符，这需要很长的时间才能生成一个靓号！")
	add("estimate_progress", "%d times have been tried, estimated progress: %.2f%%\n", "已经进行了%d次尝试，估计的完成度为%.2f%%")
	add("mnemonic", "Mnemonic:", "助记词：")
	add("address", "Address:", "地址：")
	add("title", "Generate Lucky Account Address", "靓号生成器")
	add("prefix", "Prefix:", "前缀：")
	add("suffix", "Suffix:", "后缀：")
	add("progress", "Progress:", "进展：")
	add("generate", "Generate!", "生成！")
	add("line1", "Please enter the prefix and suffix of your desired address below.",
	"请在下方输入您所期待的地址的前缀和后缀。")
	add("line2", "Prefix is the characters coming immediately after \"coinex1\".",
	"前缀是指紧接着\"coinex1\"出现的若干字符。")
	add("line3", "Suffix is the last few characters at the ending of an address.",
	"后缀是指地址末尾的若干字符。")
	add("line4", "Valid characters in a bech32 address are '023456789acdefghjklmnpqrstuvwxyz'.",
	"请注意在bech32格式的地址中，只有这些字符是合法的：'023456789acdefghjklmnpqrstuvwxyz'")
	add("line5", "Please note 'b', 'i', 'o' and '1' are not valid after \"coinex1\".",
	"也就是说，这几个字符不能出现在\"coinex1\"之后：b, i, o, 1")
}

func run() {
	if !runBtn.Enabled() {
		return
	}
	prefix := prefixEntry.Text()
	s, ok := accgen.CheckValid(prefix)
	if !ok {
		ui.MsgBoxError(mainwin, T("invalid_prefix"), fmt.Sprintf(T("invalid_char"), s))
		return
	}
	prefix = accgen.AddrPrefix + prefix

	suffix := suffixEntry.Text()
	s, ok = accgen.CheckValid(suffix)
	if !ok {
		ui.MsgBoxError(mainwin, T("invalid_suffix"), fmt.Sprintf(T("invalid_char"), s))
		return
	}

	if n := len(prefix + suffix); n > len(accgen.AddrPrefix)+7 {
		s := fmt.Sprintf(T("long_run_time"), n)
		ui.MsgBox(mainwin, T("warn"), s)
	}
	coreCount := runtime.NumCPU()
	runBtn.Disable()
	go func() {
		addr, mnemonic := accgen.TryAddressParallel(prefix, suffix, func(count uint64, percent float64) {
			ui.QueueMain(func() {
				s := fmt.Sprintf(T("estimate_progress"), count, percent)
				multiEntry.Append(s)
			})
		}, coreCount)
		ui.QueueMain(func() {
			multiEntry.Append(fmt.Sprintf("%s %s\n", T("mnemonic"), mnemonic))
			multiEntry.Append(fmt.Sprintf("%s %s\n", T("address"), addr))
			runBtn.Enable()
		})
	}()
}

func setupUI() {
	mainwin = ui.NewWindow(T("title"), 640, 480, true)
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
	vbox.Append(ui.NewLabel(T("line1")), false)
	vbox.Append(ui.NewLabel(T("line2")), false)
	vbox.Append(ui.NewLabel(T("line3")), false)
	vbox.Append(ui.NewLabel(T("line4")), false)
	vbox.Append(ui.NewLabel(T("line5")), false)

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	prefixEntry = ui.NewEntry()
	entryForm.Append(T("prefix"), prefixEntry, false)
	suffixEntry = ui.NewEntry()
	entryForm.Append(T("suffix"), suffixEntry, false)
	multiEntry = ui.NewMultilineEntry()
	entryForm.Append(T("progress"), multiEntry, true)
	multiEntry.SetReadOnly(true)
	vbox.Append(entryForm, false)
	runBtn = ui.NewButton(T("generate"))
	vbox.Append(runBtn, false)
	runBtn.OnClicked(func(_ *ui.Button) {
		run()
	})

	mainwin.SetChild(vbox)
	mainwin.SetMargined(true)

	mainwin.Show()
}

func main() {
	LC, _ = jibber_jabber.DetectIETF()
	sdk.GetConfig().SetBech32PrefixForAccount("coinex", "coinexpub")
	ui.Main(setupUI)
}

