package nmconv

import (
	"bytes"
	"strings"
	"unicode"
)

const (
	Lisp  = "-"
	Snake = "_"
)

type Normalize func(string) []string

type Denormalize func([]string) string

func Convert(str string, from Normalize, to Denormalize) string {
	tmp := from(str)
	return to(tmp)
}

func Sep(separator string) Denormalize {
	return func(norm []string) string {
		return strings.Join(norm, separator)
	}
}

func SepX(xformWord func(string) string, separator string) Denormalize {
	return func(norm []string) string {
		tmp := make([]string, len(norm))
		for i := 0; i < len(norm); i++ {
			tmp[i] = xformWord(norm[i])
		}
		return strings.Join(tmp, separator)
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

func SepXConvention(xformWord func(string) string, separator string) Conversion {
	return Conversion{
		Norm:   Unsep(separator),
		Denorm: SepX(xformWord, separator),
	}
}

func capWord(w string) string {
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
			buf.WriteString(capWord(norm[i]))
		}
		return buf.String()
	}
}

func Camel1Up(norm []string) string {
	buf := bytes.NewBuffer(nil)
	for _, w := range norm {
		buf.WriteString(capWord(w))
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
	Denorm Denormalize
}

func (cnv Conversion) Convert(str string) string {
	return Convert(str, cnv.Norm, cnv.Denorm)
}
