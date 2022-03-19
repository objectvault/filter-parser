// cSpell:ignore gonic, paulo, ferreira
package main

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
	"os"

	"github.com/objectvault/filter-parser/ast"
	"github.com/objectvault/filter-parser/builder"
	"github.com/objectvault/filter-parser/lexer"
	"github.com/objectvault/filter-parser/parser"
	"github.com/objectvault/filter-parser/syntax"
	"github.com/objectvault/filter-parser/token"
	"github.com/objectvault/filter-parser/transpiler"
)

func main() {
	// f := `and(eq(type,1),contains(alias,"org"))`
	// f := `and(type, eq(type,1))`
	// f := `not(eq(Type, 1))`
	// f := `and(gt(type,1),neq(alias,"org"))`
	f := `and(gt(type,1),neq(alias,"*\" ' \*org*"))`

	fmt.Printf("INPUT [%q] - [%d]\n", f, len(f))

	// Create New Lexer (for Input)
	l := lexer.NewLexer(f)

	for tok := l.NextToken(); ; tok = l.NextToken() {
		fmt.Printf("TOKEN [%q] - [%q]\n", tok.Type, tok.Literal)
		if tok.Type == token.EOL {
			break
		}
	}

	// Parse Input
	p := parser.NewParser(l.Reset())
	rp := p.ParseFilter()
	a, ok := rp.(ast.Node)
	if !ok {
		fmt.Println("FILTER Has Parse Error")
		os.Exit(1)
	}

	// Check Resultant AST
	s := syntax.NewSyntaxChecker(a)
	e := s.Verify()
	if e != nil {
		fmt.Printf("SYNTAX ERROR: %s\n", e.ToString())
		os.Exit(2)
	}

	fmt.Println("FILTER Has Passed Syntax Checker")
	fmt.Printf("AST [%q]\n", a.ToString())

	// Transpile Resultant AST
	t := transpiler.NewTranspileToMysqlWhere(rp.(*ast.Filter), nil)
	rt := t.Transpile()
	rts, ok := rt.(string)
	if ok {
		fmt.Printf("WHERE %s\n", rts)
	} else {
		fmt.Printf("TRANSPILER Error [%s]\n", rt.(*transpiler.TranspilerError).Message)
		os.Exit(3)
	}

	hf := builder.ASTFilter(
		builder.ASTAND(
			builder.ASTGT("type", builder.ASTValue(token.NUMBER, "1")),
			builder.ASTNEQ("alias", builder.ASTValue(token.STRING, `*\" ' \*org*`)),
		),
	)

	// Transpile Resultant AST
	t = transpiler.NewTranspileToMysqlWhere(hf, nil)
	rt = t.Transpile()
	rts, ok = rt.(string)
	if ok {
		fmt.Printf("WHERE %s\n", rts)
	} else {
		fmt.Printf("TRANSPILER Error [%s]\n", rt.(*transpiler.TranspilerError).Message)
		os.Exit(3)
	}
}
