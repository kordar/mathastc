package mathastc

type TokenType int32

const (
	IdentifierType TokenType = iota // 标识符(常量、函数、变量)
	LiteralType                     // 字面文字
	OperatorType                    // 操作符号
	CommaType                       // 逗号
)

type Token struct {
	// raw characters
	Value string
	// type with Identifier/Literal/Operator/Comma/Variable
	Type   TokenType
	Flag   int
	Offset int
}
