package mathastc

import (
	"context"
	"errors"
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
