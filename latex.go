package mathastc

import (
	"context"
	"log"
)

// ToLaTex latex支持
func ToLaTex(expr ExprNode, ctx context.Context) string {

	var l, r string
	switch node := expr.(type) {

	case OperatorExprNode:
		operator := GetOperator(node.Op[0])
		l = ToLaTex(node.Lhs, ctx)
		r = ToLaTex(node.Rhs, ctx)
		if node.Flag {
			return "(" + operator.ToLaTex(l, r) + ")"
		}
		return operator.ToLaTex(l, r)

	case NumberExprNode:
		return node.Str

	case ConstExprNode:
		if GetDefConstLaTex(node.Name) != "" {
			return GetDefConstLaTex(node.Name)
		}
		return node.Name

	case VariableExprNode:
		val := node.Val
		parameter, err := GetCtxParameter(ctx)
		if err != nil {
			log.Panicf("no parameter found, err = %v\n", err)
		}

		value, ok := parameter.Vars[val]
		if !ok {
			return val
		}

		switch t := value.(type) {
		case string:
			expression, err2 := ParseExpression(t)
			if err2 != nil {
				return val
			}
			return ToLaTex(expression, ctx)
		case ExprNode:
			return ToLaTex(t, ctx)
		default:
			return node.Val
		}

	case FunCallerExprNode:
		def := GetDefFunc(node.Name)
		switch t := def.(type) {
		case LaTexFunc:
			return t.LaTex(ctx, node.Arg...)
		default:
			return node.Name
		}
	}

	return ""
}

