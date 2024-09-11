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

// ClearContents removes all children of the node.
func (n *Node) ClearContents() *Node {
	bn := (*html.Node)(n)
	for c := bn.FirstChild; c != nil; c = bn.FirstChild {
		bn.RemoveChild(c)
	}
	return n
}

// C sets the children of the node to the given nodes.
func (n *Node) C(childs ...*Node) *Node {
	n.ClearContents()
	n.AppendC(childs...)
	return n
}

// AppendC appends the given nodes to the children of the node.
func (n *Node) AppendC(childs ...*Node) *Node {
	for _, c := range childs {
		(*html.Node)(n).AppendChild((*html.Node)(c))
	}
	return n
}

// Attr creates a new attribute with the given key and value.
func Attr(key, value string) Attribute {
	return (Attribute)(html.Attribute{
		Key: lower(key),
		Val: value,
	})
}

// A replaces all attributes with the attributes specified in the arguments. 
// If there are duplicate keys, it sets the latter value.
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

// SetA appends the given attributes to the attributes of the node. 
// If keys are already in the attributes of the node, the values are overwritten.
func (n *Node) SetA(attr ...Attribute) *Node {
	var attrs []Attribute
	for _, a := range n.Attr {
		attrs = append(attrs, Attribute(a))
	}
	return n.A(slices.Concat(attrs, attr)...)
}

// AppendA appends the given attributes to the attributes of the node.
func (n *Node) AppendA(attr Attribute) *Node {
	n.Attr = append(n.Attr, html.Attribute{Key: attr.Key, Val: attr.Val})
	return n
}

// AttrHref creates a new attribute with the key "href" and the value of the given URL.
func AttrHref(u url.URL) Attribute {
	return Attr("href", u.String())
}

func AttrID(id string) Attribute {
	return Attr("id", id)
}

func (n *Node) SetClass(class string) *Node {
	c := n.GetAttr("class")
	return n.SetA(Attr("class", strings.Join([]string{c, class}, " ")))
}

// JsLetString creates a JavaScript let statement with the given name and value.
func JsLetString(name, val string) string {
	return "let " + template.JSEscapeString(name) + " = \"" + template.JSEscapeString(val) + "\";"
}

// JsConstString creates a JavaScript const statement with the given name and value.
func JsConstString(name, val string) string {
	return "const " + template.JSEscapeString(name) + " = \"" + template.JSEscapeString(val) + "\";"
}

// Element creates a new element node with the given atom.
func Element(a atom.Atom) *Node {
	return (*Node)(
		&html.Node{
			Type:     html.ElementNode,
			DataAtom: a,
			Data:     a.String(),
		})
}

// Text creates a new text node with the given text.
func Text(text string) *Node {
	return (*Node)(
		&html.Node{
			Type: html.TextNode,
			Data: text,
		},
	)
}

type Selector css.Selector

// CssMustParse parses the given CSS selector and panics if it fails.
func CssMustParse(s string) *Selector {
	return (*Selector)(css.MustParse(s))
}

// Select selects the nodes that match the selector.
func (s *Selector) Select(n *Node) []*Node {
	nodes := (*css.Selector)(s).Select((*html.Node)(n))
	nArray := make([]*Node, len(nodes))
	for i := range len(nodes) {
		nArray[i] = (*Node)(nodes[i])
	}
	return nArray
}

// HtmlParsePage parses the HTML page from the given reader.
func HtmlParsePage(s io.Reader) (*Node, error) {
	n, err := html.Parse(s)
	return (*Node)(n), err
}

// HtmlParsePageString parses the HTML page from the given string.
func HtmlParsePageString(s string) (*Node, error) {
	return HtmlParsePage(strings.NewReader(s))
}

// HtmlParseFragment parses the HTML fragment with node context from the given reader.
func HtmlParseFragment(s io.Reader, node *Node) (*Node, error) {
	n, err := html.ParseFragment(s, (*html.Node)(node))
	return (*Node)(n[0]), err
}

// HtmlParseFragmentString parses the HTML fragment with node context from the given string.
func HtmlParseFragmentString(s string, node *Node) (*Node, error) {
	return HtmlParseFragment(strings.NewReader(s), node)
}

// HasRoot returns true if the node has the given root node as an ancestor.
func (n *Node) HasRoot(r *Node) bool {
	for p := n.Parent; p != nil; p = p.Parent {
		if p == (*html.Node)(r) {
			return true
		}
	}
	return false
}

// Query return iter.Seq[*Node] delivers the nodes that match the selector.
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

// InputText return iter.Seq[*Node] delivers the input text nodes that match the selector.
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

// HasAttrValue returns true if the node has an attribute with the given key and value.
func (n *Node) HasAttrValueLower(key, val string) bool {
	for _, a := range n.Attr {
		if a.Key == lower(key) && lower(a.Val) == lower(val) {
			return true
		}
	}
	return false
}

// Render renders the node to the given writer.
func (n *Node) Render(w io.Writer, checker ...Checker) error {
	for _, c := range checker {
		if err := c(n); err != nil {
			return err
		}
	}
	return html.Render(w, (*html.Node)(n))
}

func (n *Node) GetAttr(key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

func (n *Node) ID() string {
	for _, a := range n.Attr {
		if a.Key == "id" {
			return a.Val
		}
	}
	return ""
}

// Checker is a function that checks the node.
type Checker func(*Node) error

func IDDuplicateCheck(n *Node) error {
	ids := map[string]struct{}{}
	for e := range n.Query("[id]") {
		id := e.ID()
		if _, ok := ids[id]; ok {
			return fmt.Errorf("duplicate id: %s", id)
		}
		ids[id]=struct{}{}
	}
	return nil
}

func IDMissingCheck(n *Node) error {
	for e := range n.Query("[id]") {
		if e.ID() == "" {
			return fmt.Errorf("missing id")
		}
	}
	return nil
}

func IDHasBlankCheck(n *Node) error {
	for e := range n.Query("[id]") {
		id := e.ID()
		if strings.Contains(id, " ") {
			return fmt.Errorf("id has blank: %s", id)
		}
	}
	return nil
}

// TypeString returns the string representation of the node type.
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

// DumpNode prints the node to the standard output.
func DumpNode(n *html.Node, indent int, mark string) {
	if n == nil {
		return
	}
	fmt.Printf("T55: %s%*s", mark, indent, "")
	fmt.Println(TypeString(n), ">"+strings.ReplaceAll(n.Data, "\n", "<LF>")+"<")
	DumpNode(n.FirstChild, indent+2, "C")
	DumpNode(n.NextSibling, indent, "S")
}
