/*
 * Copyright(C) 2015 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 * 
 */


package template

import "github.com/yuin/gopher-lua"
import "github.com/maxymania/langtransf/ast"
import "fmt"
import lib_json "github.com/layeh/gopher-json"

const luaAST = "ast"

const icode = `
P.select = function (tab,ast,...)
	local fu = tab[ast.head]
	if not fu then fu = tab.not_found end
	if not fu then fu = function() end end
	return fu(ast,...)
end
`

func astPush(L *lua.LState, a *ast.AST) {
	ud := L.NewUserData()
	ud.Value = a
	L.SetMetatable(ud, L.GetTypeMetatable(luaAST))
	L.Push(ud)
}

func getAst(L *lua.LState, a *ast.AST) lua.LValue {
	ud := L.NewUserData()
	ud.Value = a
	L.SetMetatable(ud, L.GetTypeMetatable(luaAST))
	return ud
}

func astIndex(L *lua.LState) int {
	a := L.CheckUserData(1).Value.(*ast.AST)
	v := L.Get(2)
	switch v.Type() {
		case lua.LTNumber:
			i := L.CheckInt(2)
			if i>=0 && i<len(a.Childs) {
				astPush(L,a.Childs[i])
				return 1
			}
		case lua.LTString:
		switch L.CheckString(2) {
			case "head":
				L.Push(lua.LString(a.Head))
				return 1
			case "content":
				L.Push(lua.LString(a.Content))
				return 1
			case "string":
				L.Push(lua.LString(a.String()))
				return 1
			case "position":
				L.Push(lua.LString(fmt.Sprint(a.Token.Position())))
				return 1
			case "len":
				L.Push(lua.LNumber(float64(len(a.Childs))))
				return 1
			case "last":
				L.Push(lua.LNumber(float64(len(a.Childs)-1)))
				return 1
		}
	}
	L.Push(lua.LNil)
	return 1
}

func registerAST(L *lua.LState){
	mt := L.NewTypeMetatable(luaAST)
	L.SetField(mt, "__index", L.NewFunction(astIndex))
}

func i_process(a *ast.AST,fu func(L *lua.LState) error) (string,error){
	L := lua.NewState()
	defer L.Close()
	lua.OpenBase(L)
	lua.OpenPackage(L)
	registerAST(L)
	registerBuffer(L)
	registerTemplate(L)
	registerStrconv(L)
	lib_json.Preload(L)
	L.SetGlobal("P",L.NewTable())
	L.SetGlobal("U",L.NewTable())
	t := L.NewTable()
	dse := L.DoString(icode)
	if dse!=nil {
		fmt.Println(dse); panic("Can't Run Internal Code.")
	}
	L.SetGlobal("call",t)
	L.SetField(t,"source",getAst(L,a))
	e := fu(L)
	s := L.GetField(t,"result").String()
	return s,e
}

func Process(a *ast.AST,fs...string) (string,error){
	return i_process(a,func (L *lua.LState) error{
		for _,f := range fs {
			e := L.DoFile(f)
			if e!=nil { return e }
		}
		return nil
	})
}

