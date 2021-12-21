package parser

import (
	"testing"
)

func TestA(t *testing.T) {
	actual := true
	expected := true
	if actual != expected {
		t.Errorf("%v, %v", actual, expected)
	}
}
