package nodestack

import (
	"sort"

	"github.com/google/go-jsonnet/ast"
)

type NodeStack struct {
	From  ast.Node
	Stack []ast.Node
}

func NewNodeStack(from ast.Node) *NodeStack {
	return &NodeStack{
		From:  from,
		Stack: []ast.Node{from},
	}
}

func (s *NodeStack) Clone() *NodeStack {
	return &NodeStack{
		From:  s.From,
		Stack: append([]ast.Node{}, s.Stack...),
	}
}

func (s *NodeStack) Push(n ast.Node) *NodeStack {
	s.Stack = append(s.Stack, n)
	return s
}

func (s *NodeStack) Pop() (*NodeStack, ast.Node) {
	l := len(s.Stack)
	if l == 0 {
		return s, nil
	}
	n := s.Stack[l-1]
	s.Stack = s.Stack[:l-1]
	return s, n
}

func (s *NodeStack) IsEmpty() bool {
	return len(s.Stack) == 0
}

func (s *NodeStack) BuildIndexList() []string {
	var indexList []string
	for !s.IsEmpty() {
		_, curr := s.Pop()
		switch curr := curr.(type) {
		case *ast.SuperIndex:
			s = s.Push(curr.Index)
			indexList = append(indexList, "super")
		case *ast.Index:
			s = s.Push(curr.Index)
			s = s.Push(curr.Target)
		case *ast.LiteralString:
			indexList = append(indexList, curr.Value)
		case *ast.Self:
			indexList = append(indexList, "self")
		case *ast.Var:
			indexList = append(indexList, string(curr.Id))
		case *ast.Import:
			indexList = append(indexList, curr.File.Value)
		}
	}
	return indexList
}

func (s *NodeStack) ReorderDesugaredObjects() *NodeStack {
	sort.SliceStable(s.Stack, func(i, j int) bool {
		_, iIsDesugared := s.Stack[i].(*ast.DesugaredObject)
		_, jIsDesugared := s.Stack[j].(*ast.DesugaredObject)
		if !iIsDesugared && !jIsDesugared {
			return false
		}

		iLoc, jLoc := s.Stack[i].Loc(), s.Stack[j].Loc()
		if iLoc.Begin.Line < jLoc.Begin.Line && iLoc.End.Line > jLoc.End.Line {
			return true
		}

		return false
	})
	return s
}
