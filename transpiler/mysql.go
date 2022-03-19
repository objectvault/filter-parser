package transpiler

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

type TMapIdentityToField = func(string) string

func reflectIdentityToFieldMapper(identity string) string {
	return identity
}

type TranspileToMysqlWhere struct {
	Transpiler
	Filter      *ast.Filter
	FieldMapper TMapIdentityToField
}

func NewTranspileToMysqlWhere(a *ast.Filter, mapper TMapIdentityToField) *TranspileToMysqlWhere {
	t := &TranspileToMysqlWhere{Filter: a, FieldMapper: reflectIdentityToFieldMapper}
	if mapper != nil {
		t.FieldMapper = mapper
	}

	return t
}

func (c *TranspileToMysqlWhere) Transpile() interface{} {
	f := c.Filter.F
	return c.mysqlFunctionToStatement(nil, f)
}

func (c *TranspileToMysqlWhere) mysqlFunctionToStatement(p *ast.Function, f *ast.Function) interface{} {
	// ASSUMPTION: Filter has been run through Syntax Checker so AST is Correct
	fname := f.Name.Literal
	switch fname {
	case "NOT":
		return c.mysqlLogicalNOT(f)
	case "AND":
		return c.mysqlLogicalAND(f)
	case "OR":
		return c.mysqlLogicalOR(f)
	case "EQ":
		return c.mysqlOperatorEQ(f)
	case "NEQ":
		return c.mysqlOperatorNEQ(f)
	case "GT":
		return c.mysqlOperatorGT(f)
	case "GTE":
		return c.mysqlOperatorGTE(f)
	case "LT":
		return c.mysqlOperatorLT(f)
	case "LTE":
		return c.mysqlOperatorLTE(f)
	case "CONTAINS":
		return c.mysqlOperatorCONTAINS(f)
	case "IN":
		return c.mysqlOperatorIN(f)
	default:
		return &TranspilerError{Message: fmt.Sprintf("Unsupported Funcion [%s]", fname)}
	}
}

// Logical NOT
func (c *TranspileToMysqlWhere) mysqlLogicalNOT(fnot *ast.Function) interface{} {
	// 1st Parameter to NOT Should be a Logical Function or Operator Function
	pf1 := (fnot.Parameters[0]).(*ast.Function)
	r := c.mysqlFunctionToStatement(fnot, pf1)

	// Converted Function?
	rs, ok := r.(string)
	if !ok { // NO: Abort
		return r
	}

	return fmt.Sprintf("NOT(%s)", rs)
}

// Logical AND
func (c *TranspileToMysqlWhere) mysqlLogicalAND(fand *ast.Function) interface{} {
	// There should be 2 Parameter
	// Both Should be ast.Function
	pf1 := (fand.Parameters[0]).(*ast.Function)
	pf2 := (fand.Parameters[1]).(*ast.Function)

	return c.mysqlBinaryLogical(fand, "AND", pf1, pf2)
}

// Logical OR
func (c *TranspileToMysqlWhere) mysqlLogicalOR(fop *ast.Function) interface{} {
	// There should be 2 Parameter
	// Both Should be ast.Function
	pf1 := (fop.Parameters[0]).(*ast.Function)
	pf2 := (fop.Parameters[1]).(*ast.Function)

	return c.mysqlBinaryLogical(fop, "OR", pf1, pf2)
}

// Operator EQ
func (c *TranspileToMysqlWhere) mysqlOperatorEQ(f *ast.Function) interface{} {
	// There should be 2 Parameter
	// Both Should be ast.Value
	pv1 := (f.Parameters[0]).(*ast.Value) // Identifier
	pv2 := (f.Parameters[1]).(*ast.Value)

	return c.mysqlBinaryOperator("=", pv1, pv2)
}

func (c *TranspileToMysqlWhere) mysqlOperatorNEQ(f *ast.Function) interface{} {
	// There should be 2 Parameter
	// Both Should be ast.Value
	pv1 := (f.Parameters[0]).(*ast.Value) // Identifier
	pv2 := (f.Parameters[1]).(*ast.Value)

	return c.mysqlBinaryOperator("!=", pv1, pv2)
}

func (c *TranspileToMysqlWhere) mysqlOperatorGT(f *ast.Function) interface{} {
	// There should be 2 Parameter
	// Both Should be ast.Value
	pv1 := (f.Parameters[0]).(*ast.Value) // Identifier
	pv2 := (f.Parameters[1]).(*ast.Value)

	return c.mysqlBinaryOperator(">", pv1, pv2)
}

