/*
      Copyright(C) 2015 Simon Schmidt

      This Source Code Form is subject to the
      terms of the Mozilla Public License, v.
      2.0. If a copy of the MPL was not
      distributed with this file, You can
      obtain one at
      http://mozilla.org/MPL/2.0/.
*/

/*
 An SSA to code package. Intended for lua target for now.
*/
package ssacomp

type Statement struct {
	Dest []string
	Args []string
	Expr string
	Jump string
	Mark string
}
func Jump(d,ex string, args...string) Statement {
	return Statement{Args:args,Expr:ex,Jump:d}
}
func Cmd(ex string, args...string) Statement {
	return Statement{Args:args,Expr:ex}
}
func Expr(d,ex string, args...string) Statement {
	return Statement{Dest: []string{d}, Args:args,Expr:ex}
}
func Expr2(d []string,ex string, args...string) Statement {
	return Statement{Dest:d, Args:args,Expr:ex}
}
func Marker(d string) Statement {
	return Statement{Mark:d}
}


