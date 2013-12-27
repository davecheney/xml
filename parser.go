package xml

import (
	"fmt"
	"io"
)

type Parser struct {
	r io.Reader
}

func (p *Parser) Parse() (Node, error) {
	const (
		CHARACTERS = iota
		STAG_NAME_START
		STAG_NAME
		DECLARATION_START
		PROCESSING_INSTRUCTION_START
		ETAG_NAME_START
		ETAG_NAME
		ATTRIBUTE_NAME_START
		ELEMENT_EMPTY_END
		ELEMENT_END
	)
	state := CHARACTERS
	builder := &builder{
		declaredNamespaces: map[string]string{
			"":      "",
			"xml":   "http://www.w3.org/XML/1998/namespace",
			"xmlns": "http://www.w3.org/2000/xmlns/",
		},
	}

	buf := make([]byte, 4096)

	// fill buffer
	n, err := p.r.Read(buf)
	var offset, limit int
	reset := func() {
		offset = limit
		offset++
	}
	subsequence := func() string {
		return string(buf[offset:limit])
	}
	doCharacters := func() {
		chars := subsequence()
		reset()
		if isBlank(chars) {
			return
		}
		builder.doCharacters(chars)
	}
	doElementStart := func() {
		builder = builder.doElementStart(subsequence())
		reset()
	}
	doElementEnd := func() {
		builder = builder.doElementEnd()
		reset()
	}
	for ; limit < n; limit++ {
		c := buf[limit]
		switch state {
		case CHARACTERS:
			if c == '<' {
				doCharacters()
				state = STAG_NAME_START
			}
		case STAG_NAME_START:
			if isNameStartChar(rune(c)) {
				state = STAG_NAME
			} else if c == '!' {
				reset()
				state = DECLARATION_START
			} else if c == '?' {
				reset()
				state = PROCESSING_INSTRUCTION_START
			} else if c == '/' {
				reset()
				state = ETAG_NAME_START
			} else {
				return nil, fmt.Errorf("%q not allowed in state %v", c, STAG_NAME_START)
			}
		case STAG_NAME:
			if isNameChar(rune(c)) {
				// consume
			} else if isWhitespace(rune(c)) {
				doElementStart()
				state = ATTRIBUTE_NAME_START
			} else if c == '>' {
				doElementStart()
				state = CHARACTERS
			} else if c == '/' {
				doElementStart()
				state = ELEMENT_EMPTY_END
			} else {
				return nil, fmt.Errorf("%q not allowed in state %v", c, STAG_NAME)
			}
		case ETAG_NAME_START:
			if isNameStartChar(rune(c)) {
				state = ETAG_NAME
			} else {
				return nil, fmt.Errorf("%q not allowed in state %v", c, ETAG_NAME_START)
			}
		case ETAG_NAME:
			if isNameChar(rune(c)) {
				// consume
			} else if isWhitespace(rune(c)) {
				reset()
				state = ELEMENT_END
			} else if c == '>' {
				doElementEnd()
				state = CHARACTERS
			} else {
				return nil, fmt.Errorf("%q not allowed in state %v", c, ETAG_NAME)
			}
		case ELEMENT_EMPTY_END:
			if c == '>' {
				doElementEnd()
				state = CHARACTERS
			} else {
				return nil, fmt.Errorf("%q not allowed in state %v", c, ELEMENT_EMPTY_END)
			}
		default:
			return nil, fmt.Errorf("unhandled state: %v", state)
		}
	}
	if err != nil {
		return nil, err
	}
	return builder.content[0], nil
}

// TODO(dfc)
func isBlank(s string) bool {
	for _, c := range s {
		if !isWhitespace(c) {
			return false
		}
	}
	return true
}

func isNameStartChar(c rune) bool {
	return c == ':' || (c >= 'A' && c <= 'Z') || c == '_' || (c >= 'a' && c <= 'z') || (c >= '\u00C0' && c <= '\u00D6') || (c >= '\u00D8' && c <= '\u00F6') || (c >= '\u00F8' && c <= '\u02FF') || (c >= '\u0370' && c <= '\u037D') || (c >= '\u037F' && c <= '\u1FFF') || (c >= '\u200C' && c <= '\u200D') || (c >= '\u2070' && c <= '\u218F') || (c >= '\u2C00' && c <= '\u2FEF') || (c >= '\u3001' && c <= '\uD7FF') || (c >= '\uF900' && c <= '\uFDCF') || (c >= '\uFDF0' && c <= '\uFFFD')
}

func isNameChar(c rune) bool {
	return isNameStartChar(c) || c == '-' || c == '.' || (c >= '0' && c <= '9') || c == '\u00B7' || (c >= '\u0300' && c <= '\u036F') || (c >= '\u023F' && c <= '\u2040')
}

func isWhitespace(c rune) bool {
	return (c == ' ' || c == '\n' || c == '\r' || c == '\t')
}
