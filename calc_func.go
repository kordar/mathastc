package mathastc

import (
	"context"
	"fmt"
)

type numericDiffFunc struct {
}

func (n numericDiffFunc) Calculate(ctx context.Context, args ...ExprNode) float64 {
	if len(args) < 3 || len(args) > 4 {
		panic(fmt.Errorf("diff(varName, expr, x, [h]) expect 3 or 4 args"))
	}
	nameNode, ok := args[0].(VariableExprNode)
	if !ok {
		panic(fmt.Errorf("first arg of diff must be variable name"))
	}
	varName := nameNode.Val
	expr := args[1]
	xValExpr := args[2]
	x := Calculate(xValExpr, ctx)
	h := 1e-6
	if len(args) == 4 {
		h = Calculate(args[3], ctx)
	}
	return DiffAt(ctx, expr, varName, x, h)
}

func (n numericDiffFunc) ToExprStr(ctx context.Context, args ...ExprNode) string {
	return "diff(...)"
}

func (n numericDiffFunc) Argc() int {
	return -1
}

type numericIntegralFunc struct {
}

func (n numericIntegralFunc) Calculate(ctx context.Context, args ...ExprNode) float64 {
	if len(args) < 4 || len(args) > 5 {
		panic(fmt.Errorf("integral(varName, expr, a, b, [n]) expect 4 or 5 args"))
	}
	nameNode, ok := args[0].(VariableExprNode)
	if !ok {
		panic(fmt.Errorf("first arg of integral must be variable name"))
	}
	varName := nameNode.Val
	expr := args[1]
	a := Calculate(args[2], ctx)
	b := Calculate(args[3], ctx)
	steps := 100.0
	if len(args) == 5 {
		steps = Calculate(args[4], ctx)
	}
	return IntegralSimpson(ctx, expr, varName, a, b, int(steps))
}

func (n numericIntegralFunc) ToExprStr(ctx context.Context, args ...ExprNode) string {
	return "integral(...)"
}

func (n numericIntegralFunc) Argc() int {
	return -1
}

func init() {
	_ = RegDefFunc("diff", numericDiffFunc{})
	_ = RegDefFunc("integral", numericIntegralFunc{})
}
