package mathastc

import (
	"errors"
	"fmt"
	"strings"
)

type Parser struct {
	Source string
	ch     byte
	offset int
	err    error
}

func Parse(s string) ([]*Token, error) {
	p := &Parser{
		Source: s,
		err:    nil,
		ch:     s[0],
	}
	toks := p.parse()
	if p.err != nil {
		return nil, p.err
	}
	return toks, nil
}

func (p *Parser) parse() []*Token {
	toks := make([]*Token, 0)
	for {
		tok := p.nextTok()
		if tok == nil {
			break
		}
		toks = append(toks, tok)
	}
	return toks
}

func (p *Parser) nextTok() *Token {
	if p.offset >= len(p.Source) || p.err != nil {
		return nil
	}
	var err error
	for p.isWhitespace(p.ch) && err == nil {
		err = p.nextCh()
	}
	start := p.offset
	var tok *Token

	// 判断是否操作符号, []()+-*/%^
	if ounit, ok := Operators[p.ch]; ok == true {
		tok = &Token{
			Value: string(ounit.Name()),
			Type:  OperatorType,
		}
		tok.Offset = start
		err = p.nextCh()
		return tok
	}

	// 判断是否字面数字
	if p.IsLiteral(p.ch) {
		tok = &Token{
			Value: strings.ReplaceAll(p.Source[start:p.offset], "_", ""),
			Type:  LiteralType,
		}
		tok.Offset = start
		return tok
	}

	// 判断是否逗号
	if p.ch == ',' {
		tok = &Token{
			Value: string(p.ch),
			Type:  CommaType,
		}
		tok.Offset = start
		err = p.nextCh()
		return tok
	}

	// 判断是否为变量
	//if p.isVar(p.ch) {
	//	tok = &Token{
	//		Value: p.Source[start:p.offset],
	//		Type:  VariableType,
	//	}
	//	tok.Offset = start
	//	return tok
	//}

	if p.isChar(p.ch) {
		for p.isWordChar(p.ch) && p.nextCh() == nil {
		}
		tok = &Token{
			Value: p.Source[start:p.offset],
			Type:  IdentifierType,
		}
		tok.Offset = start
	} else if p.ch != ' ' {
		s := fmt.Sprintf("symbol error: unknown '%v', pos [%v:]\n%s",
			string(p.ch),
			start,
			ErrPos(p.Source, start))
		p.err = errors.New(s)
	}

	return tok
}

func (p *Parser) IsLiteral(v byte) bool {
	switch v {
	case
		'0',
		'1',
		'2',
		'3',
		'4',
		'5',
		'6',
		'7',
		'8',
		'9':
		for p.isDigitNum(p.ch) && p.nextCh() == nil {
			if (p.ch == '-' || p.ch == '+') && p.Source[p.offset-1] != 'e' {
				break
			}
		}
		return true
	default:
		return false
	}
}

func (p *Parser) nextCh() error {
	p.offset++
	if p.offset < len(p.Source) {
		p.ch = p.Source[p.offset]
		return nil
	}
	return errors.New("EOF")
}

func (p *Parser) isWhitespace(c byte) bool {
	return c == ' ' ||
		c == '\t' ||
		c == '\n' ||
		c == '\v' ||
		c == '\f' ||
		c == '\r'
}

func (p *Parser) isDigitNum(c byte) bool {
	return '0' <= c && c <= '9' || c == '.' || c == '_' || c == 'e' || c == '-' || c == '+'
}

func (p *Parser) isChar(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '\'' || c == '$' || c == '#'
}

func (p *Parser) isWordChar(c byte) bool {
	return p.isChar(c) || '0' <= c && c <= '9' || c == '$' || c == '#'
}

//func (p *Parser) isVar(v byte) bool {
//	switch v {
//	case
//		'$',
//		'#':
//		for p.nextCh() == nil {
//			c := p.Source[p.offset]
//			if !p.isChar(c) {
//				break
//			}
//		}
//		return true
//	default:
//		return false
//	}
//}
