package evaluator

import (
	"testing"

	"github.com/nazeemnato/sloth/lexer"
	"github.com/nazeemnato/sloth/object"
	"github.com/nazeemnato/sloth/parser"
)

func TestEvalIntegerExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"-10", -10},
		{"10", 10},
		{"10 + 10", 20},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"1 + 1", 2},
	}

	for _, tt := range tests {
		evaluted := testEval(tt.input)
		testIntegerObject(t, evaluted, tt.expected)
	}
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluted := testEval(tt.input)
		testBooleanObject(t, evaluted, tt.expected)
	}
}

func TestBangExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"10 > 20", false},
		{"10 == 10", true},
		{"10 != 10", false},
		{"10 < 20", true},
		{"true == true", true},
	}

	for _, tt := range tests {
		evaluted := testEval(tt.input)
		testBooleanObject(t, evaluted, tt.expected)
	}
}

func TestIfElseExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{"if(true) { 19 }", 19},
		{"if(false) { 19 }", nil},
		{"if( 1 < 2) { 1 }", 1},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		integer, ok := tt.expected.(int)
		if ok {
			testIntegerObject(t, evaluated, int64(integer))
		} else {
			testNullObject(t, evaluated)
		}
	}
}
func TestReturnExpression(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{`
		if(10 > 1) {
			if (10 > 1) {
				return 10;
			}
			return 1;
		}
		`, 10},
		{"return 10;", 10},
		{"return -1; 2;", -1},
		{"return 5 * 2; 4;", 10},
	}

	for _, tt := range tests {
		evaluted := testEval(tt.input)
		testIntegerObject(t, evaluted, tt.expected)
	}
}

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{`if(10 > 1) {
			if(10 > 1) {
				return true + false;
			}
			return 1
		}`, "unknown operator: BOOLEAN + BOOLEAN"},
		{
			"foobar",
			"identifier not found: foobar",
		},
	}

	for _, tt := range tests {
		evaluted := testEval(tt.input)
		errObj, ok := evaluted.(*object.Error)

		if !ok {
			t.Errorf("no error object returned got=%T (%+v)", evaluted, evaluted)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong message! exptected=%q got=%q", tt.expectedMessage, errObj.Message)
		}
	}
}

func TestVarStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"var x = 10; x;", 10},
		{"var x = 10 * 10; x;", 100},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionObject(t *testing.T) {
	input := "fun(x) { x + 2; }"

	evaluated := testEval(input)

	fun, ok := evaluated.(*object.Function)

	if !ok {
		t.Fatalf("object is not Function got=%T got=%T", evaluated, evaluated)
	}

	if len(fun.Parameters) !=  1 {
		t.Fatalf("function has wrong parameters. Parameters=%+v", fun.Parameters)
	}

	if fun.Parameters[0].String() != "x" {
		t.Fatalf("parameter is not 'x' got=%q", fun.Parameters[0])
	}

	expectedBody := "(x + 2)"

	if fun.Body.String() != expectedBody {
		t.Fatalf("body is not %q got=%q", expectedBody, fun.Body.String())
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	} {
		{"var k = fun(x) { x;}; k(4)", 4},
		{"var add = fun(x,y) { return x + y;}; add(1,2);",3},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestClosures(t *testing.T) {
	input := `
	var add = fun(x) {
		fun(y) {
			x + y;
		};
	};
	var addTwo = add(2)
	addTwo(2)
	`
	testIntegerObject(t, testEval(input), 4)
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	env := object.NewEnviroment()
	program := p.ParseProgram()
	return Eval(program,env)
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)

	if !ok {
		t.Errorf("object is not integer got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("pbject as wrong value got=%d, want=%d", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)

	if !ok {
		t.Errorf("object is not integer got=%T (%+v)", obj, obj)
		return false
	}

	if result.Value != expected {
		t.Errorf("pbject as wrong value got=%t, want=%t", result.Value, expected)
		return false
	}

	return true
}

func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NULL {
		t.Errorf("object is not null got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
