package js

import (
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/plugin"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var (
	NEW       = token.RegisterType("new")
	STRICT_EQ = token.RegisterType("===")
)

func Parse(input []byte) (*js.Program, error) {
	p := PluginBuilder().Build(input)
	return js.ParseProgram(p)
}

func Print(result ast.Node, opts ...printer.Option) (string, error) {
	pr := PrinterBuilder().Build(opts...)
	pr.Print(result)
	return pr.Output()
}

func PluginBuilder() *plugin.Builder {
	return xjs.PluginBuilder().Install(func(b *plugin.Builder) {
		token.RegisterUnaryType(NEW)
		token.RegisterBinaryType(STRICT_EQ, token.EQ.Precedence())

		b.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
			if tok, err = next(); err != nil {
				return
			}
			switch tok.Type {
			case token.IDENT:
				switch tok.Literal {
				case "const":
					tok.Type = CONST
				case "var":
					tok.Type = VAR
				case "try":
					tok.Type = TRY
				case "catch":
					tok.Type = CATCH
				case "finally":
					tok.Type = FINALLY
				case "switch":
					tok.Type = SWITCH
				case "case":
					tok.Type = CASE
				case "default":
					tok.Type = DEFAULT
				case "throw":
					tok.Type = THROW
				case "new":
					tok.Type = NEW
				}
			case token.EQ:
				if sc.CurrentChar() == '=' {
					sc.AdvanceChar()
					tok.Type = STRICT_EQ
					tok.Literal = "==="
				}
			}
			return
		})
		b.UseUnaryParser(func(p *parser.Parser, next func() (ast.Expr, error)) (ast.Expr, error) {
			switch p.CurrentToken.Type {
			case NEW:
				return ParseNewExpr(p)
			}
			return next()
		})
		b.UseBinaryParser(func(p *parser.Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (ast.Expr, error) {
			switch p.CurrentToken.Type {
			case STRICT_EQ:
				return js.ParseBinaryExpr(p, left)
			}
			return next(left)
		})
		b.UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (ast.Stmt, error) {
			switch p.CurrentToken.Type {
			case CONST:
				return ParseConstStmt(p)
			case VAR:
				return ParseVarStmt(p)
			case TRY:
				return ParseTryStmt(p)
			case SWITCH:
				return ParseSwitchStmt(p)
			case THROW:
				return ParseThrowStmt(p)
			}
			return next()
		})
	})
}

func PrinterBuilder() *printer.Builder {
	return xjs.PrinterBuilder().UsePrinter(func(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
		switch v := node.(type) {
		case *ConstStmt:
			PrintConstStmt(pr, v)
		case *VarStmt:
			PrintVarStmt(pr, v)
		case *TryStmt:
			PrintTryStmt(pr, v)
		case *SwitchStmt:
			PrintSwitchStmt(pr, v)
		case *ThrowStmt:
			PrintThrowStmt(pr, v)
		case *NewExpr:
			PrintNewExpr(pr, v)
		default:
			return next(node)
		}
		return nil
	})
}
