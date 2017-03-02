package example

import (
	"testing"
)

func TestExample(t *testing.T) {
	xx := New(Name("haha"), Age(11), Addr("some address"))
	t.Log(xx)
}
