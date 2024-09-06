package haat

import (
	"fmt"
	"html/template"
	"io"
	"net/url"
	"slices"
	"strings"

	"github.com/ericchiang/css"
	"golang.org/x/net/html"
	"golang.org/x/net/html/atom"
)

func lower(s string) string {
	b := []byte(s)
	for i, c := range b {
		if 'A' <= c && c <= 'Z' {
			b[i] = c + 'a' - 'A'
		}
	}
	return string(b)
}

type Node html.Node
type Attribute html.Attribute

func (n *Node) ClearContents() *Node {
	bn := (*html.Node)(n)
	for c := bn.FirstChild; c != nil; c = bn.FirstChild {
		bn.RemoveChild(c)
	}
	return n
}

func (n *Node) C(childs ...*Node) *Node {
	n.ClearContents()
	n.AppendC(childs...)
	return n
}

func (n *Node) AppendC(childs ...*Node) *Node {
	for _, c := range childs {
		(*html.Node)(n).AppendChild((*html.Node)(c))
	}
	return n
}

func Attr(key, value string) Attribute {
	return (Attribute)(html.Attribute{
		Key: lower(key),
		Val: value,
	})
}

func (n *Node) A(attrs ...Attribute) *Node {
	slices.SortStableFunc(attrs, func(a, b Attribute) int {
		return strings.Compare(lower(a.Key), lower(b.Key))
	})

	lastKey := ""
	var newAttrs []html.Attribute
	for _, a := range attrs {
		if a.Key != "" {
			if a.Key == lastKey {
				newAttrs = slices.Delete(newAttrs, len(newAttrs)-1, len(newAttrs))
			}
			lastKey = a.Key
			newAttrs = append(newAttrs, html.Attribute(a))
		}
	}

	n.Attr = newAttrs
	return n
}

func (n *Node) SetA(attr ...Attribute) *Node {
	var attrs []Attribute
	for _, a := range n.Attr {
		attrs = append(attrs, Attribute(a))
	}
	return n.A(slices.Concat(attrs, attr)...)
}

func (n *Node) AppendA(attr Attribute) *Node {
	n.Attr = append(n.Attr, html.Attribute{Key: attr.Key, Val: attr.Val})
	return n
}

func AttrHref(u url.URL) Attribute {
	return Attr("href", u.String())
}

// func (n *Node) SetAHref(u url.URL) *Node {
// 	return n.AppendA(
// 		html.Attribute{
// 			Key: "href",
// 			Val: u.String(),
// 		},
// 	)
// }

// func jsStringVar(name, val string) string {
// 	return "var " + template.JSEscapeString(name) + " = \"" + template.JSEscapeString(val) + "\";"
// }

func JsLetString(name, val string) string {
	return "let " + template.JSEscapeString(name) + " = \"" + template.JSEscapeString(val) + "\";"
}

func JsConstString(name, val string) string {
	return "const " + template.JSEscapeString(name) + " = \"" + template.JSEscapeString(val) + "\";"
}

func Element(a atom.Atom) *Node {
	return (*Node)(
		&html.Node{
			Type:     html.ElementNode,
			DataAtom: a,
			Data:     a.String(),
		})
}

func Text(text string) *Node {
	return (*Node)(
		&html.Node{
			Type: html.TextNode,
			Data: text,
		},
	)
}

type Selector css.Selector

func CssMustParse(s string) *Selector {
	return (*Selector)(css.MustParse(s))
}

func (s *Selector) Select(n *Node) []*Node {
	nodes := (*css.Selector)(s).Select((*html.Node)(n))
	nArray := make([]*Node, len(nodes))
	for i := range len(nodes) {
		nArray[i] = (*Node)(nodes[i])
	}
	return nArray
}

func HtmlParsePage(s io.Reader) (*Node, error) {
	n, err := html.Parse(s)
	return (*Node)(n), err
}

func HtmlParsePageString(s string) (*Node, error) {
	return HtmlParsePage(strings.NewReader(s))
}

func (n *Node) HasRoot(r *Node) bool {
	for p := n.Parent; p != nil; p = p.Parent {
		if p == (*html.Node)(r) {
			return true
		}
	}
	return false
}

func (n *Node) Query(selector string) func(yield func(c *Node) bool) {
	s := CssMustParse(selector)
	return func(yield func(c *Node) bool) {
		for _, e := range s.Select(n) {
			if !yield((*Node)(e)) {
				return
			}
		}
	}
}

func (n *Node) InputText(selector string) func(yield func(c *Node) bool) {
	sel := CssMustParse(selector)
	return func(yield func(c *Node) bool) {
		for _, e := range sel.Select(n) {
			if e.DataAtom == atom.Input && e.HasAttrValueLower("type", "text") {
				if !yield((*Node)(e)) {
					return
				}
			}
		}
	}
}

func (n *Node) HasAttrValueLower(key, val string) bool {
	for _, a := range n.Attr {
		if a.Key == lower(key) && lower(a.Val) == lower(val) {
			return true
		}
	}
	return false
}

func (n *Node) Render(w io.Writer) error {
	return html.Render(w, (*html.Node)(n))
}

func TypeString(n *html.Node) string {
	switch n.Type {
	case html.ErrorNode:
		return "Error"
	case html.TextNode:
		return "Text"
	case html.DocumentNode:
		return "Document"
	case html.ElementNode:
		return "Element"
	case html.CommentNode:
		return "Comment"
	case html.DoctypeNode:
		return "Doctype"
	case html.RawNode:
		return "Raw"
	}
	return "NoType"
}

func DumpNode(n *html.Node, indent int, mark string) {
	if n == nil {
		return
	}
	fmt.Printf("T55: %s%*s", mark, indent, "")
	fmt.Println(TypeString(n), ">"+strings.ReplaceAll(n.Data, "\n", "<LF>")+"<")
	DumpNode(n.FirstChild, indent+2, "C")
	DumpNode(n.NextSibling, indent, "S")
}
