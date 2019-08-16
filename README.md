# fmthtml - HTML formatter for Go

fmthtml is an HTML formatter for Go. You can format HTML source codes by using this package.

## Example

Example Go source code:

```go
package main

import (
	"fmt"

	"github.com/yosssi/gohtml"
)

func main() {
	h := `<!DOCTYPE html><html><head><title>This is a title.</title><script type="text/javascript">
alert('aaa');
if (0 < 1) {
	alert('bbb');
}
</script><style type="text/css">
body {font-size: 14px;}
h1 {
	font-size: 16px;
	font-weight: bold;
}
</style></head><body><form><input type="name"><p>AAA<br>BBB></p></form><!-- This is a comment. --></body></html>`
	fmt.Println(gohtml.Format(h))
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
    <script type="text/javascript">
      alert("aaa");
      if (0 < 1) {
        alert("bbb");
      }
    </script>
    <style type="text/css">
      body {
        font-size: 14px;
      }
      h1 {
        font-size: 16px;
        font-weight: bold;
      }
    </style>
  </head>
  <body>
    <form>
      <input type="name" />
      <p>
        AAA
        <br />
        BBB>
      </p>
    </form>
    <!-- This is a comment. -->
  </body>
</html>
```

## Format Go html/template Package's Template's Execute Result

You can format [Go html/template package](http://golang.org/pkg/html/template/)'s template's execute result by passing `Writer` to the `tpl.Execute`:

```go
package main

import (
	"os"
	"text/template"

	"github.com/yosssi/gohtml"
)

func main() {

	tpl, err := template.New("test").Parse("<html><head></head><body>{{.Msg}}</body></html>")

	if err != nil {
		panic(err)
	}

	data := map[string]interface{}{"Msg": "Hello!"}

	err = tpl.Execute(gohtml.NewWriter(os.Stdout), data)

	if err != nil {
		panic(err)
	}
}
```

Output:

```html
<html>
  <head> </head>
  <body>
    Hello!
  </body>
</html>
```

## Docs

- [GoDoc](https://godoc.org/github.com/kjk/fmthtml)
