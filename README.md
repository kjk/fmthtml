# fmthtml - Go package for formatting HTML

`fmthtml` is an HTML formatter (pretty-printer) for Go.

## Example

```go
package main

import (
	"fmt"

	"github.com/kjk/fmthtml"
)

func main() {
  d := `<!DOCTYPE html><html><head><title>This is a title.</title></head><body>foo</body></html>`
  d = fmthtml.Format([]byte(d))
	fmt.Printf("%s\n", d)
}
```

Output:

```html
<!DOCTYPE html>
<html>
  <head>
    <title>
      This is a title.
    </title>
  </head>
  <body>
    foo
  </body>
</html>
```

## Docs

- [GoDoc](https://godoc.org/github.com/kjk/fmthtml)
