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

import "github.com/maxymania/language-coder/template/parser"
import "encoding/json"
import "bytes"
import "os"


type compscript struct{
	BasePath string
	Parser   string
	Compiler string
}

func LoadAll(f string) (*parser.Parser,*Compiler,error){
	c   := new(Compiler)
	p   := new(parser.Parser)
	buf := new(bytes.Buffer)
	var cs compscript
	z,e := os.Open(f)
	if e!=nil { return nil,nil,e }
	buf.ReadFrom(z)
	z.Close()
	buf.ReadString('\n')
	json.Unmarshal(buf.Bytes(),&cs)
	e = p.LoadFile(join(cs.BasePath,cs.Parser))
	if e!=nil { return nil,nil,e }
	e = c.LoadFile(join(cs.BasePath,cs.Compiler))
	if e!=nil { return nil,nil,e }
	return p,c,nil
}
