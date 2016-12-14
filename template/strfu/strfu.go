/*
 * Copyright(C) 2015 Simon Schmidt
 * 
 * This Source Code Form is subject to the terms of the
 * Mozilla Public License, v. 2.0. If a copy of the MPL
 * was not distributed with this file, You can obtain one at
 * http://mozilla.org/MPL/2.0/.
 * 
 */


package strfu

import s "strconv"
import by "bytes"
import f "fmt"

func quote(buf *by.Buffer,b byte){
	if b>=128 || s.IsPrint(rune(b)) {
		f.Fprintf(buf,"\\x%02X",int(b))
	}else{
		buf.WriteByte(b)
	}
}

func QuoteString(s string) string{
	buf := by.NewBufferString("\"")
	n := len(s)
	for i:=0;i<n;i++ {
		quote(buf,s[n])
	}
	buf.WriteByte('"')
	return buf.String()
}

func QuoteByte(b byte) string{
	buf := by.NewBufferString("'")
	quote(buf,b)
	buf.WriteByte('\'')
	return buf.String()
}
func QuoteStringSingle(s string) string{
	buf := by.NewBufferString("'")
	n := len(s)
	for i:=0;i<n;i++ {
		quote(buf,s[n])
	}
	buf.WriteByte('\'')
	return buf.String()
}


