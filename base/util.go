package parser

import (
	"strings"
)

//
func GetLineAndColPosition(src string, pos SourcePosition) LineAndColPosition {
	prevLineIndex := 0
	lineIndex := 0
	line := 1
	col := 1

	for i := 0; i < pos.Position; i++ {
		switch src[i] {
		case '\r':
			line++
			if i+1 < pos.Position && src[i+1] == '\n' {
				i++
			}
			prevLineIndex = lineIndex
			lineIndex = i + 1
			col = 1
		case '\n':
			line++
			prevLineIndex = lineIndex
			lineIndex = i + 1
			col = 1
		default:
			col++
		}
	}

	srcLen := len(src)
	endIndex := pos.Position
OUTER:
	for i := pos.Position; i < srcLen; i++ {
		switch src[i] {
		case '\r', '\n':
			break OUTER
		default:
			endIndex++
		}
	}

	errGuideCol := col - 1
	if errGuideCol < 0 {
		errGuideCol = 0
	}

	errSource := src[prevLineIndex:endIndex]
	if strings.Contains(errSource, "\n") {
		errSource = "   | " + strings.Replace(errSource, "\n", "\n > | ", 1) + "\n"
	} else if strings.Contains(errSource, "\r") {
		errSource = "   | " + strings.Replace(errSource, "\r", "\r > | ", 1) + "\r"
	} else {
		errSource = " > | " + errSource + "\n"
	}
	errSource += "   | " + strings.Repeat(" ", errGuideCol) + "^^^^"

	return LineAndColPosition{
		LineIndex: lineIndex,
		Line:      line,
		Col:       col,
		Position:  pos.Position,
		ErrSource: errSource,
	}
}
