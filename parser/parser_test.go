package parser_test

import (
	"fmt"
	"testing"

	"github.com/nazeemnato/sloth/ast"
	"github.com/nazeemnato/sloth/lexer"
	"github.com/nazeemnato/sloth/parser"
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
		integerValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
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
			t.Fatalf("Program.statement[0] is not .EastxpressionStatement. got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.PrefixExpression)

		if !ok {
			t.Fatalf("stmt is not ast.PreflixExpress. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %s got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInflixExpression(t *testing.T) {
	inflixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"false == true", false, "==", true},
		{"false != true", false, "!=", true},
	}

	for _, tt := range inflixTests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParseErrrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("Program statements does not contain 1  statments. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Program.statement[0] is not .EastxpressionStatement. got=%T", program.Statements[0])
		}
		exp, ok := stmt.Expression.(*ast.InflixExpression)

		if !ok {
			t.Fatalf("exp is not ast.InflixExpress got=%T", stmt.Expression)
		}

		if !testLiteralExpression(t, exp.Left, tt.leftValue) {
			return
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s' got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.rightValue) {
			return
		}

	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"a + add(b * c) + d", "((a + add((b * c))) + d)"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a+b+c", "((a + b) + c)"},
		{"true", "true"},
		{"false", "false"},
		{"3 > 5 == false", "((3 > 5) == false)"},
		{"a * [1,2,3,5]", "((a * [1, 2, 3, 5]))"},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParseErrrors(t, p)

		actual := program.String()

		if actual != tt.expected {
			t.Errorf("excpected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	inputr := `true;`

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

	ident, ok := stmt.Expression.(*ast.Boolean)

	if !ok {
		t.Fatalf("Exp no *ast.Identifier got=%T", stmt.Expression)
	}
	if ident.Value != true {
		t.Errorf("Ident value not footbar got=%v", ident.Value)
	}
	if ident.TokenLiteral() != "true" {
		t.Errorf("Ident token literal not foobar got=%s", ident.TokenLiteral())
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParseErrrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements does not contain  statments. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Program.statement[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.IfExpression)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression, got=%T", stmt.Expression)
	}

	if !testInflixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("Consequence is not 1 statements got=%d\n", len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement got=%T", exp.Consequence.Statements[0])
	}

	if !testIdentifer(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative == nil {
		t.Errorf("exp.Alternative.Statements was not nil got=%+v", exp.Alternative)
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `fun(x, y) { x + y; }`

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParseErrrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements does not contain  statments. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Program.statement[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)

	if !ok {
		t.Fatalf("stmt.Expression is not ast.FuncionLiteral got=%T", stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("functin literal parameters wrong want 2 got=%d\n", len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")

	if len(function.Body.Statements) != 1 {
		t.Fatalf("Function body statements has not 1 statements got=%d\n", len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)

	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement got=%T", function.Body.Statements[0])
	}
	testInflixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestFunctionLiteralParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []string
	}{
		{input: "fun(){};", expectedParams: []string{}},
		{input: "fun(x){};", expectedParams: []string{"x"}},
		{input: "fun(x, y, z){};", expectedParams: []string{"x", "y", "z"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		checkParseErrrors(t, p)

		stmt := program.Statements[0].(*ast.ExpressionStatement)
		function := stmt.Expression.(*ast.FunctionLiteral)

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Errorf("length parameter wrong want %d got %d", len(tt.expectedParams), len(function.Parameters))
		}

		for i, ident := range tt.expectedParams {
			testLiteralExpression(t, function.Parameters[i], ident)
		}
	}
}

func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := parser.New(l)

	program := p.ParseProgram()
	checkParseErrrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("Program statements does not contain  statments. got=%d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Program.statement[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)

	if !ok {
		t.Fatalf("stmt.Expression is unknow got=%T", stmt.Expression)
	}

	if !testIdentifer(t, exp.Function, "add") {
		return
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("Length argument is not 3 got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0],1)
	testInflixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInflixExpression(t, exp.Arguments[2], 4, "+", 5)

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

func testIdentifer(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("Exp not *ast.Identifier got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("Ident.Value not %s got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s got =%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testInegerLiteral(t, exp, int64(v))
	case int64:
		return testInegerLiteral(t, exp, v)
	case string:
		return testIdentifer(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("Type of exp not handled got=%T", exp)
	return false
}

func testInflixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InflixExpression)

	if !ok {
		t.Errorf("exp is not ast.OperatorExpression got=%T, %s", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("Exp operator is not %s got =%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t got=%s", value, bo.TokenLiteral())
		return false
	}

	return true
}
