// +build js,wasm

package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	js.Global().Get("document").Set("wasmFunc", js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		js.Global().Get("window").Call("fromWasm", fmt.Sprintf("insidewasm-%s", args[0].String()))
		return nil
	}))

	select {}
}
