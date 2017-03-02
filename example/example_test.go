package example

import (
	"testing"
)

func TestExample(t *testing.T) {
	xx := New(Name("haha"), Option1("ok"))
	t.Log(xx)
}
