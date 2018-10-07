package nmconv

import (
	"bytes"
	"strings"
	"unicode"
)

const (
	// Lisp is the separator for lisp-like (kebab case) naming
	// convetions
	Lisp = "-"
	// Snake is the separator for C-like (snake case) naming
	// convetions
	Snake = "_"
)

// Normalize functions convert names from some naming convetion to the
// normalized form, i.e. an array of strings - each containing one word of the
// name.
type Normalize func(string) []string

// Denormalie functions convert names from the normalized form to a naming
// convention, i.e. it reverses the effect of a Normalize function.
type Denormalize func([]string) string

// Transform the normalized name. Can be used for up-/downcase conversion or
// adding pre- or postfixes.
type Transform func([]string) []string

// ChainX creates a Transform function that applies all xs functions in order.
func ChainX(xs ...Transform) Transform {
	return func(segs []string) []string {
		for _, x := range xs {
			segs = x(segs)
		}
		return segs
	}
}

// PerSegment creates a Transform function that replaces each name segment s
// in place with the result of x(s).
func PerSegment(x func(string) string) Transform {
	return func(segs []string) []string {
		for i := range segs {
			segs[i] = x(segs[i])
		}
		return segs
	}
}

// Prefix creates a Transform function that adds segmens at the beginning of a
// noramlized name.
func Prefix(pfs ...string) Transform {
	return func(segs []string) []string {
		res := make([]string, len(segs)+len(pfs))
		copy(res, pfs)
		copy(res[len(pfs):], segs)
		return res
	}
}

// Postfix creates a Transform function that adds segmens at the end of a
// noramlized name.
func Postfix(pfs ...string) Transform {
	return func(segs []string) []string {
		res := make([]string, len(segs)+len(pfs))
		copy(res, segs)
		copy(res[len(segs):], pfs)
		return res
	}
}

// Convert converts a given name by first normlizing it from its current naming
// convetion and then denormalizing its to the target naming convetnion.
func Convert(name string, from Normalize, to Denormalize) string {
	tmp := from(name)
	return to(tmp)
}

// ConvertX converts a given name by first normlizing it from its current naming
// convetion then transforming the normalized name with x and eventually
// denormalizing its to the target naming convention.
func ConvertX(name string, from Normalize, x Transform, to Denormalize) string {
	tmp := from(name)
	tmp = x(tmp)
	return to(tmp)
}

func NormX(n Normalize, x Transform) Normalize {
	return func(s string) []string {
		tmp := n(s)
		tmp = x(tmp)
		return tmp
	}
}

func XDenorm(x Transform, d Denormalize) Denormalize {
	return func(normal []string) string {
		normal = x(normal)
		return d(normal)
	}
}

func Sep(separator string) Denormalize {
	return func(norm []string) string {
		return strings.Join(norm, separator)
	}
}

func Unsep(separator string) Normalize {
	return func(str string) []string {
		return strings.Split(str, separator)
	}
}

func SepConvention(separator string) Conversion {
	return Conversion{
		Norm:   Unsep(separator),
		Denorm: Sep(separator),
	}
}

func SepXConvention(x Transform, separator string) Conversion {
	return Conversion{
		Norm:   Unsep(separator),
		Xform:  x,
		Denorm: Sep(separator),
	}
}

func CapWord(w string) string {
	if len(w) == 0 {
		return ""
	}
	buf := bytes.NewBuffer(make([]byte, 0, len(w)))
	buf.WriteString(strings.ToUpper(w[:1]))
	if w = w[1:]; len(w) > 0 {
		buf.WriteString(strings.ToLower(w))
	}
	return buf.String()
}

func Camel1Low(norm []string) string {
	if len(norm) == 0 {
		return ""
	} else {
		buf := bytes.NewBufferString(strings.ToLower(norm[0]))
		for i := 1; i < len(norm); i++ {
			buf.WriteString(CapWord(norm[i]))
		}
		return buf.String()
	}
}

func Camel1Up(norm []string) string {
	buf := bytes.NewBuffer(nil)
	for _, w := range norm {
		buf.WriteString(CapWord(w))
	}
	return buf.String()
}

func Uncamel(str string) (norm []string) {
	sep := strings.IndexFunc(str[1:], unicode.IsUpper)
	for sep >= 0 {
		sep++
		norm = append(norm, str[:sep])
		str = str[sep:]
		sep = strings.IndexFunc(str[1:], unicode.IsUpper)
	}
	norm = append(norm, str)
	return norm
}

type Conversion struct {
	Norm   Normalize
	Xform  Transform
	Denorm Denormalize
}

func (cnv *Conversion) Convert(str string) string {
	if cnv.Xform == nil {
		return Convert(str, cnv.Norm, cnv.Denorm)
	} else {
		return ConvertX(str, cnv.Norm, cnv.Xform, cnv.Denorm)
	}
}
