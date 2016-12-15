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
import "fmt"
import "math/rand"

const luaALLOC = "alloc"

var allocMethods = map[string]lua.LGFunction{
    "alloc": allocAlloc,
}

// Getter and setter for the Person#Name
func allocAlloc(L *lua.LState) int {
	m := L.CheckUserData(1).Value.(map[string]string)
	u := L.Get(2).String()
	s := "v_"+u
	for{
		_,ok := m[s]
		if !ok {
			m[s]=s
			L.Push(lua.LString(s))
			return 1
		}
		s = fmt.Sprintf("v_%x",rand.Int())
	}
	panic("unreachable")
}


func fuAlloc(L *lua.LState) int {
	ud := L.NewUserData()
	ud.Value = make(map[string]string)
	L.SetMetatable(ud, L.GetTypeMetatable(luaALLOC))
	L.Push(ud)
	return 1
}
func cloneStruct(L *lua.LState) int {
	v := L.CheckTable(1)
	mt := L.NewTable()
	L.SetField(mt, "__index", v)
	nt := L.NewTable()
	L.SetMetatable(nt, mt)
	L.Push(nt)
	return 1
}


func registerAlloc(L *lua.LState){
	mt := L.NewTypeMetatable(luaALLOC)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), allocMethods))
	L.SetGlobal("clone",L.NewFunction(cloneStruct))
	L.SetGlobal("alloc",L.NewFunction(fuAlloc))
}

//------------

