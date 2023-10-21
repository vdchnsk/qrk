package evaluator

import (
	"fmt"

	"github.com/vdchnsk/i-go/src/ast"
	"github.com/vdchnsk/i-go/src/error"
	"github.com/vdchnsk/i-go/src/memory"
	"github.com/vdchnsk/i-go/src/object"
	"github.com/vdchnsk/i-go/src/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node, env *memory.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements)

	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements)

	case *ast.ExpressionStatement:
		return Eval(node.Value, env)

	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}

	case *ast.Boolean:
		return nativeBoolToBooleanObject(node.Value)

	case *ast.PrefixExpression:
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		right := Eval(node.Right, env)
		if isError(right) {
			return right
		}
		return evalInfixExpression(node.Operator, left, right)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Put(node.Identifier.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node.Value, env)

	case *ast.IfExpression:
		return evalIfExpression(node.Condition, node.Consequence, node.Alternative)

	case *ast.ReturnStatement:
		returningVal := Eval(node.Value, nil)
		if isError(returningVal) {
			return returningVal
		}
		return &object.ReturnWrapper{Value: returningVal}
	}
	return nil
}

func newError(format string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, args...)}
}

func nativeBoolToBooleanObject(input bool) *object.Boolean {
	if input {
		return TRUE
	} else {
		return FALSE
	}

}

func evalProgram(statements []ast.Statement) object.Object {
	var result object.Object
	env := memory.NewEnvironment()

	for _, statement := range statements {
		result = Eval(statement, env)

		switch result := result.(type) {
		case *object.ReturnWrapper:
			return result.Value
		case *object.Error:
			return result
		}
	}
	return result
}

func evalBlockStatements(statements []ast.Statement) object.Object {
	var result object.Object
	env := memory.NewEnvironment()

	for _, statement := range statements {
		result = Eval(statement, env)

		_, isReturnWrapper := result.(*object.ReturnWrapper)
		_, isError := result.(*object.Error)

		if result != nil && (isReturnWrapper || isError) {
			return result
		}
	}
	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case token.BANG:
		return evalBangOperatorExpression(right)
	case token.MINUS:
		return evalMinusOperatorExpression(right)
	default:
		return newError(
			"%s: %s%s",
			error.UNKNOWN_OPERATOR, operator, right.Type(),
		)
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	lType := left.Type()
	rType := right.Type()

	if lType != rType {
		return newError(
			"%s: %s %s %s",
			error.TYPE_MISMATCH, left.Type(), operator, right.Type(),
		)
	}

	bothOperandsInts := lType == object.INTEGER_OBJ && rType == object.INTEGER_OBJ
	if bothOperandsInts {
		return evalInfixIntExpression(operator, left, right)
	}

	switch operator {
	case token.EQ:
		return nativeBoolToBooleanObject(left == right)
	case token.NOT_EQ:
		return nativeBoolToBooleanObject(left != right)
	default:
		return newError(
			"%s: %s %s %s",
			error.UNKNOWN_OPERATOR, lType, operator, rType,
		)
	}
}

func evalInfixIntExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value

	switch operator {
	case token.PLUS:
		return &object.Integer{Value: leftVal + rightVal}
	case token.MINUS:
		return &object.Integer{Value: leftVal - rightVal}
	case token.SLASH:
		return &object.Integer{Value: leftVal / rightVal}
	case token.ASTERISK:
		return &object.Integer{Value: leftVal * rightVal}
	case token.LT:
		return nativeBoolToBooleanObject(leftVal < rightVal)
	case token.GT:
		return nativeBoolToBooleanObject(leftVal > rightVal)
	case token.EQ:
		return nativeBoolToBooleanObject(leftVal == rightVal)
	case token.NOT_EQ:
		return nativeBoolToBooleanObject(leftVal != rightVal)
	default:
		return newError(
			"%s: %s %s %s",
			error.UNKNOWN_OPERATOR, left.Type(), operator, right.Type(),
		)
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE: // * !true == false
		return FALSE
	case FALSE: // * !false == true
		return TRUE
	case NULL: // * !null == true
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError(
			"%s: -%s",
			error.UNKNOWN_OPERATOR, right.Type(),
		)
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIfExpression(condition ast.Expression, consequence, alternative *ast.BlockStatement) object.Object {
	env := memory.NewEnvironment()
	conditionResult := Eval(condition, env)

	if isError(conditionResult) {
		return conditionResult
	}
	if isTruthy(conditionResult) {
		return Eval(consequence, env)
	}
	if alternative != nil {
		return Eval(alternative, env)
	}
	return NULL
}

func evalIdentifier(identifier string, env *memory.Environment) object.Object {
	val, ok := env.Get(identifier)
	if !ok {
		return newError(
			"%s: %s",
			error.IDENTIFIER_NOT_FOUND, identifier,
		)
	}
	return val
}

func isTruthy(obj object.Object) bool {
	if obj == TRUE {
		return true
	}
	if obj == NULL || obj == FALSE {
		return false
	}
	return true
}

func isError(obj object.Object) bool {
	if obj == nil {
		return false
	}
	return obj.Type() == object.ERROR_OBJ
}
