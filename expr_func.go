package mathastc

func SetFlag(expr ExprNode, prec int) ExprNode {
	switch node := expr.(type) {
	case OperatorExprNode:
		operator := GetOperator(node.Op[0])
		precedence := operator.Precedence()
		lhs := SetFlag(node.Lhs, precedence)
		rhs := SetFlag(node.Rhs, precedence)
		flag := false
		if prec > precedence {
			flag = true
		}
		return OperatorExprNode{
			Op: node.Op, Lhs: lhs, Rhs: rhs, Flag: flag,
		}
	default:
		return expr
	}
}

func ClearZero(expr ExprNode) ExprNode {
	switch node := expr.(type) {
	case OperatorExprNode:
		lhs := ClearZero(node.Lhs)
		rhs := ClearZero(node.Rhs)
		l, lok := lhs.(NumberExprNode)
		r, rok := rhs.(NumberExprNode)
		if node.Op == "*" && (lok && l.Str == "0" || rok && r.Str == "0") {
			return NumberExprNode{Val: 0, Str: "0"}
		}
		if node.Op == "/" && lok && l.Str == "0" {
			return NumberExprNode{Val: 0, Str: "0"}
		}
		if node.Op == "%" && lok && l.Str == "0" {
			return NumberExprNode{Val: 0, Str: "0"}
		}
		if (node.Op == "+" || node.Op == "-") && (lok && l.Str == "0") {
			return rhs
		}
		if (node.Op == "+" || node.Op == "-") && (rok && r.Str == "0") {
			return lhs
		}
		if node.Op == "^" && lok && l.Str == "0" {
			return NumberExprNode{Val: 0, Str: "0"}
		}
		if node.Op == "^" && rok && r.Str == "0" {
			return NumberExprNode{Val: 1, Str: "1"}
		}
		return OperatorExprNode{
			Op: node.Op, Lhs: lhs, Rhs: rhs, Flag: node.Flag,
		}
	default:
		return expr
	}
}

func MergeNode(expr ExprNode) ExprNode {
	switch node := expr.(type) {
	case OperatorExprNode:
		return node
	default:
		return expr
	}
}

