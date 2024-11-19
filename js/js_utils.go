package js

import (
	"encoding/json"
	"fmt"
	"html/template"
)

// LetString creates a JavaScript let statement with the given name and string value.
func LetString(name, val string) string {
	return fmt.Sprintf("let %s = \"%s\";", template.JSEscapeString(name), template.JSEscapeString(val))
}

// LetInt creates a JavaScript let statement with the given name and int value.
func LetInt(name string, val int) string {
	return fmt.Sprintf("let %s = %d;", template.JSEscapeString(name), val)
}

// LetBool creates a JavaScript let statement with the given name and bool value.
func LetBool(name string, val bool) string {
	return fmt.Sprintf("let %s = %t;", template.JSEscapeString(name), val)
}

// LetJson creates a JavaScript let statement with the given name and JSON value.
func LetJson(name string, val any) (string, error) {
	// val を json.Marshal で文字列に変換する
	json, err := json.Marshal(val)
	if err == nil {
		return fmt.Sprintf("let %s = %s;", template.JSEscapeString(name), json), nil
	}

	return "", err
}

// ConstString creates a JavaScript const statement with the given name and string value.
func ConstString(name, val string) string {
	return fmt.Sprintf("const %s = \"%s\";", template.JSEscapeString(name), template.JSEscapeString(val))
}

// ConstInt creates a JavaScript const statement with the given name and int value.
func ConstInt(name string, val int) string {
	return fmt.Sprintf("const %s = %d;", template.JSEscapeString(name), val)
}

// ConstBool creates a JavaScript const statement with the given name and bool value.
func ConstBool(name string, val bool) string {
	return fmt.Sprintf("const %s = %t;", template.JSEscapeString(name), val)
}

// ConstJson creates a JavaScript const statement with the given name and JSON value.
func ConstJson(name string, val any) (string, error) {
	// val を json.Marshal で文字列に変換する
	json, err := json.Marshal(val)
	if err == nil {
		return fmt.Sprintf("const %s = %s;", template.JSEscapeString(name), json), nil
	}
	return "", err
}
