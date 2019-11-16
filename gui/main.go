package main

import (
	"github.com/andlabs/ui"
	_ "github.com/andlabs/ui/winmanifest"
)

var mainwin *ui.Window

func setupUI() {
	mainwin = ui.NewWindow("libui Control Gallery", 640, 480, true)
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
	vbox.Append(ui.NewLabel("This is a label. Right now, labels can only span one line."), false)

	entryForm := ui.NewForm()
	entryForm.SetPadded(true)
	entryForm.Append("Entry", ui.NewEntry(), false)
	entryForm.Append("Password Entry", ui.NewPasswordEntry(), false)
	multiEntry := ui.NewMultilineEntry()
	entryForm.Append("Multiline Entry", multiEntry, true)
	multiEntry.SetReadOnly(true)
	vbox.Append(entryForm, false)
	vbox.Append(ui.NewButton("Button"), false)

	mainwin.SetChild(vbox)
	mainwin.SetMargined(true)

	mainwin.Show()
}

// func MsgBoxError(w *Window, title string, description string) 

func main() {
	ui.Main(setupUI)
}

