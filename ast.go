package mathastc

import (
	"errors"
	"fmt"
	"strconv"
)

// AST 抽象语法树
type AST struct {
	source    string
	currTok   *Token
	currIndex int
	depth     int

	Tokens []*Token
	Err    error
}

func NewAST(toks []*Token, s string) *AST {
	a := &AST{
		Tokens: toks,
		source: s,
	}
	if a.Tokens == nil || len(a.Tokens) == 0 {
		a.Err = errors.New("empty token")
	} else {
		// 设置第一个token为首节点
		a.currIndex = 0
		a.currTok = a.Tokens[0]
	}
	return a
}

func (a *AST) ParseExpression() ExprNode {
	a.depth++ // called depth
	lhs := a.parsePrimary()
	r := a.parseBinOpRHS(0, lhs)
	a.depth--
	if a.depth == 0 && a.currIndex != len(a.Tokens) && a.Err == nil {
		a.Err = errors.New(
			fmt.Sprintf("bad expression, reaching the end or missing the operator\n%s",
				ErrPos(a.source, a.currTok.Offset)))
	}
	return r
}

func (a *AST) getNextToken() *Token {
	a.currIndex++
	if a.currIndex < len(a.Tokens) {
		a.currTok = a.Tokens[a.currIndex]
		return a.currTok
	}
	return nil
}

func (a *AST) getTokPrecedence() int {
	key := a.currTok.Value[0]
	if p, ok := Operators[key]; ok {
		return p.Precedence()
	}
	return -1
}

// 解析Number值
func (a *AST) parseNumber() NumberExprNode {
	f64, err := strconv.ParseFloat(a.currTok.Value, 64)
	if err != nil {
		a.Err = errors.New(
			fmt.Sprintf("%v\nwant '(' or '0-9' but get '%s'\n%s",
				err.Error(),
				a.currTok.Value,
				ErrPos(a.source, a.currTok.Offset)))
		return NumberExprNode{}
	}
	n := NumberExprNode{
		Val: f64,
		Str: a.currTok.Value,
	}
	a.getNextToken()
	return n
}

// 解析函数或常量
func (a *AST) parseFunCallerOrConst() ExprNode {
	name := a.currTok.Value
	a.getNextToken()
	// call func，如果下一个节点为"("表示该节点为函数，否则为常量值
	if a.currTok.Value == "(" {
		f := FunCallerExprNode{}
		if _, ok := defFunc[name]; !ok {
			a.Err = errors.New(
				fmt.Sprintf("function `%s` is undefined\n%s",
					name,
					ErrPos(a.source, a.currTok.Offset)))
			return f
		}
		a.getNextToken()
		exprs := make([]ExprNode, 0)
		if a.currTok.Value == ")" {
			// function call without parameters
			// ignore the process of parameter resolution
		} else {
			exprs = append(exprs, a.ParseExpression())
			for a.currTok.Value != ")" && a.getNextToken() != nil {
				if a.currTok.Type == CommaType {
					continue
				}
				exprs = append(exprs, a.ParseExpression())
			}
		}
		def := defFunc[name]
		// 校验函数参数
		if def.Argc() >= 0 && len(exprs) != def.Argc() {
			a.Err = errors.New(
				fmt.Sprintf("wrong way calling function `%s`, parameters want %d but get %d\n%s",
					name,
					def.Argc(),
					len(exprs),
					ErrPos(a.source, a.currTok.Offset)))
		}
		a.getNextToken()
		f.Name = name
		f.Arg = exprs
		return f
	}

	// call const
	if v, ok := defConst[name]; ok {
		return ConstExprNode{
			Name: name,
			Val:  v,
			Str:  strconv.FormatFloat(v, 'f', 0, 64),
		}
	} else {
		return VariableExprNode{Val: name}
		//a.Err = errors.New(
		//	fmt.Sprintf("const `%s` is undefined\n%s",
		//		name,
		//		ErrPos(a.source, a.currTok.Offset)))
		//return NumberExprNode{}
	}

	//if v, ok := defConst[name]; ok {
	//	return ConstExprNode{
	//		Name: name,
	//		Val:  v,
	//		Str:  strconv.FormatFloat(v, 'f', 0, 64),
	//	}
	//} else {
	//	a.Err = errors.New(
	//		fmt.Sprintf("const `%s` is undefined\n%s",
	//			name,
	//			ErrPos(a.source, a.currTok.Offset)))
	//	return NumberExprNode{}
	//}
}

// 解析操作符
func (a *AST) parseOperator() ExprNode {
	if a.currTok.Value == "(" {
		t := a.getNextToken()
		if t == nil {
			a.Err = errors.New(
				fmt.Sprintf("want '(' or '0-9' but get EOF\n%s",
					ErrPos(a.source, a.currTok.Offset)))
			return nil
		}
		e := a.ParseExpression()
		if e == nil {
			return nil
		}
		if a.currTok.Value != ")" {
			a.Err = errors.New(
				fmt.Sprintf("want ')' but get %s\n%s",
					a.currTok.Value,
					ErrPos(a.source, a.currTok.Offset)))
			return nil
		}
		a.getNextToken()
		return e
	} else if a.currTok.Value == "-" {
		if a.getNextToken() == nil {
			a.Err = errors.New(
				fmt.Sprintf("want '0-9' but get '-'\n%s",
					ErrPos(a.source, a.currTok.Offset)))
			return nil
		}
		bin := OperatorExprNode{
			Op:  "-",
			Lhs: NumberExprNode{},
			Rhs: a.parsePrimary(),
		}
		return bin
	} else {
		return a.parseNumber()
	}
}

// 解析变量
func (a *AST) parseVariable() ExprNode {
	n := VariableExprNode{
		Val: a.currTok.Value,
	}
	a.getNextToken()
	return n
}

func (a *AST) parsePrimary() ExprNode {
	switch a.currTok.Type {
	case IdentifierType:
		return a.parseFunCallerOrConst()
	case LiteralType:
		return a.parseNumber()
	case OperatorType:
		return a.parseOperator()
	case CommaType:
		a.Err = errors.New(
			fmt.Sprintf("want '(' or '0-9' but get %s\n%s",
				a.currTok.Value,
				ErrPos(a.source, a.currTok.Offset)))
		return nil
	default:
		return nil
	}
}

func (a *AST) parseBinOpRHS(execPrec int, lhs ExprNode) ExprNode {
	for {
		tokPrec := a.getTokPrecedence()
		if tokPrec < execPrec {
			return lhs
		}
		binOp := a.currTok.Value
		if a.getNextToken() == nil {
			a.Err = errors.New(
				fmt.Sprintf("want '(' or '0-9' but get EOF\n%s",
					ErrPos(a.source, a.currTok.Offset)))
			return nil
		}
		rhs := a.parsePrimary()
		if rhs == nil {
			return nil
		}
		nextPrec := a.getTokPrecedence()
		if tokPrec < nextPrec {
			rhs = a.parseBinOpRHS(tokPrec+1, rhs)
			if rhs == nil {
				return nil
			}
		}
		lhs = OperatorExprNode{
			Op:   binOp,
			Lhs:  lhs,
			Rhs:  rhs,
			Flag: false,
		}
	}
}
