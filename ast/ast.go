package ast

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
	"strings"

	"github.com/objectvault/filter-parser/token"
)

/*
  Filter ::= Function
  Function ::= <IDENTIFIER> "(" Parameters ")"
  Parameters ::= Function |
                 <IDENTIFIER> |
                 <IDENTIFIER> "," ValueList
  ParameterList :== Value |
                    Value "," ParameterList
  Value ::= <STRING> | <INT> | <NUMBER>

  PARSE RULES:
  - There are 2 types of functions (LOGICAL OPERATOR, FIELD OPERATORS)
  - They can be further subdivided into
  -- UNARY : Single Parameter
  -- BINARY : 2 Parameters
  - They can have an optional final parameter, that would allow for passing
    extra conditions to be passed into the interpreted so as to allow for
    different types of processing

  SYNTAX CHECKER
  - FUNCTIONS can be:
  -- LOGICAL UNARY not(eq(..,..))
  -- LOGICAL BINARY and(eq(...), eq(...))
  -- FIELD OPERATORS eq(...,...), neq, contains, gt, lt, gte, lte, etc

  FUNCTION UNARY: example not
  - ACCEPTS one parameter of type function
*/
type Node interface {
	ToString() string
}

type Value struct {
	Node
	V token.Token
}

type Function struct {
	Node
	Name       token.Token
	Parameters []interface{}
}

type Filter struct {
	Node
	F *Function
}

type ParseError struct {
	Node
	Message string
}

func (vs *Value) ToString() string {
	if vs.V.Type == token.STRING {
		s := strings.ReplaceAll(vs.V.Literal, "*", `\*`)
		s = strings.ReplaceAll(s, "�", "*")
		return fmt.Sprintf("\"%s\"", s)
	}
	return vs.V.Literal
}

func (fs *Filter) ToString() string {
	if fs.F == nil {
		return "nil"
	} else {
		return fs.F.ToString()
	}
}

func (fs *Function) ToString() string {
	comma := false
	var buffer strings.Builder

	// Loop over Parameters
	for _, ni := range fs.Parameters {
		node := ni.(Node)
		// Need to Append ","
		if comma { // YES
			buffer.WriteString(", ")
		}
		// ELSE: NO
		buffer.WriteString(node.ToString())
		comma = true
	}

	return fmt.Sprintf("%s ( %s )", fs.Name.Literal, buffer.String())
}

func (pes *ParseError) ToString() string {
	return fmt.Sprintf("ERROR [%s]", pes.Message)
}
