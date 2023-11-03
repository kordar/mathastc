package mathastc

import (
	"errors"
	"fmt"
	"math"
	"math/big"
)

const (
	NonePrecedence = -1 // 权重值
	NoneResult     = 0.0
)

type OperatorItem interface {
	Name() byte
	Precedence() int
	Result(a float64, b float64) float64
	ToLaTex(a string, b string) string
	ToExprStr(a string, b string) string
}

var Operators = map[byte]OperatorItem{
	'(': &LBrackets{},
	')': &RBrackets{},
	//'[': &LMBrackets{},
	//']': &RMBrackets{},
	'+': &Plus{},
	'-': &Minus{},
	'*': &Mul{},
	'/': &Div{},
	'^': &Pow{},
	'%': &Mod{},
}

// LBrackets 左括号
type LBrackets struct {
}

func (L *LBrackets) Name() byte {
	return '('
}

func (L *LBrackets) Precedence() int {
	return NonePrecedence
}

func (L *LBrackets) Result(a float64, b float64) float64 {
	return NoneResult
}

func (L *LBrackets) ToExprStr(a string, b string) string {
	return ""
}

func (L *LBrackets) ToLaTex(a string, b string) string {
	return ""
}

// RBrackets 右括号
type RBrackets struct {
}

func (R *RBrackets) Name() byte {
	return ')'
}

func (R *RBrackets) Precedence() int {
	return NonePrecedence
}

func (R *RBrackets) Result(a float64, b float64) float64 {
	return NoneResult
}

func (R *RBrackets) ToExprStr(a string, b string) string {
	return ""
}

func (R *RBrackets) ToLaTex(a string, b string) string {
	return ""
}

// LMBrackets 左中括号
type LMBrackets struct {
	*LBrackets
}

func (L *LMBrackets) Name() byte {
	return '['
}

// RMBrackets 右中括号
type RMBrackets struct {
	*RBrackets
}

func (R *RMBrackets) Name() byte {
	return ']'
}

// Div 两数相除
type Div struct {
}

func (d *Div) Name() byte {
	return '/'
}

func (d *Div) Precedence() int {
	return 40
}

func (d *Div) Result(a float64, b float64) float64 {
	if b == 0 {
		panic(errors.New(
			fmt.Sprintf("violation of arithmetic specification: a division by zero in ExprASTResult: [%g/%g]",
				a,
				b)))
	}
	f, _ := new(big.Float).Quo(new(big.Float).SetFloat64(a), new(big.Float).SetFloat64(b)).Float64()
	return f
}

func (d *Div) ToExprStr(a string, b string) string {
	return fmt.Sprintf("%s/%s", a, b)
}

func (d *Div) ToLaTex(a string, b string) string {
	return fmt.Sprintf("\\frac{%s}{%s}", a, b)
}

// Minus 两数相减
type Minus struct {
}

func (m *Minus) Name() byte {
	return '-'
}

func (m *Minus) Precedence() int {
	return 20
}

func (m *Minus) Result(a float64, b float64) float64 {
	lh := big.NewFloat(a)
	rh := big.NewFloat(b)
	f, _ := new(big.Float).Sub(lh, rh).Float64()
	return f
}

func (m *Minus) ToExprStr(a string, b string) string {
	return fmt.Sprintf("%s - %s", a, b)
}

func (m *Minus) ToLaTex(a string, b string) string {
	return fmt.Sprintf("%s - %s", a, b)
}

// Mod 两数取模
type Mod struct {
}

func (m *Mod) Name() byte {
	return '%'
}

func (m *Mod) Precedence() int {
	return 40
}

func (m *Mod) Result(a float64, b float64) float64 {
	if b == 0 {
		panic(errors.New(
			fmt.Sprintf("violation of arithmetic specification: a division by zero in ExprASTResult: [%g%%%g]",
				a,
				b)))
	}
	return float64(int(a) % int(b))
}

func (m *Mod) ToExprStr(a string, b string) string {
	return fmt.Sprintf("(%s %% %s)", a, b)
}

func (m *Mod) ToLaTex(a string, b string) string {
	return fmt.Sprintf("(%s %% %s)", a, b)
}

// Mul 两数相乘
type Mul struct {
}

func (m *Mul) Name() byte {
	return '*'
}

func (m *Mul) Precedence() int {
	return 40
}

func (m *Mul) Result(a float64, b float64) float64 {
	f, _ := new(big.Float).Mul(new(big.Float).SetFloat64(a), new(big.Float).SetFloat64(b)).Float64()
	return f
}

func (m *Mul) ToExprStr(a string, b string) string {
	return fmt.Sprintf("%s * %s", a, b)
}

func (m *Mul) ToLaTex(a string, b string) string {
	return fmt.Sprintf("%s \\times %s", a, b)
}

// Plus 两数相加
type Plus struct {
}

func (p *Plus) Name() byte {
	return '+'
}

func (p *Plus) Precedence() int {
	return 20
}

func (p *Plus) Result(a float64, b float64) float64 {
	lh := big.NewFloat(a)
	rh := big.NewFloat(b)
	f, _ := new(big.Float).Add(lh, rh).Float64()
	return f
}

func (p *Plus) ToExprStr(a string, b string) string {
	return fmt.Sprintf("%s + %s", a, b)
}

func (p *Plus) ToLaTex(a string, b string) string {
	return fmt.Sprintf("%s + %s", a, b)
}

// Pow 指数运算
type Pow struct {
}

func (p *Pow) Name() byte {
	return '^'
}

func (p *Pow) Precedence() int {
	return 60
}

func (p *Pow) Result(a float64, b float64) float64 {
	return math.Pow(a, b)
}

func (p *Pow) ToExprStr(a string, b string) string {
	return fmt.Sprintf("%s^%s", a, b)
}

func (p *Pow) ToLaTex(a string, b string) string {
	return fmt.Sprintf("%s^{%s}", a, b)
}
