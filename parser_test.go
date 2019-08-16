package fmthtml

import (
	"strings"
	"testing"
)

func TestHtmlEscape(t *testing.T) {
	Condense = false
	s := `<!DOCTYPE html><html><body><div>0 &lt; 1. great insight! &lt;/sarcasm&gt; over&amp;out.&</div></body></html>`
	expected := `<!DOCTYPE html>
<html>
  <body>
    <div>
      0 &lt; 1. great insight! &lt;/sarcasm&gt; over&amp;out.&
    </div>
  </body>
</html>`
	htmlDoc := parse(strings.NewReader(s))
	actual := string(htmlDoc.toHTML())
	if actual != expected {
		t.Errorf("Invalid result. [expected: %s][actual: %s]", expected, actual)
	}
}