func (c *TranspileToMysqlWhere) mysqlOperatorGTE(f *ast.Function) interface{} {
	// There should be 2 Parameter
	// Both Should be ast.Value
	pv1 := (f.Parameters[0]).(*ast.Value) // Identifier
	pv2 := (f.Parameters[1]).(*ast.Value)

	return c.mysqlBinaryOperator(">=", pv1, pv2)
}

func (c *TranspileToMysqlWhere) mysqlOperatorLT(f *ast.Function) interface{} {
	// There should be 2 Parameter
	// Both Should be ast.Value
	pv1 := (f.Parameters[0]).(*ast.Value) // Identifier
	pv2 := (f.Parameters[1]).(*ast.Value)

	return c.mysqlBinaryOperator("<", pv1, pv2)
}

func (c *TranspileToMysqlWhere) mysqlOperatorLTE(f *ast.Function) interface{} {
	// There should be 2 Parameter
	// Both Should be ast.Value
	pv1 := (f.Parameters[0]).(*ast.Value) // Identifier
	pv2 := (f.Parameters[1]).(*ast.Value)

	return c.mysqlBinaryOperator("<=", pv1, pv2)
}

func (c *TranspileToMysqlWhere) mysqlOperatorCONTAINS(f *ast.Function) interface{} {
	// There should be 2 Parameter
	// Both Should be ast.Value
	pv1 := (f.Parameters[0]).(*ast.Value) // Identifier
	pv2 := (f.Parameters[1]).(*ast.Value) // STRING

	// Is Valid Field?
	field := c.FieldMapper(pv1.V.Literal)
	if field == "" { // NO
		return &TranspilerError{Message: fmt.Sprintf("Invalid Field [%s]", pv1.V.Literal)}
	}

	return fmt.Sprintf("%s LIKE %q", field, mysqlEscapeValue(pv2))
}

func (c *TranspileToMysqlWhere) mysqlOperatorIN(f *ast.Function) interface{} {
	// There should be 2 Parameter
	// Both Should be ast.Value
	pv1 := (f.Parameters[0]).(*ast.Value) // Identifier
	pv2 := (f.Parameters[1]).(*ast.Value) // STRING

	// Is Valid Field?
	field := c.FieldMapper(pv1.V.Literal)
	if field == "" { // NO
		return &TranspilerError{Message: fmt.Sprintf("Invalid Field [%s]", pv1.V.Literal)}
	}

	return fmt.Sprintf("%s IN %q", field, mysqlEscapeValue(pv2))
}

// HELPERS //
func (c *TranspileToMysqlWhere) mysqlBinaryLogical(p *ast.Function, op string, f1 *ast.Function, f2 *ast.Function) interface{} {
	// 1st Parameter to NOT Should be a Logical Function or Operator Function
	r1 := c.mysqlFunctionToStatement(p, f1)
	r2 := c.mysqlFunctionToStatement(p, f2)

	// Converted 1st Function?
	rs1, ok := r1.(string)
	if !ok { // NO: Abort
		return rs1
	}

	// Converted 2nd Function?
	rs2, ok := r2.(string)
	if !ok { // NO: Abort
		return rs2
	}

	return fmt.Sprintf("(%s) %s (%s)", rs1, op, rs2)
}

func (c *TranspileToMysqlWhere) mysqlBinaryOperator(op string, f *ast.Value, v *ast.Value) interface{} {
	// 1st Parameter should be an Identifier (Field Name)
	field := c.FieldMapper(f.V.Literal)

	// Is Valid Field?
	if field == "" { // NO
		return &TranspilerError{Message: fmt.Sprintf("Invalid Field [%s]", f.V.Literal)}
	}

	value := mysqlEscapeValue(v)
	if v.V.Type == token.STRING {
		return fmt.Sprintf("%s %s %q", field, op, value)
	}

	return fmt.Sprintf("%s %s %s", field, op, value)
}

func mysqlEscapeValue(v *ast.Value) string {
	if v.V.Type != token.STRING {
		return v.V.Literal
	}

	// Escape the Escape Character
	s := strings.ReplaceAll(v.V.Literal, "\\", "\\\\")

	// Make sure Embedded Quotes are Escaped
	s = strings.ReplaceAll(s, `"`, `\"`)
	s = strings.ReplaceAll(s, `'`, `\'`)

	// Escape Special Characters
	s = strings.ReplaceAll(s, `%`, `\%`)

	// Convert '\uFFFD' (replacement for '*') to '%'
	s = strings.ReplaceAll(s, "�", "%") // Temporarily Store \* as something else

	return s
}
