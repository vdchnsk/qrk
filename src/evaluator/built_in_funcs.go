package evaluator

import "github.com/vdchnsk/i-go/src/object"

var builtInFuncs = map[string]*object.BuiltInFunction{
	"len": &object.BuiltInFunction{
		Fn: func(args ...object.Object) object.Object {
			maxAllowedArgs := 1

			if len(args) > maxAllowedArgs {
				return newError(
					"wrong number arguments passed to len func, got=%d, expected=%d",
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
		},
	},
}
