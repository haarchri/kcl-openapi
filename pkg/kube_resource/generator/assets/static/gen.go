// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package static

//go:generate go run makestatic.go

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"strings"
	"unicode"
)

var files = []string{
	"../files/api_spec/k8s/k8s.json",
}

// Generate reads a set of files and returns a file buffer that declares
// a map of string constants containing contents of the input files.
func Generate() ([]byte, error) {
	buf := new(bytes.Buffer)
	fmt.Fprintf(buf, "%v\n\n%v\n\npackage static\n\n", license, warning)
	fmt.Fprintf(buf, "var Files = map[string]string{\n")
	for _, fn := range files {
		b, err := os.ReadFile(fn)
		if err != nil {
			return b, err
		}
		fmt.Fprintf(buf, "\t%q: ", strings.Replace(fn, "../files/", "", 1))
		appendQuote(buf, b)
		fmt.Fprintf(buf, ",\n\n")
	}
	fmt.Fprintln(buf, "}")
	return format.Source(buf.Bytes())
}

// appendQuote is like strconv.AppendQuote, but we avoid the latter
// because it changes when Unicode evolves, breaking gen_test.go.
func appendQuote(out *bytes.Buffer, data []byte) {
	out.WriteByte('"')
	for _, b := range data {
		if b == '\\' || b == '"' {
			out.WriteByte('\\')
			out.WriteByte(b)
		} else if b <= unicode.MaxASCII && unicode.IsPrint(rune(b)) && !unicode.IsSpace(rune(b)) {
			out.WriteByte(b)
		} else {
			fmt.Fprintf(out, "\\x%02x", b)
		}
	}
	out.WriteByte('"')
}

const warning = `// Code generated by "makestatic"; DO NOT EDIT.`

const license = `// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.`
