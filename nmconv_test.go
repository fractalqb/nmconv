package nmconv

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stvp/assert"
)

func TestUncamel_1wordLow(t *testing.T) {
	res := Uncamel("foo")
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "foo", res[0])
}

func TestUncamel_1wordUp(t *testing.T) {
	res := Uncamel("Foo")
	assert.Equal(t, 1, len(res))
	assert.Equal(t, "Foo", res[0])
}

func TestUncamel_2wordLow(t *testing.T) {
	res := Uncamel("fooBar")
	assert.Equal(t, 2, len(res))
	assert.Equal(t, "foo", res[0])
	assert.Equal(t, "Bar", res[1])
}

func TestUncamel_2wordUp(t *testing.T) {
	res := Uncamel("FooBar")
	assert.Equal(t, 2, len(res))
	assert.Equal(t, "Foo", res[0])
	assert.Equal(t, "Bar", res[1])
}

func TestUncamel_short(t *testing.T) {
	res := Uncamel("FB")
	assert.Equal(t, 2, len(res), res)
	assert.Equal(t, "F", res[0])
	assert.Equal(t, "B", res[1])
}

func TestCamelLow(t *testing.T) {
	res := Camel1Low([]string{"FOO", "bar"})
	assert.Equal(t, "fooBar", res)
}

func TestCamelUp(t *testing.T) {
	res := Camel1Up([]string{"foo", "BAR"})
	assert.Equal(t, "FooBar", res)
}

func TestCamel_shortLow(t *testing.T) {
	res := Camel1Low([]string{"f", "b"})
	assert.Equal(t, "fB", res)
}

func TestCamel_shortUp(t *testing.T) {
	res := Camel1Up([]string{"f", "b"})
	assert.Equal(t, "FB", res)
}

func TestFunConv(t *testing.T) {
	conv := Conversion{Norm: Unsep("_"), Denorm: Sep("-")}
	res := conv.Convert("foo_bar_baz")
	if res != "foo-bar-baz" {
		t.Error(res)
	}
}

func BenchmarkFunConv(b *testing.B) {
	conv := Conversion{Norm: Unsep("_"), Denorm: Sep("-")}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conv.Convert("foo_bar_baz")
	}
}

func ExamplePassConversion() {
	user := func(nm string, conv func(string) string) {
		fmt.Printf("'%s' → '%s'", nm, conv(nm))
	}
	conv := Conversion{Norm: Unsep("_"), Denorm: Sep("-")}
	user("foo_bar_baz", conv.Convert)
	// Output:
	// 'foo_bar_baz' → 'foo-bar-baz'
}

type IfConvn interface {
	Norm(str string) []string
	Denorm(n []string) string
}

type IfSep string

func (s IfSep) Norm(str string) []string {
	return strings.Split(str, string(s))
}

func (s IfSep) Denorm(n []string) string {
	return strings.Join(n, string(s))
}

type IfCnv struct {
	From IfConvn
	To   IfConvn
}

func (cnv IfCnv) Convert(str string) string {
	tmp := cnv.From.Norm(str)
	return cnv.To.Denorm(tmp)
}

func BenchmarkIfConv(b *testing.B) {
	conv := IfCnv{From: IfSep("_"), To: IfSep("-")}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		conv.Convert("foo_bar_baz")
	}
}

// TODO new tests for the Transform concept
