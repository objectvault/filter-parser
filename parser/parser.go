package parser

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
	"fmt"

	"github.com/objectvault/filter-parser/ast"
	"github.com/objectvault/filter-parser/lexer"
	"github.com/objectvault/filter-parser/token"
)

// Single Token Look ahead
type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	peekToken token.Token
}

func NewParser(l *lexer.Lexer) *Parser {
	p := &Parser{l: l}

	p.nextToken() // Set 1st Token as Peek Token
	p.nextToken() // Set 1st Token as Current Token
	return p
}

func (p *Parser) ParseFilter() interface{} {
	// Is 1st Token an Identifier?
	if p.curToken.Type == token.IDENT { // YES
		// Looks like the start of a function
		name := p.curToken
		p.nextToken()
		f := p.parseFunction(name)

		// Parsed Function Parameters without Errors?
		if _, ok := f.(*ast.ParseError); ok { // NO: Stop Parsing
			return f
		}

		// Reached End of Line?
		if p.curToken.Type != token.EOL { // NO: Error
			return &ast.ParseError{Message: "FILTER: End-of-line Expected"}
		}

		if fast, ok := f.(*ast.Function); ok {
			return &ast.Filter{F: fast}
		}

		return f
	}

	return &ast.ParseError{Message: "FILTER: expecting function name"}
}

func (p *Parser) parseFunction(name token.Token) interface{} {
	// Expecting "("
	if p.curToken.Type != token.LPAREN { // NOT FOUND
		return &ast.ParseError{Message: "FUNCTION: expecting function \"(\""}
	}

	// Consume LPAREN
	p.nextToken()

	// Expecting IDENTIFIER (1st Parameter is alwasy an IDENTIFIER)
	if p.curToken.Type != token.IDENT { // NOT FOUND
		return &ast.ParseError{Message: "FUNCTION: expecting function IDENTIFIER"}
	}

	// Create Function AST
	f := &ast.Function{Name: name}

	// Able to Parse Parameters?
	params := p.parseFunctionParameters()

	// Parsed Function Parameters without Errors?
	if _, ok := params.(*ast.ParseError); ok { // NO: Stop Parsing
		return params
	}

	// ELSE: Parameters Parsed Okay
	f.Parameters = params.([]interface{})

	// Expecting ")"
	if p.curToken.Type != token.RPAREN { // NOT FOUND
		return &ast.ParseError{Message: "FUNCTION: expecting function \")\""}
	}

	// Consume RPAREN
	p.nextToken()
	return f
}

func (p *Parser) parseFunctionParameters() interface{} {
	// p.curToken.Type == token.IDENT
	params := make([]interface{}, 0)

	// LET Syntax Checker Worry about if the Parameter Type is Valid for it's use
	var n interface{}
	finished := false
	for current := p.nextToken(); ; current = p.nextToken() {

		switch p.curToken.Type {
		case token.LPAREN:
			n = p.parseFunction(current)

			// Parsed Function without Errors?
			if _, ok := n.(*ast.ParseError); ok { // NO: Stop Parsing
				finished = true
			} else                              // Any More Parameters?
			if p.curToken.Type == token.COMMA { // YES: Position at Start of Next Parameter
				p.nextToken()
			}
		case token.COMMA:
			n = &ast.Value{V: current}
			p.nextToken() // Consume ','
		case token.RPAREN:
			n = &ast.Value{V: current}
			finished = true
		default:
			return &ast.ParseError{Message: fmt.Sprintf("FUNCTION PARAMS: unexpected token type [%q]\n", p.peekToken.Type)}
		}

		params = append(params, n)

		// Finished Parsing Parameter List?
		if finished || p.curToken.Type == token.RPAREN { // YES
			break
		}
	}

	return params
}

func (p *Parser) nextToken() token.Token {
	current := p.curToken
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
	return current
}
