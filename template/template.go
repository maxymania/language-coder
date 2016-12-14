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
import "text/template"
import "bytes"
import "github.com/maxymania/language-coder/template/strfu"

const luaTEMPLATE = "template"

const luaTDATA = "tdata"

var templateMethods = map[string]lua.LGFunction{
    "render":    templateRender,
    "renderToS": templateRenderToString,
}

func templateRender(L *lua.LState) int {
	te   := L.CheckUserData(1).Value.(*template.Template)
	buf  := L.CheckUserData(2).Value.(*bytes.Buffer)
	data := L.CheckUserData(3).Value
	e := te.Execute(buf,data)
	if e!=nil { panic(e) }
	return 0
}

func templateRenderToString(L *lua.LState) int {
	te   := L.CheckUserData(1).Value.(*template.Template)
	buf  := new(bytes.Buffer)
	data := L.CheckUserData(2).Value
	e := te.Execute(buf,data)
	if e!=nil { panic(e) }
	return 0
}

func templateNew(L *lua.LState) int {
	ud := L.NewUserData()
	temp := template.Must(template.New("literal").Parse(L.CheckString(1)))
	temp.Funcs(template.FuncMap{
		"quote":strfu.QuoteString,
		"quotesingle":strfu.QuoteStringSingle,
	})
	ud.Value = temp
	L.SetMetatable(ud, L.GetTypeMetatable(luaTEMPLATE))
	L.Push(ud)
	return 1
}

func tdataIndex(L *lua.LState) int {
	mp := L.CheckUserData(1).Value.(map[string]string)
	name := L.Get(2).String()
	s,ok := mp[name]
	if ok {
		L.Push(lua.LString(s))
	}else{
		L.Push(lua.LNil)
	}
	return 1
}
func tdataNewIndex(L *lua.LState) int {
	mp := L.CheckUserData(1).Value.(map[string]string)
	name := L.Get(2).String()
	value := L.Get(3)
	if value==lua.LNil {
		delete(mp,name)
	}else{
		mp[name] = value.String()
	}
	return 0
}
func tdataNew(L *lua.LState) int {
	ud := L.NewUserData()
	ud.Value = make(map[string]string)
	L.SetMetatable(ud, L.GetTypeMetatable(luaTDATA))
	L.Push(ud)
	return 1
}

func registerTemplate(L *lua.LState){
	mt := L.NewTypeMetatable(luaTDATA)
	L.SetField(mt, "__index", L.NewFunction(tdataIndex))
	L.SetField(mt, "__newindex", L.NewFunction(tdataNewIndex))
	L.SetGlobal("tdata",L.NewFunction(tdataNew))
	
	mt = L.NewTypeMetatable(luaTEMPLATE)
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), templateMethods))
	L.SetGlobal("template",L.NewFunction(templateNew))
}


