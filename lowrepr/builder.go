/*
      Copyright(C) 2015 Simon Schmidt

      This Source Code Form is subject to the
      terms of the Mozilla Public License, v.
      2.0. If a copy of the MPL was not
      distributed with this file, You can
      obtain one at
      http://mozilla.org/MPL/2.0/.
*/

package lowrepr

import "errors"
import "fmt"

type refcReg struct{
	Spec ValueSpec
	Rc int64
}
func (r *refcReg)Grab(){ r.Rc++ }
func (r *refcReg)Drop(){ r.Rc-- }
func (r *refcReg)wrap() rcSpec {
	return rcSpec{r,r.Spec}
}

type noop struct{}
func (r *noop)Grab(){ }
func (r *noop)Drop(){ }

type rcSpec struct {
	Refc
	Spec ValueSpec
}

type Refc interface{
	Grab()
	Drop()
}

type Builder struct{
	user Helper
	registers map[ValueType][]*refcReg
	variables map[string]ValueSpec
	names map[string]bool
	nc int
	f *Function
}
func (b *Builder) Init(h Helper){
	b.user = h
	b.f = new(Function)
	b.registers = make(map[ValueType][]*refcReg)
	b.variables = make(map[string]ValueSpec)
	b.names = make(map[string]bool)
}
func (b *Builder) DefineArg(n string, t ValueType) error {
	if _,ok := b.variables[n]; ok {
		return errors.New("Variable '"+n+"' already exists")
	}
	if _,ok := b.names[n]; ok {
		return errors.New("The name '"+n+"' already used")
	}
	b.names[n] = true
	vs := t.VSpec(n)
	b.variables[n] = vs
	b.f.Params = append(b.f.Params,vs)
	return nil
}
func (b *Builder) DefineVar(n string, t ValueType) error {
	if _,ok := b.variables[n]; ok {
		return errors.New("Variable '"+n+"' already exists")
	}
	if _,ok := b.names[n]; ok {
		return errors.New("The name '"+n+"' already used")
	}
	b.names[n] = true
	vs := t.VSpec(n)
	b.variables[n] = vs
	b.f.Locals = append(b.f.Locals,vs)
	return nil
}
func (b *Builder) GetVar(n string) Refc {
	if r,ok := b.variables[n]; ok { return rcSpec{&noop{},r} }
	return nil
}
func (b *Builder) getName() string {
	for{
		n := fmt.Sprintf("r%d",b.nc)
		if _,ok := b.names[n]; !ok { return n }
		b.nc++
	}
	panic("unreachable")
}
func (b *Builder) getReg(dt ValueType) Refc {
	rt := b.user.GprType(dt)
	if r,ok := b.registers[rt]; ok {
		for _,reg := range r {
			if reg.Rc<=0 {
				reg.Rc = 1
				return reg.wrap()
			}
		}
		name := b.getName()
		b.names[name] = false
		res := &refcReg{rt.VSpec(name),1}
		b.registers[rt] = append(r,res)
		return res.wrap()
	}
	name := b.getName()
	b.names[name] = false
	b.registers[rt] = []*refcReg{
		&refcReg{rt.VSpec(name),1},
	}
	return b.registers[rt][0].wrap()
}
func (b *Builder) Perform(op int,v interface{},args []Refc) ([]Refc,error){
	return b.perform(op,v,args,false,nil)
}
func (b *Builder) PerformEx(op int,v interface{},args []Refc,res []Refc) ([]Refc,error){
	return b.perform(op,v,args,true,res)
}
func (b *Builder) perform(oc int,v interface{},args []Refc,dest bool,res []Refc) ([]Refc,error){
	//rcSpec
	opD,err := b.user.Op(oc,v)
	if err!=nil { return nil,err }
	vt := make([]ValueType,len(args))
	for i,arg := range args { vt[i] = arg.(rcSpec).Spec.VType() }
	op,vtn,err := opD(vt)
	if err!=nil { return nil,err }
	oplist := []Op{op}
	oplist[0].In = make([]ValueSpec,len(vt))
	for i,arg := range args { oplist[0].In[i] = arg.(rcSpec).Spec }
	
	oplist[0].Out = make([]ValueSpec,len(vtn))
	result := make([]Refc,len(vtn))
	if !dest {
		for i,v := range vtn {
			reg := b.getReg(v).(rcSpec)
			reg.Spec.SetVType(v)
			oplist[0].Out[i] = reg.Spec
			result[i] = reg
		}
	}else{
		copy(result,res)
		for i,v := range vtn {
			conv,err := b.user.Check(res[i].(rcSpec).Spec.VType(),v)
			if err!=nil { return nil,err }
			if conv!=nil {
				co,_,err := conv([]ValueType{v})
				if err!=nil { return nil,err }
				reg := b.getReg(v)
				co.In[0] =reg.(rcSpec).Spec
				co.Out[0] = res[i].(rcSpec).Spec
				oplist[0].Out[i] = reg.(rcSpec).Spec
				oplist = append(oplist,co)
				reg.Drop()
			}else{
				oplist[0].Out[i] = res[i].(rcSpec).Spec
			}
		}
	}
	for _,a := range args { a.Drop() }
	
	b.f.Ops = append(b.f.Ops,oplist...)
	return result,nil
}
func (b *Builder) Function() *Function { return b.f }

