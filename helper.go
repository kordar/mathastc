package mathastc

import (
	"context"
	"errors"
	"reflect"
	"strconv"
	"strings"
)

type Int interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64
}

type Uint interface {
	~uint | ~uint8 | ~uint16 | ~uint32
}

type Float interface {
	~float32 | ~float64
}

type Integer interface {
	Int | Uint
}

func ErrPos(s string, pos int) string {
	r := strings.Repeat("-", len(s)) + "\n"
	s += "\n"
	for i := 0; i < pos; i++ {
		s += " "
	}
	s += "^\n"
	return r + s + r
}

func Float64ToStr(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func typeof(v any) string {
	return reflect.TypeOf(v).Name()
}

func cloneCtxWithVar(ctx context.Context, name string, v float64) context.Context {
	p, err := GetCtxParameter(ctx)
	if err != nil {
		vars := map[string]any{name: v}
		np := NewParameter(vars, nil)
		return context.WithValue(ctx, "parameter", np)
	}
	vars := make(map[string]any, len(p.Vars)+1)
	for k, val := range p.Vars {
		vars[k] = val
	}
	vars[name] = v
	np := NewParameter(vars, p.Diff)
	return context.WithValue(ctx, "parameter", np)
}

func DiffAt(ctx context.Context, expr ExprNode, varName string, x float64, h float64) float64 {
	if h <= 0 {
		h = 1e-6
	}
	ctx1 := cloneCtxWithVar(ctx, varName, x+h)
	ctx2 := cloneCtxWithVar(ctx, varName, x-h)
	f1 := Calculate(expr, ctx1)
	f2 := Calculate(expr, ctx2)
	return (f1 - f2) / (2 * h)
}

func DiffExprAt(ctx context.Context, expr string, varName string, x float64, h float64) (float64, error) {
	e, err := ParseExpression(expr)
	if err != nil {
		return 0, err
	}
	return DiffAt(ctx, e, varName, x, h), nil
}

func IntegralSimpson(ctx context.Context, expr ExprNode, varName string, a float64, b float64, n int) float64 {
	if n <= 0 {
		n = 100
	}
	if n%2 == 1 {
		n++
	}
	h := (b - a) / float64(n)
	sum := 0.0
	for i := 0; i <= n; i++ {
		x := a + float64(i)*h
		c := cloneCtxWithVar(ctx, varName, x)
		fx := Calculate(expr, c)
		coef := 1.0
		if i == 0 || i == n {
			coef = 1
		} else if i%2 == 1 {
			coef = 4
		} else {
			coef = 2
		}
		sum += coef * fx
	}
	return sum * h / 3.0
}

func IntegralExprSimpson(ctx context.Context, expr string, varName string, a float64, b float64, n int) (float64, error) {
	e, err := ParseExpression(expr)
	if err != nil {
		return 0, err
	}
	return IntegralSimpson(ctx, e, varName, a, b, n), nil
}

// ParseExpression 解析表达式
func ParseExpression(s string) (ExprNode, error) {
	toks, err := Parse(s)
	if err != nil {
		return nil, err
	}
	ast := NewAST(toks, s)
	if ast.Err != nil {
		return nil, ast.Err
	}
	ar := ast.ParseExpression()
	if ast.Err != nil {
		return nil, ast.Err
	}
	defer func() {
		if e := recover(); e != nil {
			err = e.(error)
		}
	}()
	return ar, nil
}

// RegDefFunc 注册函数
func RegDefFunc(name string, df DefFunc) error {
	if len(name) == 0 {
		return errors.New("RegFunction name is not empty")
	}
	if df.Argc() < -1 {
		return errors.New("RegFunction argc should be -1, 0, or a positive integer")
	}
	if _, ok := defFunc[name]; ok {
		return errors.New("RegFunction name is already exist")
	}
	defFunc[name] = df
	return nil
}

// RegConst 注册全局常量
func RegConst(name string, value float64) error {
	if len(name) == 0 {
		return errors.New("RegConst name is not empty")
	}
	if _, ok := defConst[name]; ok {
		return errors.New("RegConst name is already exist")
	}
	defConst[name] = value
	return nil
}

// RegConstLaTex 注册全局latex
func RegConstLaTex(name string, value string) error {
	if len(name) == 0 {
		return errors.New("RegConstLaTex name is not empty")
	}
	if _, ok := defConstLaTex[name]; ok {
		return errors.New("RegConstLaTex name is already exist")
	}
	defConstLaTex[name] = value
	return nil
}

// GetDefFunc 获取函数
func GetDefFunc(name string) DefFunc {
	return defFunc[name]
}

// GetOperator 获取操作单元
func GetOperator(name byte) OperatorItem {
	return Operators[name]
}

// GetDefConstLaTex 获取全局latex
func GetDefConstLaTex(name string) string {
	return defConstLaTex[name]
}

// GetDefConst 获取全局常量
func GetDefConst(name string) float64 {
	return defConst[name]
}

// GetCtxParameter 解析上下文Parameter对象
func GetCtxParameter(ctx context.Context) (*Parameter, error) {
	value := ctx.Value("parameter")
	if value == nil {
		return nil, errors.New("no parameter found")
	}

	parameter, exists := value.(*Parameter)
	if !exists {
		return nil, errors.New("no parameter found")
	}
	return parameter, nil
}
