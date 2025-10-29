package stdlib

import (
	"fmt"

	"github.com/vdchnsk/qrk/src/object"
)

var NULL = &object.Null{}

func lenBuiltin(args ...object.Object) object.Object {
	maxAllowedArgs := 1

	if len(args) > maxAllowedArgs {
		return newError(
			"wrong number of arguments passed to len func, got=%d, supported=%d",
			len(args), maxAllowedArgs,
		)
	}
	switch arg := args[0].(type) {
	case *object.String:
		return &object.Integer{
			Value: int64(len(arg.Value)),
		}
	default:
		return newError(
			"argument to `len` is not supported, got %s",
			arg.Type(),
		)
	}
}

func print(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return NULL
}

var FuncsMap = map[string]*object.BuiltInFunction{
	"len":   {Fn: lenBuiltin, Name: "len", ParamsCount: 1},
	"print": {Fn: print, Name: "print", ParamsCount: 1},
}

var Funcs = []*object.BuiltInFunction{
	FuncsMap["len"],
	FuncsMap["print"],
}
