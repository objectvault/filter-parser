package builder

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

	"github.com/objectvault/filter-parser/ast"
	"github.com/objectvault/filter-parser/token"
)

func ASTFilter(f *ast.Function) *ast.Filter {
	filter := &ast.Filter{F: f}
	return filter
}

func ASTNOT(f *ast.Function, rhs *ast.Function) *ast.Function {
	return buildUnaryFunction("NOT", f)
}

func ASTAND(lhs *ast.Function, rhs *ast.Function) *ast.Function {
	return buildBinaryFunction("AND", lhs, rhs)
}

func ASTOR(lhs *ast.Function, rhs *ast.Function) *ast.Function {
	return buildBinaryFunction("OR", lhs, rhs)
}

func ASTValue(t token.TokenType, v string) *ast.Value {
	return &ast.Value{V: buildToken(t, v)}
}

func ASTEQ(field string, value *ast.Value) *ast.Function {
	return buildOperatorFunction("EQ", field, value)
}

func ASTNEQ(field string, value *ast.Value) *ast.Function {
	return buildOperatorFunction("NEQ", field, value)
}

func ASTGT(field string, value *ast.Value) *ast.Function {
	return buildOperatorFunction("GT", field, value)
}

func ASTGTE(field string, value *ast.Value) *ast.Function {
	return buildOperatorFunction("GTE", field, value)
}

func ASTLT(field string, value *ast.Value) *ast.Function {
	return buildOperatorFunction("LT", field, value)
}

func ASTLTE(field string, value *ast.Value) *ast.Function {
	return buildOperatorFunction("LTE", field, value)
}

func ASTCONTAINS(field string, value *ast.Value) *ast.Function {
	return buildOperatorFunction("CONTAINS", field, value)
}

func ASTIN(field string, value *ast.Value) *ast.Function {
	return buildOperatorFunction("IN", field, value)
}

func buildUnaryFunction(name string, p interface{}) *ast.Function {
	fname := buildToken(token.IDENT, name)
	f := &ast.Function{Name: fname, Parameters: []interface{}{p}}
	return f
}

func buildBinaryFunction(name string, lhs interface{}, rhs interface{}) *ast.Function {
	fname := buildToken(token.IDENT, name)
	f := &ast.Function{Name: fname, Parameters: []interface{}{lhs, rhs}}
	return f
}

func buildOperatorFunction(name string, field string, value *ast.Value) *ast.Function {
	fname := buildToken(token.IDENT, name)
	lhs := ASTValue(token.IDENT, strings.ToLower(field))
	f := &ast.Function{Name: fname, Parameters: []interface{}{lhs, value}}
	return f
}

func buildToken(t token.TokenType, v string) token.Token {
	return token.Token{Type: t, Literal: v}
}
