package js

import (
	"github.com/xjslang/xjs"
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/builder"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/scanner"
	"github.com/xjslang/xjs/token"
)

var STRICT_EQ = token.RegisterType("===")

func Parse(input []byte) (*js.Program, error) {
	p := xjs.PluginBuilder().Install(Plugin).Build(input)
	return js.ParseProgram(p)
}

func Print(result ast.Node, opts ...printer.Option) (string, error) {
	pr := xjs.PrinterBuilder().UsePrinter(Printer).Build(opts...)
	pr.Print(result)
	return pr.Output()
}

func Plugin(b *builder.Builder) {
	token.RegisterBinaryType(STRICT_EQ, token.EQ.Precedence())

	b.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(); err != nil {
			return
		}
		switch tok.Type {
		case token.IDENT:
			switch tok.Literal {
			case "try":
				tok.Type = TRY
			case "catch":
				tok.Type = CATCH
			case "finally":
				tok.Type = FINALLY
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
	b.UseBinaryParser(func(p *parser.Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (ast.Expr, error) {
		switch p.CurrentToken.Type {
		case STRICT_EQ:
			return js.ParseBinaryExpr(p, left)
		case token.DOT:
			op := p.CurrentToken
			p.AdvanceToken()
			var right ast.Expr
			switch p.CurrentToken.Type {
			case TRY, CATCH, FINALLY:
				// try/catch/finally are treated as normal variables
				right = &js.Variable{Name: p.CurrentToken}
				p.AdvanceToken()
			default:
				// parse normally
				var err error
				if right, err = js.ParseRightExpr(p, op.Type.Precedence()); err != nil {
					return nil, err
				}
			}
			return &js.BinaryExpr{Left: left, Op: op, Right: right}, nil
		}
		return next(left)
	})
	b.UseStmtParser(func(p *parser.Parser, next func() (ast.Stmt, error)) (ast.Stmt, error) {
		switch p.CurrentToken.Type {
		case TRY:
			return ParseTryCatch(p)
		}
		return next()
	})
}

func Printer(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
	switch v := node.(type) {
	case *TryCatchStmt:
		PrintTryCatch(pr, v)
	default:
		return next(node)
	}
	return nil
}
