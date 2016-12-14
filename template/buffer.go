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
import "bytes"

const luaSTRBUF = "strbuf"

var strbufMethods = map[string]lua.LGFunction{
    "add": strbufAdd,
    "str": strbufStr,
}

// Getter and setter for the Person#Name
func strbufAdd(L *lua.LState) int {
	a := L.CheckUserData(1).Value.(*bytes.Buffer)
	b := L.Get(2).String()
	a.WriteString(b)
	return 0
}

func strbufStr(L *lua.LState) int {
	a := L.CheckUserData(1).Value.(*bytes.Buffer)
	L.Push(lua.LString(a.String()))
	return 1
}

func fuBuffer(L *lua.LState) int {
	ud := L.NewUserData()
	ud.Value = new(bytes.Buffer)
	L.SetMetatable(ud, L.GetTypeMetatable(luaSTRBUF))
	L.Push(ud)
	return 1
}

func registerBuffer(L *lua.LState){
	mt := L.NewTypeMetatable(luaSTRBUF)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), strbufMethods))
	L.SetField(mt, "__tostring", L.NewFunction(strbufStr))
	L.SetGlobal("buffer",L.NewFunction(fuBuffer))
}

//------------

