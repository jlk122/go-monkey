package parser

import (
	"fmt"
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

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)

	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not match test case (1 statement), got %d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement) // type assertion
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression not *ast.Identifier, got %T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Fatalf("ident.Value not foobar, got %s", ident.Value)
	}

	if ident.TokenLiteral() != "foobar" {
		t.Fatalf("ident.TokenLiteral() not foobar, got %s", ident.Value)
	}
}

func TestIntegerLiteralExpressions(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not match test case (1 statement), got %d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement) // type assertion
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression not *ast.IntegerLiteral, got %T", stmt.Expression)
	}

	if literal.Value != 5 {
		t.Fatalf("literal.Value not 5, got %d", literal.Value)
	}

	if literal.TokenLiteral() != "5" {
		t.Fatalf("ident.TokenLiteral() not 5, got %d", literal.Value)
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	literal, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not ast.IntegerLiteral, got %T", il)
		return false
	}

	if literal.Value != value {
		t.Errorf("literal.Value not %d, got %d", value, literal.Value)
		return false
	}

	if literal.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("literal.TokenLiteral() not %d, got %s", value, literal.TokenLiteral())
		return false
	}

	return true
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-13", "-", 13},
	}

	for _, tt := range prefixTests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not match test case (1 statement), got %d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement) // type assertion
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("statement not *ast.PrefixExpression, got %T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator not '!', got %s", exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"2 + 2;", 2, "+", 2},
		{"2 - 2;", 2, "-", 2},
		{"2 * 2;", 2, "*", 2},
		{"2 / 2;", 2, "/", 2},
		{"2 > 2;", 2, ">", 2},
		{"2 < 2;", 2, "<", 2},
		{"2 == 2;", 2, "==", 2},
		{"2 != 2;", 2, "!=", 2},
	}

	for _, tt := range infixTests {

		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not match test case (1 statement), got %d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement) // type assertion
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement, got %T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("statement not *ast.InfixExpression, got %T", stmt.Expression)
		}

		if !testIntegerLiteral(t, exp.Left, tt.leftValue) {
			return
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator not '%s', got %s", tt.operator, exp.Operator)
		}

		if !testIntegerLiteral(t, exp.Right, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4; -5 * 5", "(3 + 4)((-5) * 5)"},
		{"5 > 4 == 3 < 4", "((5 > 4) == (3 < 4))"},
		{"5 > 4 != 3 < 4", "((5 > 4) != (3 < 4))"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}
