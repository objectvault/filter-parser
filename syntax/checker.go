package syntax

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

	"github.com/objectvault/filter-parser/ast"
	"github.com/objectvault/filter-parser/token"
)

// Syntax Error Object
type SyntaxError struct {
	Message string
}

func (e *SyntaxError) ToString() string {
	return e.Message
}

// Syntax Checker Object
type SyntaxChecker struct {
	AST ast.Node
}

func NewSyntaxChecker(root ast.Node) *SyntaxChecker {
	c := &SyntaxChecker{AST: root}

	return c
}

func (c *SyntaxChecker) Verify() *SyntaxError {
	f, ok := c.AST.(*ast.Filter)

	if !ok {
		e := &SyntaxError{Message: "Invalid AST Object"}
		return e
	}
	return c.verifyFilter(f)
}

func (c *SyntaxChecker) verifyFilter(f *ast.Filter) *SyntaxError {
	fu := f.F
	return c.verifyFunction(nil, fu)
}

func (c *SyntaxChecker) verifyFunction(p *ast.Function, f *ast.Function) *SyntaxError {
	var e *SyntaxError

	// Normalize Function Name (ALL UPPERCASE)
	fname := f.Name.Literal
	fname = strings.ToUpper(fname)
	f.Name.Literal = fname

	t := functionType(fname)
	switch t {
	case "logical-unary":
		if len(f.Parameters) != 1 {
			e = &SyntaxError{Message: fmt.Sprintf("Function [%s] should have 1 parameter, found [%d]", fname, len(f.Parameters))}
			break
		}

		p0, ok := (f.Parameters[0]).(*ast.Function)
		if !ok {
			e = &SyntaxError{Message: fmt.Sprintf("Function [%s] should have a function as only parameter", fname)}
			break
		}

		e = c.verifyFunction(f, p0)
	case "logical-binary":
		if len(f.Parameters) != 2 {
			e = &SyntaxError{Message: fmt.Sprintf("Function [%s] should have 2 parameter, found [%d]", fname, len(f.Parameters))}
			break
		}

		p1, ok := (f.Parameters[0]).(*ast.Function)
		if !ok {
			e = &SyntaxError{Message: fmt.Sprintf("Function [%s] 1st parameter should be a function", fname)}
			break
		}

		p2, ok := (f.Parameters[1]).(*ast.Function)
		if !ok {
			e = &SyntaxError{Message: fmt.Sprintf("Function [%s] 2nd parameter should be a function", fname)}
			break
		}

		// Valid 1st Parameter?
		e = c.verifyFunction(f, p1)
		if e == nil { // YES: Validate Second Parameter
			e = c.verifyFunction(f, p2)
		}
	case "operator":
		// CHECK :Number of Parameters
		if len(f.Parameters) != 2 {
			e = &SyntaxError{Message: fmt.Sprintf("Function [%s] should have 2 parameter, found [%d]", fname, len(f.Parameters))}
			break
		}

		// CHECK: Parameter 1 should be an identifier
		pv1, ok := f.Parameters[0].(*ast.Value)
		if !ok {
			e = &SyntaxError{Message: fmt.Sprintf("Function [%s] invalid type for Parameter 1", fname)}
			break
		}

		if pv1.V.Type != token.IDENT {
			e = &SyntaxError{Message: fmt.Sprintf("Function [%s] Parameter 1 is not a Field Identifier, not [%s]", fname, pv1.V.Type)}
			break
		}

		// Field Names should always be Lower Case
		pv1.V.Literal = strings.ToLower(pv1.V.Literal)

		// CHECK: Parameter 2 should be a Non Identifier Value
		pv2, ok := f.Parameters[1].(*ast.Value)
		if !ok {
			e = &SyntaxError{Message: fmt.Sprintf("Function [%s] invalid type for Parameter 2", fname)}
			break
		}

		if pv2.V.Type == token.IDENT {
			e = &SyntaxError{Message: fmt.Sprintf("Function [%s] Parameter 2 should not be an Identifier", fname)}
			break
		}

		if fname == "CONTAINS" || fname == "IN" {
			if pv2.V.Type != token.STRING {
				e = &SyntaxError{Message: fmt.Sprintf("Function [%s] Parameter 2 should be a String no [%s]", fname, pv2.V.Type)}
				break
			}
		}

	default:
		e = &SyntaxError{Message: fmt.Sprintf("Function [%s] is not recognized", f.Name.Literal)}
	}

	return e
}

func functionType(name string) string {
	switch name {
	case "NOT":
		return "logical-unary"
	case "OR", "AND":
		return "logical-binary"
	case "EQ", "NEQ", "GT", "GTE", "LT", "LTE", "CONTAINS", "IN":
		return "operator"
	}
	return "unknown"
}
