package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var THROW = token.RegisterType("throw")

type ThrowStmt struct {
	ast.BaseStmt
	Layout struct {
		Throw token.Token
	}
	Expr ast.Expr
}

func ParseThrowStmt(p *parser.Parser) (node *ThrowStmt, err error) {
	node = &ThrowStmt{}
	if node.Layout.Throw, err = p.Expect(THROW); err != nil {
		return
	}
	if node.Expr, err = p.ParseExpr(); err != nil {
		return
	}
	return
}

func PrintThrowStmt(p *printer.Printer, node *ThrowStmt) {
	p.LnPrint(node.Layout.Throw)
	p.SpPrint(node.Expr)
}
