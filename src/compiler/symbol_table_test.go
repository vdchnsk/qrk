package compiler

import "testing"

func TestDefine(t *testing.T) {
	expected := map[string]Symbol{
		"a": Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		"b": Symbol{Name: "b", Scope: GlobalScope, Index: 1},
	}

	global := NewSymbolTable()

	a := global.Define("a")
	if a != expected["a"] {
		t.Errorf("a is not defined. got=%+v", a)
	}

	b := global.Define("b")
	if b != expected["b"] {
		t.Errorf("b is not defined. got=%+v", b)
	}
}

func TestResolveGlobal(t *testing.T) {
	global := NewSymbolTable()
	global.Define("a")
	global.Define("b")

	expected := []Symbol{
		Symbol{Name: "a", Scope: GlobalScope, Index: 0},
		Symbol{Name: "b", Scope: GlobalScope, Index: 1},
	}

	for _, symbol := range expected {
		result, ok := global.Resolve(symbol.Name)
		if !ok {
			t.Errorf("name %s not resolvable", symbol.Name)
			continue
		}
		if result != symbol {
			t.Errorf("expected %s to resolve to %+v, got %+v", symbol.Name, symbol, result)
		}
	}
}
