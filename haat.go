package haat

import (
	"encoding/json"
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

type Document html.Node
type Element html.Node
type Text html.Node
type Doctype html.Node
type Comment html.Node

type Node interface {
	NodeType() html.NodeType
}

type ElementChild interface {
	ParentElement() *Element
}

type Attribute html.Attribute

// NodeType returns the node type.
func (d *Document) NodeType() html.NodeType {
	return d.Type
}

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

// ParentElement returns the parent element of the node.
func (e *Element) ParentElement() *Element {
	return (*Element)(e.Parent)
}

func (t *Text) ParentElement() *Element {
	return (*Element)(t.Parent)
}

func (c *Comment) ParentElement() *Element {
	return (*Element)(c.Parent)
}

// ClearContents removes all children of the node.
func (e *Element) ClearContents() *Element {
	bn := (*html.Node)(e)
	for c := bn.FirstChild; c != nil; c = bn.FirstChild {
		bn.RemoveChild(c)
	}
	return e
}

// C sets the children of the node to the given nodes.
func (e *Element) C(childs ...ElementChild) *Element {
	e.ClearContents()
	e.AppendC(childs...)
	return e
}

func convertNode(n Node) *html.Node {
	var n0 *html.Node
	switch c := n.(type) {
	case *Document:
		n0 = (*html.Node)(c)
	case *Element:
		n0 = (*html.Node)(c)
	case *Text:
		n0 = (*html.Node)(c)
	case *Doctype:
		n0 = (*html.Node)(c)
	case *Comment:
		n0 = (*html.Node)(c)
	}

	return n0
}

// AppendC appends the given nodes to the children of the node.
func (e *Element) AppendC(childs ...ElementChild) *Element {
	for _, c := range childs {
		n := c.(Node)
		(*html.Node)(e).AppendChild(convertNode(n))
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
		if a.Key != lower(key) {
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
func (d *Document) Clone() *Document {
	return (*Document)(nodeClone((*html.Node)(d)))
}

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
	newDoctype := &Doctype{
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
func (e *Element) A(attrs ...Attribute) *Element {
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

	e.Attr = newAttrs
	return e
}

// SetA appends the given attributes to the attributes of the node.
// If keys are already in the attributes of the node, the values are overwritten.
func (e *Element) SetA(attr ...Attribute) *Element {
	var attrs []Attribute
	for _, a := range e.Attr {
		attrs = append(attrs, Attribute(a))
	}
	return e.A(slices.Concat(attrs, attr)...)
}

// SetBoolA appends the given boolean attributes to the attributes of the node or removes them.
func (e *Element) SetBoolA(key string, v bool) *Element {
	if v {
		return e.SetA(Attr(key, ""))
	}
	return e.RemoveAttr(key)
}

// AttrHref creates a new attribute with the key "href" and the value of the given URL.
func AttrHref(u url.URL) Attribute {
	return Attr("href", u.String())
}

// AttrID creates a new attribute with the key "id" and the given value.
func AttrID(id string) Attribute {
	return Attr("id", id)
}

// SetClasses sets the given class to the class attribute of the node.
// Class duplication check is performed.
func (e *Element) SetClasses(classes ...string) *Element {
	old := e.GetAttr("class")
	if old == "" {
		return e.SetA(Attr("class", strings.Join(classes, " ")))
	}
	newClasses := slices.Concat(strings.Split(old, " "), classes)
	slices.Sort(newClasses)
	return e.SetA(Attr("class", strings.Join(slices.Compact(newClasses), " ")))
}

// JsLetString creates a JavaScript let statement with the given name and string value.
func JsLetString(name, val string) *Text {
	return T(fmt.Sprintf("let %s = \"%s\";", template.JSEscapeString(name), template.JSEscapeString(val)))
}

// JsLetInt creates a JavaScript let statement with the given name and int value.
func JsLetInt(name string, val int) *Text {
	return T(fmt.Sprintf("let %s = %d;", template.JSEscapeString(name), val))
}

// JsLetBool creates a JavaScript let statement with the given name and bool value.
func JsLetBool(name string, val bool) *Text {
	return T(fmt.Sprintf("let %s = %t;", template.JSEscapeString(name), val))
}

// JsLetJson creates a JavaScript let statement with the given name and JSON value.
func JsLetJson(name string, val any) (*Text, error) {
	// val を json.Marshal で文字列に変換する
	json, err := json.Marshal(val)
	if err == nil {
		return T(fmt.Sprintf("let %s = %s;", template.JSEscapeString(name), json)), nil
	}

	return nil, err
}

// JsConstString creates a JavaScript const statement with the given name and string value.
func JsConstString(name, val string) *Text {
	return T(fmt.Sprintf("const %s = \"%s\";", template.JSEscapeString(name), template.JSEscapeString(val)))
}

// JsConstInt creates a JavaScript const statement with the given name and int value.
func JsConstInt(name string, val int) *Text {
	return T(fmt.Sprintf("const %s = %d;", template.JSEscapeString(name), val))
}

// JsConstBool creates a JavaScript const statement with the given name and bool value.
func JsConstBool(name string, val bool) *Text {
	return T(fmt.Sprintf("const %s = %t;", template.JSEscapeString(name), val))
}

// JsConstJson creates a JavaScript const statement with the given name and JSON value.
func JsConstJson(name string, val any) (*Text, error) {
	// val を json.Marshal で文字列に変換する
	json, err := json.Marshal(val)
	if err == nil {
		return T(fmt.Sprintf("const %s = %s;", template.JSEscapeString(name), json)), nil
	}
	return nil, err
}

// Lf creates a new text node with a line feed.
func Lf() *Text {
	return T("\n")
}

// E creates a new element node with the given atom.
func E(a atom.Atom) *Element {
	return &Element{
		Type:     html.ElementNode,
		DataAtom: a,
		Data:     a.String(),
	}
}

// T creates a new text node with the given text.
func T(text string) *Text {
	return &Text{
		Type: html.TextNode,
		Data: text,
	}
}

type Selector css.Selector

// SelectorMustParse parses the given CSS selector and panics if it fails.
func SelectorMustParse(s string) *Selector {
	return (*Selector)(css.MustParse(s))
}

// SelectorParse parse the givent CSS selector
func SelectorParse(s string) (*Selector, error) {
	sel, err := css.Parse(s)
	return (*Selector)(sel), err
}

// ParseHtml parses the HTML page from the given reader.
func ParseHtml(s io.Reader) (*Document, error) {
	n, err := html.Parse(s)
	return (*Document)(n), err
}

// ParseHtmlFragment parses the HTML fragment with node context from the given reader.
func ParseHtmlFragment(s io.Reader, node *Element) ([]*Element, error) {
	n, err := html.ParseFragment(s, (*html.Node)(node))

	nodes := make([]*Element, len(n))
	for i := range len(n) {
		nodes[i] = (*Element)(n[i])
	}
	return nodes, err
}

// HasRoot returns true if the node has the given root node as parent.
func (e *Element) HasRoot(root *Element) bool {
	for p := e.Parent; p != nil; p = p.Parent {
		if p == (*html.Node)(root) {
			return true
		}
	}
	return false
}

func queryNode(n *html.Node, selector string) []*Element {

	return queryNodeBySelector(n, (*css.Selector)(SelectorMustParse(selector)))
}

func queryNodeBySelector(n *html.Node, sel *css.Selector) []*Element {
	nodes := sel.Select(n)
	elements := make([]*Element, len(nodes))
	for i, e := range nodes {
		elements[i] = (*Element)(e)
	}
	return elements
}

// Query return the nodes that match the selector.
func (d *Document) Query(selector string) []*Element {
	return queryNode((*html.Node)(d), selector)
}

func (e *Element) Query(selector string) []*Element {
	return queryNode((*html.Node)(e), selector)
}

// QeurySelector return the all nodes that match the selector.
func (d *Document) QuerySelector(selector *Selector) []*Element {
	return queryNodeBySelector((*html.Node)(d), (*css.Selector)(selector))
}

func (e *Element) QuerySelector(selector *Selector) []*Element {
	return queryNodeBySelector((*html.Node)(e), (*css.Selector)(selector))
}

func inputText(n *html.Node, selector string) []*Element {
	elements := make([]*Element, 0)
	for _, i := range queryNode(n, selector) {
		if i.DataAtom == atom.Input && i.HasAttrValueLower("type", "text") {
			elements = append(elements, i)
		}
	}
	return elements
}

// InputText return the input text nodes that match the selector.
func (d *Document) InputText(selector string) []*Element {
	return inputText((*html.Node)(d), selector)
}

func (e *Element) InputText(selector string) []*Element {
	return inputText((*html.Node)(e), selector)
}

// HasAttrValue returns true if the node has an attribute with the given key and value.
func (e *Element) HasAttrValueLower(key, val string) bool {
	for _, a := range e.Attr {
		if a.Key == lower(key) && lower(a.Val) == lower(val) {
			return true
		}
	}
	return false
}

// Render renders the node to the given writer.
func (d *Document) Render(w io.Writer, checker ...Checker) error {
	for _, html := range d.Query("html") {
		for _, c := range checker {
			if err := c(html); err != nil {
				return err
			}
		}
	}
	return html.Render(w, (*html.Node)(d))
}

func (e *Element) Render(w io.Writer, checker ...Checker) error {
	for _, c := range checker {
		if err := c(e); err != nil {
			return err
		}
	}
	return html.Render(w, (*html.Node)(e))
}

// GetAttr returns the value of the attribute with the given key.
func (e *Element) GetAttr(key string) string {
	for _, a := range e.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// ID returns the value of the id attribute.
func (e *Element) ID() string {
	for _, a := range e.Attr {
		if a.Key == "id" {
			return a.Val
		}
	}
	return ""
}

// Checker is a function that checks the node.
type Checker func(*Element) error

// IDDuplicateCheck checks if the node has duplicate id attributes.
func IDDuplicateCheck(e *Element) error {
	ids := map[string]struct{}{}
	for _, e := range e.Query("[id]") {
		id := e.ID()
		if _, ok := ids[id]; ok {
			return fmt.Errorf("duplicate id: %s", id)
		}
		ids[id] = struct{}{}
	}
	return nil
}

// IDMissingCheck checks if the node has id without value.
func IDMissingCheck(e *Element) error {
	for _, e := range e.Query("[id]") {
		if e.ID() == "" {
			return fmt.Errorf("missing id")
		}
	}
	return nil
}

// IDHasBlankCheck checks if the node has id with blank.
func IDHasBlankCheck(e *Element) error {
	for _, e := range e.Query("[id]") {
		id := e.ID()
		if strings.Contains(id, " ") {
			return fmt.Errorf("id has blank: %s", id)
		}
	}
	return nil
}

// typeString returns the string representation of the node type.
func typeString(t html.NodeType) string {
	switch t {
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

func dumpNode(n *html.Node, indent int, mark string) {
	if n == nil {
		return
	}
	fmt.Printf("T55: %s%*s", mark, indent, "")
	fmt.Println(typeString(n.Type), ">"+strings.ReplaceAll(n.Data, "\n", "<LF>")+"<")
	dumpNode(n.FirstChild, indent+2, "C")
	dumpNode(n.NextSibling, indent, "S")
}

// DumpDocument prints the node to the standard output.
func DumpDocument(d *Document, indent int, mark string) {
	dumpNode((*html.Node)(d), indent, mark)
}

func DumpElement(e *Element, indent int, mark string) {
	dumpNode((*html.Node)(e), indent, mark)
}

func dumpNode2(n *html.Node, indent int, mark string) {
	if n == nil {
		return
	}
	fmt.Printf("T55: %s%*s", mark, indent, "")
	fmt.Println(typeString(n.Type), "<"+strings.ReplaceAll(n.Data, "\n", "<LF>")+">")
	for _, a := range n.Attr {
		fmt.Printf("T468: %s%*s", mark, indent+2, "")
		fmt.Println("Attr", a.Key, a.Val)
	}
	dumpNode2(n.FirstChild, indent+2, "C")
	dumpNode2(n.NextSibling, indent, "S")
}

// DumpNode2 prints the node to the standard output.
func DumpDocument2(d *Document, indent int, mark string) {
	dumpNode2((*html.Node)(d), indent, mark)
}
