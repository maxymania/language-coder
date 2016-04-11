/*
      Copyright(C) 2015 Simon Schmidt

      This Source Code Form is subject to the
      terms of the Mozilla Public License, v.
      2.0. If a copy of the MPL was not
      distributed with this file, You can
      obtain one at
      http://mozilla.org/MPL/2.0/.
*/

package ssacomp

import "fmt"
import "io"
import "strconv"
import "os"
import "bytes"

/*
type Statement struct {
	Dest []string
	Args []string
	Expr string
	Jump string
	Mark string
}
*/
type cset map[string]bool
func (c cset) set(s string){
	c[s]=true
}
func (c cset) get(s string) bool{
	r,_ := c[s]
	return r
}

type ctable map[string]int
func (c ctable)incr(s string){
	i,_ := c[s]
	c[s] = i+1
}
func (c ctable)decr(s string) int{
	i,_ := c[s]
	c[s] = i-1
	return i-1
}
func (c ctable)get(s string) int{
	i,_ := c[s]
	return i
}

type rfile struct {
	rmp map[string]string
	c int
	stck []string
	used []string
}
func (r *rfile) register(s string) string {
	sz := len(r.stck)
	if sz>0 {
		r.rmp[s] = r.stck[sz-1]
		r.stck = r.stck[:sz-1]
	} else {
		r.rmp[s] = fmt.Sprintf("r%d",r.c)
		r.c++
		r.used = append(r.used,r.rmp[s])
	}
	return r.rmp[s]
}
func (r *rfile) remove(s string) {
	l,ok := r.rmp[s]
	if !ok { return }
	r.stck = append(r.stck,l)
}

func Compile(sts []Statement,d io.Writer){
	B := new(bytes.Buffer)
	S := Statement{}
	ct := make(ctable)
	ct2 := make(ctable)
	cs := make(cset)
	inln := make(map[string]string)
	rf := rfile{rmp:make(map[string]string)}
	sub := func(s string)string {
		if s=="$" { return s }
		i,_ := strconv.Atoi(s)
		if i==0 { return "nil" }
		i--
		if i>=len(S.Args) { return "nil" }
		a := S.Args[i]
		r,ok := rf.rmp[a]
		if ok { return r }
		r,ok = inln[a]
		if ok { return r }
		return "nil"
	}
	
	for _,s := range sts {
		for _,arg := range s.Args { ct.incr(arg); ct2.incr(arg) }
		if len(s.Dest)>1 {
			for _,dst := range s.Dest{ cs.set(dst) }
		}
	}
	for _,s := range sts {
		S = s
		val := os.Expand(s.Expr,sub)
		for _,arg := range s.Args {
			if ct.decr(arg)==0 {
				rf.remove(arg)
			}
		}
		if s.Jump!="" {
			if val!="" { fmt.Fprintf(B,"if (%s) then goto m_%s_m end\n",val,s.Jump) }
			fmt.Fprintf(B,"goto m_%s_m\n",s.Jump)
		} else if s.Mark!="" {
			fmt.Fprintf(B,"::m_%s_m::\n",s.Mark)
		} else if len(s.Dest)==1 && ct2.get(s.Dest[0])==1 {
			inln[s.Dest[0]] = val
		} else {
			for i,dst := range s.Dest {
				dst2 := rf.register(dst)
				if i==0 {
					fmt.Fprintf(B,"%s",dst2)
				} else {
					fmt.Fprintf(B,", %s",dst2)
				}
			}
			if len(s.Dest)==0 {
				fmt.Fprintf(B,"%s\n",val)
			} else {
				fmt.Fprintf(B," = %s\n",val)
			}
		}
	}
	for _,r := range rf.used {
		fmt.Fprintf(d,"local %s\n",r)
	}
	B.WriteTo(d)
}


