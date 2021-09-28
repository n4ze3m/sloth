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
		t.Errorf("Ident token literal not foobar got=%s",ident.TokenLiteral())
	}
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
