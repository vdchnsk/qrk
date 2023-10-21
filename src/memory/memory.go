package memory

import "github.com/vdchnsk/i-go/src/object"

type Environment struct {
	store map[string]object.Object
}

func NewEnvironment() *Environment {
	s := make(map[string]object.Object)
	return &Environment{store: s}
}

func (env *Environment) Get(ident string) (object.Object, bool) {
	val, ok := env.store[ident]
	return val, ok
}

func (env *Environment) Put(ident string, val object.Object) object.Object {
	env.store[ident] = val
	return val
}
