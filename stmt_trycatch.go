package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var (
	TRY     = token.RegisterType("try")
	CATCH   = token.RegisterType("catch")
	FINALLY = token.RegisterType("finally")
)

type TryCatchStmt struct {
	ast.BaseStmt
	Layout struct {
		Try     token.Token
		Catch   token.Token
		Lparen  token.Token
		Rparen  token.Token
		Finally token.Token
	}
	Try        *js.BlockStmt
	Catch      *js.BlockStmt
	CatchParam *js.Ident
	Finally    *js.BlockStmt
}

// A semicolon at the end is not necessary.
func (node *TryCatchStmt) SelfClosing() bool {
	return true
}

func ParseTryCatch(p *parser.Parser) (node *TryCatchStmt, err error) {
	node = &TryCatchStmt{}
	if node.Layout.Try, err = p.Expect(TRY); err != nil {
		return
	}
	if node.Try, err = js.ParseBlockStmt(p); err != nil {
		return
	}
	if p.CurrentToken.Type == CATCH {
		node.Layout.Catch = p.CurrentToken
		p.AdvanceToken()
		if p.CurrentToken.Type == token.LPAREN {
			node.Layout.Lparen = p.CurrentToken
			p.AdvanceToken()
			if node.CatchParam, err = js.ParseIdent(p); err != nil {
				return
			}
			if node.Layout.Rparen, err = p.Expect(token.RPAREN); err != nil {
				return
			}
		}
		if node.Catch, err = js.ParseBlockStmt(p); err != nil {
			return
		}
	}
	if p.CurrentToken.Type == FINALLY {
		node.Layout.Finally = p.CurrentToken
		p.AdvanceToken()
		if node.Finally, err = js.ParseBlockStmt(p); err != nil {
			return
		}
	}
	if node.Catch == nil && node.Finally == nil {
		err = p.Error("missing catch or finally after try")
	}
	return
}

func PrintTryCatch(p *printer.Printer, node *TryCatchStmt) {
	p.LnPrint(node.Layout.Try)
	p.SpPrint(node.Try)
	if node.Catch != nil {
		p.SpPrint(node.Layout.Catch)
		if node.CatchParam != nil {
			p.SpPrint(node.Layout.Lparen)
			p.IncreaseIndent()
			p.Print(node.CatchParam)
			p.DecreaseIndent()
			p.Print(node.Layout.Rparen)
		}
		p.SpPrint(node.Catch)
	}
	if node.Finally != nil {
		p.SpPrint(node.Layout.Finally)
		p.SpPrint(node.Finally)
	}
}
