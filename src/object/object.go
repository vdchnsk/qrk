package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/vdchnsk/i-go/src/ast"
)

type ObjectType string

const (
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	RETURN_OBJ   = "RETURN_OBJECT"
	ERROR_OBJ    = "ERROR"
	FUNC_OBJ     = "FUNCTION"
	BUILT_IN_OBJ = "BUILT_IN"
	STRING_OBJ   = "STRING"
	ARRAY_OBJ    = "ARRAY"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Error struct {
	Message string
}

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

func NewEnclosedEnv(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer

	return env
}

func (env *Environment) Get(ident string) (Object, bool) {
	val, ok := env.store[ident]
	if !ok && env.outer != nil {
		return env.outer.Get(ident)
	}
	return val, ok
}

func (env *Environment) Put(ident string, val Object) Object {
	env.store[ident] = val
	return val
}

func (err *Error) Type() ObjectType { return RETURN_OBJ }
func (err *Error) Inspect() string  { return "ERROR:" + err.Message }

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return fmt.Sprintf(s.Value) }

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

type ReturnWrapper struct {
	Value Object
}

func (rw *ReturnWrapper) Type() ObjectType { return RETURN_OBJ }
func (rw *ReturnWrapper) Inspect() string  { return fmt.Sprintf("%d", rw.Value) }

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (fn *Function) Type() ObjectType { return FUNC_OBJ }
func (fn *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range fn.Parameters {
		params = append(params, p.ToString())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(fn.Body.ToString())
	out.WriteString("}\n")

	return out.String()
}

type BuiltInFunction struct {
	Fn func(args ...Object) Object
}

func (fn *BuiltInFunction) Type() ObjectType { return BUILT_IN_OBJ }
func (fn *BuiltInFunction) Inspect() string {
	return "builtIn function"
}

type Array struct {
	Elements []Object
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer

	elements := []string{}
	for _, el := range a.Elements {
		elements = append(elements, el.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")

	return out.String()
}
