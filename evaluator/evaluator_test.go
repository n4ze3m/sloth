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

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()
	return Eval(program)
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
