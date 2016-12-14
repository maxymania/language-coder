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
	"io"
	"encoding/xml"
	"github.com/maxymania/langtransf/syntaxfile"
	"errors"
	"fmt"
	"bytes"
	"path/filepath"
	"os"
)

var E_NO_START_RULE = errors.New("Start rule has not been found")

func join(a, b string) string {
	if a=="" { return b }
	if b=="" { return a }
	return filepath.Join(a,b)
}

type Parser struct{
	BasePath   string
	Syntaxfile []string
	Start      string
}

type ParserInstance struct{
	file syntaxfile.SyntaxFile
	main syntaxfile.Rule
	kwds map[string]string
}

func (p *Parser) LoadFile(f string) error{
	r,e := os.Open(f)
	if e!=nil { return e }
	defer r.Close()
	return p.Load(r)
}

func (p *Parser) Load(r io.Reader) error{
	d := xml.NewDecoder(r)
	return d.Decode(p)
}

func (p *Parser) parse() (f syntaxfile.SyntaxFile,e error){
	buf := new(bytes.Buffer)
	for _,sf := range p.Syntaxfile {
		sfp := join(p.BasePath,sf)
		stream,ee := os.Open(sfp)
		if ee!=nil { e=ee; return }
		_,ee = buf.ReadFrom(stream)
		stream.Close()
		if ee!=nil { e=ee; return }
	}
	defer func(){
		c := recover()
		if c==nil { return }
		if ee,ok := c.(error); ok {
			e = ee
		} else {
			e = errors.New(fmt.Sprint(c))
		}
	}()
	f = syntaxfile.Parse(buf)
	return
}
func (p *Parser) CreateInstance() (*ParserInstance,error){
	f,e := p.parse()
	if e!=nil { return nil,e }
	r,ok := f[p.Start]
	if !ok { return nil,E_NO_START_RULE }
	kw := make(map[string]string)
	f.ScanForKeyWords(kw)
	return &ParserInstance{f,r,kw},nil
}

