package haat

import (
	"bytes"
	"testing"

	"golang.org/x/net/html/atom"
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

func TestSetClasses(t *testing.T) {
	ht, err := HtmlParseFragmentString(`<p id="foo" class="bar baz">Hello</p>`, E(atom.Div))
	if err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}
	ht[0].SetClasses("bar", "baz")

	var buf bytes.Buffer
	if err := ht[0].Render(&buf); err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}

	expected := `<p class="bar baz bar baz" id="foo">Hello</p>`
	actual := buf.String()
	if actual != expected {
		t.Errorf("got: %v\nwant: %v", actual, expected)
	}
}

func TestSetA(t *testing.T) {
	ht, err := HtmlParseFragmentString(`<p xxx="yyy" id="foo" xxx="zzz">Hello</p>`, E(atom.Div))
	if err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}
	ht[0].SetA(Attr("xxx", "aaa"))

	var buf bytes.Buffer
	if err := ht[0].Render(&buf); err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}

	expected := `<p id="foo" xxx="aaa">Hello</p>`
	actual := buf.String()
	if actual != expected {
		t.Errorf("got: %v\nwant: %v", actual, expected)
	}
}

func TestRemoveAttr(t *testing.T) {
	ht, err := HtmlParseFragmentString(`<p id="foo" xxx="aaa" xxx="bbb">Hello</p>`, E(atom.Div))
	if err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}
	ht[0].RemoveAttr("xxx")

	var buf bytes.Buffer
	if err := ht[0].Render(&buf); err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}

	expected := `<p id="foo">Hello</p>`
	actual := buf.String()
	if actual != expected {
		t.Errorf("got: %v\nwant: %v", actual, expected)
	}
}

func TestRemoveClass(t *testing.T) {
	ht, err := HtmlParseFragmentString(`<p id="foo" class="bar baz bar">Hello</p>`, E(atom.Div))
	if err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}
	ht[0].RemoveClass("bar")

	var buf bytes.Buffer
	if err := ht[0].Render(&buf); err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}

	expected := `<p class="baz" id="foo">Hello</p>`
	actual := buf.String()
	if actual != expected {
		t.Errorf("got: %v\nwant: %v", actual, expected)
	}
}
