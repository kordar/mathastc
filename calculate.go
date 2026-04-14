package mathastc

import (
	"context"
	"github.com/spf13/cast"
	"log"
)

// Calculate 计算节点
func Calculate(expr ExprNode, ctx context.Context) float64 {

	var l, r float64
	switch node := expr.(type) {

	case OperatorExprNode:
		l = Calculate(node.Lhs, ctx)
		r = Calculate(node.Rhs, ctx)
		return GetOperator(node.Op[0]).Result(l, r)

	case NumberExprNode:
		return node.Val

	case ConstExprNode:
		return node.Val

	case VariableExprNode:
		val := node.Val
		parameter, err := GetCtxParameter(ctx)
		if err != nil {
			log.Panicf("no parameter found, err = %v\n", err)
		}

		value, ok := parameter.Vars[val]
		if !ok {
			log.Panicf("no parameter value found for %s\n", val)
		}

		switch t := value.(type) {
		case string:
			expression, err2 := ParseExpression(t)
			if err2 != nil {
				log.Panicf("Convert to expression node exception, err = %s\n", t)
			}
			return Calculate(expression, ctx)
		case ExprNode:
			return Calculate(t, ctx)
		case int, int8, int64, int16, int32, uint, uint8, uint16, uint32, uint64, float32, float64:
			return cast.ToFloat64(t)
		default:
			log.Panicln("unknown expr type")
		}

	case FunCallerExprNode:
		def := GetDefFunc(node.Name)
		return def.Calculate(ctx, node.Arg...)
	}

	return 0.0
}

// ToExprStr 打印节点
func ToExprStr(expr ExprNode, ctx context.Context) string {

	var l, r string
	switch node := expr.(type) {

	case OperatorExprNode:
		l = ToExprStr(node.Lhs, ctx)
		r = ToExprStr(node.Rhs, ctx)
		operator := GetOperator(node.Op[0])
		if node.Flag {
			return "(" + operator.ToExprStr(l, r) + ")"
		}
		return operator.ToExprStr(l, r)

	case NumberExprNode:
		return node.Str

	case ConstExprNode:
		return node.Str

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
			return ToExprStr(expression, ctx)
		case ExprNode:
			return ToExprStr(t, ctx)
		case int, int8, int64, int16, int32, uint, uint8, uint16, uint32, uint64, float32, float64:
			return cast.ToString(t)
		default:
			return node.Val
		}

	case FunCallerExprNode:
		def := GetDefFunc(node.Name)
		return def.ToExprStr(ctx, node.Arg...)
	}

	return ""
}
