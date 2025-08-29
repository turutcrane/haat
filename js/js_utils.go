package js

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"strings"
	"unicode"
)

func makeDeclaration(decl, funcName, name string, val any, format string) string {
	if !IsIdentifire(name) {
		name = "</script> add By " + funcName + ": " + html.EscapeString(name) // cause ECMAscript syntax error
	}
	return fmt.Sprintf("%s %s = "+format+";", decl, name, val)
}

// LetString creates a JavaScript let statement with the given name and string value.
func LetString(name, val string) string {
	return makeDeclaration("let", "js.LetString", name, template.JSEscapeString(val), `"%s"`)
}

// LetInt creates a JavaScript let statement with the given name and int value.
func LetInt(name string, val int) string {
	return makeDeclaration("let", "js.LetInt", name, val, `%d`)
}

// LetBool creates a JavaScript let statement with the given name and bool value.
func LetBool(name string, val bool) string {
	return makeDeclaration("let", "js.LetBool", name, val, `%t`)
}

func makeJsonDeclaration(decl, name string, val any) (string, error) {
	if !IsIdentifire(name) {
		return "", fmt.Errorf("invalid identifire: %s", name)
	}
	// val を json.Marshal で文字列に変換する
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s %s = %s;", decl, name, jsonBytes), nil
}

// LetJson creates a JavaScript let statement with the given name and JSON value.
func LetJson(name string, val any) (string, error) {
	return makeJsonDeclaration("let", name, val)
}

// ConstString creates a JavaScript const statement with the given name and string value.
func ConstString(name, val string) string {
	return makeDeclaration("const", "js.ConstString", name, template.JSEscapeString(val), `"%s"`)
}

// ConstInt creates a JavaScript const statement with the given name and int value.
func ConstInt(name string, val int) string {
	return makeDeclaration("const", "js.ConstInt", name, val, `%d`)
}

// ConstBool creates a JavaScript const statement with the given name and bool value.
func ConstBool(name string, val bool) string {
	return makeDeclaration("const", "js.ConstBool", name, val, `%t`)
}

// ConstJson creates a JavaScript const statement with the given name and JSON value.
func ConstJson(name string, val any) (string, error) {
	return makeJsonDeclaration("const", name, val)
}

// check name as ECMAScript IsIdentifire
func IsIdentifire(name string) bool {
	// name を rune 配列に変換する
	runes := []rune(name)
	if len(runes) == 0 {
		return false
	}

	if !IsIDStart(runes[0]) {
		return false
	}
	for _, c := range runes[1:] {
		if !IsIDContinue(c) {
			return false
		}
	}
	return true
}

func IsIDStart(c rune) bool {
	return c == rune('_') || c == rune('$') ||
		((unicode.Is(unicode.L, c) || unicode.Is(unicode.Nl, c) || unicode.Is(unicode.Other_ID_Start, c)) &&
			!unicode.Is(unicode.Pattern_Syntax, c) &&
			!unicode.Is(unicode.Pattern_White_Space, c))

}

func IsIDContinue(c rune) bool {
	return c == rune('$') ||
		((unicode.Is(unicode.L, c) ||
			unicode.Is(unicode.Nl, c) ||
			unicode.Is(unicode.Other_ID_Start, c) ||
			unicode.Is(unicode.Mn, c) ||
			unicode.Is(unicode.Mc, c) ||
			unicode.Is(unicode.Nd, c) ||
			unicode.Is(unicode.Pc, c) ||
			unicode.Is(unicode.Other_ID_Continue, c) ||
			c == 0x0200c || c == 0x0200d || // ZERO WIDTH NON-JOINER..ZERO WIDTH JOINER
			c == 0x030fb || c == 0xff65) && // KATAKANA MIDDLE DOT, HALFWIDTH KATAKANA MIDDLE DOT
			!unicode.Is(unicode.Pattern_Syntax, c) &&
			!unicode.Is(unicode.Pattern_White_Space, c))
}

func isPropertyAccess(name string) bool {
	for _, id := range strings.Split(name, ".") {
		if !IsIdentifire(id) {
			return false
		}
	}
	return true
}

func makeAssignment(funcName, name string, val any, format string) string {
	if !isPropertyAccess(name) {
		name = "</script> add By " + funcName + ": " + html.EscapeString(name) // cause ECMAscript syntax error
	}
	return fmt.Sprintf("%s = "+format+";", name, val)
}

// AssignString creates a JavaScript assignment statement with the given name and string value.
func AssignString(name, val string) string {
	return makeAssignment("js.AssignString", name, template.JSEscapeString(val), `"%s"`)
}

// AssignInt creates a JavaScript assignment statement with the given name and int value.
func AssignInt(name string, val int) string {
	return makeAssignment("js.AssignInt", name, val, `%d`)
}

// AssignBool creates a JavaScript assignment statement with the given name and bool value.
func AssignBool(name string, val bool) string {
	return makeAssignment("js.AssignBool", name, val, `%t`)
}

func makeJsonAssignment(name string, val any) (string, error) {
	if !isPropertyAccess(name) {
		return "", fmt.Errorf("invalid identifire: %s", name)
	}
	jsonBytes, err := json.Marshal(val)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s = %s;", name, jsonBytes), nil
}

// AssignJson creates a JavaScript assignment statement with the given name and JSON value.
func AssignJson(name string, val any) (string, error) {
	return makeJsonAssignment(name, val)
}
