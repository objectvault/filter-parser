package token

/*
 * This file is part of the ObjectVault Project.
 * Copyright (C) 2020-2022 Paulo Ferreira <vault at sourcenotes.org>
 *
 * This work is published under the GNU AGPLv3.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

type TokenType string

type Token struct {
	Type    TokenType
	Literal string
}

const (
	// State
	ILLEGAL = "ILLEGAL"
	EOL     = "EOL"

	// Literals
	IDENT  = "IDENT"
	STRING = "STRING"
	INT    = "INT"
	NUMBER = "NUMBER"

	// Delimiters
	COMMA  = ","
	LPAREN = "("
	RPAREN = ")"
)
