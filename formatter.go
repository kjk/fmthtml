package fmthtml

import (
	"bytes"
	"io"
	"strconv"
	"strings"
)

const (
	defaultIndentString = "  "
	startIndent         = 0
	defaultLastElement  = "</html>"
)

// A Writer represents a formatted HTML source codes writer.
type Writer struct {
	writer      io.Writer
	lastElement string
	bf          *bytes.Buffer
}

// SetLastElement set the lastElement to the Writer.
func (wr *Writer) SetLastElement(lastElement string) *Writer {
	wr.lastElement = lastElement
	return wr
}

// Write writes the parameter.
func (wr *Writer) Write(p []byte) (n int, err error) {
	wr.bf.Write(p)
	if bytes.HasSuffix(p, []byte(wr.lastElement)) {
		return wr.writer.Write([]byte(Format(wr.bf.String()) + "\n"))
	}
	return 0, nil
}

// NewWriter generates a Writer and returns it.
func NewWriter(wr io.Writer) *Writer {
	return &Writer{writer: wr, lastElement: defaultLastElement, bf: &bytes.Buffer{}}
}

// A textElement represents a text element of an HTML document.
type textElement struct {
	text string
}

// write writes a text to the buffer.
func (e *textElement) write(bf *bytes.Buffer, indent int) {
	lines := strings.Split(strings.Trim(unifyLineFeed(e.text), "\n"), "\n")
	for _, line := range lines {
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

// html generates an HTML source code and returns it.
func (htmlDoc *htmlDocument) html() string {
	return string(htmlDoc.bytes())
}

// bytes reads from htmlDocument's internal array of elements and returns HTML source code
func (htmlDoc *htmlDocument) bytes() []byte {
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

// Format parses the input HTML string, formats it and returns the result.
func Format(s string) string {
	return parse(strings.NewReader(s)).html()
}

// FormatBytes parses input HTML as bytes, formats it and returns the result.
func FormatBytes(b []byte) []byte {
	return parse(bytes.NewReader(b)).bytes()
}

// Format parses the input HTML string, formats it and returns the result with line no.
func FormatWithLineNo(s string) string {
	return AddLineNo(Format(s))
}

func AddLineNo(s string) string {
	lines := strings.Split(s, "\n")
	maxLineNoStrLen := len(strconv.Itoa(len(lines)))
	bf := &bytes.Buffer{}
	for i, line := range lines {
		lineNoStr := strconv.Itoa(i + 1)
		if i > 0 {
			bf.WriteString("\n")
		}
		bf.WriteString(strings.Repeat(" ", maxLineNoStrLen-len(lineNoStr)) + lineNoStr + "  " + line)
	}
	return bf.String()

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
