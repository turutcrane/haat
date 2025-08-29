package js

import (
	"strings"
	"testing"
)

func TestIdentifire(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want bool
	}{
		{"valid simple", "myVar", true},
		{"valid with underscore start", "_myVar", true},
		{"valid with dollar start", "$myVar", true},
		{"valid with numbers", "myVar123", true},
		{"valid unicode letter", "変数", true},
		{"valid with underscore continue", "my_var", true},
		{"valid with dollar continue", "my$var", true},
		{"invalid starts with number", "1myVar", false},
		{"invalid with hyphen", "my-var", false},
		{"invalid with space", "my var", false},
		{"empty string", "", false},
		{"just a number", "123", false},
		{"keyword (should be true for identifier check)", "let", true},
		{"just underscore", "_", true},
		{"just dollar", "$", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsIdentifire(tt.in); got != tt.want {
				t.Errorf("IsIdentifire(%q) = %v, want %v", tt.in, got, tt.want)
			}
		})
	}
}

func TestLetString(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		val     string
		want    string
	}{
		{"valid name", "myVar", "hello", `let myVar = "hello";`},
		{"value with quotes", "myVar", `he"llo`, `let myVar = "he\"llo";`},
		{"value with slashes", "myVar", `he\llo`, `let myVar = "he\\llo";`},
		{"value with script tag", "myVar", `</script>`, `let myVar = "\u003C/script\u003E";`},
		{"invalid name", "1-invalid", "hello", `let </script> add By js.LetString: 1-invalid = "hello";`},
		{"invalid name with html chars", "1<invalid", "hello", `let </script> add By js.LetString: 1&lt;invalid = "hello";`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LetString(tt.varName, tt.val)
			if got != tt.want {
				t.Errorf("LetString(%q, %q) = %q, want %q", tt.varName, tt.val, got, tt.want)
			}
		})
	}
}

func TestAssignString(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		val     string
		want    string
	}{
		{"simple assignment", "myVar", "value", `myVar = "value";`},
		{"property assignment", "myObj.prop", "value", `myObj.prop = "value";`},
		{"nested property", "a.b.c", "value", `a.b.c = "value";`},
		{"invalid identifier part", "a.1b.c", "value", `</script> add By js.AssignString: a.1b.c = "value";`},
		{"invalid with html chars", "a.b<d.c", "value", `</script> add By js.AssignString: a.b&lt;d.c = "value";`},
		{"value with quotes", "myVar", `he"llo`, `myVar = "he\"llo";`},
		{"value with slashes", "myVar", `he\llo`, `myVar = "he\\llo";`},
		{"value with script tag", "myVar", `</script>`, `myVar = "\u003C/script\u003E";`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssignString(tt.varName, tt.val)
			if got != tt.want {
				t.Errorf("AssignString(%q, %q) = %q, want %q", tt.varName, tt.val, got, tt.want)
			}
		})
	}
}

func TestAssignInt(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		val     int
		want    string
	}{
		{"simple assignment", "myVar", 123, `myVar = 123;`},
		{"property assignment", "myObj.prop", -45, `myObj.prop = -45;`},
		{"invalid identifier part", "a.1b.c", 0, `</script> add By js.AssignInt: a.1b.c = 0;`},
		{"invalid with html chars", "a.b<d.c", 99, `</script> add By js.AssignInt: a.b&lt;d.c = 99;`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssignInt(tt.varName, tt.val)
			if got != tt.want {
				t.Errorf("AssignInt(%q, %d) = %q, want %q", tt.varName, tt.val, got, tt.want)
			}
		})
	}
}

func TestAssignBool(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		val     bool
		want    string
	}{
		{"simple assignment true", "myVar", true, `myVar = true;`},
		{"property assignment false", "myObj.prop", false, `myObj.prop = false;`},
		{"invalid identifier part", "a.1b.c", true, `</script> add By js.AssignBool: a.1b.c = true;`},
		{"invalid with html chars", "a.b<d.c", false, `</script> add By js.AssignBool: a.b&lt;d.c = false;`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := AssignBool(tt.varName, tt.val)
			if got != tt.want {
				t.Errorf("AssignBool(%q, %v) = %q, want %q", tt.varName, tt.val, got, tt.want)
			}
		})
	}
}

