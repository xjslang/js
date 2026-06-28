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
	p := xjs.NewBuilder().
		Install(Plugin).
		Build(input)
	return js.ParseProgram(p)
}

func Print(result ast.Node, opts ...printer.Option) (string, error) {
	pr := &printer.Printer{}
	pr.UsePrinter(xjs.Printer)
	pr.UsePrinter(Printer)
	pr.Init(opts...)
	pr.Print(result)
	return pr.Output()
}

func Plugin(b *builder.Builder) {
	token.RegisterBinaryType(STRICT_EQ, token.EQ.Precedence())

	b.UseScanner(func(sc *scanner.Scanner, next func() (token.Token, error)) (tok token.Token, err error) {
		if tok, err = next(); err != nil {
			return
		}
		if tok.Type == token.EQ && sc.CurrentChar() == '=' {
			sc.AdvanceChar()
			tok.Type = STRICT_EQ
			tok.Literal = "==="
		}
		return
	})
	b.UseBinaryParser(func(p *parser.Parser, left ast.Expr, next func(left ast.Expr) (ast.Expr, error)) (ast.Expr, error) {
		if p.CurrentToken.Type == STRICT_EQ {
			return js.ParseBinaryExpr(p, left)
		}
		return next(left)
	})
}

func Printer(pr *printer.Printer, node ast.Node, next func(node ast.Node) error) error {
	return next(node)
}
