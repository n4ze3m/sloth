package parser_test

import (
	"fmt"
	"github.com/nazeemnato/sloth/ast"
	"github.com/nazeemnato/sloth/lexer"
	"github.com/nazeemnato/sloth/parser"
	"testing"
)

func TestVarStatement(t *testing.T) {
	input := `
	var x = 5;
	var y = 10;
	var foobar = 2232;
	`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParseErrrors(t, p)

	if program == nil {
		t.Fatalf("ParserProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Program statements does not contain 3 statments. got=%d", len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		fmt.Println(program.Statements)
		stmt := program.Statements[i]

		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func TestReturnStatement(t *testing.T) {
	inputr := `return 5;
	return 10;
	return 10000;
	`

	l := lexer.New(inputr)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParseErrrors(t, p)

	if program == nil {
		t.Fatalf("ParserProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("Program statements does not contain 3 statments. got=%d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnsStatement, got=%T", stmt)
			continue
		}
		if returnStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not return got %q", returnStmt.TokenLiteral())
		}

	}
}

func TestIdentifierExpress(t *testing.T) {
	inputr := `foobar;`

	l := lexer.New(inputr)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParseErrrors(t, p)

	if program == nil {
		t.Fatalf("ParserProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements does not contain  statments. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Program.statement[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)

	if !ok {
		t.Fatalf("Exp no *ast.Identifier got=%T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("Ident value not footbar got=%s", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("Ident token literal not foobar got=%s", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpress(t *testing.T) {
	inputr := `5;`

	l := lexer.New(inputr)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParseErrrors(t, p)

	if program == nil {
		t.Fatalf("ParserProgram() returned nil")
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements does not contain  statments. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Program.statement[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.IntegerLiteral)

	if !ok {
		t.Fatalf("Exp no *ast.IntegerLiteral got=%T", stmt.Expression)
	}
	if ident.Value != 5 {
		t.Errorf("Ident value not 5  got=%d", ident.Value)
	}
	if ident.TokenLiteral() != "5" {
		t.Errorf("Ident token literal not 5 got=%s", ident.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTest := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
	}

	for _, tt := range prefixTest {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParseErrrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Program statements does not contain 1  statments. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Program.statement[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt is not ast.PreflixExpress. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s got=%s", tt.operator, exp.Operator)
		}

		if !testInegerLiteral(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func testInegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)

	if !ok {
		t.Errorf("il not *ast.IntegerLiteral, got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d got %d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d got=%s", value, integ.TokenLiteral())
	}

	return true
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "var" {
		t.Errorf("s.Tokenliteral not 'var' got=%q", s.TokenLiteral())
		return false
	}
	varStmt, ok := s.(*ast.VarStatement)

	if !ok {
		t.Errorf("s not *ast.VarStatement got=%q", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s' got=%s", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name.Value not '%s' got=%s", name, varStmt.Name)
		return false
	}

	return true
}

func checkParseErrrors(t *testing.T, p *parser.Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}
	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("Parser error :%q", msg)
	}
	t.FailNow()
}
