package evaluator

import (
	"fmt"

	"github.com/vdchnsk/quasark/src/ast"
	"github.com/vdchnsk/quasark/src/object"
	"github.com/vdchnsk/quasark/src/token"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

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

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)

	case *ast.BlockStatement:
		return evalBlockStatements(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Value, env)

	case *ast.IntegerLiteral:
		return &object.Integer{
			Value: node.Value,
		}

	case *ast.StringLiteral:
		return &object.String{
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

	case *ast.AssignStatement:
		val := Eval(node.Value, env)
		if isError(val) {
			return val
		}
		env.Put(node.Identifier.Value, val)

	case *ast.Identifier:
		return evalIdentifier(node.Value, env)

	case *ast.IfExpression:
		return evalIfExpression(node.Condition, node.Consequence, node.Alternative, env)

	case *ast.ReturnStatement:
		returningVal := Eval(node.Value, env)
		if isError(returningVal) {
			return returningVal
		}
		return &object.ReturnWrapper{Value: returningVal}

	case *ast.FuncLiteral:
		params := node.Parameters
		body := node.Body

		return &object.Function{
			Parameters: params,
			Body:       body,
			Env:        env,
		}

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalExpressions(node.Argments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args)

	case *ast.ArrayLiteral:
		elements := evalExpressions(node.Elements, env)
		if len(elements) == 1 && isError(elements[0]) {
			return elements[0]
		}
		return &object.Array{
			Elements: elements,
		}

	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if isError(left) {
			return left
		}
		index := Eval(node.Index, env)
		if isError(index) {
			return index
		}
		return evalIndexExpression(left, index)

	case *ast.HashMapLiteral:
		return evalHashMap(node, env)
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

func evalProgram(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

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

func evalBlockStatements(statements []ast.Statement, env *object.Environment) object.Object {
	var result object.Object

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

func evalExpressions(expressions []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, e := range expressions {
		evalueated := Eval(e, env)
		if isError(evalueated) {
			return []object.Object{evalueated}
		}
		result = append(result, evalueated)
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
			UNKNOWN_OPERATOR, operator, right.Type(),
		)
	}
}

func evalInfixExpression(operator string, left, right object.Object) object.Object {
	lType := left.Type()
	rType := right.Type()

	if lType != rType {
		return newError(
			"%s: %s %s %s",
			TYPE_MISMATCH, left.Type(), operator, right.Type(),
		)
	}

	bothOperandsInts := lType == object.INTEGER_OBJ && rType == object.INTEGER_OBJ
	if bothOperandsInts {
		return evalInfixIntExpression(operator, left, right)
	}

	bothOperandsStrings := lType == object.STRING_OBJ && rType == object.STRING_OBJ
	if bothOperandsStrings {
		return evalInfixStringExpression(operator, left, right)
	}

	bothOperandsBools := lType == object.BOOLEAN_OBJ && rType == object.BOOLEAN_OBJ
	if bothOperandsBools {
		return evalInfixBooleanExpression(operator, left, right)
	}

	return newError(
		"%s: %s %s %s",
		UNKNOWN_OPERATOR, lType, operator, rType,
	)
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
			UNKNOWN_OPERATOR, left.Type(), operator, right.Type(),
		)
	}
}

func evalInfixStringExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value

	switch operator {
	case token.PLUS:
		return &object.String{Value: leftVal + rightVal}
	default:
		return newError(
			"%s: %s %s %s",
			UNKNOWN_OPERATOR, left.Type(), operator, right.Type(),
		)
	}
}

func evalInfixBooleanExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value

	switch operator {
	case token.EQ:
		return nativeBoolToBooleanObject(left == right)
	case token.NOT_EQ:
		return nativeBoolToBooleanObject(left != right)
	case token.AND:
		return nativeBoolToBooleanObject(leftVal && rightVal)
	case token.OR:
		return nativeBoolToBooleanObject(leftVal || rightVal)
	default:
		return newError(
			"%s: %s %s %s",
			UNKNOWN_OPERATOR, left.Type(), operator, right.Type(),
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
			UNKNOWN_OPERATOR, right.Type(),
		)
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func evalIfExpression(condition ast.Expression, consequence, alternative *ast.BlockStatement, env *object.Environment) object.Object {
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

func evalIdentifier(identifier string, env *object.Environment) object.Object {
	if envFunc, ok := env.Get(identifier); ok {
		return envFunc
	}

	if builtInFunc, ok := builtInFuncs[identifier]; ok {
		return builtInFunc
	}

	return newError(
		"%s: %s",
		IDENTIFIER_NOT_FOUND, identifier,
	)
}

func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		extendedEnv := extendFuncEnv(fn, args)
		evaluatedBody := Eval(fn.Body, extendedEnv)

		return unwrapReturnWrapper(evaluatedBody)
	case *object.BuiltInFunction:
		return fn.Fn(args...)
	default:
		return newError(
			"%s %s",
			NOT_A_FUNCTION, fn.Type(),
		)
	}
}

func extendFuncEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnv(fn.Env)

	for paramId, paramData := range fn.Parameters {
		env.Put(paramData.Value, args[paramId])
	}

	return env
}

func unwrapReturnWrapper(obj object.Object) object.Object {
	if returnWrapper, ok := obj.(*object.ReturnWrapper); ok {
		return returnWrapper.Value
	}
	return obj
}

func evalIndexExpression(left, index object.Object) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ && index.Type() == object.INTEGER_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.HASH_MAP_OBJ:
		return evalHashMapIndexExpression(left, index)
	default:
		return newError(
			"%s %s",
			INDEX_OPERATOR_NOT_SUPPORTED, left.Type(),
		)
	}
}

func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObject := array.(*object.Array)
	guestIndex := index.(*object.Integer).Value

	maxIndex := int64(len(arrayObject.Elements)) - 1
	if guestIndex > maxIndex {
		return NULL
	}

	var hostIndex int64
	if guestIndex >= 0 {
		hostIndex = guestIndex
	} else {
		hostIndex = int64(len(arrayObject.Elements)) + int64(guestIndex)
	}

	return arrayObject.Elements[hostIndex]
}

func evalHashMap(node *ast.HashMapLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.HashPair)

	for nodeKey, nodeValue := range node.Pairs {
		key := Eval(nodeKey, env)
		if isError(key) {
			return key
		}

		hashableKey, ok := key.(object.Hashable)
		if !ok {
			return newError(
				fmt.Sprintf(
					"%s %s",
					KEY_IS_NOT_HASHABLE, key.Type(),
				),
			)
		}

		value := Eval(nodeValue, env)
		if isError(key) {
			return key
		}

		pairs[hashableKey.HashKey()] = object.HashPair{Key: key, Value: value}
	}

	return &object.HashMap{Pairs: pairs}
}

func evalHashMapIndexExpression(hashMap, index object.Object) object.Object {
	hashMapObject := hashMap.(*object.HashMap)
	hashMapIndex, ok := index.(object.Hashable)
	if !ok {
		return newError(PROVIDED_INDEX_CANNOT_BE_USED_AS_HASHMAP_KEY)
	}

	pair, ok := hashMapObject.Pairs[hashMapIndex.HashKey()]
	if !ok {
		return NULL
	}
	return pair.Value
}
