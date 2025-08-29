package js

import (
	"encoding/json"
	"fmt"
	"html"
	"html/template"
	"unicode"
)

// LetString creates a JavaScript let statement with the given name and string value.
func LetString(name, val string) string {
	if !IsIdentifire(name) {
		name = "</script> add By js.LetString: " + html.EscapeString(name) // cause ECMAscript syntax error
	}
	return fmt.Sprintf("let %s = \"%s\";", name, template.JSEscapeString(val))
}

// LetInt creates a JavaScript let statement with the given name and int value.
func LetInt(name string, val int) string {
	if !IsIdentifire(name) {
		name = "</script> add By js.LetInt: " + html.EscapeString(name) // cause ECMAscript syntax error
	}
	return fmt.Sprintf("let %s = %d;", name, val)
}

// LetBool creates a JavaScript let statement with the given name and bool value.
func LetBool(name string, val bool) string {
	if !IsIdentifire(name) {
		name = "</script> add By js.LetBool: " + html.EscapeString(name) // cause ECMAscript syntax error
	}
	return fmt.Sprintf("let %s = %t;", name, val)
}

// LetJson creates a JavaScript let statement with the given name and JSON value.
func LetJson(name string, val any) (string, error) {
	if !IsIdentifire(name) {
		return "", fmt.Errorf("invalid identifire: %s", name)
	}
	// val を json.Marshal で文字列に変換する
	json, err := json.Marshal(val)
	if err == nil {
		return fmt.Sprintf("let %s = %s;", name, json), nil
	}

	return "", err
}

// ConstString creates a JavaScript const statement with the given name and string value.
func ConstString(name, val string) string {
	if !IsIdentifire(name) {
		name = "</script> add By js.ConstString: " + html.EscapeString(name) // cause ECMAscript syntax error
	}
	return fmt.Sprintf("const %s = \"%s\";", name, template.JSEscapeString(val))
}

// ConstInt creates a JavaScript const statement with the given name and int value.
func ConstInt(name string, val int) string {
	if !IsIdentifire(name) {
		name = "</script> add By js.ConstInt: " + html.EscapeString(name) // cause ECMAscript syntax error
	}
	return fmt.Sprintf("const %s = %d;", name, val)
}

// ConstBool creates a JavaScript const statement with the given name and bool value.
func ConstBool(name string, val bool) string {
	if !IsIdentifire(name) {
		name = "</script> add By js.ConstBool: " + html.EscapeString(name) // cause ECMAscript syntax error
	}
	return fmt.Sprintf("const %s = %t;", name, val)
}

// ConstJson creates a JavaScript const statement with the given name and JSON value.
func ConstJson(name string, val any) (string, error) {
	if !IsIdentifire(name) {
		return "", fmt.Errorf("invalid identifire: %s", name)
	}
	// val を json.Marshal で文字列に変換する
	json, err := json.Marshal(val)
	if err == nil {
		return fmt.Sprintf("const %s = %s;", name, json), nil
	}
	return "", err
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
