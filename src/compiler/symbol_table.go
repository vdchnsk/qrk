package compiler

type SymbolScope string

const (
	GlobalScope SymbolScope = "GLOBAL"
	LocalScope  SymbolScope = "LOCAL"
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

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]

	if !ok && s.Outer != nil {
		return s.Outer.Resolve(name)
	}

	return obj, ok
}
