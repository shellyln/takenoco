package parser

// An interface that provides a debug trace callback.
type ParserTracer interface {
	// before event
	Before(scope string, trNo int, className string, ctx *ParserContext)
	// after event
	After(scope string, trNo int, className string, ctx *ParserContext, err error)
	// error event
	Panic(scope string, trNo int, className string, ctx *ParserContext, r interface{})
}

// A debug trace callback that does nothing.
type nullParserTracer struct{}

// before event
func (s nullParserTracer) Before(scope string, trNo int, className string, ctx *ParserContext) {
}

// after event
func (s nullParserTracer) After(scope string, trNo int, className string, ctx *ParserContext, err error) {
}

// error event
func (s nullParserTracer) Panic(scope string, trNo int, className string, ctx *ParserContext, r interface{}) {
}

var (
	//
	traceEnabled bool
	//
	traceScope string
	//
	traceParserTrackingNo int
	//
	parserTracer ParserTracer = nullParserTracer{}
)

// Debug tracing is enabled for child parsers.
// NOTE: It's not thread safe.
func DebugTrace(scope string, pt ParserTracer) func(children ...ParserFn) ParserFn {
	const ClassName = ":Base:DebugTrace"

	savedTraceEnabled := traceEnabled
	savedTraceScope := traceScope
	savedParserTracer := parserTracer

	traceEnabled = true
	traceScope = traceScope + "/" + scope
	parserTracer = pt

	return func(children ...ParserFn) ParserFn {
		defer func() {
			traceEnabled = savedTraceEnabled
			traceScope = savedTraceScope
			parserTracer = savedParserTracer
		}()
		parser := BaseParser(ClassName, nil, nil, children, nil)
		return parser
	}
}

// Handler for debug traces.
func tracer(scope string, pt ParserTracer, trNo int, className string, parser ParserFn) ParserFn {
	return func(ctx ParserContext) (ParserContext, error) {
		defer func() {
			if r := recover(); r != nil {
				pt.Panic(scope, trNo, className, &ctx, r)
				panic(r)
			}
		}()

		pt.Before(scope, trNo, className, &ctx)
		ctx, err := parser(ctx)
		pt.After(scope, trNo, className, &ctx, err)
		return ctx, err
	}
}
