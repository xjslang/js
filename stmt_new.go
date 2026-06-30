package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

type NewExpr struct {
	ast.BaseExpr
	Layout struct {
		New token.Token
	}
	Value ast.Expr
}

func ParseNewExpr(p *parser.Parser) (node *NewExpr, err error) {
	node = &NewExpr{}
	if node.Layout.New, err = p.Expect(NEW); err != nil {
		return
	}
	if node.Value, err = js.ParseValue(p); err != nil {
		return
	}
	return
}

func PrintNewExpr(p *printer.Printer, node *NewExpr) {
	p.Print(node.Layout.New)
	p.SpPrint(node.Value)
}
