package mathastc

type Parameter struct {
	Vars map[string]any // number | string
	Diff []string
}

func (p Parameter) HasDiffVar(v string) bool {
	for _, s := range p.Diff {
		if s == v {
			return true
		}
	}
	return false
}

func NewParameter(vars map[string]any, diff []string) *Parameter {
	if vars == nil {
		vars = make(map[string]any)
	}
	if diff == nil {
		diff = make([]string, 0)
	}
	return &Parameter{Vars: vars, Diff: diff}
}
