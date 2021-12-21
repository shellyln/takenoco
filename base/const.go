package parser

// Don't change these variables at runtime.
var (
	// quantifier {1,1}
	qtyOnce = Times{
		Min: 1,
		Max: 1,
	}
	// quantifier {0,1}
	qtyZeroOrOnce = Times{
		Min: 0,
		Max: 1,
	}
	// quantifier {0,}
	qtyZeroOrMoreTimes = Times{
		Min: 0,
		Max: -1,
	}
	// quantifier {1,}
	qtyOneOrMoreTimes = Times{
		Min: 1,
		Max: -1,
	}
	// quantifier {0,0}
	tempEmptyAstSlice = AstSlice{}
)
