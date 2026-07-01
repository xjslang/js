package js

import (
	"github.com/xjslang/xjs/ast"
	"github.com/xjslang/xjs/js"
	"github.com/xjslang/xjs/parser"
	"github.com/xjslang/xjs/printer"
	"github.com/xjslang/xjs/token"
)

var IMPORT = token.RegisterType("import")

type ImportStmt struct {
	ast.BaseStmt
	Layout struct {
		Import token.Token
		From   token.Token
	}
	Imports []ast.Node
	Source  token.Token
}

type DefaultImport struct {
	ast.BaseNode
	Name *js.Ident
}

type NamespaceImport struct {
	ast.BaseNode
	Layout struct {
		Multiply token.Token
		As       token.Token
	}
	Name *js.Ident
}

type NamedImport struct {
	ast.BaseNode
	Layout struct {
		As token.Token
	}
	Name  *js.Ident
	Alias *js.Ident
}

type NamedImportList struct {
	ast.BaseNode
	Layout struct {
		Lbrace token.Token
		Rbrace token.Token
	}
	Imports []*NamedImport
}

func ParseImportStmt(p *parser.Parser) (node *ImportStmt, err error) {
	node = &ImportStmt{}
	if node.Layout.Import, err = p.Expect(IMPORT); err != nil {
		return
	}
	if p.CurrentToken.Type == js.STRING {
		node.Source = p.CurrentToken
		p.AdvanceToken()
	} else {
		if p.CurrentToken.Type == token.IDENT {
			n := &DefaultImport{Name: &js.Ident{Token: p.CurrentToken}}
			p.AdvanceToken()
			node.Imports = append(node.Imports, n)
			if p.CurrentToken.Type == token.COMMA {
				p.AdvanceToken()
				var n ast.Node
				if n, err = parseImport(p); err != nil {
					return
				}
				node.Imports = append(node.Imports, n)
			}
		} else if typ := p.CurrentToken.Type; typ == token.MULTIPLY || typ == token.LBRACE {
			var n ast.Node
			if n, err = parseImport(p); err != nil {
				return
			}
			node.Imports = append(node.Imports, n)
		} else {
			err = p.Error("syntax error")
			return
		}
		if node.Layout.From, err = p.ExpectString("from"); err != nil {
			return
		}
		if node.Source, err = p.Expect(js.STRING); err != nil {
			return
		}
	}
	return
}

func parseNamespaceImport(p *parser.Parser) (node *NamespaceImport, err error) {
	node = &NamespaceImport{}
	node.Layout.Multiply = p.CurrentToken
	p.AdvanceToken() // consume *
	if node.Layout.As, err = p.ExpectString("as"); err != nil {
		return
	}
	if node.Name, err = js.ParseIdent(p); err != nil {
		return
	}
	return
}

func parseNamedImport(p *parser.Parser) (node *NamedImportList, err error) {
	node = &NamedImportList{}
	node.Layout.Lbrace = p.CurrentToken
	for p.AdvanceToken(); p.CurrentToken.Type != token.RBRACE; p.AdvanceToken() {
		n := &NamedImport{}
		if n.Name, err = js.ParseIdent(p); err != nil {
			return
		}
		if p.CurrentToken.Literal == "as" {
			n.Layout.As = p.CurrentToken
			p.AdvanceToken()
			if n.Alias, err = js.ParseIdent(p); err != nil {
				return
			}
		}
		node.Imports = append(node.Imports, n)
		if p.CurrentToken.Type != token.COMMA {
			break
		}
	}
	if node.Layout.Rbrace, err = p.Expect(token.RBRACE); err != nil {
		return
	}
	return
}

func parseImport(p *parser.Parser) (node ast.Node, err error) {
	switch p.CurrentToken.Type {
	case token.MULTIPLY:
		if node, err = parseNamespaceImport(p); err != nil {
			return
		}
	case token.LBRACE:
		if node, err = parseNamedImport(p); err != nil {
			return
		}
	default:
		err = p.Error("syntax error")
	}
	return
}

func PrintImportStmt(p *printer.Printer, node *ImportStmt) {
	p.LnPrint(node.Layout.Import)
	if len(node.Imports) > 0 {
		for i, imp := range node.Imports {
			if i > 0 {
				p.Print(",")
				p.EnsureSpace()
			}
			switch v := imp.(type) {
			case *DefaultImport:
				p.SpPrint(v.Name)
			case *NamespaceImport:
				p.SpPrint(v.Layout.Multiply).SpPrint(v.Layout.As).SpPrint(v.Name)
			case *NamedImportList:
				p.SpPrint(v.Layout.Lbrace)
				p.IncreaseIndent()
				for j, namedImp := range v.Imports {
					if j > 0 {
						p.Print(",")
						p.EnsureSpace()
					}
					p.SpPrint(namedImp.Name)
					if namedImp.Alias != nil {
						p.SpPrint(namedImp.Layout.As).SpPrint(namedImp.Alias)
					}
				}
				if len(v.Imports) > 0 {
					p.EnsureSpace()
				}
				p.DecreaseIndent()
				p.Print(v.Layout.Rbrace)
			}
		}
		p.SpPrint(node.Layout.From)
	}
	p.SpPrint(node.Source)
}
