package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var CONST = token.RegisterType("const")

type ConstStmt struct {
	ast.BaseStmt
	Layout struct {
		Const  token.Token
		Assign token.Token
	}
	Name  *js.Ident
	Value ast.Expr
}

func ParseConstStmt(p *parser.Parser) (node *ConstStmt, err error) {
	node = &ConstStmt{}
	if node.Layout.Const, err = p.Expect(CONST); err != nil {
		return
	}
	if node.Name, err = js.ParseIdent(p); err != nil {
		return
	}
	if node.Layout.Assign, err = p.Expect(token.ASSIGN); err != nil {
		return
	}
	node.Value, err = p.ParseExpr()
	if err != nil {
		return
	}
	return
}

func PrintConstStmt(p *printer.Printer, node *ConstStmt) {
	p.LnPrint(node.Layout.Const)
	p.SpPrint(node.Name)
	p.SpPrint(node.Layout.Assign)
	p.SpPrint(node.Value)
}
