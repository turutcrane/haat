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

type Element html.Node
type Text html.Node
type Doctype html.Node
type Comment html.Node

type Node interface {
	NodeType() html.NodeType
}

type Attribute html.Attribute

// NodeType returns the node type.
func (e *Element) NodeType() html.NodeType {
	return e.Type
}

func (t *Text) NodeType() html.NodeType {
	return t.Type
}

func (d *Doctype) NodeType() html.NodeType {
	return d.Type
}

func (c *Comment) NodeType() html.NodeType {
	return c.Type
}

// ClearContents removes all children of the node.
func (n *Element) ClearContents() *Element {
	bn := (*html.Node)(n)
	for c := bn.FirstChild; c != nil; c = bn.FirstChild {
		bn.RemoveChild(c)
	}
	return n
}

// C sets the children of the node to the given nodes.
func (n *Element) C(childs ...Node) *Element {
	n.ClearContents()
	n.AppendC(childs...)
	return n
}

func convertNode(n Node) *html.Node {
	var n0 *html.Node
	switch c := n.(type) {
	case *Element:
		n0 = (*html.Node)(c)
	case *Text:
		n0 = (*html.Node)(c)
	}

	return n0
}

// AppendC appends the given nodes to the children of the node.
func (e *Element) AppendC(childs ...Node) *Element {
	for _, c := range childs {
		(*html.Node)(e).AppendChild(convertNode(c))
	}
	return e
}

// Remove removes the node from the parent node.
func Remove(n Node) {
	n0 := convertNode(n)
	if n0.Parent == nil {
		return
	}
	n0.Parent.RemoveChild(n0)
}

// RemoveAttr removes the attribute with the given key from the node.
func (e *Element) RemoveAttr(key string) *Element {
	attr := make([]html.Attribute, 0, len(e.Attr))
	for _, a := range e.Attr {
		if a.Key != key {
			attr = append(attr, a)
		}
	}
	e.Attr = attr
	return e
}

// RemoveClass removes the class from the class attribute of the node.
func (e *Element) RemoveClass(class string) *Element {
	classes := strings.Split(e.GetAttr("class"), " ")
	cs := make([]string, 0, len(classes))
	for _, c := range classes {
		if c != class {
			cs = append(cs, c)
		}
	}
	e.SetA(Attr("class", strings.Join(cs, " ")))
	return e
}

// Clone node
func nodeClone(n *html.Node) *html.Node {
	m := &html.Node{
		Type:      n.Type,
		DataAtom:  n.DataAtom,
		Data:      n.Data,
		Namespace: n.Namespace,
		Attr:      make([]html.Attribute, len(n.Attr)),
	}
	copy(m.Attr, n.Attr)
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		d := nodeClone(c)
		m.AppendChild(d)
	}
	return m
}

// Clone creates a deep copy of the node.
func (e *Element) Clone() *Element {
	return (*Element)(nodeClone((*html.Node)(e)))
}

func (t *Text) Clone() *Text {
	return &Text{
		Type: t.Type,
		Data: t.Data,
	}
}

func (d *Doctype) Clone() *Doctype {
	newDoctype :=  &Doctype{
		Type:      d.Type,
		Data:      d.Data,
		Namespace: d.Namespace,
		Attr:      make([]html.Attribute, len(d.Attr)),
	}
	copy(newDoctype.Attr, d.Attr)
	return newDoctype
}

