package lexer

/*
 * This file is part of the ObjectVault Project.
 * Copyright (C) 2020-2022 Paulo Ferreira <vault at sourcenotes.org>
 *
 * This work is published under the GNU AGPLv3.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

import (
	"testing"

	"github.com/objectvault/filter-parser/token"
)

func TestLeadingWhiteSpace(t *testing.T) {
	input := " ,"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.COMMA, ","},
		{token.EOL, "\x00"},
	}

	// Create New Lexer (for Input)
	l := NewLexer(input)

	// Run Tests
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestTrailingWhiteSpace(t *testing.T) {
	input := ", "

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.COMMA, ","},
		{token.EOL, "\x00"},
	}

	// Create New Lexer (for Input)
	l := NewLexer(input)

	// Run Tests
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestEmbeddedWhiteSpace(t *testing.T) {
	input := " ( , ) "

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LPAREN, "("},
		{token.COMMA, ","},
		{token.RPAREN, ")"},
		{token.EOL, "\x00"},
	}

	// Create New Lexer (for Input)
	l := NewLexer(input)

	// Run Tests
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestEscapeCharacter(t *testing.T) {
	input := `" \ \\ \\\ \" * \* \\* \\\* \\\\* "`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.STRING, ` \ \ \\ " � * \� \* \\� `},
	}

	// Create New Lexer (for Input)
	l := NewLexer(input)

	// Run Tests
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestLimiters(t *testing.T) {
	input := " ( , ) "

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LPAREN, "("},
		{token.COMMA, ","},
		{token.RPAREN, ")"},
		{token.EOL, "\x00"},
	}

	// Create New Lexer (for Input)
	l := NewLexer(input)

	// Run Tests
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestIdentifiers(t *testing.T) {
	input := "and, or, a_b b"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.IDENT, "and"},
		{token.COMMA, ","},
		{token.IDENT, "or"},
		{token.COMMA, ","},
		{token.IDENT, "a_b"},
		{token.IDENT, "b"},
		{token.EOL, "\x00"},
	}

	// Create New Lexer (for Input)
	l := NewLexer(input)

	// Run Tests
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestValidNumbers(t *testing.T) {
	input := "1,12 34 4.5 .5"

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "1"},
		{token.COMMA, ","},
		{token.INT, "12"},
		{token.INT, "34"},
		{token.NUMBER, "4.5"},
		{token.NUMBER, ".5"},
		{token.EOL, "\x00"},
	}

	// Create New Lexer (for Input)
	l := NewLexer(input)

	// Run Tests
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestInvalidNumbers(t *testing.T) {
	input := " . 5.5.0 .. "

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.ILLEGAL, "."},
		{token.ILLEGAL, "5.5."},
		{token.INT, "0"},
		{token.ILLEGAL, ".."},
		{token.EOL, "\x00"},
	}

	// Create New Lexer (for Input)
	l := NewLexer(input)

	// Run Tests
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestValidStrings(t *testing.T) {
	input := ` "" " " "a" "a\"b" `

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.STRING, ""},
		{token.STRING, " "},
		{token.STRING, "a"},
		{token.STRING, "a\"b"},
		{token.EOL, "\x00"},
	}

	// Create New Lexer (for Input)
	l := NewLexer(input)

	// Run Tests
	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
