# Formula to RPN converter

## Usage

```go
package main

import (
    "fmt"
    "log"
    "github.com/shellyln/takenoco/_examples/torpn"
)

func main() {
    data, err := torpn.Parse("1+2+3") // []any
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(data)
}
```
