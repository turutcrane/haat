# haat (Html as a Template)

Wrapper Library for x/net/html and github.com/ericchiang/css

```go
package main

import (
	"log"
	"os"

	"github.com/turutcrane/haat"
)

func main() {

	h, err := haat.HtmlParsePageString(`
<!DOCTYPE html>
<html>
<head>
<title>Hello haat</title>
</head>
<body>
Hello <span id="pkgname"></span>!!
</body>
</html>
`)
	if err != nil {
		log.Panicln(err)
	}
	for span := range h.Query("span#pkgname") {
		span.C(haat.T("haat"))
	}
	h.Render(os.Stdout)
}
```

Output:
```
<!DOCTYPE html><html><head>
<title>Hello haat</title>
</head>
<body>
Hello <span id="pkgname">haat</span>!!


</body></html>
```
