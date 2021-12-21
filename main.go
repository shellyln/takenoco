package main

import (
	"fmt"
	"os"

	csv "github.com/shellyln/takenoco/_examples/csv"
)

func main() {
	x, err := csv.Parse("0,1,2,3,4,5,6,7,8,9\n0,1,2,3,4,5,6,7,8,9")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(-1)
	}
	fmt.Println(x)

	y := csv.ToCsv(x)
	fmt.Println(y)

	os.Exit(0)
}
