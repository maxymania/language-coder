package main;

import "github.com/maxymania/language-coder/template"
import "github.com/maxymania/language-coder/template/parser"
import "fmt"
import "os"

func pe(e error){
	if e==nil { return }
	fmt.Println(e)
	os.Exit(1)
}

func main(){
	var p parser.Parser
	pe(p.LoadFile("parser.xml"))
	i,e := p.CreateInstance()
	pe(e)
	f,e := os.Open("example")
	pe(e)
	defer f.Close()
	tok := i.Scan(f)
	a,ee := i.Parse(tok)
	fmt.Println(a)
	fmt.Println(ee)
	s,e := template.Process(a,"process.lua")
	fmt.Println(s,e)
}

