# Formula parser

## Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/shellyln/takenoco/_examples/formula"
)

func main() {
    data, err := formula.Parse("1+2+3") // int64
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("%v", data)
}
```
