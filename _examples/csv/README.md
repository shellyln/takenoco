# CSV parser

## Usage

### Parse

```go
package main

import (
    "fmt"
    "log"
    "github.com/shellyln/takenoco/_examples/csv"
)

func main() {
    data, err := csv.Parse("a,b,c\n1,2,3") // [][]string
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println("%v", data)
}
```

### Unmarshal

```go
package main

import (
    "fmt"
    "log"
    "github.com/shellyln/takenoco/_examples/csv"
)

type Foobar struct {
    XFoo string `csv:"Foo"`
    Bar  string
    XBaz int `csv:"Baz"`
}

func main() {
    data := []Foobar{}

    err := csv.Unmarshal(&data, "Foo,Bar,Baz\n1,2,3") // [][]string
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(data)
}
```

### Marshal

```go
package main

import (
    "fmt"
    "log"
    "github.com/shellyln/takenoco/_examples/csv"
)

type Foobar struct {
    XFoo string `csv:"Foo"`
    Bar  string
    XBaz int `csv:"Baz"`
}

func main() {
    data := []Foobar{
        {
            XFoo: "foo0",
            Bar:  "bar0",
            XBaz: 0,
        },
        {
            XFoo: "foo1",
            Bar:  "bar1",
            XBaz: 1,
        },
    }

    csv, err := csv.Marshal([]string{"Foo", "Bar", "Baz"}, data) // string
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(csv)
}
```

### Convert to CSV

```go
package main

import (
    "fmt"
    "log"
    "github.com/shellyln/takenoco/_examples/csv"
)

func main() {
    data := [][]string{
        {"a", "b", "c"},
    }

    fmt.Println(csv.ToCsv(data))
}
```
