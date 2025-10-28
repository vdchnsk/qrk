package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
	StdlibScope SymbolScope = "STD_LIB"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolTable struct {
	Outer *SymbolTable

	store            map[string]Symbol
	definitionsCount int
}

func NewSymbolTable() *SymbolTable {
	return &SymbolTable{
		store: make(map[string]Symbol),
	}
}

func NewEnclosedSymbolTable(outerSymbolTable *SymbolTable) *SymbolTable {
	store := NewSymbolTable()
	store.Outer = outerSymbolTable

	return store
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{
		Name:  name,
		Index: s.definitionsCount,
	}

	isGlobalScope := s.Outer == nil

	if isGlobalScope {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}

	s.store[name] = symbol
	s.definitionsCount++

	return symbol
}

func (s *SymbolTable) DefineStdlibFunc(index int, name string) Symbol {
	sym := Symbol{
		Name:  name,
		Scope: StdlibScope,
		Index: index,
	}

	s.store[name] = sym

	return sym
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, inCurScope := s.store[name]

	if !inCurScope && s.Outer != nil {
		return s.Outer.Resolve(name)
	}

	return obj, inCurScope
}
