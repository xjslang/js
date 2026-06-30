package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var VAR = token.RegisterType("const")

type VarStmt struct {
	ast.BaseStmt
	Layout struct {
		Var    token.Token
		Assign token.Token
	}
	Name  *js.Ident
	Value ast.Expr
}

func ParseVarStmt(p *parser.Parser) (node *VarStmt, err error) {
	node = &VarStmt{}
	if node.Layout.Var, err = p.Expect(VAR); err != nil {
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

func PrintVarStmt(p *printer.Printer, node *VarStmt) {
	p.LnPrint(node.Layout.Var)
	p.SpPrint(node.Name)
	p.SpPrint(node.Layout.Assign)
	p.SpPrint(node.Value)
}
