/*
 * Copyright(C) 2015 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 * 
 */


package compiler

import (
	"io"
	"encoding/xml"
	"github.com/maxymania/langtransf/ast"
	"github.com/maxymania/language-coder/template"
	"path/filepath"
	"os"
)

func join(a, b string) string {
	if a=="" { return b }
	if b=="" { return a }
	return filepath.Join(a,b)
}

type Compiler struct{
	BasePath string
	Script   []string
	scripts  []string
}

func (p *Compiler) LoadFile(f string) error{
	r,e := os.Open(f)
	if e!=nil { return e }
	defer r.Close()
	return p.Load(r)
}

func (p *Compiler) Load(r io.Reader) error{
	d := xml.NewDecoder(r)
	e := d.Decode(p)
	if e==nil {
		p.scripts = make([]string,len(p.Script))
		for i,s := range p.Script {
			p.scripts[i] = join(p.BasePath,s)
		}
	}
	return e
}


func (p *Compiler) Compile(a *ast.AST) (string,error){
	return template.Process(a,p.scripts...)
}

