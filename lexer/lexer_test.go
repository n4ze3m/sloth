package lexer

import (
	"testing"

	"github.com/nazeemnato/sloth/token"
)

func TestNextToken(t *testing.T) {
	input := `var five = 5;
	var two = 2;
	var add = fun(x,y) {
		x + y;
	};

	var result = add(five,two);

	!-/*5;
	5 < 10 > 5;

	if (5 < 10) {
		return true;
	}else {
		return false;
	}

	5 == 6;
	5 != 6;
	"foobar"
	"foo bar"
	[1,2];
	`

	tests := []struct {
		expectedType    token.TokenType
		exceptedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLN, ";"},

		{token.VAR, "var"},
		{token.IDENT, "two"},
		{token.ASSIGN, "="},
		{token.INT, "2"},
		{token.SEMICOLN, ";"},

		{token.VAR, "var"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fun"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLN, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLN, ";"},

		{token.VAR, "var"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "two"},
		{token.RPAREN, ")"},
		{token.SEMICOLN, ";"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLN, ";"},

		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLN, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLN, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLN, ";"},
		{token.RBRACE, "}"},


		{token.INT, "5"},
		{token.EQ, "=="},
		{token.INT, "6"},
		{token.SEMICOLN, ";"},

		{token.INT, "5"},
		{token.NOT_EQ, "!="},
		{token.INT, "6"},
		{token.SEMICOLN, ";"},

		{token.STRING, "foobar"},
		{token.STRING, "foo bar"},

		{token.LBRACKET, "["},
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "2"},
		{token.RBRACKET, "]"},
		{token.SEMICOLN, ";"},
		
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - token type wrong. exected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.exceptedLiteral {
			t.Fatalf("test[%d] - literal  wrong. exected=%q, got=%q", i, tt.exceptedLiteral, tok.Literal)
		}
	}
}

