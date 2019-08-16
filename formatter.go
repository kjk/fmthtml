package fmthtml

import (
	"bytes"
	"strings"
)

const (
	defaultIndentString = "  "
	startIndent         = 0
)

// A textElement represents a text element of an HTML document.
type textElement struct {
	text string
}

// write writes a text to the buffer.
func (e *textElement) write(bf *bytes.Buffer, indent int) {
	s := unifyLineFeed(e.text)
	s = strings.Trim(s, "\n")
	lines := strings.Split(s, "\n")
	for _, line := range lines {
		line = strings.TrimLeft(line, " ")
		if line == "" {
			continue
		}
		writeLineFeed(bf)
		writeIndent(bf, indent)
		bf.WriteString(line)
	}
}

// A tagElement represents a tag element of an HTML document.
type tagElement struct {
	tagName     string
	startTagRaw string
	endTagRaw   string
	children    []element
}

// Condense any tag with no child tags (only text or nothing) onto a single line
var Condense bool

// write writes a tag to the buffer.
func (e *tagElement) write(bf *bytes.Buffer, indent int) {
	if Condense {
		l := len(e.children)
		if l == 0 {
			writeLine(bf, indent, e.startTagRaw, e.endTagRaw)
			return
		} else if l == 1 && e.endTagRaw != "" {
			if c, ok := e.children[0].(*textElement); ok {
				writeLine(bf, indent, e.startTagRaw, c.text, e.endTagRaw)
				return
			}
		}
	}

	writeLine(bf, indent, e.startTagRaw)
	for _, c := range e.children {
		var childIndent int
		if e.endTagRaw != "" {
			childIndent = indent + 1
		} else {
			childIndent = indent
		}
		c.write(bf, childIndent)
	}
	if e.endTagRaw != "" {
		writeLine(bf, indent, e.endTagRaw)
	}
}

// appendChild append an element to the element's children.
func (e *tagElement) appendChild(child element) {
	e.children = append(e.children, child)
}

// An element represents an HTML element.
type element interface {
	write(*bytes.Buffer, int)
}

// An htmlDocument represents an HTML document.
type htmlDocument struct {
	elements []element
}

// bytes reads from htmlDocument's internal array of elements and returns HTML source code
func (htmlDoc *htmlDocument) toHTML() []byte {
	bf := &bytes.Buffer{}
	for _, e := range htmlDoc.elements {
		e.write(bf, startIndent)
	}
	return bf.Bytes()
}

// append appends an element to the htmlDocument.
func (htmlDoc *htmlDocument) append(e element) {
	htmlDoc.elements = append(htmlDoc.elements, e)
}

// Format pretty-prints HTML
func Format(s []byte) []byte {
	doc := parse(bytes.NewReader(s))
	res := doc.toHTML()
	return res
}

// writeLine writes an HTML line to the buffer.
func writeLine(bf *bytes.Buffer, indent int, strs ...string) {
	writeLineFeed(bf)
	writeIndent(bf, indent)
	for _, s := range strs {
		bf.WriteString(s)
	}
}

// writeLineFeed writes a line feed to the buffer.
func writeLineFeed(bf *bytes.Buffer) {
	if bf.Len() > 0 {
		bf.WriteString("\n")
	}
}

// writeIndent writes indents to the buffer.
func writeIndent(bf *bytes.Buffer, indent int) {
	bf.WriteString(strings.Repeat(defaultIndentString, indent))
}

// unifyLineFeed unifies line feeds.
func unifyLineFeed(s string) string {
	return strings.Replace(strings.Replace(s, "\r\n", "\n", -1), "\r", "\n", -1)
}
