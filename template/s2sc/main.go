/*
 * Copyright(C) 2015 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 * 
 */


package main;

import "fmt"
import "os"
import "github.com/maxymania/language-coder/template/compiler"
import fp "path/filepath"

func pe(e error){
	if e==nil { return }
	fmt.Println(e)
	os.Exit(1)
}

func x(s string) string{
	d,_ := fp.Abs(s)
	return d
}

func main(){
	if len(os.Args) < 4 {
		fmt.Println("usage:",os.Args[0]," CompilerScript SourceFile TargetFile")
		os.Exit(1)
	}
	if x(os.Args[2]) == x(os.Args[3]) {
		fmt.Println("SourceFile == TargetFile")
		os.Exit(1)
	}
	p,c,e := compiler.LoadAll(os.Args[1])
	pe(e)
	i,e := p.CreateInstance()
	pe(e)
	f,e := os.Open(os.Args[2])
	pe(e)
	defer f.Close()
	tok := i.Scan(f)
	a,ee := i.Parse(tok)
	if len(ee.Errors)>0 {
		fmt.Println("Parse error")
		for _,eee := range ee.Errors {
			fmt.Println(os.Args[2],eee.Position,eee.What)
		}
		os.Exit(1)
	}
	s,e := c.Compile(a)
	if e!=nil {
		fmt.Println("Compile error:",e)
		os.Exit(1)
	}
	f2,e := os.Create(os.Args[3])
	pe(e)
	defer f2.Close()
	f2.Write([]byte(s))
}

