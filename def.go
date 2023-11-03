package mathastc

import (
	"context"
	"math"
)

// 定义全局常量
var defConst = map[string]float64{
	"pi":    math.Pi,
	"e":     math.E,
	"infty": 0,
}

// 定义全局LaTex常量
var defConstLaTex = map[string]string{
	"pi":    "π",
	"e":     "e",
	"infty": "\\infty",
}

// FuncExprNode处理对象
var defFunc map[string]DefFunc = map[string]DefFunc{}

// DefFunc 节点运算
type DefFunc interface {
	Calculate(ctx context.Context, args ...ExprNode) float64
	ToExprStr(ctx context.Context, args ...ExprNode) string
	Argc() int
}

// LaTexFunc LaTex生成
type LaTexFunc interface {
	LaTex(ctx context.Context, args ...ExprNode) string
}

// DiffExprNodeFunc 微分运算
type DiffExprNodeFunc interface {
	DiffExprNode(ctx context.Context, args ...ExprNode) ExprNode
}
