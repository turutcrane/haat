package haat

import (
	"bytes"
	"strings"
	"testing"

	"golang.org/x/net/html/atom"
)

// Clone メソッド をテストする
func TestClone(t *testing.T) {
	html := strings.NewReader(`<!doctype html><html><head>
<title>Hello haat</title>
</head>
<body>
Hello <span id="pkgname"></span>!!
</body></html>`,
	)

	ht, err := ParseHtml(html)
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
	ht, err := ParseHtmlFragment(strings.NewReader(`<p id="foo" class="baz bar">Hello</p>`), Elem(atom.Div))
	if err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}
	ht[0].SetClasses("bar", "baz")

	var buf bytes.Buffer
	if err := ht[0].Render(&buf); err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}

	expected := `<p class="bar baz" id="foo">Hello</p>`
	actual := buf.String()
	if actual != expected {
		t.Errorf("got: %v\nwant: %v", actual, expected)
	}
}

func TestSetA(t *testing.T) {
	ht, err := ParseHtmlFragment(strings.NewReader(`<p xxx="yyy" id="foo" xxx="zzz">Hello</p>`), Elem(atom.Div))
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
	ht, err := ParseHtmlFragment(strings.NewReader(`<p id="foo" xxx="aaa" xxx="bbb">Hello</p>`), Elem(atom.Div))
	if err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}
	ht[0].RemoveAttr("Xxx")

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
	ht, err := ParseHtmlFragment(strings.NewReader(`<p id="foo" class="bar baz bar">Hello</p>`), Elem(atom.Div))
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

func TestRawTest(t *testing.T) {
	ht, err := ParseHtmlFragment(strings.NewReader(`<p>aaa</p>`), Elem(atom.Div))
	if err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}
	ht[0].AppendC(RawT("<i>bbb</i>"))

	var buf bytes.Buffer
	if err := ht[0].Render(&buf); err != nil {
		t.Errorf("got: %v\nwant: %v", err, nil)
	}

	expected := `<p>aaa<i>bbb</i></p>`
	actual := buf.String()
	if actual != expected {
		t.Errorf("got: %v\nwant: %v", actual, expected)
	}
}