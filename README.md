# 🎋 Takenoco

A parser combinator library for Go.

<img src="https://raw.githubusercontent.com/shellyln/takenoco/master/_assets/logo-go-takenoco.svg" alt="logo" style="width:250px;" width="250">


---


## 👋 Examples

* [CSV parser](https://github.com/shellyln/takenoco/tree/master/_examples/csv)
* [Dust - toy scripting language](https://github.com/shellyln/dust-lang)
    * [Parsers](https://github.com/shellyln/dust-lang/tree/master/scripting/parser)
    * [Production rules](https://github.com/shellyln/dust-lang/tree/master/scripting/rules)


## 🚀 Getting started

### Define the parser:

```go
package csv

import (
    "errors"
    "strconv"

    . "github.com/shellyln/takenoco/base"
    . "github.com/shellyln/takenoco/string"
)

var (
    // Comma and line break characters
    cellBreakCharacters []string
    documentParser      ParserFn
)

func init() {
    cellBreakCharacters = make([]string, 0, len(LineBreakCharacters)+1)
    cellBreakCharacters = append(cellBreakCharacters, ",")
    cellBreakCharacters = append(cellBreakCharacters, LineBreakCharacters...)
    documentParser = document()
}

// Remove the resulting AST.
func erase(fn ParserFn) ParserFn {
    return Trans(fn, Erase)
}

// Whitespaces
func sp() ParserFn {
    return erase(ZeroOrMoreTimes(WhitespaceNoLineBreak()))
}

func quotedCell() ParserFn {
    return Trans(
        OneOrMoreTimes(
            FlatGroup(
                sp(),
                erase(Seq("\"")),
                ZeroOrMoreTimes(
                    First(
                        erase(Seq("\"\"")),
                        CharClassN("\""),
                    ),
                ),
                First(
                    erase(Seq("\"")),
                    FlatGroup(End(), Error("Unexpected EOF")),
                ),
                sp(),
            ),
        ),
        Concat,
    )
}

func cell() ParserFn {
    return Trans(
        ZeroOrMoreTimes(CharClassN(cellBreakCharacters...)),
        Trim,
    )
}

// Convert AST to array data. (line)
func lineTransform(_ ParserContext, asts AstSlice) (AstSlice, error) {
    w := make([]string, len(asts))
    length := len(asts)

    for i := 0; i < length; i++ {
        w[i] = asts[i].Value.(string)
    }

    return AstSlice{{
        ClassName: "*Line",
        Type:      AstType_Any,
        Value:     w,
    }}, nil
}

func line() ParserFn {
    return Trans(
        FlatGroup(
            ZeroOrMoreTimes(
                First(quotedCell(), cell()),
                erase(Seq(",")),
            ),
            First(quotedCell(), cell()),
        ),
        lineTransform,
    )
}

// Convert AST to array data. (Entire document)
func documentTransform(_ ParserContext, asts AstSlice) (AstSlice, error) {
    length := len(asts)
    w := make([][]string, length)

    for i := 0; i < length; i++ {
        w[i] = asts[i].Value.([]string)
    }
    for i := length - 1; i >= 0; i-- {
        if len(w[i]) == 0 || len(w[i]) == 1 && w[i][0] == "" {
            w = w[:i]
        } else {
            break
        }
    }

    return AstSlice{{
        ClassName: "*Document",
        Type:      AstType_Any,
        Value:     w,
    }}, nil
}

func document() ParserFn {
    return Trans(
        FlatGroup(
            ZeroOrMoreTimes(
                line(),
                erase(OneOrMoreTimes(LineBreak())),
            ),
            line(),
            End(),
        ),
        documentTransform,
    )
}

func Parse(s string) ([][]string, error) {
    out, err := documentParser(*NewStringParserContext(s))
    if err != nil {
        return nil, err
    } else {
        if out.MatchStatus == MatchStatus_Matched {
            return out.AstStack[0].Value.([][]string), nil
        } else {
            return nil, errors.New("Parse failed at " + strconv.Itoa(out.SourcePosition.Position))
        }
    }
}
```

### Use the parser:

```go
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
```


## 📦 Build example

### Build to native executable

```bash
make
```

### Build to WebAssembly
#### Windows prerequirements:

```bash
scoop install tinygo
scoop install binaryen
```

#### Build:

```bash
make wasm
```


## ⚖️ License

MIT  
Copyright (c) 2021 Shellyl_N and Authors.