func (c *Comment) Clone() *Comment {
	return &Comment{
		Type: c.Type,
		Data: c.Data,
	}
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
func (n *Element) A(attrs ...Attribute) *Element {
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
func (n *Element) SetA(attr ...Attribute) *Element {
	var attrs []Attribute
	for _, a := range n.Attr {
		attrs = append(attrs, Attribute(a))
	}
	return n.A(slices.Concat(attrs, attr)...)
}

// AppendA appends the given attributes to the attributes of the node.
// No attribute duplication check is performed.
func (n *Element) AppendA(attr Attribute) *Element {
	n.Attr = append(n.Attr, html.Attribute{Key: attr.Key, Val: attr.Val})
	return n
}

// AttrHref creates a new attribute with the key "href" and the value of the given URL.
func AttrHref(u url.URL) Attribute {
	return Attr("href", u.String())
}

// AttrID creates a new attribute with the key "id" and the given value.
func AttrID(id string) Attribute {
	return Attr("id", id)
}

// SetClasses setsthe given class to the class attribute of the node.
// No class duplication check is performed.
func (n *Element) SetClasses(class ...string) *Element {
	old := n.GetAttr("class")
	if old == "" {
		return n.SetA(Attr("class", strings.Join(class, " ")))
	}
	return n.SetA(Attr("class", strings.Join(slices.Concat(strings.Split(old, " "), class), " ")))
}

// JsLetExpr creates a JavaScript let statement with the given name and value.
func JsLetExpr(name, val string) *Element {
	return T("let " + template.JSEscapeString(name) + " = \"" + template.JSEscapeString(val) + "\";")
}

// JsConstExpr creates a JavaScript const statement with the given name and value.
func JsConstExpr(name, val string) *Element {
	return T("const " + template.JSEscapeString(name) + " = \"" + template.JSEscapeString(val) + "\";")
}

// Lf creates a new text node with a line feed.
func Lf() *Element {
	return T("\n")
}

// E creates a new element node with the given atom.
func E(a atom.Atom) *Element {
	return (*Element)(
		&html.Node{
			Type:     html.ElementNode,
			DataAtom: a,
			Data:     a.String(),
		})
}

// T creates a new text node with the given text.
func T(text string) *Element {
	return (*Element)(
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
func (s *Selector) Select(n *Element) []*Element {
	nodes := (*css.Selector)(s).Select((*html.Node)(n))
	nArray := make([]*Element, len(nodes))
	for i := range len(nodes) {
		nArray[i] = (*Element)(nodes[i])
	}
	return nArray
}

// HtmlParsePage parses the HTML page from the given reader.
func HtmlParsePage(s io.Reader) (*Element, error) {
	n, err := html.Parse(s)
	return (*Element)(n), err
}

// HtmlParsePageString parses the HTML page from the given string.
func HtmlParsePageString(s string) (*Element, error) {
	return HtmlParsePage(strings.NewReader(s))
}

// HtmlParseFragment parses the HTML fragment with node context from the given reader.
func HtmlParseFragment(s io.Reader, node *Element) ([]*Element, error) {
	n, err := html.ParseFragment(s, (*html.Node)(node))

	nodes := make([]*Element, len(n))
	for i := range len(n) {
		nodes[i] = (*Element)(n[i])
	}
	return nodes, err
}

// HtmlParseFragmentString parses the HTML fragment with node context from the given string.
func HtmlParseFragmentString(s string, node *Element) ([]*Element, error) {
	return HtmlParseFragment(strings.NewReader(s), node)
}

// HasRoot returns true if the node has the given root node as an ancestor.
func (n *Element) HasRoot(root *Element) bool {
	for p := n.Parent; p != nil; p = p.Parent {
		if p == (*html.Node)(root) {
			return true
		}
	}
	return false
}

// Query return iter.Seq[*Element] delivers the nodes that match the selector.
func (n *Element) Query(selector string) func(yield func(c *Element) bool) {
	s := CssMustParse(selector)
	return func(yield func(c *Element) bool) {
		for _, e := range s.Select(n) {
			if !yield((*Element)(e)) {
				return
			}
		}
	}
}

// InputText return iter.Seq[*Element] delivers the input text nodes that match the selector.
func (n *Element) InputText(selector string) func(yield func(c *Element) bool) {
	sel := CssMustParse(selector)
	return func(yield func(c *Element) bool) {
		for _, e := range sel.Select(n) {
			if e.DataAtom == atom.Input && e.HasAttrValueLower("type", "text") {
				if !yield((*Element)(e)) {
					return
				}
			}
		}
	}
}

// HasAttrValue returns true if the node has an attribute with the given key and value.
func (n *Element) HasAttrValueLower(key, val string) bool {
	for _, a := range n.Attr {
		if a.Key == lower(key) && lower(a.Val) == lower(val) {
			return true
		}
	}
	return false
}

// Render renders the node to the given writer.
func (n *Element) Render(w io.Writer, checker ...Checker) error {
	for _, c := range checker {
		if err := c(n); err != nil {
			return err
		}
	}
	return html.Render(w, (*html.Node)(n))
}

// GetAttr returns the value of the attribute with the given key.
func (n *Element) GetAttr(key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// ID returns the value of the id attribute.
func (n *Element) ID() string {
	for _, a := range n.Attr {
		if a.Key == "id" {
			return a.Val
		}
	}
	return ""
}

// Checker is a function that checks the node.
type Checker func(*Element) error

// IDDuplicateCheck checks if the node has duplicate id attributes.
func IDDuplicateCheck(n *Element) error {
	ids := map[string]struct{}{}
	for e := range n.Query("[id]") {
		id := e.ID()
		if _, ok := ids[id]; ok {
			return fmt.Errorf("duplicate id: %s", id)
		}
		ids[id] = struct{}{}
	}
	return nil
}

// IDMissingCheck checks if the node has id without value.
func IDMissingCheck(n *Element) error {
	for e := range n.Query("[id]") {
		if e.ID() == "" {
			return fmt.Errorf("missing id")
		}
	}
	return nil
}

// IDHasBlankCheck checks if the node has id with blank.
func IDHasBlankCheck(n *Element) error {
	for e := range n.Query("[id]") {
		id := e.ID()
		if strings.Contains(id, " ") {
			return fmt.Errorf("id has blank: %s", id)
		}
	}
	return nil
}

// TypeString returns the string representation of the node type.
func TypeString(n *Element) string {
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
func DumpNode(n *Element, indent int, mark string) {
	if n == nil {
		return
	}
	fmt.Printf("T55: %s%*s", mark, indent, "")
	fmt.Println(TypeString(n), ">"+strings.ReplaceAll(n.Data, "\n", "<LF>")+"<")
	DumpNode((*Element)(n.FirstChild), indent+2, "C")
	DumpNode((*Element)(n.NextSibling), indent, "S")
}