func TestAssignJson(t *testing.T) {
	type myStruct struct {
		Foo string `json:"foo"`
	}

	tests := []struct {
		name    string
		varName string
		val     any
		want    string
		wantErr bool
		errText string
	}{
		{"valid struct", "myJson", myStruct{Foo: "bar"}, `myJson = {"foo":"bar"};`, false, ""},
		{"property assignment", "myObj.prop", myStruct{Foo: "bar"}, `myObj.prop = {"foo":"bar"};`, false, ""},
		{"invalid identifier", "1-invalid", myStruct{Foo: "bar"}, "", true, "invalid identifire: 1-invalid"},
		{"invalid property part", "a.1b.c", myStruct{Foo: "bar"}, "", true, "invalid identifire: a.1b.c"},
		{"json marshal error", "myJson", make(chan int), "", true, "json: unsupported type: chan int"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AssignJson(tt.varName, tt.val)
			if (err != nil) != tt.wantErr {
				t.Fatalf("AssignJson() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errText) {
					t.Errorf("AssignJson() error = %q, want to contain %q", err.Error(), tt.errText)
				}
			}
			if got != tt.want {
				t.Errorf("AssignJson() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestLetInt(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		val     int
		want    string
	}{
		{"valid name, positive int", "myVar", 123, `let myVar = 123;`},
		{"valid name, negative int", "myVar", -123, `let myVar = -123;`},
		{"valid name, zero", "myVar", 0, `let myVar = 0;`},
		{"invalid name", "1-invalid", 42, `let </script> add By js.LetInt: 1-invalid = 42;`},
		{"invalid name with html chars", "1<invalid", 42, `let </script> add By js.LetInt: 1&lt;invalid = 42;`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LetInt(tt.varName, tt.val)
			if got != tt.want {
				t.Errorf("LetInt(%q, %d) = %q, want %q", tt.varName, tt.val, got, tt.want)
			}
		})
	}
}

func TestLetBool(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		val     bool
		want    string
	}{
		{"valid name, true", "myVar", true, `let myVar = true;`},
		{"valid name, false", "myVar", false, `let myVar = false;`},
		{"invalid name", "1-invalid", true, `let </script> add By js.LetBool: 1-invalid = true;`},
		{"invalid name with html chars", "1<invalid", false, `let </script> add By js.LetBool: 1&lt;invalid = false;`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := LetBool(tt.varName, tt.val)
			if got != tt.want {
				t.Errorf("LetBool(%q, %v) = %q, want %q", tt.varName, tt.val, got, tt.want)
			}
		})
	}
}

func TestLetJson(t *testing.T) {
	type myStruct struct {
		Foo string `json:"foo"`
		Bar int    `json:"bar"`
	}

	tests := []struct {
		name    string
		varName string
		val     any
		want    string
		wantErr bool
		errText string
	}{
		{"valid struct", "myJson", myStruct{Foo: "baz", Bar: 42}, `let myJson = {"foo":"baz","bar":42};`, false, ""},
		{"invalid identifier", "1-invalid", myStruct{Foo: "baz", Bar: 42}, "", true, "invalid identifire: 1-invalid"},
		{"json marshal error", "myJson", make(chan int), "", true, "json: unsupported type: chan int"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LetJson(tt.varName, tt.val)
			if (err != nil) != tt.wantErr {
				t.Fatalf("LetJson() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errText) {
					t.Errorf("LetJson() error = %q, want to contain %q", err.Error(), tt.errText)
				}
			}
			if got != tt.want {
				t.Errorf("LetJson() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestConstString(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		val     string
		want    string
	}{
		{"valid name", "myVar", "hello", `const myVar = "hello";`},
		{"invalid name", "1-invalid", "hello", `const </script> add By js.ConstString: 1-invalid = "hello";`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConstString(tt.varName, tt.val)
			if got != tt.want {
				t.Errorf("ConstString(%q, %q) = %q, want %q", tt.varName, tt.val, got, tt.want)
			}
		})
	}
}

func TestConstInt(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		val     int
		want    string
	}{
		{"valid name", "myVar", 123, `const myVar = 123;`},
		{"invalid name", "1-invalid", 42, `const </script> add By js.ConstInt: 1-invalid = 42;`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConstInt(tt.varName, tt.val)
			if got != tt.want {
				t.Errorf("ConstInt(%q, %d) = %q, want %q", tt.varName, tt.val, got, tt.want)
			}
		})
	}
}

func TestConstBool(t *testing.T) {
	tests := []struct {
		name    string
		varName string
		val     bool
		want    string
	}{
		{"valid name, true", "myVar", true, `const myVar = true;`},
		{"invalid name", "1-invalid", true, `const </script> add By js.ConstBool: 1-invalid = true;`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ConstBool(tt.varName, tt.val)
			if got != tt.want {
				t.Errorf("ConstBool(%q, %v) = %q, want %q", tt.varName, tt.val, got, tt.want)
			}
		})
	}
}

func TestConstJson(t *testing.T) {
	type myStruct struct {
		Foo string `json:"foo"`
	}

	tests := []struct {
		name    string
		varName string
		val     any
		want    string
		wantErr bool
		errText string
	}{
		{"valid struct", "myJson", myStruct{Foo: "bar"}, `const myJson = {"foo":"bar"};`, false, ""},
		{"invalid identifier", "1-invalid", myStruct{Foo: "bar"}, "", true, "invalid identifire: 1-invalid"},
		{"json marshal error", "myJson", make(chan int), "", true, "json: unsupported type: chan int"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConstJson(tt.varName, tt.val)
			if (err != nil) != tt.wantErr {
				t.Fatalf("ConstJson() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr {
				if !strings.Contains(err.Error(), tt.errText) {
					t.Errorf("ConstJson() error = %q, want to contain %q", err.Error(), tt.errText)
				}
			}
			if got != tt.want {
				t.Errorf("ConstJson() = %q, want %q", got, tt.want)
			}
		})
	}
}
