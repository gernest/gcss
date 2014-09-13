package gcss

import (
	"bytes"
	"fmt"
	"io"
	"strings"
)

// declaration represents a declaration of CSS.
type declaration struct {
	elementBase
	property string
	value    string
}

// WriteTo writes the declaration to the writer.
func (dec *declaration) WriteTo(w io.Writer) (int64, error) {
	bf := new(bytes.Buffer)

	bf.WriteString(dec.property)
	bf.WriteString(colon)

	if strings.HasPrefix(dec.value, dollarMark) {
		v, ok := dec.Context().vars[strings.TrimPrefix(dec.value, dollarMark)]
		if ok {
			// Writing to the bytes.Buffer never returns an error.
			v.WriteTo(bf)
		}
	} else {
		bf.WriteString(dec.value)
	}

	bf.WriteString(semicolon)

	n, err := w.Write(bf.Bytes())

	return int64(n), err
}

// declarationPV extracts a declaration property and value
// from the line.
func declarationPV(ln *line) (string, string, error) {
	pv := strings.SplitN(strings.TrimSpace(ln.s), space, 2)

	if len(pv) < 2 {
		return "", "", fmt.Errorf("declaration's property and value should be divided by a space [line: %d]", ln.no)
	}

	if !strings.HasSuffix(pv[0], colon) {
		return "", "", fmt.Errorf("property should end with a colon [line: %d]", ln.no)
	}

	return strings.TrimSuffix(pv[0], colon), pv[1], nil
}

// newDeclaration creates and returns a declaration.
func newDeclaration(ln *line, parent element) (*declaration, error) {
	property, value, err := declarationPV(ln)

	if err != nil {
		return nil, err
	}

	if strings.HasSuffix(value, semicolon) {
		return nil, fmt.Errorf("declaration must not end with %q [line: %d]", semicolon, ln.no)
	}

	return &declaration{
		elementBase: newElementBase(ln, parent),
		property:    property,
		value:       value,
	}, nil
}
