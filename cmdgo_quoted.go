// Code generated by golang.org/x/tools/cmd/bundle. DO NOT EDIT.
//go:generate bundle -o cmdgo_quoted.go -prefix cmdgoQuoted cmd/internal/quoted

// Package quoted provides string manipulation utilities.
//

package main

import (
	"flag"
	"fmt"
	"strings"
	"unicode"
)

func cmdgoQuotedisSpaceByte(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}

// Split splits s into a list of fields,
// allowing single or double quotes around elements.
// There is no unescaping or other processing within
// quoted fields.
func cmdgoQuotedSplit(s string) ([]string, error) {
	// Split fields allowing '' or "" around elements.
	// Quotes further inside the string do not count.
	var f []string
	for len(s) > 0 {
		for len(s) > 0 && cmdgoQuotedisSpaceByte(s[0]) {
			s = s[1:]
		}
		if len(s) == 0 {
			break
		}
		// Accepted quoted string. No unescaping inside.
		if s[0] == '"' || s[0] == '\'' {
			quote := s[0]
			s = s[1:]
			i := 0
			for i < len(s) && s[i] != quote {
				i++
			}
			if i >= len(s) {
				return nil, fmt.Errorf("unterminated %c string", quote)
			}
			f = append(f, s[:i])
			s = s[i+1:]
			continue
		}
		i := 0
		for i < len(s) && !cmdgoQuotedisSpaceByte(s[i]) {
			i++
		}
		f = append(f, s[:i])
		s = s[i:]
	}
	return f, nil
}

// Join joins a list of arguments into a string that can be parsed
// with Split. Arguments are quoted only if necessary; arguments
// without spaces or quotes are kept as-is. No argument may contain both
// single and double quotes.
func cmdgoQuotedJoin(args []string) (string, error) {
	var buf []byte
	for i, arg := range args {
		if i > 0 {
			buf = append(buf, ' ')
		}
		var sawSpace, sawSingleQuote, sawDoubleQuote bool
		for _, c := range arg {
			switch {
			case c > unicode.MaxASCII:
				continue
			case cmdgoQuotedisSpaceByte(byte(c)):
				sawSpace = true
			case c == '\'':
				sawSingleQuote = true
			case c == '"':
				sawDoubleQuote = true
			}
		}
		switch {
		case !sawSpace && !sawSingleQuote && !sawDoubleQuote:
			buf = append(buf, []byte(arg)...)

		case !sawSingleQuote:
			buf = append(buf, '\'')
			buf = append(buf, []byte(arg)...)
			buf = append(buf, '\'')

		case !sawDoubleQuote:
			buf = append(buf, '"')
			buf = append(buf, []byte(arg)...)
			buf = append(buf, '"')

		default:
			return "", fmt.Errorf("argument %q contains both single and double quotes and cannot be quoted", arg)
		}
	}
	return string(buf), nil
}

// A Flag parses a list of string arguments encoded with Join.
// It is useful for flags like cmd/link's -extldflags.
type cmdgoQuotedFlag []string

var _ flag.Value = (*cmdgoQuotedFlag)(nil)

func (f *cmdgoQuotedFlag) Set(v string) error {
	fs, err := cmdgoQuotedSplit(v)
	if err != nil {
		return err
	}
	*f = fs[:len(fs):len(fs)]
	return nil
}

func (f *cmdgoQuotedFlag) String() string {
	if f == nil {
		return ""
	}
	s, err := cmdgoQuotedJoin(*f)
	if err != nil {
		return strings.Join(*f, " ")
	}
	return s
}
