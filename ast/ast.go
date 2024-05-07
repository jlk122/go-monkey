package ast

import "monkey/token"

type Node interface {
	TokenLiteral() string // for testing and debugging
}

// We use these dummy methods behind the interfaces, help the Go compiler with throwing errors if wrong one is used
type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program contain statements, makes sense'
// is the root node of every AST the parser produces
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

type Identifier struct {
	Token token.Token // IDENT token
	Value Expression
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }

type LetStatements struct {
	Token token.Token // LET token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatements) statementNode()       {}
func (ls *LetStatements) TokenLiteral() string { return ls.Token.Literal }
