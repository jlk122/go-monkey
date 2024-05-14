package parser

import (
	"monkey/ast"
	"monkey/lexer"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("Parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestLetStatments(t *testing.T) {
	// 	inputFail := `
	// let x  5;
	// let  = 10;
	// let 8554;
	// `
	inputSucced := `
let x = 5;
let y = 10;
let foo = 8554;
`
	l := lexer.New(inputSucced)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)
	if program == nil {
		t.Fatalf("ParseProgram failed")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not match test case (3 statements), got %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foo"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

// s ast.Statement is an interface so it can take any type which implemented said interface?
func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral is not `let`, got %q", s)
		return false
	}

	// Type assertion
	// (or converting an interface value to another type, in this case a pointer to ast.LetStatement)
	// s - interface value, checks out
	// *ast.LetStatement - the type were trying to assert s to
	// letStmt - variable holding the asserted value if assertion succesful (nil of assertion fails)
	// ok - bool that indicates assertion success
	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement, got %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s', got %s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s', got %s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 1337;
`
	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram failed")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not match test case (3 statements), got %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("Statment not *ast.ReturnStatement, got %T", stmt)
			continue
		}

		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("TokenLiteral() not 'return', got %s", returnStmt.TokenLiteral())
		}
	}
}
