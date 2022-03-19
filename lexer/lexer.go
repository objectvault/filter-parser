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
	"strings"
	"unicode"

	"github.com/objectvault/filter-parser/token"
)

type Lexer struct {
	input        []rune
	position     int  // current position in input (points to current char)
	readPosition int  // current reading position in input (after current char)
	ch           rune // current char under examination
}

func NewLexer(input string) *Lexer {
	// Create Lexer Object
	l := &Lexer{input: []rune(input)}

	// Load 1st Unicode Character
	l.nextChar()
	return l
}

func (l *Lexer) Reset() *Lexer {
	// Reset Positions
	l.position = 0
	l.readPosition = 0

	// Load 1st Unicode Character
	l.nextChar()
	return l
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	// Skip Leading Whitespaces
	l.skipWhiteSpaces()

	// See what we have as the current character
	// start := l.position
	if l.ch == '(' {
		tok = newToken(token.LPAREN, l.ch)
	} else if l.ch == ')' {
		tok = newToken(token.RPAREN, l.ch)
	} else if l.ch == ',' {
		tok = newToken(token.COMMA, l.ch)
	} else if l.ch == 0 { // EOL: Marker
		// NOTE: l.ch contain run '\x00'
		tok = newToken(token.EOL, l.ch)
	} else if isValidNumberRune(l.ch) {
		tok = l.nextTokenNumber()
	} else if isValidFirstIdentifierRune(l.ch) {
		tok = l.nextTokenIdentifier()
	} else if l.ch == '"' {
		tok = l.nextTokenString()
	} else {
		tok = newToken(token.ILLEGAL, l.ch)
	}
	// fmt.Printf("SLICE [%d:%d] [%q] [%d] - REMAIN [%q]\n", start, l.position, string(l.input[start:l.position+1]), l.position-start+1, string(l.input[l.position+1:]))

	// Move Forward in Stream
	l.nextChar()
	return tok
}

func (l *Lexer) nextTokenNumber() token.Token {
	matchedPeriod := false    // Found a Period in the Number
	requireNextDigit := false // Next Character has to be a Digit?
	tt := token.INT           // Set Default Number Type

	// Did we start with a decimal?
	if l.ch == '.' {
		tt = token.NUMBER
		requireNextDigit = true
		matchedPeriod = true
	}

	// MARK Start of Number
	start := l.position
	for l.nextChar(); isValidNumberRune(l.ch); l.nextChar() {
		// Is Next Character a Period
		if l.ch == '.' { // YES: Set Token Type and Flags
			// Another Period?
			if matchedPeriod || requireNextDigit { // YES: Illegal Number
				tt = token.ILLEGAL
				l.nextChar()
				break
			}

			tt = token.NUMBER
			requireNextDigit = true
			matchedPeriod = true
		} else { // NO: Clear Flag
			requireNextDigit = false
		}
	}
	// MARK End of Number +1
	end := l.position

	// Is the Next Digit Flags Still Set?
	if requireNextDigit { // YES: Then this is an illegal number
		tt = token.ILLEGAL
	}

	// BACKUP One Letter (nextToken will move forward one)
	l.undoChar()
	return token.Token{Type: token.TokenType(tt), Literal: string(l.input[start:end])}
}

func (l *Lexer) nextTokenIdentifier() token.Token {
	// MARK Start of Idenitifer
	start := l.position
	for l.nextChar(); isValidNextIdentifierRune(l.ch); l.nextChar() {
	}
	// MARK End of Identifier + 1
	end := l.position
	// BACKUP One Letter (nextToken will move forward one)
	l.undoChar()
	return token.Token{Type: token.IDENT, Literal: string(l.input[start:end])}
}

func (l *Lexer) nextTokenString() token.Token {
	// Peek at Next Character
	nch := l.peekChar(l.readPosition)

	// Is EOL?
	if nch == '\x00' { // YES: Missing Closing "
		return token.Token{Type: token.ILLEGAL, Literal: ""}
	} else if nch == '"' { // ELSE: Empty String
		l.nextChar() // Skip Closing Quote
		return token.Token{Type: token.STRING, Literal: ""}
	}

	foundEOS := false
	var s strings.Builder
	for l.nextChar(); isValidStringRune(l.ch) || l.ch == '"'; l.nextChar() {
		// Found End Quote?
		if l.ch == '"' { // YES: Exit
			foundEOS = true
			break
		}

		// Is Escape Character
		if l.ch == '\\' { // YES
			// Is Valid Escaped Character?
			pch := l.peekChar(l.readPosition)
			if pch == '\\' || pch == '"' { // YES: Consume 1st '\'
				l.nextChar()
			} else if pch == '*' {
				//				s.WriteRune(l.ch)
				l.nextChar()
			}
		} else if l.ch == '*' {
			l.ch = '\uFFFD'
		}

		// Add Character to Builder
		s.WriteRune(l.ch)
	}

	// Find End Quote?
	if !foundEOS { // NO: Invalid String
		return token.Token{Type: token.ILLEGAL, Literal: s.String()}
	}

	return token.Token{Type: token.STRING, Literal: s.String()}
}

func (l *Lexer) isEOL() bool {
	is := l.position >= len(l.input)
	return is
}

func (l *Lexer) nextChar() {
	// Have we reached EOL?
	if l.readPosition >= len(l.input) { // YES
		l.ch = 0

		// Reached EOL
		l.position = l.readPosition
	} else { // NO : Set Next Character
		l.ch = l.input[l.readPosition]

		// Save Current Character Position
		l.position = l.readPosition

		// Set Peek Position for Stream
		l.readPosition += 1
	}
}

func (l *Lexer) nextStringChar() {
	// Do Normal Next Character
	l.nextChar()

	// Is String Escape Caracter?
	if l.ch == '\\' { // YES
		nch := l.peekChar(l.readPosition)
		if nch == '"' {
			l.nextChar()
		}
	}
}

func (l *Lexer) undoChar() {
	// Backup one Character
	l.position -= 1
	l.ch = l.input[l.position]

	// Postion Next Character
	l.readPosition = l.position + 1
}

func (l *Lexer) peekChar(peek int) rune {
	// Looing beyond EOL?
	if peek >= len(l.input) { // YES
		return '\x00'
	} else { // NO : Return Next Unicode Character
		return l.input[peek]
	}
}

func (l *Lexer) skipWhiteSpaces() int {
	count := 0

	for ; !l.isEOL() && isWhiteSpace(l.peekChar(l.position)); l.nextChar() {
		count += 1
	}

	return count
}

func isWhiteSpace(ch rune) bool {
	is := unicode.IsSpace(ch)
	return is
}

func isLetter(ch rune) bool {
	is := unicode.IsLetter(ch)
	return is
}

func isDigit(ch rune) bool {
	is := unicode.IsDigit(ch)
	return is
}

func isValidFirstIdentifierRune(ch rune) bool {
	return isLetter(ch)
}

func isValidNextIdentifierRune(ch rune) bool {
	return isLetter(ch) || (ch == '_')
}

func isValidNumberRune(ch rune) bool {
	return isDigit(ch) || ch == '.'
}

func isValidStringRune(ch rune) bool {
	return unicode.IsPrint(ch)
}

func newToken(tokenType token.TokenType, ch rune) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
