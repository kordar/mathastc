package mathastc

type MergePiece struct {
	op   string
	node ExprNode
}

func NewMergePiece(op string, node ExprNode) MergePiece {
	return MergePiece{node: node, op: op}
}

func ExprBreakUp(expr ExprNode) []MergePiece {
	data := make([]MergePiece, 0)
	breakup(expr, &data, "+")
	return data
}

func breakup(expr ExprNode, data *[]MergePiece, op string) ExprNode {
	switch node := expr.(type) {
	case OperatorExprNode:
		if node.Op == "+" || node.Op == "-" {
			lhs := breakup(node.Lhs, data, node.Op)
			if lhs != nil {
				*data = append(*data, NewMergePiece(op, lhs))
			}
			rhs := breakup(node.Rhs, data, node.Op)
			if rhs != nil {
				*data = append(*data, NewMergePiece(op, rhs))
			}
			return nil
		}
		if node.Op == "*" || node.Op == "/" {
			node.Lhs = breakup(node.Lhs, data, op)
			node.Rhs = breakup(node.Rhs, data, op)
			return node
		}
	default:
	}
	return expr
}

