//go:build wasm
// +build wasm

package main

import (
	"syscall/js"

	"github.com/shellyln/takenoco/_examples/csv"
	"github.com/shellyln/takenoco/_examples/formula"
)

func parseCsv(this js.Value, args []js.Value) interface{} {
	src := ""
	if 0 < len(args) {
		src = args[0].String()
	}

	data, err := csv.Parse(src)
	if err != nil {
		println(err)
	}

	// NOTE: `js.ValueOf()` can't accept `[][]string`.

	rows := make([]interface{}, len(data))

	for i := 0; i < len(data); i++ {
		row := data[i]
		cells := make([]interface{}, len(row))
		for j := 0; j < len(row); j++ {
			cells[j] = row[j]
		}
		rows[i] = cells
	}

	return rows
}

func parseFormula(this js.Value, args []js.Value) interface{} {
	src := ""
	if 0 < len(args) {
		src = args[0].String()
	}

	data, err := formula.Parse(src)
	if err != nil {
		println(err)
	}

	return data
}

//export parseCsvWithTinyGo
func parseCsvWithTinyGo(src string) [][]string {
	// https://tinygo.org/docs/guides/webassembly/
	// NOTE: BUG: At this time, non-numeric type parameters are not accepted.
	data, err := csv.Parse(src)
	if err != nil {
		println(err)
	}
	// NOTE: BUG: At this time, non-numeric type parameters are not accepted.
	return data
}

func main() {
	ch := make(chan struct{}, 0)
	println("Go WebAssembly Initialized")

	// for golang/go
	js.Global().Set("parseCsv", js.FuncOf(parseCsv))

	// for golang/go
	js.Global().Set("parseFormula", js.FuncOf(parseFormula))

	// for golang/go
	<-ch
}
