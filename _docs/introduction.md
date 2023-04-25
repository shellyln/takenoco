# Introduction to takenoco, a parser combinator library for Go

This document describes [takenoco](https://github.com/shellyln/takenoco), a parser combinator library for the Go language.

Parser combinator is a parser construction method that combines functions of small parsers to build the desired parser.
It can be implemented rather easily by the "higher-order function" mechanism of programming languages.

[takenoco](https://github.com/shellyln/takenoco) provides a foundation for parsing strings or slices of any type
and converting them into ASTs (abstract syntax trees), a set of general-purpose small parsers,
and AST conversion functions with production rules.

# Packages in takenoco
[takenoco](https://github.com/shellyln/takenoco) is provided by the following packages.

```
github.com/shellyln/takenoco/
├── base/
├── string/
├── object/
└── extra/
```
* `base/`:  
  Provides types, parser foundation, and string/arbitrary type slices' common parsers.
* `string/`:  
  Provides common parsers for string.
* `object/`:  
  Provides common parsers for arbitrary type slices'.
* `extra/`:  
  Provides additional parsers.


# Minimal parser
First, let's make the minimal parser to make a sense of [takenoco](https://github.com/shellyln/takenoco).
```go
package main

import (
    "fmt"
    "log"
    "play.ground/myparser"
)

func main() {
    data, err := myparser.Parse("foobar") // 1) Passing the text to be parsed
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(data)
}

-- go.mod --
module play.ground

-- myparser/myparser.go --
package myparser

import (
    "errors"
    "strconv"
    . "github.com/shellyln/takenoco/base"    // 2) takenoco's common functionality package
    . "github.com/shellyln/takenoco/string"  // 3) takenoco's string parser package
)

// Entire program or document to be parsed
func program() ParserFn {
    return FlatGroup( // 4) Group the parsers. However, it does not group the parsed AST structs.
                      //    (It returns multiple structs flatly)
                      //    In this case, they are grouped because a single ParserFn needs to be returned.
        Start(), // 5) A zero-width assertion indicating the start of the document.
        Trans(   // 6) Parses with the parser, and passes the parsed result to the transform function.
            OneOrMoreTimes(Alpha()), // 7) Repeats an ASCII alphabetic character one or more times.
                                     //    (One character - one AST struct)
            Concat, // 8) Combines the resulting parsed AST structs as a string
                    //    and converts them into a single AST struct.
        ),
        End(), // 9) A zero-width assertion indicating the end of the document.
    )
}

// Initialize the package
var rootParser ParserFn
func init() {
    rootParser = program() // 10) Build the parser
}

// Parser
func Parse(s string) (string, error) {
    out, err := rootParser(*NewStringParserContext(s)) // 11) Execute parsing.
    if err != nil {
        // 12) Format the error location information.
        pos := GetLineAndColPosition(s, out.SourcePosition, 4)
        return "", errors.New(
            err.Error() +
                "\n --> Line " + strconv.Itoa(pos.Line) +
                ", Col " + strconv.Itoa(pos.Col) + "\n" +
                pos.ErrSource)
    }

    // 13) If MatchStatus_Matched, parsing is successful.
    if out.MatchStatus == MatchStatus_Matched {
        return out.AstStack[0].Value.(string), nil
    } else {
        pos := GetLineAndColPosition(s, out.SourcePosition, 4)
        return "", errors.New(
            "Parse failed" +
                "\n --> Line " + strconv.Itoa(pos.Line) +
                ", Col " + strconv.Itoa(pos.Col) + "\n" +
                pos.ErrSource)
    }
}
```
Let's look at the `func program()` in `myparser.go`.  
The `Start()` and `End()` represent the start and end of the document,
and the parser `Alpha()`, which matches alphabetic character, is wrapped between them
with the parser `OneOrMoreTimes()`, which matches one or more repetitions.

It looks somewhat like a regular expression, don't you think?
If it were expressed in a regular expression, it would be `/(?:^[A-Za-z]+$)/`.

> **Note**  
> To be precise, the meaning of the regular expression non-capturing group (`(?:)`) and
> takenoco's `FlatGroup()` have different meanings, but they are only examples for contrast.

**Practice**
* Run it in [The Go Playground](https://go.dev/play/).
  The entire code above can be pasted into the playground unchanged.
* What happens if `"foobar"` in `1)` is replaced by `"foobar1"`?


# Introduction of common parsers
[takenoco](https://github.com/shellyln/takenoco) provides a set of frequently used basic and common parsers in advance.
In most cases, these can be combined to build the desired parser.

## github.com/shellyln/takenoco/base
```go
func Indirect(fn func() ParserFn) 
```
Used to parse recursive expressions. If the parser call is circular,
it will cause a stack overflow during parser construction and a run-time error.

To prevent this, parsers with circular references are wrapped in `Indirect`.
Note that the function being wrapped cannot take parameters at construction time.

```go
func If(b bool, fnT ParserFn, fnF ParserFn) ParserFn
```
Conditional branching at parser construction time. Note that this is not at parser runtime.

```go
func Error(msg string) ParserFn
```
If called, parsing fails. If a parse error occurs, it will not be backtracked.
Inserted as the last item in the argument to `First()` (see below), it can be used to indicate to the user
where the error occurred and the error message when a syntactically invalid token appears.

Inserting it where there is grammatical room for backtracking will cause an error even though it should be parsable.
Please consider carefully where to insert it.

```go
func Unmatched() ParserFn
```
When called, it is unmatched. Backtracking occurs, so the parse as a whole may still succeed.

```go
func Zero(astsToInsert ...Ast) ParserFn
```
A zero-width (not advance the source position) assertion. Matches whenever called.
Appends the argument AST struct(s) to the parse result.

```go
func Start() ParserFn
```
A zero-width assertion that matches the beginning of the source to be parsed.

```go
func FlatGroup(children ...ParserFn) ParserFn
```
Groups the argument parsers into a single parser.
The resulting AST struct(s) output by each parser are flatly added to the parse result.

```go
func Group(children ...ParserFn) ParserFn
```
Groups the argument parsers into a single parser. 
The resulting AST struct(s) output by each parser are added to the parsed result as one AST struct.
The type of `Value` of the AST will be a slice of AST structs.

```go
func First(children ...ParserFn) ParserFn
```
Matches the first match among the parsers of the arguments.

```go
func LookAhead(children ...ParserFn) ParserFn
```
A zero-width (not advance the source position) "look-ahead" assertion.
Matches if all of the parsers in the argument match in order, and then restores the source position.

```go
func LookAheadN(children ...ParserFn) ParserFn
```
A zero-width (not advance the source position) "negative look-ahead" assertion.
Matches if some of the parsers in the argument don't match in order, and then restores the source position.

```go
func LookBehind(minN, maxN int, children ...ParserFn) ParserFn
```
A zero-width (not advance the source position) "look-behind" assertion.
Matches if it matches all of the argument parsers in order, and then restores the source position.

Matching is tried from the current position, moving backward by `minN`.
If it does not match, it tries from the position +1 forward up to a maximum of `maxN`.

It does not check if the end position of the match matches the original position.

```go
func LookBehindN(minN, maxN int, children ...ParserFn) ParserFn
```
A zero-width (not advance the source position) "negative look-behind" assertion.
Matches if some of the parsers in the argument don't match in order, and then restores the source position.

Matching is tried from the current position, moving backward by `minN`.
If it does not match, it tries from the position +1 forward up to a maximum of `maxN`.

It does not check if the end position of the match matches the original position.

```go
func Repeat(times Times, children ...ParserFn) ParserFn

type Times struct {
    Min int
    Max int
}
```
Matches all argument parsers in order at least `Min` and up to `Max` times.
If more than `Max` matches are possible, it stops at `Max`.  
If `Min` and `Max` are negative, the number of matches is not limited.

```go
func Once(children ...ParserFn) ParserFn
```
Matches all argument parsers in order only once.  
Equivalent to not wrapping with `Once()`, but provided to harmonize the vocabulary
with the following `ZeroOrOnce`, `ZeroOrMoreTimes`, `ZeroOrMoreTimes`, `OneOrMoreTimes`.

```go
func ZeroOrOnce(children ...ParserFn) ParserFn
```
Matches all argument parsers in order at least 0 and up to once.  
Equivalent to `Repeat(Times{Min: 0, Max: 1}, ...)`.

```go
func ZeroOrMoreTimes(children ...ParserFn) ParserFn
```
Matches all argument parsers in order at least 0 or more times.  
Equivalent to `Repeat(Times{Min: 0, Max: -1}, ...)`.

```go
func OneOrMoreTimes(children ...ParserFn) ParserFn
```
Matches all argument parsers in order at least once or more times.  
Equivalent to `Repeat(Times{Min: 1, Max: -1}, ...)`.

```go
func Trans(child ParserFn, tr ...TransformerFn) ParserFn
```
If matched with the `child` argument, and the resulting AST struct slice is rewritten with the transform function.

> **Warning**  
> The transform function **MUST NOT change** the value of the AST struct slice given in the argument.  
> If you make changes, you **SHOULD copy** them to another slice.  
> It is ok to return the AST struct slice given in the argument as is.

## github.com/shellyln/takenoco/string

```go
func Any() ParserFn
```
Matches any single character.

```go
func End() ParserFn
```
A zero-width assertion that matches the end of the source to be parsed.

```go
func Seq(s string) ParserFn
```
Matches the argument string.

```go
func SeqI(s string) ParserFn
```
Matches the argument string. Ignore case.

```go
func CharRange(cr ...RuneRange) ParserFn

type RuneRange struct {
    Start rune
    End rune
}
```
Matches if it meets any of the character code ranges in the argument.

```go
func CharRangeN(cr ...RuneRange) ParserFn
```
Matches if it does not meet any of the character code ranges in the argument.

```go
func CharClass(cc ...string) ParserFn
```
Matches if it meets any of the argument strings.

```go
func CharClassN(cc ...string) ParserFn
```
Matches if it does not meet any of the argument strings.

```go
func CharClassFn(fn func(c rune) bool) ParserFn
```
Matches if the function determines the character to be met.

```go
func Whitespace() ParserFn
```
Matches any of `HT`, `LF`, `VT`, `FF`, `CR`, `SP`, `NEL`, `NBSP` characters.

```go
func WhitespaceNoLineBreak() ParserFn
```
Matches any of `HT`, `SP`, `NBSP` characters.

```go
func LineBreak() ParserFn
```
Matches any of `LF`, `VT`, `FF`, `CR`, `NEL` characters.

```go
func Alpha() ParserFn
```
Matches any of `/[A-Za-z]/` characters.

```go
func Number() ParserFn
```
Matches any of `/[0-9]/` characters.

```go
func Alnum() ParserFn
```
Matches any of `/[0-9A-Za-z]/` characters.

```go
func BinNumber() ParserFn
```
Matches any of `/[0-1]/` characters.

```go
func OctNumber() ParserFn
```
Matches any of `/[0-7]/` characters.

```go
func HexNumber() ParserFn
```
Matches any of `/[0-9A-Fa-f]/` characters.

```go
func WordBoundary() ParserFn 
```
A zero-width (not advance the source position) assertion.
Matches between a word construct character and a nonword construct character or the beginning or end of a source,
`/[0-9A-Za-z_$]/` as a word construct character.


## github.com/shellyln/takenoco/object

```go
func Any() ParserFn
```
Matches any single value.

```go
func End() ParserFn
```
A zero-width assertion that matches the end of the source to be parsed.

```go
func Seq(seq ...interface{}) ParserFn
```
Matches the argument values.

```go
func ObjClass(oc ...interface{}) ParserFn
```
Matches if it meets any of the argument values.

```go
func ObjClassN(oc ...interface{}) ParserFn
```
Matches if it does not meet any of the argument values.

```go
func ObjClassFn(fn func(c interface{}) bool) ParserFn
```
Matches if the function determines the value to be met.


## github.com/shellyln/takenoco/extra

```go
func BinaryNumberStr() ParserFn
```
Matches a binary number string.
Prefixes (such as `0b`) are not included. May contain `_` to delimit digits.

```go
func OctalNumberStr() ParserFn
```
Matches an octal number string.
Prefixes (such as `0o`) are not included. May contain `_` to delimit digits.

```go
func HexNumberStr() ParserFn
```
Matches a hexadecimal number string.
Prefixes (such as `0x`) are not included. May contain `_` to delimit digits.

```go
func IntegerNumberStr() ParserFn
```
Matches a string representing positive, negative, and zero integer value.
The `+` and `-` signs are not required. May contain `_` to delimit digits.

```go
func FloatNumberStr() ParserFn
```
Matches a string representing a number with a decimal point
or a number with a decimal exponent expression.
May contain `_` to delimit digits.

```go
func NumericStr() ParserFn
```
Matches a string representing an integer, a number with a decimal point
or a number with a decimal exponent expression.
May contain `_` to delimit digits.

```go
func AsciiIdentifierStr() ParserFn
```
Matches a string representing an identifier in the ASCII range, where the first character is `/[A-Za-z_$]/`
and the second and later characters are `/[0-9A-Za-z_$]/`.

```go
func UnicodeIdentifierStr() ParserFn
```
Matches a string of Unicode identifiers, where the first character is {`ID_Start`, `_`, `$`}
and the second and later characters are {`ID_Continue`, `$`, `U+200C`, `U+200D`}.

```go
func UnicodeWordBoundary() ParserFn
```
A zero-width (not advance the source position) assertion.
Matches between a word construct character and a nonword construct character or the beginning or end of a source,
{`ID_Continue`, `$`, `U+200C`, `U+200D`} as a word construct character.

```go
func DateStr() ParserFn
```
Matches an ISO 8601 extended date format string.
* `yyyy-MM-dd`

```go
func DateTimeStr() ParserFn
```
Matches an ISO 8601 extended date/time format string.
* `yyyy-MM-ddThh:mmZ` to `yyyy-MM-ddThh:mm:ss.fffffffffZ`
* `yyyy-MM-ddThh:mm+00:00` to `yyyy-MM-ddThh:mm:ss.fffffffff+00:00`

```go
func TimeStr() ParserFn
```
Matches an ISO 8601 extended time format string.
* `hh:mm` to `hh:mm:ss.fffffffff`


# Introduction of common transform functions
Common transform functions for use with `Trans()` are also provided in advance.

## github.com/shellyln/takenoco/base

```go
func Erase(_ ParserContext, asts AstSlice) (AstSlice, error)
```
Discard the parsed results.

```go
func TransformError(s string) TransformerFn
```
Cause an error.

```go
func GroupingTransform(_ ParserContext, asts AstSlice) (AstSlice, error)
```
The parsed result is grouped into a single AST struct.
The `Value` type is a slice of the AST struct.

```go
func ToSlice(ctx ParserContext, asts AstSlice) (AstSlice, error)
```
For each AST structs in the parsed result, get the `Value` of each and make it into a slice.
The type will be the same as the source passed in `NewObjectParserContext(slice)`.  
Does not work with string parsers.

```go
func SetOpCodeAndClassName(opcode AstOpCodeType, name string) TransformerFn
```
Rewrite the 0th `OpCode` and `ClassName` of the resulting AST struct slice.

```go
func SetOpCode(opcode AstOpCodeType) TransformerFn
```
Rewrite the 0th `OpCode` of the resulting AST struct slice.

```go
func ChangeClassName(name string) TransformerFn
```
Rewrite the 0th `ClassName` of the resulting AST struct slice.

```go
func SetValue(typ AstType, v interface{}) TransformerFn
```
Rewrite the 0th `Type` and `Value` of the resulting AST struct slice.

```go
func Prepend(ast Ast) TransformerFn
```
Append the argument AST struct to the beginning of the resulting AST struct slice.

```go
func Push(ast Ast) TransformerFn
```
Append the argument AST struct to the end of the resulting AST struct slice.

```go
func Pop(_ ParserContext, asts AstSlice) (AstSlice, error)
```
Remove one AST struct from the end of the resulting AST struct slice
(shorten the length of the slice by 1).

```go
func Exchange(_ ParserContext, asts AstSlice) (AstSlice, error)
```
Exchange `asts[len(asts)-2]` and `asts[len(asts)-1]` in the resulting AST struct slice.

```go
func Roll(n int) TransformerFn
```
Shift the resulting AST struct slice by n.
The overflowed struct(s) comes in from the other end.

* examples
    * n == 2
        ```
        asts bottom |0  |1  |2  |3  |...|n-2|n-1| top
                    <=  <=
                    |2  |3  |...|n-2|n-1|0  |1  |
        ```
    * n == -2
        ```
        asts bottom |0  |1  |...|n-4|n-3|n-2|n-1| top
                                        =>  =>
                    |n-2|n-1|0  |1  |...|n-4|n-3|
        ```

## github.com/shellyln/takenoco/string

```go
func Concat(_ ParserContext, asts AstSlice) (AstSlice, error)
```
For each resulting AST struct(s) in the slice, get its `Value` and concatenate the strings.  
If the `Value` are not string, an error occurs.

```go
func Trim(ctx ParserContext, asts AstSlice) (AstSlice, error)
```
For each resulting AST struct(s) in the slice, get its `Value` and concatenate the strings,
and Remove consecutive whitespaces from the beginning and end of the string.  
If the `Value` are not string, an error occurs.

```go
func TrimStart(ctx ParserContext, asts AstSlice) (AstSlice, error)
```
For each resulting AST struct(s) in the slice, get its `Value` and concatenate the strings,
and Remove consecutive whitespaces from the beginning of the string.  
If the `Value` are not string, an error occurs.

```go
func TrimEnd(ctx ParserContext, asts AstSlice) (AstSlice, error)
```
For each resulting AST struct(s) in the slice, get its `Value` and concatenate the strings,
and Remove consecutive whitespaces from the end of the string.  
If the `Value` are not string, an error occurs.

```go
func ParseIntRadix(base int) TransformerFn
```
For each resulting AST struct(s) in the slice, get its `Value` and concatenate the strings,
and convert to `int64` specifying the radix.  
If the `Value` are not string, or if it cannot be converted, an error occurs.

```go
func ParseInt(ctx ParserContext, asts AstSlice) (AstSlice, error)
```
For each resulting AST struct(s) in the slice, get its `Value` and concatenate the strings,
and convert it to `int64` with radix 10.  
If the `Value` are not string, or if it cannot be converted, an error occurs.

```go
func ParseUintRadix(base int) TransformerFn
```
For each resulting AST struct(s) in the slice, get its `Value` and concatenate the strings,
and convert to `uint64` specifying the radix.  
If the `Value` are not string, or if it cannot be converted, an error occurs.

```go
func ParseUint(ctx ParserContext, asts AstSlice) (AstSlice, error)
```
For each resulting AST struct(s) in the slice, get its `Value` and concatenate the strings,
and convert it to `uint64` with radix 10.  
If the `Value` are not string, or if it cannot be converted, an error occurs.

```go
func ParseFloat(ctx ParserContext, asts AstSlice) (AstSlice, error)
```
For each resulting AST struct(s) in the slice, get its `Value` and concatenate the strings,
and convert it to `float64`.  
If the `Value` are not string, or if it cannot be converted, an error occurs.

```go
func RuneFromInt(ctx ParserContext, asts AstSlice) (AstSlice, error)
```
Converts the 0th `Value` of the resulting AST struct slice from `int64` to `rune`.  
If `Value` is not `int64`, it panics.

```go
func IntFromRune(ctx ParserContext, asts AstSlice) (AstSlice, error)
```
Converts the 0th `Value` of the resulting AST struct slice from `rune` to `int64`.  
If `Value` is not `rune`, it panics.

```go
func StringFromInt(ctx ParserContext, asts AstSlice) (AstSlice, error)
```
Converts the 0th `Value` of the resulting AST struct slice from `int64` to `rune` and then to `string`.  
If `Value` is not `int64`, it panics.

```go
func StringFromRune(ctx ParserContext, asts AstSlice) (AstSlice, error)
```
Converts the 0th `Value` of the resulting AST struct slice from `rune` to `string`.  
If `Value` is not `rune`, it panics.


## github.com/shellyln/takenoco/extra

```go
func ParseDate(ctx ParserContext, asts AstSlice) (AstSlice, error)
```
Converts the `[len(asts) - 1]`th `Value` of the resulting AST struct slice to `time.Time`
as a date string in ISO 8061 extended format. The resulting timezone will be UTC.  
If `Value` is not `string`, it panics.

```go
func ParseDateTime(ctx ParserContext, asts AstSlice) (AstSlice, error)
```
Converts the `[len(asts) - 1]`th `Value` of the resulting AST struct slice to `time.Time`
as a date/time string in ISO 8061 extended format. The resulting timezone will be UTC.  
If `Value` is not `string`, it panics.

```go
func ParseTime(ctx ParserContext, asts AstSlice) (AstSlice, error)
```
Converts the `[len(asts) - 1]`th `Value` of the resulting AST struct slice to `time.Time`
as a time string in ISO 8061 extended format.
The resulting date part will be a Unix epoch and the timezone will be UTC.  
If `Value` is not `string`, it panics.


# Transform parsed result ASTs by production rules

In the "Introduction of common transform functions" section, we introduced a set of functions that perform simple transformations,
such as string processing and string-to-number conversions, on certain kinds of limited and small parsed results.
User-defined conversion functions can also be used to perform type conversions not provided by the library.

But it is complicated to create transform functions across multiple tokens each time.
For example, when evaluating four arithmetic operations according to the precedence of the operators.  
In [takenoco](https://github.com/shellyln/takenoco), the framework is provided in advance to perform such transformations with rules called "production rules".

This chapter explains how to apply the production rules using the [Formula to RPN converter](https://github.com/shellyln/takenoco/tree/master/_examples/torpn) example in the takenoco repository.  
Formula to RPN converter transforms an integer four arithmetic operations formula into a Reverse Polish Notation (RPN) formula.


## Parser

The parser is an advanced of the "minimal parser" chapter.
This example differs from the "minimal parser" in that it defines parsers as functions that combine the common parsers, and calls them.
It is easier to understand if functions are defined using grammatical chunks ("Terminal symbols" and "Nonterminal symbols").

The production rules are applied by calling the `Trans()` function in the same way as the functions in the previous chapter,
"Introduction of common transform functions".

torpn.go
```go
...

// Expression enclosed in parentheses
func groupedExpresion() ParserFn {
    return FlatGroup(
        erase(CharClass("(")), // 1) Parentheses are removed with `erase` so that they are not included in the result.
                               //    Parentheses can also be converted by the production rules,
                               //    but in this example they are handled by the parser.
        First(
            FlatGroup(
                erase(sp0()),
                expression(), // 2) Expression production rules applied.
                erase(CharClass(")")),
                erase(sp0()),
            ),
            Error("Error in grouped expression"),
        ),
    )
}

// Expression before applying the production rules
func expressionInner() ParserFn {
    return FlatGroup(
        ZeroOrMoreTimes(unaryOperator()), // 3) Unary operator.
        First(
            simpleExpression(),         // 4) Expression without parentheses.
            Indirect(groupedExpresion), // 5) Expression enclosed in parentheses.
            Error("Value required"),
        ),
        ZeroOrMoreTimes(
            binaryOperator(), // 6) Binary operator.
            First(
                FlatGroup(
                    ZeroOrMoreTimes(unaryOperator()),
                    First(
                        simpleExpression(),
                        Indirect(groupedExpresion),
                    ),
                ),
                Error("Error in the expression after the binary operator"),
            ),
        ),
    )
}

// Applying production rules to expression
func expression() ParserFn {
    return Trans( // 7) `Trans()` is also used to apply the production rules.
        expressionInner(),
        formulaProductionRules(), // 8) Transform function to apply the production rules.
    )
}

// Entire program
func program() ParserFn {
    return FlatGroup(
        Start(),
        erase(sp0()),
        expression(),
        End(),
    )
}

...
```

## Production rules

A transform function using production rules can be constructed by calling `ProductionRule()` (`19)`)
with the matching conditions of the production rule and a small transform function
for each precedence order.

Expression and statement parsing often translates them into a tree structure called an abstract syntax tree (AST).  
The result of takenoco's parsing is a slice of the AST struct, that can have a tree structure under it.  
However, in this example, since we are only converting to a sequence of RPN tokens,
the slice is extended horizontally without tree structure in the transform functions
for the production rules (e.g. `transformUnaryOp`, `transformBinaryOp`).

precedence.go
```go
...

// Production rules for precedence 3 (unary prefix operator)
var expressionRule3 = Precedence{ // 1) Sets attributes to the Precedence struct.
    Rules: []ParserFn{ // 2) Describes the production rules as a slice.
                       //    This allows multiple production rules within the same precedence order.
        Trans(
            FlatGroup( // 3) Describes the match condition of the production rule as a parser.
                       //    (Use object parser)
                       //    Generative grammar:
                       //      (Expression -> UnaryOperator Expression)
                       //      (Expression -> Number)
                isOperator("UnaryOperator", []string{"-"}),
                anyOperand(),
            ),
            transformUnaryOp, // 4) Describes the production rule as transform function.
        ),
    },
    Rtol: true, // 5) Rules are scanned and applied in right-to-left order (default is left-to-right).
}

// Production rules for precedence 2 (binary multiplication and division operators)
var expressionRule2 = Precedence{
    Rules: []ParserFn{
        Trans(
            FlatGroup( // 6) Describes the match condition of the production rule as a parser.
                       //    (Use object parser)
                       //    Generative grammar:
                       //      (Expression -> Expression BinaryOperator Expression)
                       //      (Expression -> Number)
                anyOperand(),
                isOperator("BinaryOperator", []string{"*", "/"}),
                anyOperand(),
            ),
            transformBinaryOp, // 7) Describes the production rule as transform function.
        ),
    },
}

// Production rules for precedence 1 (binary addition and subtraction operators)
var expressionRule1 = Precedence{
    Rules: []ParserFn{
        Trans(
            FlatGroup( // 8) Describes the match condition of the production rule as a parser.
                       //    (Use object parser)
                       //    Generative grammar:
                       //      (Expression -> Expression BinaryOperator Expression)
                       //      (Expression -> Number)
                anyOperand(),
                isOperator("BinaryOperator", []string{"+", "-"}),
                anyOperand(),
            ),
            transformBinaryOp, // 9) Describes the production rule as transform function.
        ),
    },
}

// Production rule for unary prefix operator
func transformUnaryOp(ctx ParserContext, asts AstSlice) (AstSlice, error) {
    switch asts[1].Type {
    case AstType_Int:
        // 10) If operand is Number, applies a sign to the value.
        opcode := asts[0].Value.(string)
        op1 := asts[1].Value.(int64)

        var v int64
        switch opcode {
        case "-":
            v = -op1
        }

        return AstSlice{{
            Type:      AstType_Int,
            ClassName: "Number",
            Value:     v,
        }}, nil
    case AstType_ListOfAst:
        // 11) If the operand is an expression, add an operator at the end of the expression.
        ret := asts[1].Value.(AstSlice)
        ret = append(ret, asts[0]) // 12) Extend the slice horizontally.
        return AstSlice{{
            Type:      AstType_ListOfAst,
            ClassName: "List",
            Value:     ret,
        }}, nil
    default:
        return nil, errors.New("Cannot apply unary - operator")
    }
}

// Production rule for binary operators
func transformBinaryOp(ctx ParserContext, asts AstSlice) (AstSlice, error) {
    // 13) Swap the order of operands and operators from [0, 1, 2] to [0, 2, 1].
    opcode := asts[1].Value.(string)
    ret := make(AstSlice, 0, 3)

    switch asts[0].Type {
    case AstType_Int:
        // 14) Extend the slice horizontally.
        ret = append(ret, asts[0])
    case AstType_ListOfAst:
        // 15) Extend the slice horizontally.
        ret = append(ret, asts[0].Value.(AstSlice)...)
    default:
        return nil, errors.New("Cannot apply binary " + opcode + " operator")
    }

    switch asts[2].Type {
    case AstType_Int:
        // 16) Extend the slice horizontally.
        ret = append(ret, asts[2])
    case AstType_ListOfAst:
        // 17) Extend the slice horizontally.
        ret = append(ret, asts[2].Value.(AstSlice)...)
    default:
        return nil, errors.New("Cannot apply binary " + opcode + " operator")
    }

    ret = append(ret, asts[1])

    return AstSlice{{
        Type:      AstType_ListOfAst,
        ClassName: "List",
        Value:     ret,
    }}, nil
}

// Defines the priority of the production rules
var precedences = []Precedence{
    expressionRule3, // 18) The smallest index has higher priority.
    expressionRule2,
    expressionRule1,
}

// Constructs a transform function of the production rule
func formulaProductionRules() TransformerFn {
    return ProductionRule( // 19) Returns the transform function of the production rule.
        precedences, // 20) Production rules in order of precedence.

        // 20) Describes the pattern that must be satisfied by the results
        //     of applying all the production rules as a parser.
        //     Here we define that it must be transformed into any single AST struct.
        //     Generative grammar:
        //       (S -> Expression)
        FlatGroup(Start(), objparser.Any(), objparser.End()),
    )
}
```

> **Note**  
> Several helper functions are defined in the package (`isOperator`, `anyOperand`, etc.) .


**Practice**
* Let's add the remainder (`%`), bitwise-and (`&`), and bitwise-or (`|`) operators.

# Next step

Check out the other examples.

* [CSV parser](https://github.com/shellyln/takenoco/tree/master/_examples/csv)
* [Formula parser](https://github.com/shellyln/takenoco/tree/master/_examples/formula)
* [Formula to RPN converter](https://github.com/shellyln/takenoco/tree/master/_examples/torpn)
* [Loose JSON + TOML parsers](https://github.com/shellyln/go-loose-json-parser)
    * [JSON parser](https://github.com/shellyln/go-loose-json-parser/blob/master/jsonlp/json.go)
    * [TOML parser](https://github.com/shellyln/go-loose-json-parser/blob/master/jsonlp/toml.go)
    * [Live demo (Loose JSON | TOML normalizer)](https://shellyln.github.io/jsonlp/)
* [Dust - toy scripting language](https://github.com/shellyln/dust-lang)
    * [Parsers](https://github.com/shellyln/dust-lang/tree/master/scripting/parser)
    * [Production rules](https://github.com/shellyln/dust-lang/tree/master/scripting/rules)
* [OpenSOQL parser](https://github.com/shellyln/go-open-soql-parser)
    * [Parser](https://github.com/shellyln/go-open-soql-parser/blob/master/soql/parser/parser.go)
    * [Live demo](https://shellyln.github.io/soql/)
