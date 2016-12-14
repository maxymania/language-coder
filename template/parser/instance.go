/*
 * Copyright(C) 2015 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 * 
 */


package parser

import (
	"github.com/maxymania/langtransf/syntaxfile"
	"github.com/maxymania/langtransf/parsing"
	"github.com/maxymania/langtransf/ast"
	"io"
)

/*
type ParserInstance struct{
	file syntaxfile.SyntaxFile
	main syntaxfile.Rule
	kwds map[string]string
}
*/

func (p *ParserInstance) Scan(r io.Reader) *parsing.Token {
	return new(parsing.StdScanner).Init(p.kwds,r).FirstToken()
}

func (p *ParserInstance) Parse(head *parsing.Token) (*ast.AST,*syntaxfile.ErrorRecorder) {
	a := new(ast.AST)
	e := new(syntaxfile.ErrorRecorder)
	p.main.Parse(p.file,head,a,e)
	return a,e
}

