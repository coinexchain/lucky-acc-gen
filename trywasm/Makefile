.PHONY: all clean serve

all: main.wasm serve

%.wasm: %.go
	GOOS=js GOARCH=wasm go build -o "$@" "$<"

serve:
	goexec 'http.ListenAndServe(":8080", http.FileServer(http.Dir(".")))'
	

clean:
	rm -f *.wasm
