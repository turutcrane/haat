package haat

import (
	"bytes"
	"testing"
)

// Clone メソッド をテストする
func TestClone(t *testing.T) {
	ht, err := HtmlParsePageString(`
<!doctype html>
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
		t.Errorf("got: %v\nwant: %v", err, nil)
	}
	ht2 := ht.Clone()
	var htBuf, ht2Buf bytes.Buffer	
	if err := ht.Render(&htBuf); err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}
	if err := ht2.Render(&ht2Buf); err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}
	expected := htBuf.String()
	actual := ht2Buf.String()
	if actual != expected {
		t.Errorf("got: %v\nwant: %v", actual, expected)
	}
}
