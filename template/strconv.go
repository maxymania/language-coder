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
import "strconv"

func strconvUnquote(L *lua.LState) int {
	s := L.CheckString(1)
	n,e := strconv.Unquote(s)
	if e!=nil { panic(e) }
	L.Push(lua.LString(n))
	return 1
}

func strconvParseFloat(L *lua.LState) int {
	s := L.CheckString(1)
	n,e := strconv.ParseFloat(s, 64)
	if e!=nil { panic(e) }
	L.Push(lua.LNumber(n))
	return 1
}
func strconvParseInt(L *lua.LState) int {
	s := L.CheckString(1)
	n,e := strconv.ParseInt(s, 0, 64)
	if e!=nil { panic(e) }
	L.Push(lua.LNumber(float64(n)))
	return 1
}
func strconvParseUInt(L *lua.LState) int {
	s := L.CheckString(1)
	n,e := strconv.ParseUint(s, 0, 64)
	if e!=nil { panic(e) }
	L.Push(lua.LNumber(float64(n)))
	return 1
}

func registerStrconv(L *lua.LState){
	L.SetGlobal("unquote",L.NewFunction(strconvUnquote))
	L.SetGlobal("parse_float",L.NewFunction(strconvParseFloat))
	L.SetGlobal("parse_int",L.NewFunction(strconvParseInt))
	L.SetGlobal("parse_uint",L.NewFunction(strconvParseUInt))
}


